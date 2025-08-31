package id

import (
	"crypto/rand"
	"strings"
	"testing"

	"gitee.com/MM-Q/go-kit/pool"
)

// 当前的UUID实现（使用math/rand）
func UUIDCurrent() string {
	r := pool.GetRand()
	defer pool.PutRand(r)

	var buf strings.Builder
	buf.Grow(36)

	// 8位
	for i := 0; i < 8; i++ {
		buf.WriteByte(chars[r.Intn(62)])
	}
	buf.WriteByte('-')

	// 4位
	for i := 0; i < 4; i++ {
		buf.WriteByte(chars[r.Intn(62)])
	}
	buf.WriteByte('-')

	// 4位
	for i := 0; i < 4; i++ {
		buf.WriteByte(chars[r.Intn(62)])
	}
	buf.WriteByte('-')

	// 4位
	for i := 0; i < 4; i++ {
		buf.WriteByte(chars[r.Intn(62)])
	}
	buf.WriteByte('-')

	// 12位
	for i := 0; i < 12; i++ {
		buf.WriteByte(chars[r.Intn(62)])
	}

	return buf.String()
}

// 使用crypto/rand的UUID实现
func UUIDCrypto() string {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		// 回退到当前实现
		return UUIDCurrent()
	}

	var buf strings.Builder
	buf.Grow(36)

	byteIndex := 0

	// 8位
	for i := 0; i < 8; i++ {
		buf.WriteByte(chars[randomBytes[byteIndex]%62])
		byteIndex++
	}
	buf.WriteByte('-')

	// 4位
	for i := 0; i < 4; i++ {
		buf.WriteByte(chars[randomBytes[byteIndex]%62])
		byteIndex++
	}
	buf.WriteByte('-')

	// 4位
	for i := 0; i < 4; i++ {
		buf.WriteByte(chars[randomBytes[byteIndex]%62])
		byteIndex++
	}
	buf.WriteByte('-')

	// 4位
	for i := 0; i < 4; i++ {
		buf.WriteByte(chars[randomBytes[byteIndex]%62])
		byteIndex++
	}
	buf.WriteByte('-')

	// 12位
	for i := 0; i < 12; i++ {
		buf.WriteByte(chars[randomBytes[byteIndex]%62])
		byteIndex++
	}

	return buf.String()
}

// 性能对比基准测试
func BenchmarkUUIDCurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUIDCurrent()
	}
}

func BenchmarkUUIDCrypto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUIDCrypto()
	}
}

// 并发性能测试
func BenchmarkUUIDCurrentParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			UUIDCurrent()
		}
	})
}

func BenchmarkUUIDCryptoParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			UUIDCrypto()
		}
	})
}

// 唯一性测试
func TestUUIDCryptoUniqueness(t *testing.T) {
	const numUUIDs = 10000
	uuids := make(map[string]bool)

	for i := 0; i < numUUIDs; i++ {
		uuid := UUIDCrypto()
		if uuids[uuid] {
			t.Errorf("Duplicate UUID found: %s", uuid)
		}
		uuids[uuid] = true
	}

	if len(uuids) != numUUIDs {
		t.Errorf("Expected %d unique UUIDs, got %d", numUUIDs, len(uuids))
	}
}

// 并发唯一性测试
func TestUUIDCryptoConcurrentUniqueness(t *testing.T) {
	const numGoroutines = 10
	const numIDsPerGoroutine = 100

	results := make(chan string, numGoroutines*numIDsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numIDsPerGoroutine; j++ {
				results <- UUIDCrypto()
			}
		}()
	}

	uuids := make(map[string]bool)
	duplicates := 0
	for i := 0; i < numGoroutines*numIDsPerGoroutine; i++ {
		uuid := <-results
		if uuids[uuid] {
			duplicates++
		}
		uuids[uuid] = true
	}

	totalIDs := numGoroutines * numIDsPerGoroutine
	uniqueRatio := float64(len(uuids)) / float64(totalIDs)

	t.Logf("Crypto UUID concurrent test: %d unique IDs out of %d (%.2f%% unique, %d duplicates)",
		len(uuids), totalIDs, uniqueRatio*100, duplicates)

	if uniqueRatio < 0.999 { // 期望99.9%以上的唯一性
		t.Errorf("Crypto UUID uniqueness too low: %.3f (duplicates: %d/%d)", uniqueRatio, duplicates, totalIDs)
	}
}
