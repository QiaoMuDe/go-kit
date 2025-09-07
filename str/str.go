package str

import "strings"

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
