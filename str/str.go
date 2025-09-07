package str

import "strings"

// IsNotEmpty 检查字符串是否不为空
// 用于验证字符串在去除首尾空格后是否包含有效内容
//
// 参数:
//   - s: 待检查的字符串
//
// 返回:
//   - bool: 字符串不为空返回true，否则返回false
func IsNotEmpty(s string) bool {
	// 去除字符串两端的空格
	s = strings.TrimSpace(s)

	// 检查字符串是否为空
	return s != ""
}

// StringSuffix8 从给定字符串中获取最后8个字符。
// 如果字符串长度小于等于8，则返回原字符串。
//
// 参数:
//   - s: 输入字符串
//
// 返回:
//   - string: 字符串的最后8个字符，或原字符串（如果长度不足8），或空字符串（如果输入为空）
func StringSuffix8(s string) string {

	// 检查输入字符串是否为空，若为空则直接返回空字符串
	if s == "" {
		return ""
	}

	// 检查输入字符串的长度是否小于等于 8，若是则直接返回该字符串本身
	if len(s) <= 8 {
		return s
	}

	// 若输入字符串长度大于 8，则截取并返回其最后 8 个字符
	return s[len(s)-8:]
}

// SafeDeref 安全地解引用字符串指针
// 用于安全地获取字符串指针的值，避免空指针异常
//
// 参数:
//   - s: 字符串指针
//
// 返回:
//   - string: 解引用后的字符串，指针为nil时返回空字符串
func SafeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// BuildStr 使用默认容量的字符串构建器执行函数
//
// 参数:
//   - fn: 使用字符串构建器的函数
//
// 返回值:
//   - string: 函数执行后构建的字符串结果
//
// 说明:
//   - 创建新的字符串构建器（不使用对象池）
//   - 执行用户提供的函数
//   - 返回构建的字符串结果
//   - 适用于不需要对象池优化的简单字符串构建场景
//
// 使用示例:
//
//	result := str.BuildStr(func(buf *strings.Builder) {
//	    buf.WriteString("Hello")
//	    buf.WriteByte(' ')
//	    buf.WriteString("World")
//	})
func BuildStr(fn func(*strings.Builder)) string {
	var buf strings.Builder
	fn(&buf)
	return buf.String()
}

// BuildStrCap 使用指定容量的字符串构建器执行函数
//
// 参数:
//   - cap: 字符串构建器初始容量
//   - fn: 使用字符串构建器的函数
//
// 返回值:
//   - string: 函数执行后构建的字符串结果
//
// 说明:
//   - 创建指定容量的新字符串构建器（不使用对象池）
//   - 执行用户提供的函数
//   - 返回构建的字符串结果
//   - 适用于已知字符串长度且不需要对象池优化的场景
//
// 使用示例:
//
//	result := str.BuildStrCap(64, func(buf *strings.Builder) {
//	    buf.WriteString("Hello")
//	    buf.WriteByte(' ')
//	    buf.WriteString("World")
//	})
func BuildStrCap(cap int, fn func(*strings.Builder)) string {
	var buf strings.Builder
	if cap > 0 {
		buf.Grow(cap)
	}
	fn(&buf)
	return buf.String()
}
