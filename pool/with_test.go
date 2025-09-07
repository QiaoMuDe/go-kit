package pool

import (
	"bytes"
	"strings"
	"testing"
)

func TestWithString(t *testing.T) {
	t.Run("Basic usage", func(t *testing.T) {
		result := WithStr(func(buf *strings.Builder) {
			buf.WriteString("Hello")
			buf.WriteByte(' ')
			buf.WriteString("World")
		})

		expected := "Hello World"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("Empty function", func(t *testing.T) {
		result := WithStr(func(buf *strings.Builder) {
			// 不写入任何内容
		})

		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})

	t.Run("Large content", func(t *testing.T) {
		result := WithStr(func(buf *strings.Builder) {
			for i := 0; i < 100; i++ {
				buf.WriteString("test")
			}
		})

		expected := strings.Repeat("test", 100)
		if result != expected {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
	})
}

func TestWithBuffer(t *testing.T) {
	t.Run("Basic usage", func(t *testing.T) {
		result := WithBuf(func(buf *bytes.Buffer) {
			buf.WriteString("Hello")
			buf.WriteByte(' ')
			buf.WriteString("World")
		})

		expected := []byte("Hello World")
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("Binary data", func(t *testing.T) {
		result := WithBuf(func(buf *bytes.Buffer) {
			buf.Write([]byte{0x01, 0x02, 0x03, 0x04})
		})

		expected := []byte{0x01, 0x02, 0x03, 0x04}
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Empty buffer", func(t *testing.T) {
		result := WithBuf(func(buf *bytes.Buffer) {
			// 不写入任何内容
		})

		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})
}

// 基准测试对比传统方式和新方式
func BenchmarkTraditionalString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := GetStr()
		buf.WriteString("Hello")
		buf.WriteByte(' ')
		buf.WriteString("World")
		result := buf.String()
		PutStr(buf)
		_ = result
	}
}

func BenchmarkWithString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := WithStr(func(buf *strings.Builder) {
			buf.WriteString("Hello")
			buf.WriteByte(' ')
			buf.WriteString("World")
		})
		_ = result
	}
}

func BenchmarkTraditionalBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := GetBuf()
		buf.WriteString("Hello")
		buf.WriteByte(' ')
		buf.WriteString("World")
		result := make([]byte, buf.Len())
		copy(result, buf.Bytes())
		PutBuf(buf)
		_ = result
	}
}

func BenchmarkWithBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := WithBuf(func(buf *bytes.Buffer) {
			buf.WriteString("Hello")
			buf.WriteByte(' ')
			buf.WriteString("World")
		})
		_ = result
	}
}

// 并发安全测试
func TestWithStringConcurrent(t *testing.T) {
	const numGoroutines = 10
	const numOperations = 100

	results := make(chan string, numGoroutines*numOperations)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numOperations; j++ {
				result := WithStr(func(buf *strings.Builder) {
					buf.WriteString("goroutine-")
					buf.WriteString(string(rune('0' + id)))
					buf.WriteString("-op-")
					buf.WriteString(string(rune('0' + j%10)))
				})
				results <- result
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
