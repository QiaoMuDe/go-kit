package pool

import (
	"fmt"
	"runtime"
	"testing"
)

// 基准测试：对比使用字节切片对象池和不使用对象池的性能差异

// BenchmarkByteWithPool 使用对象池的基准测试
func BenchmarkByteWithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetByte()
		buf = append(buf[:0], "Hello"...)
		buf = append(buf, ' ')
		buf = append(buf, "World"...)
		buf = append(buf, fmt.Sprintf(" %d", i)...)
		_ = buf
		PutByte(buf)
	}
}

// BenchmarkByteWithoutPool 不使用对象池的基准测试
func BenchmarkByteWithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 256)
		buf = append(buf, "Hello"...)
		buf = append(buf, ' ')
		buf = append(buf, "World"...)
		buf = append(buf, fmt.Sprintf(" %d", i)...)
		_ = buf
		// 不归还，让GC处理
	}
}

// BenchmarkByteWithPoolLarge 使用对象池处理大数据的基准测试
func BenchmarkByteWithPoolLarge(b *testing.B) {
	data := make([]byte, 1024) // 1KB数据
	for i := range data {
		data[i] = byte(i % 256)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetByteCap(10240) // 10KB容量
		buf = buf[:0]
		for j := 0; j < 10; j++ {
			buf = append(buf, data...)
		}
		_ = buf
		PutByte(buf)
	}
}

// BenchmarkByteWithoutPoolLarge 不使用对象池处理大数据的基准测试
func BenchmarkByteWithoutPoolLarge(b *testing.B) {
	data := make([]byte, 1024) // 1KB数据
	for i := range data {
		data[i] = byte(i % 256)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 10240) // 10KB容量
		for j := 0; j < 10; j++ {
			buf = append(buf, data...)
		}
		_ = buf
		// 不归还，让GC处理
	}
}

// BenchmarkByteEmptyWithPool 使用对象池获取空缓冲区的基准测试
func BenchmarkByteEmptyWithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetByteEmpty(256)
		buf = append(buf, "Hello"...)
		buf = append(buf, ' ')
		buf = append(buf, "World"...)
		buf = append(buf, fmt.Sprintf(" %d", i)...)
		_ = buf
		PutByte(buf)
	}
}

// BenchmarkByteEmptyWithoutPool 不使用对象池获取空缓冲区的基准测试
func BenchmarkByteEmptyWithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 256)
		buf = append(buf, "Hello"...)
		buf = append(buf, ' ')
		buf = append(buf, "World"...)
		buf = append(buf, fmt.Sprintf(" %d", i)...)
		_ = buf
	}
}

// 内存分配测试：对比内存分配次数和大小

// BenchmarkByteMemoryWithPool 测试使用对象池的内存分配
func BenchmarkByteMemoryWithPool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetByteEmpty(256)
		buf = append(buf, "Hello World"...)
		buf = append(buf, fmt.Sprintf(" %d", i)...)
		_ = buf
		PutByte(buf)
	}
}

// BenchmarkByteMemoryWithoutPool 测试不使用对象池的内存分配
func BenchmarkByteMemoryWithoutPool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 256)
		buf = append(buf, "Hello World"...)
		buf = append(buf, fmt.Sprintf(" %d", i)...)
		_ = buf
	}
}

// 并发测试：测试对象池在并发环境下的性能

// BenchmarkByteConcurrentWithPool 并发使用对象池的基准测试
func BenchmarkByteConcurrentWithPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := GetByteEmpty(256)
			buf = append(buf, "Hello"...)
			buf = append(buf, ' ')
			buf = append(buf, "World"...)
			_ = buf
			PutByte(buf)
		}
	})
}

// BenchmarkByteConcurrentWithoutPool 并发不使用对象池的基准测试
func BenchmarkByteConcurrentWithoutPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := make([]byte, 0, 256)
			buf = append(buf, "Hello"...)
			buf = append(buf, ' ')
			buf = append(buf, "World"...)
			_ = buf
		}
	})
}

// 数据处理场景测试

// BenchmarkByteDataProcessingWithPool 使用对象池进行数据处理
func BenchmarkByteDataProcessingWithPool(b *testing.B) {
	input := []byte("The quick brown fox jumps over the lazy dog")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetByteEmpty(len(input) * 2)
		// 模拟数据处理：转换为大写并添加前缀
		buf = append(buf, "PROCESSED: "...)
		for _, c := range input {
			if c >= 'a' && c <= 'z' {
				buf = append(buf, c-32) // 转大写
			} else {
				buf = append(buf, c)
			}
		}
		_ = buf
		PutByte(buf)
	}
}

// BenchmarkByteDataProcessingWithoutPool 不使用对象池进行数据处理
func BenchmarkByteDataProcessingWithoutPool(b *testing.B) {
	input := []byte("The quick brown fox jumps over the lazy dog")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, len(input)*2)
		// 模拟数据处理：转换为大写并添加前缀
		buf = append(buf, "PROCESSED: "...)
		for _, c := range input {
			if c >= 'a' && c <= 'z' {
				buf = append(buf, c-32) // 转大写
			} else {
				buf = append(buf, c)
			}
		}
		_ = buf
	}
}

// 网络数据包处理场景

// BenchmarkBytePacketWithPool 使用对象池处理网络数据包
func BenchmarkBytePacketWithPool(b *testing.B) {
	header := []byte{0x01, 0x02, 0x03, 0x04} // 4字节头部
	payload := make([]byte, 1024)            // 1KB载荷
	for i := range payload {
		payload[i] = byte(i % 256)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetByteEmpty(1100) // 预分配足够空间
		buf = append(buf, header...)
		buf = append(buf, payload...)
		// 添加校验和（简单示例）
		checksum := byte(0)
		for _, b := range buf {
			checksum ^= b
		}
		buf = append(buf, checksum)
		_ = buf
		PutByte(buf)
	}
}

// BenchmarkBytePacketWithoutPool 不使用对象池处理网络数据包
func BenchmarkBytePacketWithoutPool(b *testing.B) {
	header := []byte{0x01, 0x02, 0x03, 0x04} // 4字节头部
	payload := make([]byte, 1024)            // 1KB载荷
	for i := range payload {
		payload[i] = byte(i % 256)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 1100) // 预分配足够空间
		buf = append(buf, header...)
		buf = append(buf, payload...)
		// 添加校验和（简单示例）
		checksum := byte(0)
		for _, b := range buf {
			checksum ^= b
		}
		buf = append(buf, checksum)
		_ = buf
	}
}

// JSON编码场景

// BenchmarkByteJSONWithPool 使用对象池构建JSON
func BenchmarkByteJSONWithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetByteEmpty(512)
		buf = append(buf, `{"id":`...)
		buf = append(buf, fmt.Sprintf("%d", i)...)
		buf = append(buf, `,"name":"user`...)
		buf = append(buf, fmt.Sprintf("%d", i)...)
		buf = append(buf, `","email":"user`...)
		buf = append(buf, fmt.Sprintf("%d", i)...)
		buf = append(buf, `@example.com","active":true}`...)
		_ = buf
		PutByte(buf)
	}
}

// BenchmarkByteJSONWithoutPool 不使用对象池构建JSON
func BenchmarkByteJSONWithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 512)
		buf = append(buf, `{"id":`...)
		buf = append(buf, fmt.Sprintf("%d", i)...)
		buf = append(buf, `,"name":"user`...)
		buf = append(buf, fmt.Sprintf("%d", i)...)
		buf = append(buf, `","email":"user`...)
		buf = append(buf, fmt.Sprintf("%d", i)...)
		buf = append(buf, `@example.com","active":true}`...)
		_ = buf
	}
}

// GC压力测试：测试对象池对GC的影响

// TestByteGCPressure 测试GC压力
func TestByteGCPressure(t *testing.T) {
	// 测试使用对象池的GC压力
	t.Run("WithPool", func(t *testing.T) {
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		for i := 0; i < 10000; i++ {
			buf := GetByteEmpty(256)
			buf = append(buf, "Hello World"...)
			buf = append(buf, fmt.Sprintf(" %d", i)...)
			_ = buf
			PutByte(buf)
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		t.Logf("使用字节对象池 - 分配次数: %d, 总分配: %d bytes",
			m2.Mallocs-m1.Mallocs, m2.TotalAlloc-m1.TotalAlloc)
	})

	// 测试不使用对象池的GC压力
	t.Run("WithoutPool", func(t *testing.T) {
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		for i := 0; i < 10000; i++ {
			buf := make([]byte, 0, 256)
			buf = append(buf, "Hello World"...)
			buf = append(buf, fmt.Sprintf(" %d", i)...)
			_ = buf
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		t.Logf("不使用字节对象池 - 分配次数: %d, 总分配: %d bytes",
			m2.Mallocs-m1.Mallocs, m2.TotalAlloc-m1.TotalAlloc)
	})
}

// 复用效率测试

// TestBytePoolReuse 测试对象池的复用效率
func TestBytePoolReuse(t *testing.T) {
	pool := NewBytePool(256, 1024)

	// 获取一个缓冲区并使用
	buf1 := pool.GetEmpty(256)
	buf1 = append(buf1, "test"...)
	originalCap := cap(buf1)
	pool.Put(buf1)

	// 再次获取，应该复用同一个对象
	buf2 := pool.GetEmpty(256)
	if cap(buf2) != originalCap {
		t.Errorf("期望复用相同容量的缓冲区，原容量: %d, 新容量: %d", originalCap, cap(buf2))
	}
	if len(buf2) != 0 {
		t.Errorf("复用的缓冲区应该被重置为空，当前长度: %d", len(buf2))
	}
	pool.Put(buf2)
}

// 容量管理测试

// TestBytePoolCapacityManagement 测试容量管理
func TestBytePoolCapacityManagement(t *testing.T) {
	pool := NewBytePool(256, 1024)

	t.Run("超大容量不回收", func(t *testing.T) {
		buf := pool.GetCap(2048) // 超过maxCap(1024)
		if cap(buf) < 2048 {
			t.Errorf("缓冲区容量应至少为2048，实际为 %d", cap(buf))
		}
		if len(buf) != 2048 {
			t.Errorf("缓冲区长度应为2048，实际为 %d", len(buf))
		}
		pool.Put(buf) // 应该不会回收，因为容量过大
	})

	t.Run("容量扩展", func(t *testing.T) {
		buf := pool.GetCap(512)
		if cap(buf) < 512 {
			t.Errorf("缓冲区容量应至少为512，实际为 %d", cap(buf))
		}
		if len(buf) != 512 {
			t.Errorf("缓冲区长度应为512，实际为 %d", len(buf))
		}
		pool.Put(buf)
	})

	t.Run("空缓冲区", func(t *testing.T) {
		buf := pool.GetEmpty(512)
		if cap(buf) < 512 {
			t.Errorf("缓冲区容量应至少为512，实际为 %d", cap(buf))
		}
		if len(buf) != 0 {
			t.Errorf("空缓冲区长度应为0，实际为 %d", len(buf))
		}
		pool.Put(buf)
	})
}

// 不同大小缓冲区性能对比

// BenchmarkByteSize64WithPool 64字节缓冲区使用对象池
func BenchmarkByteSize64WithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetByteCap(64)
		buf = buf[:0]
		buf = append(buf, "Small data"...)
		_ = buf
		PutByte(buf)
	}
}

// BenchmarkByteSize64WithoutPool 64字节缓冲区不使用对象池
func BenchmarkByteSize64WithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 64)
		buf = append(buf, "Small data"...)
		_ = buf
	}
}

// BenchmarkByteSize1KWithPool 1KB缓冲区使用对象池
func BenchmarkByteSize1KWithPool(b *testing.B) {
	data := make([]byte, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetByteCap(1024)
		buf = buf[:0]
		buf = append(buf, data...)
		_ = buf
		PutByte(buf)
	}
}

// BenchmarkByteSize1KWithoutPool 1KB缓冲区不使用对象池
func BenchmarkByteSize1KWithoutPool(b *testing.B) {
	data := make([]byte, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 1024)
		buf = append(buf, data...)
		_ = buf
	}
}

// BenchmarkByteSize64KWithPool 64KB缓冲区使用对象池
func BenchmarkByteSize64KWithPool(b *testing.B) {
	data := make([]byte, 65536)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := GetByteCap(65536)
		buf = buf[:0]
		buf = append(buf, data...)
		_ = buf
		PutByte(buf)
	}
}

// BenchmarkByteSize64KWithoutPool 64KB缓冲区不使用对象池
func BenchmarkByteSize64KWithoutPool(b *testing.B) {
	data := make([]byte, 65536)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 65536)
		buf = append(buf, data...)
		_ = buf
	}
}
