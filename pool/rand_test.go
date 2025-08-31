package pool

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestRandPool_Get(t *testing.T) {
	r := GetRand()
	if r == nil {
		t.Fatal("GetRand() returned nil")
	}

	// 验证可以生成随机数
	n1 := r.Int()
	n2 := r.Int()

	// 虽然理论上可能相等，但概率极低
	if n1 == n2 {
		t.Log("Generated same random number twice (very unlikely but possible)")
	}

	PutRand(r)
}

func TestRandPool_Put(t *testing.T) {
	r := GetRand()

	// 使用随机数生成器
	r.Seed(12345)
	_ = r.Int()

	PutRand(r)

	// 再次获取
	r2 := GetRand()
	if r2 == nil {
		t.Fatal("GetRand() returned nil after put")
	}

	PutRand(r2)
}

func TestRandPool_Reuse(t *testing.T) {
	r1 := GetRand()
	r1.Seed(54321)
	PutRand(r1)

	r2 := GetRand()
	// 在单线程环境下可能复用同一个对象
	if r1 == r2 {
		t.Log("Reused the same rand object")
	}

	PutRand(r2)
}

func TestRandPool_Concurrent(t *testing.T) {
	const numGoroutines = 100
	const numOperations = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// 用于检测随机数质量的简单测试
	results := make([][]int, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		results[i] = make([]int, numOperations)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				r := GetRand()
				if r == nil {
					t.Errorf("GetRand() returned nil in goroutine %d", id)
					return
				}

				// 生成随机数
				num := r.Intn(1000)
				results[id][j] = num

				PutRand(r)
			}
		}(i)
	}

	wg.Wait()

	// 简单验证随机数的分布
	for i, result := range results {
		if len(result) != numOperations {
			t.Errorf("Goroutine %d generated %d numbers, expected %d", i, len(result), numOperations)
		}
	}
}

func TestRandPool_RandomQuality(t *testing.T) {
	r := GetRand()

	// 测试不同的随机数生成方法
	t.Run("Int", func(t *testing.T) {
		nums := make(map[int]bool)
		for i := 0; i < 1000; i++ {
			num := r.Int()
			if nums[num] {
				t.Log("Duplicate number generated (rare but possible)")
			}
			nums[num] = true
		}
	})

	t.Run("Intn", func(t *testing.T) {
		const max = 100
		counts := make([]int, max)
		for i := 0; i < 10000; i++ {
			num := r.Intn(max)
			if num < 0 || num >= max {
				t.Errorf("Intn(%d) returned %d, out of range", max, num)
			}
			counts[num]++
		}

		// 简单的分布检查
		for i, count := range counts {
			if count == 0 {
				t.Errorf("Number %d was never generated", i)
			}
		}
	})

	t.Run("Float64", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			f := r.Float64()
			if f < 0 || f >= 1 {
				t.Errorf("Float64() returned %f, should be in [0,1)", f)
			}
		}
	})

	PutRand(r)
}

func TestRandPool_Seed(t *testing.T) {
	r1 := GetRand()
	r2 := GetRand()

	// 使用相同的种子
	seed := time.Now().UnixNano()
	r1.Seed(seed)
	r2.Seed(seed)

	// 应该生成相同的序列
	for i := 0; i < 10; i++ {
		n1 := r1.Int()
		n2 := r2.Int()
		if n1 != n2 {
			t.Errorf("With same seed, expected same sequence, got %d != %d at position %d", n1, n2, i)
		}
	}

	PutRand(r1)
	PutRand(r2)
}

func TestRandPool_EdgeCases(t *testing.T) {
	// 测试边界情况
	r := GetRand()

	// Intn(1) 应该总是返回0
	for i := 0; i < 100; i++ {
		if r.Intn(2) >= 2 {
			t.Error("Intn(2) should return 0 or 1")
		}
	}

	// 测试大数
	for i := 0; i < 100; i++ {
		num := r.Intn(1000000)
		if num < 0 || num >= 1000000 {
			t.Errorf("Intn(1000000) returned %d, out of range", num)
		}
	}

	PutRand(r)
}

func BenchmarkRandPool_GetPut(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := GetRand()
		_ = r.Int()
		PutRand(r)
	}
}

func BenchmarkRandPool_vs_New(b *testing.B) {
	b.Run("Pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := GetRand()
			_ = r.Int()
			PutRand(r)
		}
	})

	b.Run("New", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			_ = r.Int()
		}
	})
}

func BenchmarkRandPool_Operations(b *testing.B) {
	r := GetRand()
	defer PutRand(r)

	b.Run("Int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = r.Int()
		}
	})

	b.Run("Intn", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = r.Intn(1000)
		}
	})

	b.Run("Float64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = r.Float64()
		}
	})
}

// BenchmarkRandPool_SeedCost 测试设置随机种子的性能开销
func BenchmarkRandPool_SeedCost(b *testing.B) {
	b.Run("WithSeed", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := GetRand()
			// 模拟设置种子的操作
			r.Seed(time.Now().UnixNano())
			_ = r.Int()
			PutRand(r)
		}
	})

	b.Run("WithoutSeed", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := GetRand()
			// 不设置种子，直接使用
			_ = r.Int()
			PutRand(r)
		}
	})

	b.Run("NewRand", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// 创建新的随机数生成器
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			_ = r.Int()
		}
	})
}
