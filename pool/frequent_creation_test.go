package pool

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// 测试频繁创建字符串构建器的场景

// BenchmarkFrequentStringCreation_WithPool 频繁创建场景使用对象池
func BenchmarkFrequentStringCreation_WithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟在一个函数中需要多次创建字符串构建器
		for j := 0; j < 100; j++ {
			builder := GetStr()
			builder.WriteString("Request ID: ")
			fmt.Fprintf(builder, "%d-%d", i, j)
			builder.WriteString(", Status: OK")
			_ = builder.String()
			PutStr(builder)
		}
	}
}

// BenchmarkFrequentStringCreation_WithoutPool 频繁创建场景不使用对象池
func BenchmarkFrequentStringCreation_WithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟在一个函数中需要多次创建字符串构建器
		for j := 0; j < 100; j++ {
			builder := &strings.Builder{}
			builder.Grow(256) // 预分配容量
			builder.WriteString("Request ID: ")
			fmt.Fprintf(builder, "%d-%d", i, j)
			builder.WriteString(", Status: OK")
			_ = builder.String()
		}
	}
}

// BenchmarkFrequentStringCreation_WithFunction 使用With函数的频繁创建场景
func BenchmarkFrequentStringCreation_WithFunction(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟在一个函数中需要多次创建字符串构建器
		for j := 0; j < 100; j++ {
			result := WithStr(func(builder *strings.Builder) {
				builder.WriteString("Request ID: ")
				fmt.Fprintf(builder, "%d-%d", i, j)
				builder.WriteString(", Status: OK")
			})
			_ = result
		}
	}
}

// 测试不同频率的创建场景

// BenchmarkStringCreation_10Times_WithPool 10次创建使用对象池
func BenchmarkStringCreation_10Times_WithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			builder := GetStr()
			builder.WriteString("Log entry ")
			fmt.Fprintf(builder, "%d", j)
			_ = builder.String()
			PutStr(builder)
		}
	}
}

// BenchmarkStringCreation_10Times_WithoutPool 10次创建不使用对象池
func BenchmarkStringCreation_10Times_WithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			builder := &strings.Builder{}
			builder.Grow(256)
			builder.WriteString("Log entry ")
			fmt.Fprintf(builder, "%d", j)
			_ = builder.String()
		}
	}
}

// BenchmarkStringCreation_1000Times_WithPool 1000次创建使用对象池
func BenchmarkStringCreation_1000Times_WithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			builder := GetStr()
			builder.WriteString("Item ")
			fmt.Fprintf(builder, "%d", j)
			_ = builder.String()
			PutStr(builder)
		}
	}
}

// BenchmarkStringCreation_1000Times_WithoutPool 1000次创建不使用对象池
func BenchmarkStringCreation_1000Times_WithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			builder := &strings.Builder{}
			builder.Grow(256)
			builder.WriteString("Item ")
			fmt.Fprintf(builder, "%d", j)
			_ = builder.String()
		}
	}
}

// 内存分配对比测试

// BenchmarkFrequentCreation_Memory_WithPool 频繁创建的内存分配测试（使用对象池）
func BenchmarkFrequentCreation_Memory_WithPool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 50; j++ {
			builder := GetStr()
			builder.WriteString("User: ")
			fmt.Fprintf(builder, "%d", j)
			builder.WriteString(", Action: login")
			_ = builder.String()
			PutStr(builder)
		}
	}
}

// BenchmarkFrequentCreation_Memory_WithoutPool 频繁创建的内存分配测试（不使用对象池）
func BenchmarkFrequentCreation_Memory_WithoutPool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 50; j++ {
			builder := &strings.Builder{}
			builder.Grow(256)
			builder.WriteString("User: ")
			fmt.Fprintf(builder, "%d", j)
			builder.WriteString(", Action: login")
			_ = builder.String()
		}
	}
}

// 实际应用场景测试

// BenchmarkLogFormatting_WithPool 日志格式化场景使用对象池
func BenchmarkLogFormatting_WithPool(b *testing.B) {
	levels := []string{"INFO", "WARN", "ERROR", "DEBUG"}
	messages := []string{"User login", "Database query", "Cache miss", "API call"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟一次请求中需要记录多条日志
		for j := 0; j < 20; j++ {
			level := levels[j%len(levels)]
			message := messages[j%len(messages)]

			result := WithStr(func(builder *strings.Builder) {
				builder.WriteString("[")
				builder.WriteString(level)
				builder.WriteString("] ")
				builder.WriteString("2023-12-07 10:30:45 - ")
				builder.WriteString(message)
				builder.WriteString(" (request_id: ")
				fmt.Fprintf(builder, "%d-%d", i, j)
				builder.WriteString(")")
			})
			_ = result
		}
	}
}

// BenchmarkLogFormatting_WithoutPool 日志格式化场景不使用对象池
func BenchmarkLogFormatting_WithoutPool(b *testing.B) {
	levels := []string{"INFO", "WARN", "ERROR", "DEBUG"}
	messages := []string{"User login", "Database query", "Cache miss", "API call"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟一次请求中需要记录多条日志
		for j := 0; j < 20; j++ {
			level := levels[j%len(levels)]
			message := messages[j%len(messages)]

			builder := &strings.Builder{}
			builder.Grow(256)
			builder.WriteString("[")
			builder.WriteString(level)
			builder.WriteString("] ")
			builder.WriteString("2023-12-07 10:30:45 - ")
			builder.WriteString(message)
			builder.WriteString(" (request_id: ")
			fmt.Fprintf(builder, "%d-%d", i, j)
			builder.WriteString(")")
			_ = builder.String()
		}
	}
}

// BenchmarkHTTPResponseBuilding_WithPool HTTP响应构建使用对象池
func BenchmarkHTTPResponseBuilding_WithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟构建多个HTTP响应
		for j := 0; j < 30; j++ {
			result := WithStrCap(1024, func(builder *strings.Builder) {
				builder.WriteString(`{"status":"success","data":{"id":`)
				fmt.Fprintf(builder, "%d", j)
				builder.WriteString(`,"name":"user`)
				fmt.Fprintf(builder, "%d", j)
				builder.WriteString(`","email":"user`)
				fmt.Fprintf(builder, "%d", j)
				builder.WriteString(`@example.com","created_at":"2023-12-07T10:30:45Z","permissions":["read","write"]}}`)
			})
			_ = result
		}
	}
}

// BenchmarkHTTPResponseBuilding_WithoutPool HTTP响应构建不使用对象池
func BenchmarkHTTPResponseBuilding_WithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟构建多个HTTP响应
		for j := 0; j < 30; j++ {
			builder := &strings.Builder{}
			builder.Grow(1024)
			builder.WriteString(`{"status":"success","data":{"id":`)
			fmt.Fprintf(builder, "%d", j)
			builder.WriteString(`,"name":"user`)
			fmt.Fprintf(builder, "%d", j)
			builder.WriteString(`","email":"user`)
			fmt.Fprintf(builder, "%d", j)
			builder.WriteString(`@example.com","created_at":"2023-12-07T10:30:45Z","permissions":["read","write"]}}`)
			_ = builder.String()
		}
	}
}

// GC压力测试

// TestFrequentCreation_GCPressure 测试频繁创建的GC压力
func TestFrequentCreation_GCPressure(t *testing.T) {
	// 测试使用对象池的GC压力
	t.Run("WithPool", func(t *testing.T) {
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		// 模拟频繁创建场景
		for i := 0; i < 1000; i++ {
			for j := 0; j < 10; j++ {
				builder := GetStr()
				builder.WriteString("Test message ")
				fmt.Fprintf(builder, "%d-%d", i, j)
				_ = builder.String()
				PutStr(builder)
			}
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		t.Logf("频繁创建使用对象池 - 分配次数: %d, 总分配: %d bytes",
			m2.Mallocs-m1.Mallocs, m2.TotalAlloc-m1.TotalAlloc)
	})

	// 测试不使用对象池的GC压力
	t.Run("WithoutPool", func(t *testing.T) {
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		// 模拟频繁创建场景
		for i := 0; i < 1000; i++ {
			for j := 0; j < 10; j++ {
				builder := &strings.Builder{}
				builder.Grow(256)
				builder.WriteString("Test message ")
				fmt.Fprintf(builder, "%d-%d", i, j)
				_ = builder.String()
			}
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		t.Logf("频繁创建不使用对象池 - 分配次数: %d, 总分配: %d bytes",
			m2.Mallocs-m1.Mallocs, m2.TotalAlloc-m1.TotalAlloc)
	})
}

// 并发频繁创建测试

// BenchmarkFrequentCreation_Concurrent_WithPool 并发频繁创建使用对象池
func BenchmarkFrequentCreation_Concurrent_WithPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 每个goroutine中频繁创建
			for j := 0; j < 10; j++ {
				builder := GetStr()
				builder.WriteString("Concurrent message ")
				fmt.Fprintf(builder, "%d", j)
				_ = builder.String()
				PutStr(builder)
			}
		}
	})
}

// BenchmarkFrequentCreation_Concurrent_WithoutPool 并发频繁创建不使用对象池
func BenchmarkFrequentCreation_Concurrent_WithoutPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 每个goroutine中频繁创建
			for j := 0; j < 10; j++ {
				builder := &strings.Builder{}
				builder.Grow(256)
				builder.WriteString("Concurrent message ")
				fmt.Fprintf(builder, "%d", j)
				_ = builder.String()
			}
		}
	})
}

// 真实场景模拟：Web服务器请求处理

// BenchmarkWebRequestProcessing_WithPool Web请求处理使用对象池
func BenchmarkWebRequestProcessing_WithPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟处理一个Web请求，需要生成多个字符串

		// 1. 生成请求ID
		requestID := WithStr(func(builder *strings.Builder) {
			builder.WriteString("req_")
			fmt.Fprintf(builder, "%d", i)
			builder.WriteString("_")
			fmt.Fprintf(builder, "%d", i*1000)
		})

		// 2. 生成多个日志条目
		for j := 0; j < 5; j++ {
			logEntry := WithStr(func(builder *strings.Builder) {
				builder.WriteString("[INFO] Processing step ")
				fmt.Fprintf(builder, "%d", j)
				builder.WriteString(" for request ")
				builder.WriteString(requestID)
			})
			_ = logEntry
		}

		// 3. 生成响应
		response := WithStrCap(512, func(builder *strings.Builder) {
			builder.WriteString(`{"request_id":"`)
			builder.WriteString(requestID)
			builder.WriteString(`","status":"success","data":[`)
			for k := 0; k < 3; k++ {
				if k > 0 {
					builder.WriteString(",")
				}
				builder.WriteString(`{"id":`)
				fmt.Fprintf(builder, "%d", k)
				builder.WriteString(`,"value":"item`)
				fmt.Fprintf(builder, "%d", k)
				builder.WriteString(`"}`)
			}
			builder.WriteString(`]}`)
		})
		_ = response
	}
}

// BenchmarkWebRequestProcessing_WithoutPool Web请求处理不使用对象池
func BenchmarkWebRequestProcessing_WithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟处理一个Web请求，需要生成多个字符串

		// 1. 生成请求ID
		builder1 := &strings.Builder{}
		builder1.Grow(256)
		builder1.WriteString("req_")
		fmt.Fprintf(builder1, "%d", i)
		builder1.WriteString("_")
		fmt.Fprintf(builder1, "%d", i*1000)
		requestID := builder1.String()

		// 2. 生成多个日志条目
		for j := 0; j < 5; j++ {
			builder2 := &strings.Builder{}
			builder2.Grow(256)
			builder2.WriteString("[INFO] Processing step ")
			fmt.Fprintf(builder2, "%d", j)
			builder2.WriteString(" for request ")
			builder2.WriteString(requestID)
			_ = builder2.String()
		}

		// 3. 生成响应
		builder3 := &strings.Builder{}
		builder3.Grow(512)
		builder3.WriteString(`{"request_id":"`)
		builder3.WriteString(requestID)
		builder3.WriteString(`","status":"success","data":[`)
		for k := 0; k < 3; k++ {
			if k > 0 {
				builder3.WriteString(",")
			}
			builder3.WriteString(`{"id":`)
			fmt.Fprintf(builder3, "%d", k)
			builder3.WriteString(`,"value":"item`)
			fmt.Fprintf(builder3, "%d", k)
			builder3.WriteString(`"}`)
		}
		builder3.WriteString(`]}`)
		_ = builder3.String()
	}
}
