package pool

import (
	"testing"
)

// 基准测试对比传统方式和新方式
func BenchmarkTraditionalByte(b *testing.B) {
	testData := []byte("Hello World Test Data")

	for i := 0; i < b.N; i++ {
		buf := GetByte()
		copy(buf, testData)
		result := make([]byte, len(testData))
		copy(result, buf[:len(testData)])
		PutByte(buf)
		_ = result
	}
}

func BenchmarkTraditionalEmptyByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := GetEmptyByte(64)
		buf = append(buf, []byte("Hello")...)
		buf = append(buf, ' ')
		buf = append(buf, []byte("World")...)
		result := make([]byte, len(buf))
		copy(result, buf)
		PutByte(buf)
		_ = result
	}
}
