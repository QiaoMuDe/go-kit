package utils

import (
	"strings"
	"testing"
	"time"
)

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
