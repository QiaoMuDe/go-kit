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

// generateTruncatedTimestamp 生成指定长度的截断时间戳
//
// 参数:
//   - tsLen: 时间戳长度，必须大于0
//
// 返回:
//   - 截断后的时间戳字符串
func generateTruncatedTimestamp(tsLen int) string {
	// 快速失败：参数验证
	if tsLen <= 0 {
		return ""
	}

	// 计算模数：10^tsLen
	tsMod := int64(1)
	for i := 0; i < tsLen; i++ {
		tsMod *= 10
	}

	// 生成格式化字符串并应用
	tsFormat := fmt.Sprintf("%%0%dd", tsLen)
	return fmt.Sprintf(tsFormat, time.Now().UnixNano()%tsMod)
}

// genIDInternal 内部ID生成方法
// 用于生成带有时间戳和随机数的ID，支持可配置的时间戳长度
//
// 参数:
//   - tsLen: 时间戳长度，-1表示使用完整时间戳
//   - randLen: 随机部分长度
//
// 返回:
//   - 生成的ID
func genIDInternal(tsLen, randLen int) string {
	// 快速失败：参数验证
	if randLen < 0 {
		return ""
	}

	// 生成时间戳部分
	var ts string
	switch {
	case tsLen == -1:
		// 使用完整纳秒时间戳
		ts = fmt.Sprintf("%d", time.Now().UnixNano())
	case tsLen <= 0:
		// 时间戳长度无效，不生成时间戳部分
		ts = ""
	default:
		// 使用指定长度的时间戳
		ts = generateTruncatedTimestamp(tsLen)
	}

	// 快速返回：如果只需要时间戳
	if randLen == 0 {
		return ts
	}

	// 计算总长度并生成随机部分
	totalLen := len(ts) + randLen
	r := pool.GetRand()
	defer pool.PutRand(r)

	return pool.WithString(totalLen, func(buf *strings.Builder) {
		buf.WriteString(ts)
		// 生成随机数部分
		for i := 0; i < randLen; i++ {
			buf.WriteByte(chars[r.Intn(62)])
		}
	})

}

// GenID 生成ID
// 用于生成带有时间戳和随机数的ID，格式为：时间戳(16位) + 随机数(n位)
// 默认使用16位时间戳，提供更好的唯一性保证
//
// 参数:
//   - n: 随机部分长度
//
// 返回:
//   - 生成的ID
func GenID(n int) string {
	return genIDInternal(16, n)
}

// GenIDWithLen 生成指定长度的ID
// 用于生成带有自定义时间戳和随机数长度的ID
//
// 参数:
//   - tsLen: 时间戳长度，-1表示使用完整时间戳
//   - randLen: 随机部分长度
//
// 返回:
//   - 生成的ID
func GenIDWithLen(tsLen, randLen int) string {
	return genIDInternal(tsLen, randLen)
}

// GenIDs 批量生成ID
// 用于批量生成多个唯一ID，使用16位时间戳提供更好的唯一性
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
	for i := 0; i < count; i++ {
		ids[i] = genIDInternal(16, n)
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
// 用于验证ID是否符合预期格式：前16位数字 + 后n位随机字符
// 注意：此函数验证的是GenID()生成的默认格式(16位时间戳)
//
// 参数:
//   - id: 要验证的ID字符串
//   - n: 预期的随机部分长度
//
// 返回:
//   - 格式正确返回true，否则返回false
func Valid(id string, n int) bool {
	if len(id) != 16+n {
		return false
	}

	// 检查时间戳部分(前16位)
	for i := 0; i < 16; i++ {
		if id[i] < '0' || id[i] > '9' {
			return false
		}
	}

	// 检查随机部分
	for i := 16; i < len(id); i++ {
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

// ValidWithLen 验证指定长度的ID格式
// 用于验证ID是否符合指定的时间戳和随机数长度格式
//
// 参数:
//   - id: 要验证的ID字符串
//   - tsLen: 预期的时间戳长度，-1表示完整时间戳(纯数字)
//   - randLen: 预期的随机部分长度
//
// 返回:
//   - 格式正确返回true，否则返回false
func ValidWithLen(id string, tsLen, randLen int) bool {
	if tsLen == -1 {
		// 完整时间戳模式：检查是否以数字开头，后面跟随机字符
		if len(id) < randLen {
			return false
		}

		tsActualLen := len(id) - randLen
		if tsActualLen <= 0 {
			return false
		}

		// 检查时间戳部分(纯数字)
		for i := 0; i < tsActualLen; i++ {
			if id[i] < '0' || id[i] > '9' {
				return false
			}
		}

		// 检查随机部分
		for i := tsActualLen; i < len(id); i++ {
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
	} else {
		// 固定时间戳长度模式
		if len(id) != tsLen+randLen {
			return false
		}

		// 检查时间戳部分
		for i := 0; i < tsLen; i++ {
			if id[i] < '0' || id[i] > '9' {
				return false
			}
		}

		// 检查随机部分
		for i := tsLen; i < len(id); i++ {
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
