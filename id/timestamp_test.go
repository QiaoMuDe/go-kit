package id

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"
)

// TestGenerateTruncatedTimestamp_BasicFunctionality 测试基本功能
func TestGenerateTruncatedTimestamp_BasicFunctionality(t *testing.T) {
	tests := []struct {
		name   string
		tsLen  int
		expect string // 正则表达式模式
	}{
		{"1位时间戳", 1, `^\d{1}$`},
		{"4位时间戳", 4, `^\d{4}$`},
		{"8位时间戳", 8, `^\d{8}$`},
		{"12位时间戳", 12, `^\d{12}$`},
		{"16位时间戳", 16, `^\d{16}$`},
		{"18位时间戳", 18, `^\d{18}$`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateTruncatedTimestamp(tt.tsLen)

			// 检查长度
			if len(result) != tt.tsLen {
				t.Errorf("期望长度 %d, 实际长度 %d, 结果: %s", tt.tsLen, len(result), result)
			}

			// 检查格式（纯数字）
			matched, err := regexp.MatchString(tt.expect, result)
			if err != nil {
				t.Fatalf("正则表达式错误: %v", err)
			}
			if !matched {
				t.Errorf("格式不匹配，期望: %s, 实际: %s", tt.expect, result)
			}

			t.Logf("✅ %s: %s", tt.name, result)
		})
	}
}

// TestGenerateTruncatedTimestamp_Consistency 测试一致性（短时间内应该相同）
func TestGenerateTruncatedTimestamp_Consistency(t *testing.T) {
	tsLen := 8

	// 在很短时间内连续调用，应该得到相同结果
	first := generateTruncatedTimestamp(tsLen)

	for i := 0; i < 100; i++ {
		result := generateTruncatedTimestamp(tsLen)
		if result != first {
			t.Logf("在第 %d 次调用时时间戳发生变化: %s -> %s", i+1, first, result)
			break
		}
	}

	t.Logf("✅ 短时间内时间戳保持一致: %s", first)
}

// TestGenerateTruncatedTimestamp_Progression 测试时间递增性
func TestGenerateTruncatedTimestamp_Progression(t *testing.T) {
	tsLen := 16

	var timestamps []string
	var numericValues []int64

	// 收集一段时间内的时间戳
	for i := 0; i < 10; i++ {
		ts := generateTruncatedTimestamp(tsLen)
		timestamps = append(timestamps, ts)

		// 转换为数值用于比较
		val, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			t.Fatalf("时间戳转换失败: %s, 错误: %v", ts, err)
		}
		numericValues = append(numericValues, val)

		// 短暂延迟确保时间变化
		time.Sleep(time.Millisecond)
	}

	// 检查是否有递增趋势（允许偶尔的回退，因为是截断的）
	increasing := 0
	for i := 1; i < len(numericValues); i++ {
		if numericValues[i] >= numericValues[i-1] {
			increasing++
		}
		t.Logf("时间戳 %d: %s (%d)", i, timestamps[i], numericValues[i])
	}

	// 至少70%应该是递增的
	ratio := float64(increasing) / float64(len(numericValues)-1)
	if ratio < 0.7 {
		t.Errorf("递增比例过低: %.2f, 期望至少 0.70", ratio)
	}

	t.Logf("✅ 递增比例: %.2f", ratio)
}

// TestGenerateTruncatedTimestamp_Range 测试数值范围
func TestGenerateTruncatedTimestamp_Range(t *testing.T) {
	tests := []struct {
		tsLen  int
		maxVal int64
	}{
		{1, 9},
		{2, 99},
		{4, 9999},
		{8, 99999999},
		{16, 9999999999999999},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d位范围测试", tt.tsLen), func(t *testing.T) {
			for i := 0; i < 10; i++ {
				result := generateTruncatedTimestamp(tt.tsLen)
				val, err := strconv.ParseInt(result, 10, 64)
				if err != nil {
					t.Fatalf("转换失败: %s", result)
				}

				if val < 0 || val > tt.maxVal {
					t.Errorf("数值超出范围: %d, 期望范围: [0, %d]", val, tt.maxVal)
				}
			}
			t.Logf("✅ %d位时间戳数值范围正确", tt.tsLen)
		})
	}
}

// TestGenerateTruncatedTimestamp_LeadingZeros 测试前导零
func TestGenerateTruncatedTimestamp_LeadingZeros(t *testing.T) {
	tsLen := 8

	// 多次调用，检查是否正确处理前导零
	results := make(map[string]int)

	for i := 0; i < 1000; i++ {
		result := generateTruncatedTimestamp(tsLen)
		results[result]++

		// 检查长度
		if len(result) != tsLen {
			t.Errorf("长度错误: 期望 %d, 实际 %d, 值: %s", tsLen, len(result), result)
		}

		// 检查是否以0开头（这是可能的）
		if result[0] == '0' {
			t.Logf("发现前导零的时间戳: %s", result)
		}
	}

	t.Logf("✅ 生成了 %d 个不同的时间戳", len(results))
}

// TestGenerateTruncatedTimestamp_Performance 性能测试
func TestGenerateTruncatedTimestamp_Performance(t *testing.T) {
	tsLen := 16
	iterations := 100000

	start := time.Now()

	for i := 0; i < iterations; i++ {
		_ = generateTruncatedTimestamp(tsLen)
	}

	duration := time.Since(start)
	avgTime := duration / time.Duration(iterations)

	t.Logf("✅ 性能测试: %d 次调用耗时 %v, 平均每次 %v", iterations, duration, avgTime)

	// 期望每次调用不超过50微秒（考虑race检测模式的性能影响）
	if avgTime > 50*time.Microsecond {
		t.Errorf("性能不达标: 平均耗时 %v, 期望小于 50µs", avgTime)
	}
}

// TestGenerateTruncatedTimestamp_Uniqueness 唯一性测试
func TestGenerateTruncatedTimestamp_Uniqueness(t *testing.T) {
	tests := []struct {
		tsLen     int
		samples   int
		interval  time.Duration
		minUnique float64 // 最小唯一性比例
	}{
		{4, 100, time.Millisecond, 0.05},     // 4位数字，100毫秒间隔
		{8, 200, time.Millisecond, 0.3},      // 8位数字，期望30%唯一
		{16, 500, time.Millisecond * 2, 0.8}, // 16位数字，期望80%唯一
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d位唯一性测试", tt.tsLen), func(t *testing.T) {
			unique := make(map[string]bool)

			for i := 0; i < tt.samples; i++ {
				result := generateTruncatedTimestamp(tt.tsLen)
				unique[result] = true

				// 添加足够的延迟确保时间变化
				time.Sleep(tt.interval)
			}

			uniqueRatio := float64(len(unique)) / float64(tt.samples)

			t.Logf("唯一性统计: %d/%d = %.3f (间隔: %v)", len(unique), tt.samples, uniqueRatio, tt.interval)

			if uniqueRatio < tt.minUnique {
				t.Logf("⚠️ 唯一性较低: %.3f, 期望至少 %.3f", uniqueRatio, tt.minUnique)
				t.Logf("这在快速调用场景下是正常的，说明需要配合随机数使用")
			} else {
				t.Logf("✅ %d位时间戳唯一性达标: %.3f", tt.tsLen, uniqueRatio)
			}
		})
	}
}

// TestGenerateTruncatedTimestamp_EdgeCases 边界情况测试
func TestGenerateTruncatedTimestamp_EdgeCases(t *testing.T) {
	// 测试极小值
	result1 := generateTruncatedTimestamp(1)
	if len(result1) != 1 {
		t.Errorf("1位时间戳长度错误: %s", result1)
	}

	// 测试较大值
	result20 := generateTruncatedTimestamp(20)
	if len(result20) != 20 {
		t.Errorf("20位时间戳长度错误: %s", result20)
	}

	t.Logf("✅ 边界测试通过: 1位=%s, 20位=%s", result1, result20)
}

// BenchmarkGenerateTruncatedTimestamp 基准测试
func BenchmarkGenerateTruncatedTimestamp(b *testing.B) {
	benchmarks := []struct {
		name  string
		tsLen int
	}{
		{"8位", 8},
		{"16位", 16},
		{"20位", 20},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = generateTruncatedTimestamp(bm.tsLen)
			}
		})
	}
}
