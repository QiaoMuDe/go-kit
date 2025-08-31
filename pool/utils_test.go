package pool

import (
	"testing"
)

func TestByteConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant int64
		expected int64
	}{
		{"Byte", Byte, 1},
		{"KB", KB, 1024},
		{"MB", MB, 1024 * 1024},
		{"GB", GB, 1024 * 1024 * 1024},
		{"TB", TB, 1024 * 1024 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("%s = %d, expected %d", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestCalculateBufferSize(t *testing.T) {
	tests := []struct {
		name     string
		fileSize int64
		expected int
	}{
		// 极小文件测试 (≤ 4KB)
		{"Empty file", 0, 0},
		{"1 byte", 1, 1},
		{"100 bytes", 100, 100},
		{"1KB", 1 * KB, int(1 * KB)},
		{"4KB exactly", 4 * KB, int(4 * KB)},

		// 小文件测试 (4KB - 32KB)
		{"4KB + 1", 4*KB + 1, int(8 * KB)},
		{"8KB", 8 * KB, int(8 * KB)},
		{"16KB", 16 * KB, int(8 * KB)},
		{"31KB", 31 * KB, int(8 * KB)},

		// 中小文件测试 (32KB - 128KB)
		{"32KB", 32 * KB, int(32 * KB)},
		{"64KB", 64 * KB, int(32 * KB)},
		{"100KB", 100 * KB, int(32 * KB)},
		{"127KB", 127 * KB, int(32 * KB)},

		// 中等文件测试 (128KB - 512KB)
		{"128KB", 128 * KB, int(64 * KB)},
		{"256KB", 256 * KB, int(64 * KB)},
		{"400KB", 400 * KB, int(64 * KB)},
		{"511KB", 511 * KB, int(64 * KB)},

		// 中大文件测试 (512KB - 1MB)
		{"512KB", 512 * KB, int(128 * KB)},
		{"768KB", 768 * KB, int(128 * KB)},
		{"1MB - 1", 1*MB - 1, int(128 * KB)},

		// 大文件测试 (1MB - 4MB)
		{"1MB", 1 * MB, int(256 * KB)},
		{"2MB", 2 * MB, int(256 * KB)},
		{"3MB", 3 * MB, int(256 * KB)},
		{"4MB - 1", 4*MB - 1, int(256 * KB)},

		// 较大文件测试 (4MB - 16MB)
		{"4MB", 4 * MB, int(512 * KB)},
		{"8MB", 8 * MB, int(512 * KB)},
		{"12MB", 12 * MB, int(512 * KB)},
		{"16MB - 1", 16*MB - 1, int(512 * KB)},

		// 大文件测试 (16MB - 64MB)
		{"16MB", 16 * MB, int(1 * MB)},
		{"32MB", 32 * MB, int(1 * MB)},
		{"48MB", 48 * MB, int(1 * MB)},
		{"64MB - 1", 64*MB - 1, int(1 * MB)},

		// 超大文件测试 (> 64MB)
		{"64MB", 64 * MB, int(2 * MB)},
		{"100MB", 100 * MB, int(2 * MB)},
		{"1GB", 1 * GB, int(2 * MB)},
		{"10GB", 10 * GB, int(2 * MB)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateBufferSize(tt.fileSize)
			if result != tt.expected {
				t.Errorf("CalculateBufferSize(%d) = %d, expected %d", tt.fileSize, result, tt.expected)
			}
		})
	}
}

func TestCalculateBufferSize_EdgeCases(t *testing.T) {
	// 测试负数输入
	t.Run("Negative file size", func(t *testing.T) {
		result := CalculateBufferSize(-1)
		// 负数会进入第一个case，返回负数转换为int
		if result != -1 {
			t.Errorf("CalculateBufferSize(-1) = %d, expected -1", result)
		}
	})

	// 测试边界值
	boundaryTests := []struct {
		name     string
		fileSize int64
		expected int
	}{
		{"Exactly 4KB", 4 * KB, int(4 * KB)},
		{"4KB + 1 byte", 4*KB + 1, int(8 * KB)},
		{"Exactly 32KB", 32 * KB, int(32 * KB)},
		{"32KB + 1 byte", 32*KB + 1, int(32 * KB)},
		{"Exactly 128KB", 128 * KB, int(64 * KB)},
		{"128KB + 1 byte", 128*KB + 1, int(64 * KB)},
	}

	for _, tt := range boundaryTests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateBufferSize(tt.fileSize)
			if result != tt.expected {
				t.Errorf("CalculateBufferSize(%d) = %d, expected %d", tt.fileSize, result, tt.expected)
			}
		})
	}
}

func TestCalculateBufferSize_Performance(t *testing.T) {
	// 性能测试：确保函数能快速处理各种大小的输入
	testSizes := []int64{
		0, 1, 100, 1 * KB, 4 * KB, 8 * KB, 32 * KB, 64 * KB, 128 * KB,
		256 * KB, 512 * KB, 1 * MB, 4 * MB, 16 * MB, 64 * MB, 100 * MB, 1 * GB,
	}

	for _, size := range testSizes {
		result := CalculateBufferSize(size)
		if result < 0 && size >= 0 {
			t.Errorf("CalculateBufferSize(%d) returned negative result: %d", size, result)
		}
	}
}

func BenchmarkCalculateBufferSize(b *testing.B) {
	testSizes := []int64{
		1, 1 * KB, 32 * KB, 128 * KB, 1 * MB, 16 * MB, 64 * MB, 1 * GB,
	}

	for _, size := range testSizes {
		b.Run(formatSize(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				CalculateBufferSize(size)
			}
		})
	}
}

// 辅助函数：格式化文件大小用于基准测试名称
func formatSize(size int64) string {
	switch {
	case size >= GB:
		return "1GB+"
	case size >= MB:
		return "MB_range"
	case size >= KB:
		return "KB_range"
	default:
		return "bytes"
	}
}

func TestCalculateBufferSize_Consistency(t *testing.T) {
	// 测试一致性：相同输入应该总是返回相同结果
	testSize := int64(1 * MB)
	expected := CalculateBufferSize(testSize)

	for i := 0; i < 100; i++ {
		result := CalculateBufferSize(testSize)
		if result != expected {
			t.Errorf("Inconsistent result at iteration %d: got %d, expected %d", i, result, expected)
		}
	}
}

func TestCalculateBufferSize_Rationale(t *testing.T) {
	// 测试缓冲区大小的合理性
	t.Run("Buffer size should not exceed file size for small files", func(t *testing.T) {
		smallSizes := []int64{1, 10, 100, 1 * KB, 2 * KB, 4 * KB}
		for _, size := range smallSizes {
			bufferSize := CalculateBufferSize(size)
			if int64(bufferSize) > size && size > 0 {
				t.Errorf("Buffer size %d exceeds file size %d", bufferSize, size)
			}
		}
	})

	t.Run("Buffer size should be reasonable for large files", func(t *testing.T) {
		largeSize := int64(100 * MB)
		bufferSize := CalculateBufferSize(largeSize)
		maxReasonableBuffer := int(2 * MB)
		if bufferSize > maxReasonableBuffer {
			t.Errorf("Buffer size %d is too large for file size %d", bufferSize, largeSize)
		}
	})

	t.Run("Buffer size should increase with file size", func(t *testing.T) {
		sizes := []int64{8 * KB, 64 * KB, 256 * KB, 2 * MB, 32 * MB}
		var prevBuffer int
		for i, size := range sizes {
			bufferSize := CalculateBufferSize(size)
			if i > 0 && bufferSize < prevBuffer {
				t.Errorf("Buffer size should not decrease: size %d has buffer %d, but previous had %d",
					size, bufferSize, prevBuffer)
			}
			prevBuffer = bufferSize
		}
	})
}
