package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

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
			_, err := GetSize(largeFile)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ManySmallFiles", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := GetSize(manyFilesDir)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkExecuteCmd_Simple 测试简单命令执行性能
func BenchmarkExecuteCmd_Simple(b *testing.B) {
	args := getEchoCommand("hello")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ExecuteCmd(args, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkExecuteCmdWithTimeout_Various 测试不同超时设置的性能
func BenchmarkExecuteCmdWithTimeout_Various(b *testing.B) {
	timeouts := []struct {
		name    string
		timeout time.Duration
	}{
		{"1ms", time.Millisecond},
		{"10ms", 10 * time.Millisecond},
		{"100ms", 100 * time.Millisecond},
		{"1s", time.Second},
	}

	args := getEchoCommand("hello")

	for _, timeout := range timeouts {
		b.Run(timeout.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := ExecuteCmdWithTimeout(timeout.timeout, args, nil)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
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
