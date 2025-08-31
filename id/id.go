package id

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"gitee.com/MM-Q/go-kit/pool"
)

// 随机字符集(固定62位请勿修改)
const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenID 生成ID
// 用于生成带有时间戳和随机数的ID，格式为：时间戳(8位) + 随机数(n位)
//
// 参数:
//   - n: 随机部分长度
//
// 返回:
//   - 生成的ID
func GenID(n int) string {
	if n < 0 {
		return ""
	}

	// 8位时间戳
	ts := fmt.Sprintf("%08d", time.Now().UnixNano()%1e8)

	// 如果指定的随机部分长度为0，则只返回时间戳
	if n == 0 {
		return ts
	}

	// 获取随机数生成器
	r := pool.GetRand()
	defer pool.PutRand(r)

	return pool.WithString(8+n, func(buf *strings.Builder) {
		buf.WriteString(ts)
		// 生成n位随机数
		for i := 0; i < n; i++ {
			buf.WriteByte(chars[r.Intn(62)])
		}
	})
}

// GenIDs 批量生成ID
// 用于批量生成多个唯一ID，每个ID间会有纳秒级延迟确保唯一性
//
// 参数:
//   - count: 要生成的ID数量
//   - n: 每个ID随机部分的长度
//
// 返回:
//   - ID切片，参数无效时返回nil
func GenIDs(count, n int) []string {
	if count <= 0 || n < 0 {
		return nil
	}

	ids := make([]string, count)
	r := pool.GetRand()
	defer pool.PutRand(r)

	for i := 0; i < count; i++ {
		ts := fmt.Sprintf("%08d", time.Now().UnixNano()%1e8)

		if n == 0 {
			ids[i] = ts
		} else {
			ids[i] = pool.WithString(8+n, func(buf *strings.Builder) {
				buf.WriteString(ts)
				for j := 0; j < n; j++ {
					buf.WriteByte(chars[r.Intn(62)])
				}
			})
		}

		if i < count-1 {
			time.Sleep(time.Nanosecond)
		}
	}

	return ids
}

// GenWithPrefix 生成带前缀的ID
// 用于生成带有自定义前缀的ID，格式为：prefix_ID
//
// 参数:
//   - prefix: ID前缀字符串
//   - n: 随机部分长度
//
// 返回:
//   - 带前缀的ID字符串
func GenWithPrefix(prefix string, n int) string {
	if n < 0 {
		return prefix
	}

	id := GenID(n)
	if prefix == "" {
		return id
	}

	return pool.WithString(len(prefix)+len(id)+1, func(buf *strings.Builder) {
		buf.WriteString(prefix)
		buf.WriteByte('_')
		buf.WriteString(id)
	})
}

// Valid 验证ID格式
// 用于验证ID是否符合预期格式：前8位数字 + 后n位随机字符
//
// 参数:
//   - id: 要验证的ID字符串
//   - n: 预期的随机部分长度
//
// 返回:
//   - 格式正确返回true，否则返回false
func Valid(id string, n int) bool {
	if len(id) != 8+n {
		return false
	}

	// 检查时间戳部分
	for i := 0; i < 8; i++ {
		if id[i] < '0' || id[i] > '9' {
			return false
		}
	}

	// 检查随机部分
	for i := 8; i < len(id); i++ {
		found := false
		for j := 0; j < 62; j++ {
			if id[i] == chars[j] {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// UUID 生成类UUID格式
// 用于生成类似UUID的字符串，格式为：8-4-4-4-12
// 使用crypto/rand提供强随机性，确保并发安全和高唯一性
//
// 返回:
//   - 36位长度的UUID格式字符串
func UUID() string {
	// 从字节池获取32字节的缓冲区用于加密安全随机数据
	randomBytes := pool.GetByte(32)
	defer pool.PutByte(randomBytes)

	if _, err := rand.Read(randomBytes); err != nil {
		// 极少情况下crypto/rand失败时的fallback
		r := pool.GetRand()
		defer pool.PutRand(r)

		return pool.WithString(36, func(buf *strings.Builder) {
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
		})
	}

	// 使用crypto/rand生成的随机字节映射到字符集
	return pool.WithString(36, func(buf *strings.Builder) {
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
	})
}

// Short 生成短ID
// 用于生成基于当前纳秒时间戳的短ID
//
// 返回:
//   - 纳秒时间戳字符串
func Short() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Nano 生成纳秒级ID
// 用于生成基于当前纳秒时间戳的ID
//
// 返回:
//   - 纳秒时间戳字符串
func Nano() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
