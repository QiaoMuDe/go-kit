package pool

import (
	"sync"
	"testing"
)

func TestBytePool_Get(t *testing.T) {
	data := GetByteWithCapacity(1024)
	if data == nil {
		t.Fatal("GetBytes() returned nil")
	}

	// 验证返回的切片长度为请求的大小
	if len(data) != 1024 {
		t.Errorf("Expected slice length 1024, got length %d", len(data))
	}

	PutByte(data)
}

func TestBytePool_Put(t *testing.T) {
	data := GetByteWithCapacity(1024)
	// 修改数据内容
	for i := range data {
		data[i] = byte(i % 256)
	}

	if len(data) != 1024 {
		t.Fatal("Data length should be 1024")
	}

	PutByte(data)

	// 再次获取应该是指定长度的切片
	data2 := GetByteWithCapacity(512)
	if len(data2) != 512 {
		t.Errorf("Expected slice length 512, got length %d", len(data2))
	}

	PutByte(data2)
}

func TestBytePool_Reuse(t *testing.T) {
	// 测试对象池的复用机制
	data1 := GetByte()
	data1 = append(data1, 'a', 'b', 'c')
	PutByte(data1)

	data2 := GetByte()
	// 验证容量可能被保留（取决于实现）
	if cap(data2) == 0 {
		t.Log("New slice has zero capacity")
	}

	PutByte(data2)
}

func TestBytePool_Concurrent(t *testing.T) {
	const numGoroutines = 100
	const numOperations = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				data := GetByteWithCapacity(1024)
				if data == nil {
					t.Errorf("GetByte() returned nil in goroutine %d", id)
					return
				}

				// 修改数据内容
				if len(data) > 0 {
					data[0] = byte(id % 256)
				}
				if len(data) > 1 {
					data[1] = byte(j % 256)
				}

				// 验证数据长度
				if len(data) != 1024 {
					t.Errorf("Expected length 1024, got %d in goroutine %d", len(data), id)
					return
				}

				PutByte(data)
			}
		}(i)
	}

	wg.Wait()
}

func TestBytePool_LargeSlice(t *testing.T) {
	largeSize := 1024 * 1024 // 1MB
	data := GetByteWithCapacity(largeSize)

	// 验证获取的切片大小
	if len(data) != largeSize {
		t.Errorf("Expected length %d, got %d", largeSize, len(data))
	}

	// 修改数据内容
	for i := 0; i < len(data); i++ {
		data[i] = byte(i % 256)
	}

	PutByte(data)

	// 验证获取新的切片
	newData := GetByte()
	if len(newData) != 256 {
		t.Errorf("Expected new slice length 256, got length %d", len(newData))
	}

	PutByte(newData)
}

func TestBytePool_EdgeCases(t *testing.T) {
	// 测试空切片
	data := GetByte()
	PutByte(data) // 应该不会panic

	// 测试nil切片
	PutByte(nil) // 应该不会panic

	// 测试多次put同一个切片
	data2 := GetByte()
	PutByte(data2)
	PutByte(data2) // 应该不会panic，但可能导致问题
}

func TestBytePool_CapacityGrowth(t *testing.T) {
	data := GetByte()
	initialCap := cap(data)

	// 强制扩容
	for i := 0; i < 1000; i++ {
		data = append(data, byte(i%256))
	}

	if cap(data) <= initialCap {
		t.Log("Capacity did not grow as expected, but this might be implementation dependent")
	}

	PutByte(data)
}

func BenchmarkBytePool_GetPut(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data := GetByte()
		data = append(data, []byte("benchmark test data")...)
		PutByte(data)
	}
}

func BenchmarkBytePool_vs_Make(b *testing.B) {
	b.Run("Pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := GetByte()
			data = append(data, []byte("benchmark")...)
			PutByte(data)
		}
	})

	b.Run("Make", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := make([]byte, 0)
			data = append(data, []byte("benchmark")...)
			_ = data // Use the data to avoid ineffassign warning
		}
	})
}

func BenchmarkBytePool_LargeAllocation(b *testing.B) {
	const size = 1024 * 10 // 10KB

	b.Run("Pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := GetByteWithCapacity(size)
			for j := 0; j < size; j++ {
				data = append(data, byte(j%256))
			}
			_ = data // Use the data to avoid staticcheck warning
			PutByte(data)
		}
	})

	b.Run("Make", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := make([]byte, 0, size)
			for j := 0; j < size; j++ {
				data = append(data, byte(j%256))
			}
			_ = len(data) // Use the data to avoid staticcheck warning
		}
	})
}
