package pool

import (
	"bytes"
	"testing"
)

func TestWithByte(t *testing.T) {
	t.Run("Basic usage", func(t *testing.T) {
		data := WithByte(func(buf []byte) {
			copy(buf, []byte("Hello World"))
		})

		expected := []byte("Hello World")
		// 只比较有效数据部分
		if !bytes.Equal(data[:len(expected)], expected) {
			t.Errorf("Expected %q, got %q", expected, data[:len(expected)])
		}
	})

	t.Run("Binary data", func(t *testing.T) {
		testData := []byte{0x01, 0x02, 0x03, 0x04, 0xFF}
		data := WithByte(func(buf []byte) {
			copy(buf, testData)
		})

		if !bytes.Equal(data[:len(testData)], testData) {
			t.Errorf("Expected %v, got %v", testData, data[:len(testData)])
		}
	})

	t.Run("Large buffer", func(t *testing.T) {
		size := 1024
		data := WithByteCapacity(size, func(buf []byte) {
			for i := 0; i < size; i++ {
				buf[i] = byte(i % 256)
			}
		})

		if len(data) != size {
			t.Errorf("Expected length %d, got %d", size, len(data))
		}

		// 验证数据正确性
		for i := 0; i < size; i++ {
			if data[i] != byte(i%256) {
				t.Errorf("Data mismatch at index %d: expected %d, got %d", i, i%256, data[i])
				break
			}
		}
	})
}

func TestWithEmptyByte(t *testing.T) {
	t.Run("Append operations", func(t *testing.T) {
		data := WithEmptyByte(64, func(buf []byte) []byte {
			buf = append(buf, []byte("Hello")...)
			buf = append(buf, ' ')
			buf = append(buf, []byte("World")...)
			return buf
		})

		expected := []byte("Hello World")
		if !bytes.Equal(data, expected) {
			t.Errorf("Expected %q, got %q", expected, data)
		}
	})

	t.Run("Build binary data", func(t *testing.T) {
		data := WithEmptyByte(32, func(buf []byte) []byte {
			for i := 0; i < 10; i++ {
				buf = append(buf, byte(i))
			}
			return buf
		})

		expected := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		if !bytes.Equal(data, expected) {
			t.Errorf("Expected %v, got %v", expected, data)
		}
	})

	t.Run("Empty result", func(t *testing.T) {
		data := WithEmptyByte(16, func(buf []byte) []byte {
			// 不添加任何数据
			return buf
		})

		if len(data) != 0 {
			t.Errorf("Expected empty slice, got length %d", len(data))
		}
	})

	t.Run("Complex operations", func(t *testing.T) {
		data := WithEmptyByte(128, func(buf []byte) []byte {
			// 构建一个简单的协议包
			buf = append(buf, 0x01, 0x02)                // 头部
			buf = append(buf, []byte("test message")...) // 数据
			buf = append(buf, 0xFF)                      // 结尾标记
			return buf
		})

		expected := []byte{0x01, 0x02}
		expected = append(expected, []byte("test message")...)
		expected = append(expected, 0xFF)

		if !bytes.Equal(data, expected) {
			t.Errorf("Expected %v, got %v", expected, data)
		}
	})
}

func TestBytePoolWithMethods(t *testing.T) {
	pool := NewBytePool(64, 1024)

	t.Run("WithByte method", func(t *testing.T) {
		data := pool.WithByte(func(buf []byte) {
			copy(buf, []byte("test"))
		})

		expected := []byte("test")
		if !bytes.Equal(data[:len(expected)], expected) {
			t.Errorf("Expected %q, got %q", expected, data[:len(expected)])
		}
	})

	t.Run("WithEmptyByte method", func(t *testing.T) {
		data := pool.WithEmptyByte(32, func(buf []byte) []byte {
			buf = append(buf, []byte("custom pool")...)
			return buf
		})

		expected := []byte("custom pool")
		if !bytes.Equal(data, expected) {
			t.Errorf("Expected %q, got %q", expected, data)
		}
	})
}

// 基准测试对比传统方式和新方式
func BenchmarkTraditionalByte(b *testing.B) {
	testData := []byte("Hello World Test Data")

	for i := 0; i < b.N; i++ {
		buf := GetByte()
		copy(buf, testData)
		result := make([]byte, len(testData))
		copy(result, buf[:len(testData)])
		PutByte(buf)
		_ = result
	}
}

func BenchmarkWithByte(b *testing.B) {
	testData := []byte("Hello World Test Data")

	for i := 0; i < b.N; i++ {
		result := WithByte(func(buf []byte) {
			copy(buf, testData)
		})
		_ = result
	}
}

func BenchmarkTraditionalEmptyByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := GetEmptyByte(64)
		buf = append(buf, []byte("Hello")...)
		buf = append(buf, ' ')
		buf = append(buf, []byte("World")...)
		result := make([]byte, len(buf))
		copy(result, buf)
		PutByte(buf)
		_ = result
	}
}

func BenchmarkWithEmptyByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := WithEmptyByte(64, func(buf []byte) []byte {
			buf = append(buf, []byte("Hello")...)
			buf = append(buf, ' ')
			buf = append(buf, []byte("World")...)
			return buf
		})
		_ = result
	}
}

// 并发安全测试
func TestWithByteConcurrent(t *testing.T) {
	const numGoroutines = 10
	const numOperations = 100

	results := make(chan []byte, numGoroutines*numOperations)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numOperations; j++ {
				data := WithByte(func(buf []byte) {
					testData := []byte("goroutine-" + string(rune('0'+id)))
					copy(buf, testData)
				})
				results <- data
			}
		}(i)
	}

	// 收集所有结果
	for i := 0; i < numGoroutines*numOperations; i++ {
		result := <-results
		if len(result) == 0 {
			t.Error("Got empty result from concurrent operation")
		}
	}
}

func TestWithEmptyByteConcurrent(t *testing.T) {
	const numGoroutines = 10
	const numOperations = 100

	results := make(chan []byte, numGoroutines*numOperations)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numOperations; j++ {
				data := WithEmptyByte(32, func(buf []byte) []byte {
					buf = append(buf, []byte("test-")...)
					buf = append(buf, byte('0'+id))
					return buf
				})
				results <- data
			}
		}(i)
	}

	// 收集所有结果
	for i := 0; i < numGoroutines*numOperations; i++ {
		result := <-results
		if len(result) == 0 {
			t.Error("Got empty result from concurrent operation")
		}
	}
}
