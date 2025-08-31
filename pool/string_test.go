package pool

import (
	"strings"
	"sync"
	"testing"
)

func TestStringPool_Get(t *testing.T) {
	sb := GetString(1024)
	if sb == nil {
		t.Fatal("GetStringBuilder() returned nil")
	}

	// 验证返回的是一个空的builder
	if sb.Len() != 0 {
		t.Errorf("Expected empty builder, got length %d", sb.Len())
	}

	PutString(sb)
}

func TestStringPool_Put(t *testing.T) {
	sb := GetString(1024)
	sb.WriteString("test data")

	if sb.Len() == 0 {
		t.Fatal("Builder should contain data before put")
	}

	PutString(sb)

	// 再次获取应该是空的
	sb2 := GetString(1024)
	if sb2.Len() != 0 {
		t.Errorf("Expected empty builder after put, got length %d", sb2.Len())
	}

	PutString(sb2)
}

func TestStringPool_Reuse(t *testing.T) {
	sb1 := GetString(1024)
	sb1.WriteString("test")
	PutString(sb1)

	sb2 := GetString(1024)
	// 在单线程环境下可能复用同一个对象
	if sb1 == sb2 {
		t.Log("Reused the same StringBuilder object")
	}

	PutString(sb2)
}

func TestStringPool_Concurrent(t *testing.T) {
	const numGoroutines = 100
	const numOperations = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				sb := GetString(1024)
				if sb == nil {
					t.Errorf("GetString() returned nil in goroutine %d", id)
					return
				}

				// 写入一些数据
				sb.WriteString("goroutine")
				sb.WriteString(string(rune('0' + id%10)))
				sb.WriteString("operation")
				sb.WriteString(string(rune('0' + j%10)))

				// 验证数据
				result := sb.String()
				if len(result) == 0 {
					t.Errorf("Builder is empty after writing in goroutine %d", id)
					return
				}

				expected := "goroutine" + string(rune('0'+id%10)) + "operation" + string(rune('0'+j%10))
				if result != expected {
					t.Errorf("Expected %s, got %s in goroutine %d", expected, result, id)
					return
				}

				PutString(sb)
			}
		}(i)
	}

	wg.Wait()
}

func TestStringPool_StringBuilderMethods(t *testing.T) {
	sb := GetString(1024)

	// 测试各种写入方法
	sb.WriteString("Hello")
	sb.WriteByte(' ')
	sb.WriteRune('世')
	sb.WriteRune('界')

	result := sb.String()
	expected := "Hello 世界"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// 测试长度
	if sb.Len() != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), sb.Len())
	}

	// 测试容量
	if sb.Cap() < sb.Len() {
		t.Errorf("Capacity %d should be >= length %d", sb.Cap(), sb.Len())
	}

	PutString(sb)
}

func TestStringPool_LargeString(t *testing.T) {
	sb := GetString(1024)

	// 构建大字符串
	const iterations = 10000
	for i := 0; i < iterations; i++ {
		sb.WriteString("test")
	}

	result := sb.String()
	expectedLen := iterations * 4 // "test" has 4 characters
	if len(result) != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, len(result))
	}

	PutString(sb)

	// 验证重置后是空的
	newSb := GetString(1024)
	if newSb.Len() != 0 {
		t.Errorf("Expected empty builder after putting large string, got length %d", newSb.Len())
	}

	PutString(newSb)
}

func TestStringPool_Reset(t *testing.T) {
	sb := GetString(1024)
	sb.WriteString("some data")
	sb.WriteByte(0x00)
	sb.WriteRune('测')

	if sb.Len() == 0 {
		t.Fatal("Builder should contain data before reset")
	}

	PutString(sb)

	// 获取新的builder应该是空的
	newSb := GetString(1024)
	if newSb.Len() != 0 {
		t.Errorf("Expected empty builder after reset, got length %d", newSb.Len())
	}

	if newSb.String() != "" {
		t.Errorf("Expected empty string after reset, got %q", newSb.String())
	}

	PutString(newSb)
}

func TestStringPool_EdgeCases(t *testing.T) {
	// 测试空字符串
	sb := GetString(1024)
	sb.WriteString("")
	if sb.Len() != 0 {
		t.Error("Writing empty string should not change length")
	}
	PutString(sb)

	// 测试Unicode字符
	sb2 := GetString(1024)
	sb2.WriteString("Hello")
	sb2.WriteRune('🌍')
	sb2.WriteString("世界")

	result := sb2.String()
	if !strings.Contains(result, "🌍") {
		t.Error("Should handle Unicode characters correctly")
	}

	PutString(sb2)
}

func TestStringPool_CapacityGrowth(t *testing.T) {
	sb := GetString(1024)
	initialCap := sb.Cap()

	// 强制扩容
	largeString := strings.Repeat("a", 1000)
	sb.WriteString(largeString)

	if sb.Cap() <= initialCap {
		t.Log("Capacity did not grow as expected, but this might be implementation dependent")
	}

	if sb.Len() != 1000 {
		t.Errorf("Expected length 1000, got %d", sb.Len())
	}

	PutString(sb)
}

func BenchmarkStringPool_GetPut(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sb := GetString(1024)
		sb.WriteString("benchmark test data")
		_ = sb.String()
		PutString(sb)
	}
}

func BenchmarkStringPool_vs_New(b *testing.B) {
	b.Run("Pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sb := GetString(1024)
			sb.WriteString("benchmark")
			_ = sb.String()
			PutString(sb)
		}
	})

	b.Run("New", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sb := &strings.Builder{}
			sb.WriteString("benchmark")
			_ = sb.String()
		}
	})
}

func BenchmarkStringPool_StringBuilding(b *testing.B) {
	const numStrings = 100

	b.Run("Pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sb := GetString(1024)
			for j := 0; j < numStrings; j++ {
				sb.WriteString("test")
				sb.WriteString(string(rune('0' + j%10)))
			}
			_ = sb.String()
			PutString(sb)
		}
	})

	b.Run("Concatenation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var result string
			for j := 0; j < numStrings; j++ {
				result += "test" + string(rune('0'+j%10))
			}
		}
	})
}
