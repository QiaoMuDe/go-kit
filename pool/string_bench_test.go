package pool

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// 基准测试：对比使用字符串构建器对象池和不使用对象池的性能差异

// BenchmarkStringWithPool 使用对象池的基准测试
func BenchmarkStringWithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := GetStr()
		builder.WriteString("Hello")
		builder.WriteByte(' ')
		builder.WriteString("World")
		fmt.Fprintf(builder, " %d", i)
		_ = builder.String()
		PutStr(builder)
	}
}

// BenchmarkStringWithoutPool 不使用对象池的基准测试
func BenchmarkStringWithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := &strings.Builder{}
		builder.Grow(256) // 预分配容量
		builder.WriteString("Hello")
		builder.WriteByte(' ')
		builder.WriteString("World")
		fmt.Fprintf(builder, " %d", i)
		_ = builder.String()
		// 不归还，让GC处理
	}
}

// BenchmarkStringWithPoolLarge 使用对象池处理大字符串的基准测试
func BenchmarkStringWithPoolLarge(b *testing.B) {
	data := strings.Repeat("A", 1024) // 1KB数据
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := GetStrCap(10240) // 10KB容量
		for j := 0; j < 10; j++ {
			builder.WriteString(data)
		}
		_ = builder.String()
		PutStr(builder)
	}
}

// BenchmarkStringWithoutPoolLarge 不使用对象池处理大字符串的基准测试
func BenchmarkStringWithoutPoolLarge(b *testing.B) {
	data := strings.Repeat("A", 1024) // 1KB数据
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := &strings.Builder{}
		builder.Grow(10240) // 10KB容量
		for j := 0; j < 10; j++ {
			builder.WriteString(data)
		}
		_ = builder.String()
		// 不归还，让GC处理
	}
}

// BenchmarkStringWithFunction 使用With函数的基准测试
func BenchmarkStringWithFunction(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := WithStr(func(builder *strings.Builder) {
			builder.WriteString("Hello")
			builder.WriteByte(' ')
			builder.WriteString("World")
			fmt.Fprintf(builder, " %d", i)
		})
		_ = result
	}
}

// BenchmarkStringWithCapFunction 使用WithCap函数的基准测试
func BenchmarkStringWithCapFunction(b *testing.B) {
	data := strings.Repeat("A", 1024) // 1KB数据
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := WithStrCap(10240, func(builder *strings.Builder) {
			for j := 0; j < 10; j++ {
				builder.WriteString(data)
			}
		})
		_ = result
	}
}

// 内存分配测试：对比内存分配次数和大小

// BenchmarkStringMemoryWithPool 测试使用对象池的内存分配
func BenchmarkStringMemoryWithPool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := GetStr()
		builder.WriteString("Hello World")
		fmt.Fprintf(builder, " %d", i)
		_ = builder.String()
		PutStr(builder)
	}
}

// BenchmarkStringMemoryWithoutPool 测试不使用对象池的内存分配
func BenchmarkStringMemoryWithoutPool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := &strings.Builder{}
		builder.Grow(256)
		builder.WriteString("Hello World")
		fmt.Fprintf(builder, " %d", i)
		_ = builder.String()
	}
}

// 并发测试：测试对象池在并发环境下的性能

// BenchmarkStringConcurrentWithPool 并发使用对象池的基准测试
func BenchmarkStringConcurrentWithPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			builder := GetStr()
			builder.WriteString("Hello")
			builder.WriteByte(' ')
			builder.WriteString("World")
			_ = builder.String()
			PutStr(builder)
		}
	})
}

// BenchmarkStringConcurrentWithoutPool 并发不使用对象池的基准测试
func BenchmarkStringConcurrentWithoutPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			builder := &strings.Builder{}
			builder.Grow(256)
			builder.WriteString("Hello")
			builder.WriteByte(' ')
			builder.WriteString("World")
			_ = builder.String()
		}
	})
}

// 字符串拼接场景测试

// BenchmarkStringJoinWithPool 使用对象池进行字符串拼接
func BenchmarkStringJoinWithPool(b *testing.B) {
	words := []string{"apple", "banana", "cherry", "date", "elderberry"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := WithStr(func(builder *strings.Builder) {
			for j, word := range words {
				if j > 0 {
					builder.WriteString(", ")
				}
				builder.WriteString(word)
			}
		})
		_ = result
	}
}

// BenchmarkStringJoinWithoutPool 不使用对象池进行字符串拼接
func BenchmarkStringJoinWithoutPool(b *testing.B) {
	words := []string{"apple", "banana", "cherry", "date", "elderberry"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := &strings.Builder{}
		builder.Grow(256)
		for j, word := range words {
			if j > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(word)
		}
		_ = builder.String()
	}
}

// BenchmarkStringJoinBuiltIn 使用内置strings.Join进行对比
func BenchmarkStringJoinBuiltIn(b *testing.B) {
	words := []string{"apple", "banana", "cherry", "date", "elderberry"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Join(words, ", ")
	}
}

// JSON构建场景测试

// BenchmarkStringJSONWithPool 使用对象池构建JSON字符串
func BenchmarkStringJSONWithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := WithStrCap(512, func(builder *strings.Builder) {
			builder.WriteString(`{"id":`)
			fmt.Fprintf(builder, "%d", i)
			builder.WriteString(`,"name":"user`)
			fmt.Fprintf(builder, "%d", i)
			builder.WriteString(`","email":"user`)
			fmt.Fprintf(builder, "%d", i)
			builder.WriteString(`@example.com","active":true}`)
		})
		_ = result
	}
}

// BenchmarkStringJSONWithoutPool 不使用对象池构建JSON字符串
func BenchmarkStringJSONWithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := &strings.Builder{}
		builder.Grow(512)
		builder.WriteString(`{"id":`)
		fmt.Fprintf(builder, "%d", i)
		builder.WriteString(`,"name":"user`)
		fmt.Fprintf(builder, "%d", i)
		builder.WriteString(`","email":"user`)
		fmt.Fprintf(builder, "%d", i)
		builder.WriteString(`@example.com","active":true}`)
		_ = builder.String()
	}
}

// GC压力测试：测试对象池对GC的影响

// TestStringGCPressure 测试GC压力
func TestStringGCPressure(t *testing.T) {
	// 测试使用对象池的GC压力
	t.Run("WithPool", func(t *testing.T) {
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		for i := 0; i < 10000; i++ {
			builder := GetStr()
			builder.WriteString("Hello World")
			fmt.Fprintf(builder, " %d", i)
			_ = builder.String()
			PutStr(builder)
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		t.Logf("使用字符串对象池 - 分配次数: %d, 总分配: %d bytes",
			m2.Mallocs-m1.Mallocs, m2.TotalAlloc-m1.TotalAlloc)
	})

	// 测试不使用对象池的GC压力
	t.Run("WithoutPool", func(t *testing.T) {
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		for i := 0; i < 10000; i++ {
			builder := &strings.Builder{}
			builder.Grow(256)
			builder.WriteString("Hello World")
			fmt.Fprintf(builder, " %d", i)
			_ = builder.String()
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		t.Logf("不使用字符串对象池 - 分配次数: %d, 总分配: %d bytes",
			m2.Mallocs-m1.Mallocs, m2.TotalAlloc-m1.TotalAlloc)
	})
}

// 复用效率测试

// TestStringPoolReuse 测试对象池的复用效率
func TestStringPoolReuse(t *testing.T) {
	pool := NewStrPool(256, 1024)

	// 获取一个构建器并使用
	builder1 := pool.Get()
	builder1.WriteString("test")
	originalCap := builder1.Cap()
	pool.Put(builder1)

	// 再次获取，应该复用同一个对象
	builder2 := pool.Get()
	if builder2.Cap() != originalCap {
		t.Errorf("期望复用相同容量的构建器，原容量: %d, 新容量: %d", originalCap, builder2.Cap())
	}
	if builder2.Len() != 0 {
		t.Errorf("复用的构建器应该被重置，当前长度: %d", builder2.Len())
	}
	pool.Put(builder2)
}

// 容量管理测试

// TestStringPoolCapacityManagement 测试容量管理
func TestStringPoolCapacityManagement(t *testing.T) {
	pool := NewStrPool(256, 1024)

	t.Run("超大容量不回收", func(t *testing.T) {
		builder := pool.GetCap(2048) // 超过maxCap(1024)
		if builder.Cap() < 2048 {
			t.Errorf("构建器容量应至少为2048，实际为 %d", builder.Cap())
		}
		builder.WriteString("test")
		pool.Put(builder) // 应该不会回收，因为容量过大
	})

	t.Run("容量扩展", func(t *testing.T) {
		builder := pool.GetCap(512)
		if builder.Cap() < 512 {
			t.Errorf("构建器容量应至少为512，实际为 %d", builder.Cap())
		}
		pool.Put(builder)
	})
}

// 实际应用场景测试

// BenchmarkStringLogFormatWithPool 使用对象池格式化日志
func BenchmarkStringLogFormatWithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := WithStr(func(builder *strings.Builder) {
			builder.WriteString("[INFO] ")
			builder.WriteString("2023-12-07 10:30:45")
			builder.WriteString(" - User ")
			fmt.Fprintf(builder, "%d", i)
			builder.WriteString(" logged in from IP 192.168.1.")
			fmt.Fprintf(builder, "%d", i%255)
		})
		_ = result
	}
}

// BenchmarkStringLogFormatWithoutPool 不使用对象池格式化日志
func BenchmarkStringLogFormatWithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := &strings.Builder{}
		builder.Grow(256)
		builder.WriteString("[INFO] ")
		builder.WriteString("2023-12-07 10:30:45")
		builder.WriteString(" - User ")
		fmt.Fprintf(builder, "%d", i)
		builder.WriteString(" logged in from IP 192.168.1.")
		fmt.Fprintf(builder, "%d", i%255)
		_ = builder.String()
	}
}

// BenchmarkStringURLBuildWithPool 使用对象池构建URL
func BenchmarkStringURLBuildWithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := WithStr(func(builder *strings.Builder) {
			builder.WriteString("https://api.example.com/v1/users/")
			fmt.Fprintf(builder, "%d", i)
			builder.WriteString("?include=profile,settings&format=json&timestamp=")
			fmt.Fprintf(builder, "%d", 1701936645+i)
		})
		_ = result
	}
}

// BenchmarkStringURLBuildWithoutPool 不使用对象池构建URL
func BenchmarkStringURLBuildWithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := &strings.Builder{}
		builder.Grow(256)
		builder.WriteString("https://api.example.com/v1/users/")
		fmt.Fprintf(builder, "%d", i)
		builder.WriteString("?include=profile,settings&format=json&timestamp=")
		fmt.Fprintf(builder, "%d", 1701936645+i)
		_ = builder.String()
	}
}
