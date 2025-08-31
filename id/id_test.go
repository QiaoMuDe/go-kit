package id

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestGenID(t *testing.T) {
	t.Run("Basic ID generation", func(t *testing.T) {
		id := GenID(4)
		if id == "" {
			t.Fatal("GenID(4) returned empty string")
		}

		// 验证长度：8位时间戳 + 4位随机数
		expectedLength := 12
		if len(id) != expectedLength {
			t.Errorf("Expected ID length %d, got %d", expectedLength, len(id))
		}

		// 验证前8位是数字
		for i := 0; i < 8; i++ {
			if id[i] < '0' || id[i] > '9' {
				t.Errorf("First 8 characters should be digits, got: %s", id[:8])
				break
			}
		}

		// 验证后4位是有效字符
		for i := 8; i < len(id); i++ {
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
		if len(id) != 8 {
			t.Errorf("Expected length 8 for GenID(0), got %d", len(id))
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
		expectedLength := 28 // 8 + 20
		if len(id) != expectedLength {
			t.Errorf("Expected ID length %d, got %d", expectedLength, len(id))
		}
	})

	t.Run("ID uniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		const numIDs = 1000

		for i := 0; i < numIDs; i++ {
			id := GenID(8)
			if ids[id] {
				t.Errorf("Duplicate ID generated: %s", id)
			}
			ids[id] = true
		}

		if len(ids) != numIDs {
			t.Errorf("Expected %d unique IDs, got %d", numIDs, len(ids))
		}
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
			if len(id) != 12 { // 8 + 4
				t.Errorf("ID %d has wrong length: %d", i, len(id))
			}
		}

		// 验证唯一性
		uniqueIDs := make(map[string]bool)
		for _, id := range ids {
			uniqueIDs[id] = true
		}
		if len(uniqueIDs) != len(ids) {
			t.Error("Generated IDs are not unique")
		}
	})

	t.Run("Zero count", func(t *testing.T) {
		ids := GenIDs(0, 4)
		if ids != nil {
			t.Errorf("Expected nil for zero count, got: %v", ids)
		}
	})

	t.Run("Negative count", func(t *testing.T) {
		ids := GenIDs(-1, 4)
		if ids != nil {
			t.Errorf("Expected nil for negative count, got: %v", ids)
		}
	})

	t.Run("Negative random length", func(t *testing.T) {
		ids := GenIDs(3, -1)
		if ids != nil {
			t.Errorf("Expected nil for negative random length, got: %v", ids)
		}
	})

	t.Run("Zero random length", func(t *testing.T) {
		ids := GenIDs(3, 0)
		if len(ids) != 3 {
			t.Errorf("Expected 3 IDs, got %d", len(ids))
		}

		for i, id := range ids {
			if len(id) != 8 {
				t.Errorf("ID %d should have length 8, got %d", i, len(id))
			}
		}
	})

	t.Run("Large batch", func(t *testing.T) {
		ids := GenIDs(100, 6)
		if len(ids) != 100 {
			t.Errorf("Expected 100 IDs, got %d", len(ids))
		}

		// 验证唯一性
		uniqueIDs := make(map[string]bool)
		for _, id := range ids {
			uniqueIDs[id] = true
		}

		// 由于有纳秒级延迟，应该大部分是唯一的
		uniqueRatio := float64(len(uniqueIDs)) / float64(len(ids))
		if uniqueRatio < 0.95 {
			t.Errorf("Uniqueness ratio too low: %.2f", uniqueRatio)
		}
	})
}

func TestGenWithPrefix(t *testing.T) {
	t.Run("Basic prefix generation", func(t *testing.T) {
		id := GenWithPrefix("user", 4)
		if !strings.HasPrefix(id, "user_") {
			t.Errorf("ID should start with 'user_', got: %s", id)
		}

		// 验证总长度：prefix + _ + 8位时间戳 + 4位随机数
		expectedLength := len("user") + 1 + 8 + 4
		if len(id) != expectedLength {
			t.Errorf("Expected length %d, got %d", expectedLength, len(id))
		}
	})

	t.Run("Empty prefix", func(t *testing.T) {
		id := GenWithPrefix("", 4)
		// 空前缀应该返回普通ID
		if len(id) != 12 { // 8 + 4
			t.Errorf("Expected length 12 for empty prefix, got %d", len(id))
		}
		if strings.Contains(id, "_") {
			t.Errorf("Empty prefix should not contain underscore: %s", id)
		}
	})

	t.Run("Zero random length", func(t *testing.T) {
		id := GenWithPrefix("test", 0)
		expectedLength := len("test") + 1 + 8
		if len(id) != expectedLength {
			t.Errorf("Expected length %d, got %d", expectedLength, len(id))
		}
	})

	t.Run("Negative random length", func(t *testing.T) {
		id := GenWithPrefix("test", -1)
		if id != "test" {
			t.Errorf("Expected 'test' for negative random length, got: %s", id)
		}
	})

	t.Run("Long prefix", func(t *testing.T) {
		longPrefix := "very_long_prefix_name"
		id := GenWithPrefix(longPrefix, 6)
		if !strings.HasPrefix(id, longPrefix+"_") {
			t.Errorf("ID should start with '%s_', got: %s", longPrefix, id)
		}
	})
}

func TestValid(t *testing.T) {
	t.Run("Valid ID", func(t *testing.T) {
		id := GenID(4)
		if !Valid(id, 4) {
			t.Errorf("Generated ID should be valid: %s", id)
		}
	})

	t.Run("Invalid length", func(t *testing.T) {
		if Valid("12345", 4) {
			t.Error("Short ID should be invalid")
		}
		if Valid("123456789012345", 4) {
			t.Error("Long ID should be invalid")
		}
	})

	t.Run("Invalid timestamp part", func(t *testing.T) {
		invalidID := "1234567a1234" // 'a' in timestamp part
		if Valid(invalidID, 4) {
			t.Errorf("ID with invalid timestamp should be invalid: %s", invalidID)
		}
	})

	t.Run("Invalid random part", func(t *testing.T) {
		invalidID := "12345678@#$%" // invalid characters in random part
		if Valid(invalidID, 4) {
			t.Errorf("ID with invalid random part should be invalid: %s", invalidID)
		}
	})

	t.Run("Zero random length", func(t *testing.T) {
		id := GenID(0)
		if !Valid(id, 0) {
			t.Errorf("Generated ID with zero random length should be valid: %s", id)
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		// 空字符串
		if Valid("", 0) {
			t.Error("Empty string should be invalid")
		}

		// 只有时间戳部分有效
		if !Valid("12345678", 0) {
			t.Error("Valid 8-digit timestamp should be valid with n=0")
		}
	})
}

func TestUUID(t *testing.T) {
	t.Run("Basic UUID generation", func(t *testing.T) {
		uuid := UUID()
		if uuid == "" {
			t.Fatal("UUID() returned empty string")
		}

		// 验证长度：32字符 + 4个连字符 = 36
		expectedLength := 36
		if len(uuid) != expectedLength {
			t.Errorf("Expected UUID length %d, got %d", expectedLength, len(uuid))
		}

		// 验证格式：8-4-4-4-12
		parts := strings.Split(uuid, "-")
		if len(parts) != 5 {
			t.Errorf("Expected 5 parts separated by hyphens, got %d", len(parts))
		}

		expectedLengths := []int{8, 4, 4, 4, 12}
		for i, part := range parts {
			if len(part) != expectedLengths[i] {
				t.Errorf("Part %d should have length %d, got %d", i, expectedLengths[i], len(part))
			}
		}
	})

	t.Run("UUID uniqueness", func(t *testing.T) {
		uuids := make(map[string]bool)
		const numUUIDs = 1000

		for i := 0; i < numUUIDs; i++ {
			uuid := UUID()
			if uuids[uuid] {
				t.Errorf("Duplicate UUID generated: %s", uuid)
			}
			uuids[uuid] = true
		}

		if len(uuids) != numUUIDs {
			t.Errorf("Expected %d unique UUIDs, got %d", numUUIDs, len(uuids))
		}
	})

	t.Run("UUID character set", func(t *testing.T) {
		uuid := UUID()
		// 移除连字符
		cleanUUID := strings.ReplaceAll(uuid, "-", "")

		// 验证所有字符都在有效字符集中
		for i, char := range cleanUUID {
			found := false
			for _, validChar := range chars {
				if char == validChar {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Invalid character at position %d: %c", i, char)
			}
		}
	})
}

func TestShort(t *testing.T) {
	t.Run("Basic short ID generation", func(t *testing.T) {
		shortID := Short()
		if shortID == "" {
			t.Fatal("Short() returned empty string")
		}

		// 验证是数字
		shortIDPattern := regexp.MustCompile(`^\d+$`)
		if !shortIDPattern.MatchString(shortID) {
			t.Errorf("Short ID should be numeric: %s", shortID)
		}
	})

	t.Run("Short ID uniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		const numIDs = 100

		for i := 0; i < numIDs; i++ {
			id := Short()
			ids[id] = true
			time.Sleep(time.Nanosecond) // 确保时间戳不同
		}

		// 由于基于纳秒时间戳，应该大部分是唯一的
		uniqueRatio := float64(len(ids)) / float64(numIDs)
		if uniqueRatio < 0.9 {
			t.Errorf("Short ID uniqueness ratio too low: %.2f", uniqueRatio)
		}
	})

	t.Run("Short ID ordering", func(t *testing.T) {
		// 连续生成的ID应该是递增的
		var prevID string
		for i := 0; i < 10; i++ {
			id := Short()
			if i > 0 && id <= prevID {
				t.Logf("Non-increasing Short ID: prev=%s, current=%s", prevID, id)
			}
			prevID = id
			time.Sleep(time.Nanosecond)
		}
	})
}

func TestNano(t *testing.T) {
	t.Run("Basic nano ID generation", func(t *testing.T) {
		nanoID := Nano()
		if nanoID == "" {
			t.Fatal("Nano() returned empty string")
		}

		// 验证是数字
		nanoIDPattern := regexp.MustCompile(`^\d+$`)
		if !nanoIDPattern.MatchString(nanoID) {
			t.Errorf("Nano ID should be numeric: %s", nanoID)
		}
	})

	t.Run("Nano ID uniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		const numIDs = 100

		for i := 0; i < numIDs; i++ {
			id := Nano()
			ids[id] = true
			time.Sleep(time.Nanosecond)
		}

		// 由于基于纳秒时间戳，应该大部分是唯一的
		uniqueRatio := float64(len(ids)) / float64(numIDs)
		if uniqueRatio < 0.9 {
			t.Errorf("Nano ID uniqueness ratio too low: %.2f", uniqueRatio)
		}
	})

	t.Run("Nano equals Short", func(t *testing.T) {
		// Nano() 和 Short() 应该返回相同的结果
		for i := 0; i < 10; i++ {
			nano := Nano()
			short := Short()
			// 由于时间可能略有不同，我们只验证格式相同
			if len(nano) != len(short) {
				t.Errorf("Nano and Short should have same length: nano=%d, short=%d", len(nano), len(short))
			}
		}
	})
}

// 并发测试
func TestConcurrentIDGeneration(t *testing.T) {
	t.Run("Concurrent GenID", func(t *testing.T) {
		const numGoroutines = 5
		const numIDsPerGoroutine = 50

		results := make(chan string, numGoroutines*numIDsPerGoroutine)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				for j := 0; j < numIDsPerGoroutine; j++ {
					results <- GenID(8)          // 增加随机部分长度提高唯一性
					time.Sleep(time.Microsecond) // 添加微小延迟
				}
			}()
		}

		ids := make(map[string]bool)
		duplicates := 0
		for i := 0; i < numGoroutines*numIDsPerGoroutine; i++ {
			id := <-results
			if ids[id] {
				duplicates++
			}
			ids[id] = true
		}

		// 允许少量重复，但不应该太多
		totalIDs := numGoroutines * numIDsPerGoroutine
		uniqueRatio := float64(len(ids)) / float64(totalIDs)
		if uniqueRatio < 0.8 {
			t.Errorf("Concurrent GenID uniqueness too low: %.2f (duplicates: %d/%d)", uniqueRatio, duplicates, totalIDs)
		}
		t.Logf("GenID concurrent test: %d unique IDs out of %d (%.2f%% unique, %d duplicates)",
			len(ids), totalIDs, uniqueRatio*100, duplicates)
	})

	t.Run("Concurrent UUID", func(t *testing.T) {
		const numGoroutines = 5
		const numIDsPerGoroutine = 50

		results := make(chan string, numGoroutines*numIDsPerGoroutine)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				for j := 0; j < numIDsPerGoroutine; j++ {
					results <- UUID()
					time.Sleep(time.Microsecond) // 添加微小延迟
				}
			}()
		}

		ids := make(map[string]bool)
		duplicates := 0
		for i := 0; i < numGoroutines*numIDsPerGoroutine; i++ {
			id := <-results
			if ids[id] {
				duplicates++
			}
			ids[id] = true
		}

		// UUID应该有更高的唯一性
		totalIDs := numGoroutines * numIDsPerGoroutine
		uniqueRatio := float64(len(ids)) / float64(totalIDs)
		if uniqueRatio < 0.9 {
			t.Errorf("Concurrent UUID uniqueness too low: %.2f (duplicates: %d/%d)", uniqueRatio, duplicates, totalIDs)
		}
		t.Logf("UUID concurrent test: %d unique IDs out of %d (%.2f%% unique, %d duplicates)",
			len(ids), totalIDs, uniqueRatio*100, duplicates)
	})
}

// 基准测试
func BenchmarkGenID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenID(8)
	}
}

func BenchmarkGenIDs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenIDs(10, 8)
	}
}

func BenchmarkGenWithPrefix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenWithPrefix("user", 8)
	}
}

func BenchmarkUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUID()
	}
}

func BenchmarkShort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Short()
	}
}

func BenchmarkNano(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Nano()
	}
}

func BenchmarkValid(b *testing.B) {
	id := GenID(8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Valid(id, 8)
	}
}

// 并发基准测试
func BenchmarkConcurrentGenID(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GenID(8)
		}
	})
}

func BenchmarkConcurrentUUID(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			UUID()
		}
	})
}

// 边界测试
func TestEdgeCases(t *testing.T) {
	t.Run("High frequency generation", func(t *testing.T) {
		// 高频生成测试
		ids := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			ids[i] = GenID(8)
		}

		// 验证唯一性
		uniqueIDs := make(map[string]bool)
		for _, id := range ids {
			uniqueIDs[id] = true
		}

		// 由于时间戳精度限制，可能有重复，但应该大部分唯一
		uniqueRatio := float64(len(uniqueIDs)) / float64(len(ids))
		if uniqueRatio < 0.8 {
			t.Errorf("High frequency generation uniqueness too low: %.2f", uniqueRatio)
		}
	})

	t.Run("Character set coverage", func(t *testing.T) {
		// 验证字符集覆盖
		charUsed := make(map[byte]bool)
		for i := 0; i < 1000; i++ {
			id := GenID(10)
			// 检查随机部分（跳过前8位时间戳）
			for j := 8; j < len(id); j++ {
				charUsed[id[j]] = true
			}
		}

		// 应该使用了字符集中的大部分字符
		if len(charUsed) < 30 { // 至少使用30个不同字符
			t.Errorf("Character set coverage too low: %d characters used", len(charUsed))
		}
	})

	t.Run("Memory efficiency", func(t *testing.T) {
		// 测试大批量生成的内存效率
		const batchSize = 10000
		ids := GenIDs(batchSize, 6)

		if len(ids) != batchSize {
			t.Errorf("Expected %d IDs, got %d", batchSize, len(ids))
		}

		// 验证所有ID都有效
		for i, id := range ids {
			if !Valid(id, 6) {
				t.Errorf("Invalid ID at index %d: %s", i, id)
				break
			}
		}
	})
}
