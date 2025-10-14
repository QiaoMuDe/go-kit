package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/go-kit/fs"
)

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

// BenchmarkFormatBytes_AllSizes 测试不同大小的格式化性能
func BenchmarkFormatBytes_AllSizes(b *testing.B) {
	sizes := []struct {
		name  string
		bytes int64
	}{
		{"Bytes", 512},
		{"KB", 1024 * 512},
		{"MB", 1024 * 1024 * 512},
		{"GB", 1024 * 1024 * 1024 * 2},
		{"TB", 1024 * 1024 * 1024 * 1024 * 2},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				FormatBytes(size.bytes)
			}
		})
	}
}

// BenchmarkGetSize_FileVsDirectory 比较文件和目录的性能
func BenchmarkGetSize_FileVsDirectory(b *testing.B) {
	tempDir := b.TempDir()

	// 创建单个大文件
	largeFile := filepath.Join(tempDir, "large.txt")
	content := strings.Repeat("a", 1024*1024) // 1MB
	err := os.WriteFile(largeFile, []byte(content), 0644)
	if err != nil {
		b.Fatal(err)
	}

	// 创建包含多个小文件的目录
	manyFilesDir := filepath.Join(tempDir, "manyfiles")
	err = os.MkdirAll(manyFilesDir, 0755)
	if err != nil {
		b.Fatal(err)
	}

	smallContent := strings.Repeat("b", 1024) // 1KB
	for i := 0; i < 1000; i++ {
		smallFile := filepath.Join(manyFilesDir, fmt.Sprintf("file%d.txt", i))
		err := os.WriteFile(smallFile, []byte(smallContent), 0644)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.Run("SingleLargeFile", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := fs.GetSize(largeFile)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ManySmallFiles", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := fs.GetSize(manyFilesDir)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkFormatWithUnit_Comparison 比较不同单位转换的性能
func BenchmarkFormatWithUnit_Comparison(b *testing.B) {
	testCases := []struct {
		name      string
		bytes     int64
		divisor   int64
		unitIndex int
	}{
		{"KB", 1536, 1024, 0},
		{"MB", 1572864, 1048576, 1},
		{"GB", 1610612736, 1073741824, 2},
		{"TB", 1649267441664, 1099511627776, 3},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				formatWithUnit(tc.bytes, tc.divisor, tc.unitIndex)
			}
		})
	}
}
