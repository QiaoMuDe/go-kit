package pool

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// 基准测试：对比使用对象池和不使用对象池的性能差异

// BenchmarkWithPool 使用对象池的基准测试
func BenchmarkWithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetBuf()
		buf.WriteString("Hello")
		buf.WriteByte(' ')
		buf.WriteString("World")
		fmt.Fprintf(buf, " %d", i)
		_ = buf.Bytes()
		PutBuf(buf)
	}
}

// BenchmarkWithoutPool 不使用对象池的基准测试
func BenchmarkWithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(make([]byte, 0, 256))
		buf.WriteString("Hello")
		buf.WriteByte(' ')
		buf.WriteString("World")
		fmt.Fprintf(buf, " %d", i)
		_ = buf.Bytes()
		// 不归还，让GC处理
	}
}

// BenchmarkWithPoolLarge 使用对象池处理大数据的基准测试
func BenchmarkWithPoolLarge(b *testing.B) {
	data := strings.Repeat("A", 1024) // 1KB数据
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetBufCap(2048)
		for j := 0; j < 10; j++ {
			buf.WriteString(data)
		}
		_ = buf.Bytes()
		PutBuf(buf)
	}
}

// BenchmarkWithoutPoolLarge 不使用对象池处理大数据的基准测试
func BenchmarkWithoutPoolLarge(b *testing.B) {
	data := strings.Repeat("A", 1024) // 1KB数据
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(make([]byte, 0, 2048))
		for j := 0; j < 10; j++ {
			buf.WriteString(data)
		}
		_ = buf.Bytes()
		// 不归还，让GC处理
	}
}

// BenchmarkWithFunction 使用With函数的基准测试
func BenchmarkWithFunction(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := WithBuf(func(buf *bytes.Buffer) {
			buf.WriteString("Hello")
			buf.WriteByte(' ')
			buf.WriteString("World")
			fmt.Fprintf(buf, " %d", i)
		})
		_ = data
	}
}

// BenchmarkWithCapFunction 使用WithCap函数的基准测试
func BenchmarkWithCapFunction(b *testing.B) {
	data := strings.Repeat("A", 1024) // 1KB数据
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := WithBufCap(2048, func(buf *bytes.Buffer) {
			for j := 0; j < 10; j++ {
				buf.WriteString(data)
			}
		})
		_ = result
	}
}

// 内存分配测试：对比内存分配次数和大小

// BenchmarkMemoryWithPool 测试使用对象池的内存分配
func BenchmarkMemoryWithPool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetBuf()
		buf.WriteString("Hello World")
		fmt.Fprintf(buf, " %d", i)
		_ = buf.Bytes()
		PutBuf(buf)
	}
}

// BenchmarkMemoryWithoutPool 测试不使用对象池的内存分配
func BenchmarkMemoryWithoutPool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(make([]byte, 0, 256))
		buf.WriteString("Hello World")
		fmt.Fprintf(buf, " %d", i)
		_ = buf.Bytes()
	}
}

// 并发测试：测试对象池在并发环境下的性能

// BenchmarkConcurrentWithPool 并发使用对象池的基准测试
func BenchmarkConcurrentWithPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := GetBuf()
			buf.WriteString("Hello")
			buf.WriteByte(' ')
			buf.WriteString("World")
			_ = buf.Bytes()
			PutBuf(buf)
		}
	})
}

// BenchmarkConcurrentWithoutPool 并发不使用对象池的基准测试
func BenchmarkConcurrentWithoutPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := bytes.NewBuffer(make([]byte, 0, 256))
			buf.WriteString("Hello")
			buf.WriteByte(' ')
			buf.WriteString("World")
			_ = buf.Bytes()
		}
	})
}

// GC压力测试：测试对象池对GC的影响

// TestGCPressure 测试GC压力
func TestGCPressure(t *testing.T) {
	// 测试使用对象池的GC压力
	t.Run("WithPool", func(t *testing.T) {
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		for i := 0; i < 10000; i++ {
			buf := GetBuf()
			buf.WriteString("Hello World")
			fmt.Fprintf(buf, " %d", i)
			_ = buf.Bytes()
			PutBuf(buf)
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		t.Logf("使用对象池 - 分配次数: %d, 总分配: %d bytes",
			m2.Mallocs-m1.Mallocs, m2.TotalAlloc-m1.TotalAlloc)
	})

	// 测试不使用对象池的GC压力
	t.Run("WithoutPool", func(t *testing.T) {
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		for i := 0; i < 10000; i++ {
			buf := bytes.NewBuffer(make([]byte, 0, 256))
			buf.WriteString("Hello World")
			fmt.Fprintf(buf, " %d", i)
			_ = buf.Bytes()
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		t.Logf("不使用对象池 - 分配次数: %d, 总分配: %d bytes",
			m2.Mallocs-m1.Mallocs, m2.TotalAlloc-m1.TotalAlloc)
	})
}

// 功能测试：确保对象池功能正确

func TestBufPool(t *testing.T) {
	pool := NewBufPool(256, 1024)

	t.Run("基本功能", func(t *testing.T) {
		buf := pool.Get()
		if buf == nil {
			t.Fatal("Get() 返回 nil")
		}
		if buf.Len() != 0 {
			t.Errorf("新获取的缓冲区长度应为0，实际为 %d", buf.Len())
		}
		if buf.Cap() < 256 {
			t.Errorf("缓冲区容量应至少为256，实际为 %d", buf.Cap())
		}

		buf.WriteString("test")
		pool.Put(buf)
	})

	t.Run("容量指定", func(t *testing.T) {
		buf := pool.GetCap(512)
		if buf.Cap() < 512 {
			t.Errorf("缓冲区容量应至少为512，实际为 %d", buf.Cap())
		}
		pool.Put(buf)
	})

	t.Run("超大容量不回收", func(t *testing.T) {
		buf := pool.GetCap(2048) // 超过maxCap(1024)
		buf.WriteString("test")
		pool.Put(buf) // 应该不会回收
	})

	t.Run("With函数", func(t *testing.T) {
		data := pool.With(func(buf *bytes.Buffer) {
			buf.WriteString("Hello")
			buf.WriteByte(' ')
			buf.WriteString("World")
		})
		expected := "Hello World"
		if string(data) != expected {
			t.Errorf("期望 %q，实际 %q", expected, string(data))
		}
	})

	t.Run("WithCap函数", func(t *testing.T) {
		data := pool.WithCap(1024, func(buf *bytes.Buffer) {
			buf.WriteString("Hello")
			buf.WriteByte(' ')
			buf.WriteString("World")
		})
		expected := "Hello World"
		if string(data) != expected {
			t.Errorf("期望 %q，实际 %q", expected, string(data))
		}
	})
}

func TestDefaultPool(t *testing.T) {
	t.Run("全局函数", func(t *testing.T) {
		buf := GetBuf()
		buf.WriteString("test")
		PutBuf(buf)

		buf2 := GetBufCap(512)
		if buf2.Cap() < 512 {
			t.Errorf("缓冲区容量应至少为512，实际为 %d", buf2.Cap())
		}
		PutBuf(buf2)
	})

	t.Run("WithBuf函数", func(t *testing.T) {
		data := WithBuf(func(buf *bytes.Buffer) {
			buf.WriteString("Global")
			buf.WriteByte(' ')
			buf.WriteString("Test")
		})
		expected := "Global Test"
		if string(data) != expected {
			t.Errorf("期望 %q，实际 %q", expected, string(data))
		}
	})

	t.Run("WithBufCap函数", func(t *testing.T) {
		data := WithBufCap(1024, func(buf *bytes.Buffer) {
			buf.WriteString("Global")
			buf.WriteByte(' ')
			buf.WriteString("Cap")
			buf.WriteByte(' ')
			buf.WriteString("Test")
		})
		expected := "Global Cap Test"
		if string(data) != expected {
			t.Errorf("期望 %q，实际 %q", expected, string(data))
		}
	})
}

// 边界条件测试
func TestEdgeCases(t *testing.T) {
	pool := NewBufPool(256, 1024)

	t.Run("Put nil缓冲区", func(t *testing.T) {
		pool.Put(nil) // 不应该panic
	})

	t.Run("零容量", func(t *testing.T) {
		buf := pool.GetCap(0)
		if buf.Cap() < 256 { // 应该使用默认容量
			t.Errorf("零容量应使用默认容量256，实际为 %d", buf.Cap())
		}
		pool.Put(buf)
	})

	t.Run("负容量", func(t *testing.T) {
		buf := pool.GetCap(-100)
		if buf.Cap() < 256 { // 应该使用默认容量
			t.Errorf("负容量应使用默认容量256，实际为 %d", buf.Cap())
		}
		pool.Put(buf)
	})
}
