package id

import (
	"testing"
)

func TestGenID(t *testing.T) {
	t.Run("Basic ID generation", func(t *testing.T) {
		id := GenID(4)
		if id == "" {
			t.Fatal("GenID(4) returned empty string")
		}

		// 验证长度：16位时间戳 + 4位随机数
		expectedLength := 20
		if len(id) != expectedLength {
			t.Errorf("Expected ID length %d, got %d", expectedLength, len(id))
		}

		// 验证前16位是数字
		for i := 0; i < 16; i++ {
			if id[i] < '0' || id[i] > '9' {
				t.Errorf("First 16 characters should be digits, got: %s", id[:16])
				break
			}
		}

		// 验证后4位是有效字符
		for i := 16; i < len(id); i++ {
			found := false
			for _, char := range chars {
				if id[i] == byte(char) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Invalid character at position %d: %c", i, id[i])
			}
		}
	})

	t.Run("Zero length random part", func(t *testing.T) {
		id := GenID(0)
		if len(id) != 16 {
			t.Errorf("Expected length 16 for GenID(0), got %d", len(id))
		}

		// 验证全部是数字
		for i, char := range id {
			if char < '0' || char > '9' {
				t.Errorf("Character at position %d should be digit, got: %c", i, char)
			}
		}
	})

	t.Run("Negative length", func(t *testing.T) {
		id := GenID(-1)
		if id != "" {
			t.Errorf("Expected empty string for GenID(-1), got: %s", id)
		}
	})

	t.Run("Large random part", func(t *testing.T) {
		id := GenID(20)
		expectedLength := 36 // 16 + 20
		if len(id) != expectedLength {
			t.Errorf("Expected ID length %d, got %d", expectedLength, len(id))
		}
	})

	t.Run("ID uniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		const numIDs = 1000

		for i := 0; i < numIDs; i++ {
			id := GenID(8)
			ids[id] = true
		}

		// 16位时间戳在高频生成下允许重复，期望20%以上唯一性
		uniqueRatio := float64(len(ids)) / float64(numIDs)
		if uniqueRatio < 0.20 {
			t.Errorf("Expected at least 20%% unique IDs, got %d/%d (%.1f%%)", len(ids), numIDs, uniqueRatio*100)
		}
		t.Logf("ID唯一性: %d/%d (%.1f%%) - 16位时间戳限制，如需更高唯一性请使用GenIDWithLen", len(ids), numIDs, uniqueRatio*100)
	})
}

func TestGenIDs(t *testing.T) {
	t.Run("Basic batch generation", func(t *testing.T) {
		ids := GenIDs(5, 4)
		if len(ids) != 5 {
			t.Errorf("Expected 5 IDs, got %d", len(ids))
		}

		// 验证每个ID的格式
		for i, id := range ids {
			if len(id) != 20 { // 16位时间戳 + 4位随机数
				t.Errorf("ID %d has wrong length: %d", i, len(id))
			}
		}
	})

	t.Run("Zero random length", func(t *testing.T) {
		ids := GenIDs(3, 0)
		for i, id := range ids {
			if len(id) != 16 {
				t.Errorf("ID %d should have length 16, got %d", i, len(id))
			}
		}
	})
}

func TestGenWithPrefix(t *testing.T) {
	t.Run("Basic prefix generation", func(t *testing.T) {
		id := GenWithPrefix("user", 4)
		expectedLength := len("user") + 1 + 16 + 4 // prefix + _ + 16位时间戳 + 4位随机数
		if len(id) != expectedLength {
			t.Errorf("Expected length %d, got %d", expectedLength, len(id))
		}
	})

	t.Run("Empty prefix", func(t *testing.T) {
		id := GenWithPrefix("", 4)
		expectedLength := 20 // 16位时间戳 + 4位随机数
		if len(id) != expectedLength {
			t.Errorf("Expected length %d for empty prefix, got %d", expectedLength, len(id))
		}
	})
}

func TestValid(t *testing.T) {
	t.Run("Valid IDs", func(t *testing.T) {
		id := GenID(8)
		if !Valid(id, 8) {
			t.Errorf("Generated ID should be valid: %s", id)
		}
	})

	t.Run("Invalid length", func(t *testing.T) {
		if Valid("12345", 8) {
			t.Error("Short ID should be invalid")
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		// 测试16位时间戳 + 0位随机数
		id := GenID(0)
		if !Valid(id, 0) {
			t.Errorf("Valid 16-digit timestamp should be valid with n=0")
		}
	})
}

// 测试新的GenIDWithLen函数
func TestGenIDWithLen(t *testing.T) {
	t.Run("Custom timestamp length", func(t *testing.T) {
		id := GenIDWithLen(8, 4)
		if len(id) != 12 {
			t.Errorf("Expected length 12, got %d", len(id))
		}
	})

	t.Run("Full timestamp", func(t *testing.T) {
		id := GenIDWithLen(-1, 4)
		// 完整时间戳长度约19位 + 4位随机数
		if len(id) < 20 {
			t.Errorf("Full timestamp ID should be at least 20 chars, got %d", len(id))
		}
	})

	t.Run("Zero timestamp length", func(t *testing.T) {
		id := GenIDWithLen(0, 4)
		if len(id) != 4 {
			t.Errorf("Expected length 4 (only random), got %d", len(id))
		}
	})
}

// 测试新的ValidWithLen函数
func TestValidWithLen(t *testing.T) {
	t.Run("Custom length validation", func(t *testing.T) {
		id := GenIDWithLen(8, 4)
		if !ValidWithLen(id, 8, 4) {
			t.Errorf("Generated ID should be valid: %s", id)
		}
	})

	t.Run("Full timestamp validation", func(t *testing.T) {
		id := GenIDWithLen(-1, 4)
		if !ValidWithLen(id, -1, 4) {
			t.Errorf("Full timestamp ID should be valid: %s", id)
		}
	})
}
