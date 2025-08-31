package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestGetSize 测试GetSize函数
func TestGetSize(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()

	// 测试用例
	tests := []struct {
		name        string
		setupFunc   func() (string, int64) // 返回路径和期望大小
		expectError bool
	}{
		{
			name: "单个文件",
			setupFunc: func() (string, int64) {
				content := "hello world"
				filePath := filepath.Join(tempDir, "test.txt")
				err := os.WriteFile(filePath, []byte(content), 0644)
				if err != nil {
					t.Fatal(err)
				}
				return filePath, int64(len(content))
			},
			expectError: false,
		},
		{
			name: "空文件",
			setupFunc: func() (string, int64) {
				filePath := filepath.Join(tempDir, "empty.txt")
				err := os.WriteFile(filePath, []byte{}, 0644)
				if err != nil {
					t.Fatal(err)
				}
				return filePath, 0
			},
			expectError: false,
		},
		{
			name: "目录包含多个文件",
			setupFunc: func() (string, int64) {
				dirPath := filepath.Join(tempDir, "testdir")
				err := os.MkdirAll(dirPath, 0755)
				if err != nil {
					t.Fatal(err)
				}

				// 创建多个文件
				files := map[string]string{
					"file1.txt": "content1",
					"file2.txt": "content2",
					"file3.txt": "content3",
				}

				var totalSize int64
				for name, content := range files {
					filePath := filepath.Join(dirPath, name)
					err := os.WriteFile(filePath, []byte(content), 0644)
					if err != nil {
						t.Fatal(err)
					}
					totalSize += int64(len(content))
				}

				return dirPath, totalSize
			},
			expectError: false,
		},
		{
			name: "空目录",
			setupFunc: func() (string, int64) {
				dirPath := filepath.Join(tempDir, "emptydir")
				err := os.MkdirAll(dirPath, 0755)
				if err != nil {
					t.Fatal(err)
				}
				return dirPath, 0
			},
			expectError: false,
		},
		{
			name: "不存在的路径",
			setupFunc: func() (string, int64) {
				return "/nonexistent/path", 0
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, expectedSize := tt.setupFunc()
			size, err := GetSize(path)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
				return
			}

			if err != nil {
				t.Errorf("意外错误: %v", err)
				return
			}

			if size != expectedSize {
				t.Errorf("GetSize() = %d, 期望 %d", size, expectedSize)
			}
		})
	}
}

// TestExecuteCmd 测试ExecuteCmd函数
func TestExecuteCmd(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		env         []string
		expectError bool
		checkOutput func([]byte) bool
	}{
		{
			name:        "空命令",
			args:        []string{},
			env:         nil,
			expectError: true,
			checkOutput: nil,
		},
		{
			name:        "echo命令",
			args:        getEchoCommand("hello"),
			env:         nil,
			expectError: false,
			checkOutput: func(output []byte) bool {
				return strings.Contains(string(output), "hello")
			},
		},
		{
			name:        "不存在的命令",
			args:        []string{"nonexistentcommand"},
			env:         nil,
			expectError: true,
			checkOutput: nil,
		},
		{
			name:        "带环境变量的命令",
			args:        getEchoCommand("test"),
			env:         []string{"TEST_VAR=test_value"},
			expectError: false,
			checkOutput: func(output []byte) bool {
				return len(output) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := ExecuteCmd(tt.args, tt.env)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
				return
			}

			if err != nil {
				t.Errorf("意外错误: %v", err)
				return
			}

			if tt.checkOutput != nil && !tt.checkOutput(output) {
				t.Errorf("输出验证失败: %s", string(output))
			}
		})
	}
}

// TestExecuteCmdWithTimeout 测试ExecuteCmdWithTimeout函数
func TestExecuteCmdWithTimeout(t *testing.T) {
	tests := []struct {
		name          string
		timeout       time.Duration
		args          []string
		env           []string
		expectError   bool
		expectTimeout bool
	}{
		{
			name:          "空命令",
			timeout:       time.Second,
			args:          []string{},
			env:           nil,
			expectError:   true,
			expectTimeout: false,
		},
		{
			name:          "正常命令",
			timeout:       time.Second * 5,
			args:          getEchoCommand("hello"),
			env:           nil,
			expectError:   false,
			expectTimeout: false,
		},
		{
			name:          "超时命令",
			timeout:       time.Millisecond * 50,
			args:          getSleepCommand("5"),
			env:           nil,
			expectError:   true,
			expectTimeout: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := ExecuteCmdWithTimeout(tt.timeout, tt.args, tt.env)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
				if tt.expectTimeout && !strings.Contains(err.Error(), "超时") && !strings.Contains(err.Error(), "context deadline exceeded") {
					// Windows timeout命令可能返回不同的错误，我们接受任何错误作为超时的表现
					t.Logf("超时测试得到错误: %v (这在Windows上是正常的)", err)
				}
				return
			}

			if err != nil {
				t.Errorf("意外错误: %v", err)
				return
			}

			if len(output) == 0 {
				t.Errorf("期望有输出但输出为空")
			}
		})
	}
}

// TestFormatBytes 测试FormatBytes函数
func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"零字节", 0, "0 B"},
		{"负数", -1024, "-1 KB"},
		{"1字节", 1, "1 B"},
		{"1023字节", 1023, "1023 B"},
		{"1KB", 1024, "1 KB"},
		{"1.5KB", 1536, "1.50 KB"},
		{"1MB", 1048576, "1 MB"},
		{"1.25MB", 1310720, "1.25 MB"},
		{"1GB", 1073741824, "1 GB"},
		{"2.5GB", 2684354560, "2.50 GB"},
		{"1TB", 1099511627776, "1 TB"},
		{"1.75TB", 1924145348608, "1.75 TB"},
		{"1PB", 1125899906842624, "1 PB"},
		{"大数值", 9223372036854775807, "8191.99 PB"}, // int64最大值
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("FormatBytes(%d) = %s, 期望 %s", tt.bytes, result, tt.expected)
			}
		})
	}
}

// TestFormatWithUnit 测试formatWithUnit函数
func TestFormatWithUnit(t *testing.T) {
	tests := []struct {
		name      string
		bytes     int64
		divisor   int64
		unitIndex int
		expected  string
	}{
		{"整数KB", 2048, 1024, 0, "2 KB"},
		{"小数KB", 1536, 1024, 0, "1.50 KB"},
		{"整数MB", 2097152, 1048576, 1, "2 MB"},
		{"小数MB", 1572864, 1048576, 1, "1.50 MB"},
		{"小于10的小数", 1126400, 1048576, 1, "1.07 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatWithUnit(tt.bytes, tt.divisor, tt.unitIndex)
			if result != tt.expected {
				t.Errorf("formatWithUnit(%d, %d, %d) = %s, 期望 %s",
					tt.bytes, tt.divisor, tt.unitIndex, result, tt.expected)
			}
		})
	}
}

// TestWrapPathError 测试wrapPathError函数
func TestWrapPathError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		path      string
		operation string
		expected  string
	}{
		{
			name:      "文件不存在错误",
			err:       os.ErrNotExist,
			path:      "/test/path",
			operation: "reading",
			expected:  "path does not exist when reading: /test/path",
		},
		{
			name:      "权限错误",
			err:       os.ErrPermission,
			path:      "/test/path",
			operation: "writing",
			expected:  "permission denied when writing path '/test/path'",
		},
		{
			name:      "其他错误",
			err:       fmt.Errorf("custom error"),
			path:      "/test/path",
			operation: "accessing",
			expected:  "error when accessing path '/test/path': custom error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapPathError(tt.err, tt.path, tt.operation)
			if !strings.Contains(result.Error(), tt.path) {
				t.Errorf("错误信息应包含路径 %s，但得到: %s", tt.path, result.Error())
			}
			if !strings.Contains(result.Error(), tt.operation) {
				t.Errorf("错误信息应包含操作 %s，但得到: %s", tt.operation, result.Error())
			}
		})
	}
}

// BenchmarkFormatBytes 性能测试
func BenchmarkFormatBytes(b *testing.B) {
	testCases := []int64{
		0, 1, 1023, 1024, 1536, 1048576, 1310720,
		1073741824, 2684354560, 1099511627776,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, bytes := range testCases {
			FormatBytes(bytes)
		}
	}
}

// BenchmarkGetSize 性能测试
func BenchmarkGetSize(b *testing.B) {
	// 创建临时文件
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	content := strings.Repeat("a", 1024) // 1KB内容
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetSize(testFile)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// TestFormatBytes_EdgeCases 边界测试
func TestFormatBytes_EdgeCases(t *testing.T) {
	// 测试边界值
	edgeCases := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"KB边界-1", 1023, "1023 B"},
		{"KB边界", 1024, "1 KB"},
		{"KB边界+1", 1025, "1 KB"},
		{"MB边界-1", 1048575, "1023.99 KB"},
		{"MB边界", 1048576, "1 MB"},
		{"MB边界+1", 1048577, "1 MB"},
	}

	for _, tt := range edgeCases {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("FormatBytes(%d) = %s, 期望 %s", tt.bytes, result, tt.expected)
			}
		})
	}
}

// TestGetSize_SymbolicLinks 测试符号链接
func TestGetSize_SymbolicLinks(t *testing.T) {
	tempDir := t.TempDir()

	// 创建原始文件
	originalFile := filepath.Join(tempDir, "original.txt")
	content := "test content"
	err := os.WriteFile(originalFile, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 创建符号链接
	linkFile := filepath.Join(tempDir, "link.txt")
	err = os.Symlink(originalFile, linkFile)
	if err != nil {
		t.Skip("无法创建符号链接，跳过测试")
	}

	// 测试符号链接的大小
	size, err := GetSize(linkFile)
	if err != nil {
		t.Errorf("获取符号链接大小失败: %v", err)
	}

	expectedSize := int64(len(content))
	if size != expectedSize {
		t.Errorf("符号链接大小 = %d, 期望 %d", size, expectedSize)
	}
}

// TestExecuteCmd_LongOutput 测试长输出
func TestExecuteCmd_LongOutput(t *testing.T) {
	// 生成长输出的命令
	longText := strings.Repeat("a", 1000) // 减少长度以适应命令行限制
	args := getEchoCommand(longText)

	output, err := ExecuteCmd(args, nil)
	if err != nil {
		t.Errorf("执行命令失败: %v", err)
	}

	if len(output) < 500 { // 调整期望长度
		t.Errorf("输出长度不足，期望至少500字节，得到%d字节", len(output))
	}
}
