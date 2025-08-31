package pool

import (
	"bytes"
	"sync"
	"testing"
)

func TestBufferPool_Get(t *testing.T) {
	buf := GetBuffer(1024)
	if buf == nil {
		t.Fatal("GetBuffer() returned nil")
	}

	// 验证返回的是一个空的buffer
	if buf.Len() != 0 {
		t.Errorf("Expected empty buffer, got length %d", buf.Len())
	}

	PutBuffer(buf)
}

func TestBufferPool_Put(t *testing.T) {
	buf := GetBuffer(1024)
	buf.WriteString("test data")

	// Put应该重置buffer
	PutBuffer(buf)

	// 再次获取应该是空的
	buf2 := GetBuffer(1024)
	if buf2.Len() != 0 {
		t.Errorf("Expected empty buffer after put, got length %d", buf2.Len())
	}

	PutBuffer(buf2)
}

func TestBufferPool_Reuse(t *testing.T) {
	// 测试对象池的复用机制
	buf1 := GetBuffer(1024)
	buf1.WriteString("test")
	PutBuffer(buf1)

	buf2 := GetBuffer(1024)
	// 应该复用同一个对象（在单线程情况下）
	if buf1 != buf2 {
		t.Log("Buffer objects are different (this is acceptable in concurrent scenarios)")
	}

	PutBuffer(buf2)
}

func TestBufferPool_Concurrent(t *testing.T) {
	const numGoroutines = 100
	const numOperations = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				buf := GetBuffer(1024)
				if buf == nil {
					t.Errorf("GetBuffer() returned nil in goroutine %d", id)
					return
				}

				// 写入一些数据
				buf.WriteString("goroutine")
				buf.WriteString(string(rune('0' + id%10)))

				// 验证数据
				if buf.Len() == 0 {
					t.Errorf("Buffer is empty after writing in goroutine %d", id)
					return
				}

				PutBuffer(buf)
			}
		}(i)
	}

	wg.Wait()
}

func TestBufferPool_Reset(t *testing.T) {
	buf := GetBuffer(1024)
	buf.WriteString("some data")
	buf.WriteByte(0x00)

	if buf.Len() == 0 {
		t.Fatal("Buffer should contain data before reset")
	}

	PutBuffer(buf)

	// 获取新的buffer应该是空的
	newBuf := GetBuffer(1024)
	if newBuf.Len() != 0 {
		t.Errorf("Expected empty buffer after reset, got length %d", newBuf.Len())
	}

	PutBuffer(newBuf)
}

func TestBufferPool_LargeData(t *testing.T) {
	buf := GetBuffer(1024 * 1024)

	// 写入大量数据
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	buf.Write(largeData)

	if buf.Len() != len(largeData) {
		t.Errorf("Expected buffer length %d, got %d", len(largeData), buf.Len())
	}

	PutBuffer(buf)

	// 验证重置后是空的
	newBuf := GetBuffer(1024)
	if newBuf.Len() != 0 {
		t.Errorf("Expected empty buffer after putting large data, got length %d", newBuf.Len())
	}

	PutBuffer(newBuf)
}

func BenchmarkBufferPool_GetPut(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := GetBuffer(1024)
		buf.WriteString("benchmark test data")
		PutBuffer(buf)
	}
}

func BenchmarkBufferPool_vs_New(b *testing.B) {
	b.Run("Pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := GetBuffer(1024)
			buf.WriteString("benchmark")
			PutBuffer(buf)
		}
	})

	b.Run("New", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := &bytes.Buffer{}
			buf.WriteString("benchmark")
		}
	})
}
