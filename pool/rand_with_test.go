package pool

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestWithRand(t *testing.T) {
	t.Run("Generate random int", func(t *testing.T) {
		num := WithRand(func(rng *rand.Rand) int {
			return rng.Intn(100)
		})

		if num < 0 || num >= 100 {
			t.Errorf("Expected number in range [0, 100), got %d", num)
		}
	})

	t.Run("Generate random string", func(t *testing.T) {
		str := WithRand(func(rng *rand.Rand) string {
			return fmt.Sprintf("id_%d", rng.Int63())
		})

		if len(str) < 4 { // 至少包含 "id_"
			t.Errorf("Expected non-empty string with prefix, got %q", str)
		}

		if str[:3] != "id_" {
			t.Errorf("Expected string to start with 'id_', got %q", str)
		}
	})

	t.Run("Generate random slice", func(t *testing.T) {
		nums := WithRand(func(rng *rand.Rand) []int {
			result := make([]int, 5)
			for i := range result {
				result[i] = rng.Intn(10)
			}
			return result
		})

		if len(nums) != 5 {
			t.Errorf("Expected slice length 5, got %d", len(nums))
		}

		for i, num := range nums {
			if num < 0 || num >= 10 {
				t.Errorf("Expected number in range [0, 10) at index %d, got %d", i, num)
			}
		}
	})

	t.Run("Generate random float", func(t *testing.T) {
		f := WithRand(func(rng *rand.Rand) float64 {
			return rng.Float64()
		})

		if f < 0.0 || f >= 1.0 {
			t.Errorf("Expected float in range [0.0, 1.0), got %f", f)
		}
	})

	t.Run("Generate random bool", func(t *testing.T) {
		b := WithRand(func(rng *rand.Rand) bool {
			return rng.Intn(2) == 1
		})

		// 布尔值只能是true或false，这里只是确保没有panic
		_ = b
	})
}

func TestWithRandSeed(t *testing.T) {
	t.Run("Reproducible sequence", func(t *testing.T) {
		seed := int64(12345)

		// 生成第一个序列
		nums1 := WithRandSeed(seed, func(rng *rand.Rand) []int {
			result := make([]int, 10)
			for i := range result {
				result[i] = rng.Intn(100)
			}
			return result
		})

		// 使用相同种子生成第二个序列
		nums2 := WithRandSeed(seed, func(rng *rand.Rand) []int {
			result := make([]int, 10)
			for i := range result {
				result[i] = rng.Intn(100)
			}
			return result
		})

		// 两个序列应该完全相同
		if len(nums1) != len(nums2) {
			t.Errorf("Sequences have different lengths: %d vs %d", len(nums1), len(nums2))
		}

		for i := range nums1 {
			if nums1[i] != nums2[i] {
				t.Errorf("Sequences differ at index %d: %d vs %d", i, nums1[i], nums2[i])
			}
		}
	})

	t.Run("Different seeds produce different sequences", func(t *testing.T) {
		nums1 := WithRandSeed(12345, func(rng *rand.Rand) []int {
			result := make([]int, 10)
			for i := range result {
				result[i] = rng.Intn(100)
			}
			return result
		})

		nums2 := WithRandSeed(54321, func(rng *rand.Rand) []int {
			result := make([]int, 10)
			for i := range result {
				result[i] = rng.Intn(100)
			}
			return result
		})

		// 不同种子应该产生不同的序列（虽然理论上可能相同，但概率极低）
		different := false
		for i := range nums1 {
			if nums1[i] != nums2[i] {
				different = true
				break
			}
		}

		if !different {
			t.Log("Warning: Different seeds produced identical sequences (very unlikely but possible)")
		}
	})

	t.Run("Complex data structure", func(t *testing.T) {
		type TestData struct {
			ID    int
			Value float64
			Name  string
		}

		data := WithRandSeed(99999, func(rng *rand.Rand) TestData {
			return TestData{
				ID:    rng.Intn(1000),
				Value: rng.Float64() * 100,
				Name:  fmt.Sprintf("user_%d", rng.Intn(10000)),
			}
		})

		if data.ID < 0 || data.ID >= 1000 {
			t.Errorf("Expected ID in range [0, 1000), got %d", data.ID)
		}

		if data.Value < 0 || data.Value >= 100 {
			t.Errorf("Expected Value in range [0, 100), got %f", data.Value)
		}

		if len(data.Name) < 5 { // 至少包含 "user_"
			t.Errorf("Expected non-empty name with prefix, got %q", data.Name)
		}
	})
}

// 基准测试对比传统方式和便捷方式
func BenchmarkTraditionalRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rng := GetRand()
		result := rng.Intn(100)
		PutRand(rng)
		_ = result
	}
}

func BenchmarkWithRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := WithRand(func(rng *rand.Rand) int {
			return rng.Intn(100)
		})
		_ = result
	}
}

func BenchmarkTraditionalRandSeed(b *testing.B) {
	seed := int64(12345)
	for i := 0; i < b.N; i++ {
		rng := GetRandWithSeed(seed)
		result := rng.Intn(100)
		PutRand(rng)
		_ = result
	}
}

func BenchmarkWithRandSeed(b *testing.B) {
	seed := int64(12345)
	for i := 0; i < b.N; i++ {
		result := WithRandSeed(seed, func(rng *rand.Rand) int {
			return rng.Intn(100)
		})
		_ = result
	}
}

// 并发安全测试
func TestWithRandConcurrent(t *testing.T) {
	const numGoroutines = 10
	const numOperations = 100

	results := make(chan int, numGoroutines*numOperations)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numOperations; j++ {
				num := WithRand(func(rng *rand.Rand) int {
					return rng.Intn(1000)
				})
				results <- num
			}
		}()
	}

	// 收集所有结果
	for i := 0; i < numGoroutines*numOperations; i++ {
		result := <-results
		if result < 0 || result >= 1000 {
			t.Errorf("Got invalid result from concurrent operation: %d", result)
		}
	}
}

func TestWithRandSeedConcurrent(t *testing.T) {
	const numGoroutines = 10
	const numOperations = 100

	results := make(chan []int, numGoroutines*numOperations)

	for i := 0; i < numGoroutines; i++ {
		go func(seed int64) {
			for j := 0; j < numOperations; j++ {
				nums := WithRandSeed(seed, func(rng *rand.Rand) []int {
					result := make([]int, 3)
					for k := range result {
						result[k] = rng.Intn(100)
					}
					return result
				})
				results <- nums
			}
		}(int64(i))
	}

	// 收集所有结果
	for i := 0; i < numGoroutines*numOperations; i++ {
		result := <-results
		if len(result) != 3 {
			t.Errorf("Got invalid result length from concurrent operation: %d", len(result))
		}
		for j, num := range result {
			if num < 0 || num >= 100 {
				t.Errorf("Got invalid number at index %d from concurrent operation: %d", j, num)
			}
		}
	}
}

// 测试资源正确归还
func TestRandResourceManagement(t *testing.T) {
	// 测试正常情况下资源归还
	t.Run("Normal return", func(t *testing.T) {
		result := WithRand(func(rng *rand.Rand) int {
			return rng.Intn(100)
		})

		if result < 0 || result >= 100 {
			t.Errorf("Expected valid result, got %d", result)
		}
	})

	// 测试panic情况下资源归还
	t.Run("Panic recovery", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic but didn't get one")
			}
		}()

		WithRand(func(rng *rand.Rand) int {
			panic("test panic")
		})
	})
}

// 实际使用场景测试
func TestRandRealWorldUsage(t *testing.T) {
	t.Run("Generate random ID", func(t *testing.T) {
		id := WithRand(func(rng *rand.Rand) string {
			return fmt.Sprintf("user_%d_%d", rng.Int63(), rng.Intn(1000))
		})

		if len(id) < 6 { // 至少包含 "user_"
			t.Errorf("Expected valid ID, got %q", id)
		}
	})

	t.Run("Generate test data", func(t *testing.T) {
		testData := WithRandSeed(42, func(rng *rand.Rand) map[string]interface{} {
			return map[string]interface{}{
				"name":   fmt.Sprintf("test_user_%d", rng.Intn(1000)),
				"age":    rng.Intn(100),
				"score":  rng.Float64() * 100,
				"active": rng.Intn(2) == 1,
			}
		})

		if testData["name"] == nil || testData["age"] == nil {
			t.Error("Expected complete test data")
		}
	})

	t.Run("Shuffle slice", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5}
		shuffled := WithRand(func(rng *rand.Rand) []int {
			result := make([]int, len(original))
			copy(result, original)
			rng.Shuffle(len(result), func(i, j int) {
				result[i], result[j] = result[j], result[i]
			})
			return result
		})

		if len(shuffled) != len(original) {
			t.Errorf("Expected same length, got %d vs %d", len(shuffled), len(original))
		}

		// 验证所有元素都存在（虽然顺序可能不同）
		counts := make(map[int]int)
		for _, v := range shuffled {
			counts[v]++
		}

		for _, v := range original {
			if counts[v] != 1 {
				t.Errorf("Element %d appears %d times in shuffled slice", v, counts[v])
			}
		}
	})
}
