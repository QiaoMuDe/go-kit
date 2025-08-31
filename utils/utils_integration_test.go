package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gitee.com/MM-Q/go-kit/fs"
)

// TestIntegration_CommandExecutionAndSizeCheck 集成测试：命令执行和大小检查
func TestIntegration_CommandExecutionAndSizeCheck(t *testing.T) {
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.txt")

	// 使用命令创建文件
	content := "integration test content"
	args := getWriteFileCommand(content, outputFile)

	_, err := ExecuteCmd(args, nil)
	if err != nil {
		t.Skipf("跳过集成测试，命令执行失败: %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("命令执行后文件不存在")
	}

	// 获取文件大小
	size, err := fs.GetSize(outputFile)
	if err != nil {
		t.Errorf("获取文件大小失败: %v", err)
	}

	// 验证大小合理性（考虑换行符等）
	if size == 0 {
		t.Error("文件大小为0，可能命令执行失败")
	}

	// 格式化并验证
	formatted := FormatBytes(size)
	if formatted == "" {
		t.Error("格式化结果为空")
	}

	t.Logf("创建的文件大小: %s", formatted)
}

// TestIntegration_TimeoutAndErrorHandling 集成测试：超时和错误处理
func TestIntegration_TimeoutAndErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		timeout       time.Duration
		args          []string
		expectError   bool
		expectTimeout bool
	}{
		{
			name:          "快速成功命令",
			timeout:       time.Second * 5,
			args:          getEchoCommand("success"),
			expectError:   false,
			expectTimeout: false,
		},
		{
			name:          "超时命令",
			timeout:       time.Millisecond * 50,
			args:          getSleepCommand("1"),
			expectError:   true,
			expectTimeout: true,
		},
		{
			name:          "不存在的命令",
			timeout:       time.Second,
			args:          []string{"nonexistentcommand123"},
			expectError:   true,
			expectTimeout: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			output, err := ExecuteCmdWithTimeout(tt.timeout, tt.args, nil)
			duration := time.Since(start)

			if tt.expectError {
				if err == nil {
					t.Error("期望错误但没有返回错误")
				}
				if tt.expectTimeout {
					if duration > tt.timeout*2 {
						t.Errorf("超时处理不及时，期望约%v，实际%v", tt.timeout, duration)
					}
					if !strings.Contains(err.Error(), "超时") {
						t.Errorf("期望超时错误，但得到: %v", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("意外错误: %v", err)
				}
				if len(output) == 0 {
					t.Error("期望有输出但输出为空")
				}
			}
		})
	}
}

// TestIntegration_DirectoryTraversal 集成测试：目录遍历和大小计算
func TestIntegration_DirectoryTraversal(t *testing.T) {
	tempDir := t.TempDir()

	// 创建复杂的目录结构
	structure := map[string]string{
		"file1.txt":             "content1",
		"subdir1/file2.txt":     "content2",
		"subdir1/file3.txt":     "content3",
		"subdir2/file4.txt":     "content4",
		"subdir2/sub/file5.txt": "content5",
	}

	var expectedTotalSize int64
	for path, content := range structure {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)

		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatal(err)
		}

		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}

		expectedTotalSize += int64(len(content))
	}

	// 测试整个目录的大小
	totalSize, err := fs.GetSize(tempDir)
	if err != nil {
		t.Errorf("获取目录大小失败: %v", err)
	}

	if totalSize != expectedTotalSize {
		t.Errorf("目录总大小不匹配: 得到 %d, 期望 %d", totalSize, expectedTotalSize)
	}

	// 格式化总大小
	formatted := FormatBytes(totalSize)
	t.Logf("目录总大小: %s", formatted)

	// 测试子目录大小
	subdir1Size, err := fs.GetSize(filepath.Join(tempDir, "subdir1"))
	if err != nil {
		t.Errorf("获取子目录大小失败: %v", err)
	}

	expectedSubdir1Size := int64(len("content2") + len("content3"))
	if subdir1Size != expectedSubdir1Size {
		t.Errorf("子目录大小不匹配: 得到 %d, 期望 %d", subdir1Size, expectedSubdir1Size)
	}
}

// TestIntegration_ErrorPropagation 集成测试：错误传播
func TestIntegration_ErrorPropagation(t *testing.T) {
	// 测试不存在路径的错误处理
	nonExistentPath := "/absolutely/nonexistent/path/file.txt"

	size, err := fs.GetSize(nonExistentPath)
	if err == nil {
		t.Error("期望错误但没有返回错误")
	}

	if size != 0 {
		t.Errorf("错误情况下大小应为0，但得到 %d", size)
	}

	// 验证错误信息包含路径
	if !strings.Contains(err.Error(), nonExistentPath) {
		t.Errorf("错误信息应包含路径，但得到: %s", err.Error())
	}

	// 测试格式化0字节
	formatted := FormatBytes(size)
	if formatted != "0 B" {
		t.Errorf("0字节格式化应为 '0 B'，但得到 '%s'", formatted)
	}
}

// TestIntegration_ConcurrentAccess 集成测试：并发访问
func TestIntegration_ConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "concurrent_test.txt")
	content := strings.Repeat("concurrent", 1000)
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 并发测试
	const numGoroutines = 10
	results := make(chan struct {
		size      int64
		formatted string
		err       error
	}, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			size, err := fs.GetSize(testFile)
			formatted := FormatBytes(size)
			results <- struct {
				size      int64
				formatted string
				err       error
			}{size, formatted, err}
		}()
	}

	// 收集结果
	expectedSize := int64(len(content))
	for i := 0; i < numGoroutines; i++ {
		result := <-results

		if result.err != nil {
			t.Errorf("并发访问错误: %v", result.err)
		}

		if result.size != expectedSize {
			t.Errorf("并发访问大小不一致: 得到 %d, 期望 %d", result.size, expectedSize)
		}

		if result.formatted == "" {
			t.Error("并发访问格式化结果为空")
		}
	}
}
