package id

import (
	"testing"
)

// TestTimestampLengthLimit 测试时间戳长度限制
func TestTimestampLengthLimit(t *testing.T) {
	testCases := []struct {
		name        string
		tsLen       int
		expectedLen int
		description string
	}{
		{"正常长度", 8, 8, "8位时间戳"},
		{"最大长度", 16, 16, "16位时间戳(最大)"},
		{"超出限制", 20, 16, "请求20位但限制为16位"},
		{"超出限制", 25, 16, "请求25位但限制为16位"},
		{"完整时间戳", -1, 16, "完整微秒时间戳"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id := GenIDWithLen(tc.tsLen, 4) // 4位随机数

			// 分离时间戳部分
			if len(id) < 4 {
				t.Fatalf("ID长度异常: %s", id)
			}

			timestampPart := id[:len(id)-4]
			randomPart := id[len(id)-4:]

			t.Logf("%s: ID=%s, 时间戳=%s(%d位), 随机=%s",
				tc.description, id, timestampPart, len(timestampPart), randomPart)

			// 验证时间戳长度
			if len(timestampPart) != tc.expectedLen {
				t.Errorf("时间戳长度错误: 期望%d位, 实际%d位", tc.expectedLen, len(timestampPart))
			}

			// 验证时间戳都是数字
			for i, char := range timestampPart {
				if char < '0' || char > '9' {
					t.Errorf("时间戳第%d位不是数字: %c", i, char)
				}
			}

			// 验证随机部分长度
			if len(randomPart) != 4 {
				t.Errorf("随机部分长度错误: 期望4位, 实际%d位", len(randomPart))
			}
		})
	}
}

// TestLengthLimitConsistency 测试长度限制的一致性
func TestLengthLimitConsistency(t *testing.T) {
	// 测试多次调用结果的一致性
	for i := 17; i <= 25; i++ {
		id1 := GenIDWithLen(i, 0) // 只要时间戳
		id2 := GenIDWithLen(i, 0)

		if len(id1) != 16 || len(id2) != 16 {
			t.Errorf("长度%d被限制后应该都是16位: id1=%d位, id2=%d位", i, len(id1), len(id2))
		}

		t.Logf("请求%d位 -> 实际16位: %s", i, id1)
	}
}
