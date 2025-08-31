package str

import (
	"strings"
	"testing"
)

// TestIsNotEmpty 测试 IsNotEmpty 函数
func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// 正常情况
		{"非空字符串", "hello", true},
		{"包含空格的非空字符串", "hello world", true},
		{"数字字符串", "123", true},
		{"特殊字符", "!@#$%", true},
		{"中文字符串", "你好", true},

		// 边界情况
		{"空字符串", "", false},
		{"只有空格", "   ", false},
		{"只有制表符", "\t", false},
		{"只有换行符", "\n", false},
		{"混合空白字符", " \t\n ", false},
		{"单个字符", "a", true},

		// 前后有空格但中间有内容
		{"前后有空格", "  hello  ", true},
		{"前有空格", "  hello", true},
		{"后有空格", "hello  ", true},
		{"前后有制表符", "\thello\t", true},
		{"前后有换行符", "\nhello\n", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotEmpty(tt.input)
			if result != tt.expected {
				t.Errorf("IsNotEmpty(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestStringSuffix8 测试 StringSuffix8 函数
func TestStringSuffix8(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// 正常情况
		{"长度大于8的字符串", "abcdefghijk", "defghijk"},
		{"长度等于9的字符串", "123456789", "23456789"},
		{"长度远大于8的字符串", "this is a very long string", "g string"},

		// 边界情况
		{"空字符串", "", ""},
		{"长度等于8的字符串", "12345678", "12345678"},
		{"长度小于8的字符串", "hello", "hello"},
		{"单个字符", "a", "a"},
		{"长度为7的字符串", "1234567", "1234567"},

		// 特殊字符
		{"包含特殊字符", "!@#$%^&*()abcdef", "()abcdef"},

		// 中文字符串测试 - 注意：函数按字节长度处理，不是字符长度
		{"中文字符串", "这是中文测试字符串", "\xad\x97符串"},
		{"中文字符串长度小于8", "中文测试", "\x96\x87测试"},
		{"混合字符", "abc中文123defg", "\x87123defg"},

		// 简单ASCII测试避免UTF-8复杂性
		{"长ASCII字符串", "abcdefghijklmnop", "ijklmnop"},
		{"包含数字", "prefix123456789", "23456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringSuffix8(tt.input)
			if result != tt.expected {
				t.Errorf("StringSuffix8(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestSafeDeref 测试 SafeDeref 函数
func TestSafeDeref(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		// 正常情况
		{"非空指针指向非空字符串", stringPtr("hello"), "hello"},
		{"非空指针指向空字符串", stringPtr(""), ""},
		{"非空指针指向包含空格的字符串", stringPtr("  hello  "), "  hello  "},
		{"非空指针指向数字字符串", stringPtr("12345"), "12345"},
		{"非空指针指向特殊字符", stringPtr("!@#$%"), "!@#$%"},
		{"非空指针指向中文字符串", stringPtr("你好世界"), "你好世界"},

		// 边界情况
		{"nil指针", nil, ""},

		// 长字符串
		{"非空指针指向长字符串", stringPtr("this is a very long string for testing"), "this is a very long string for testing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeDeref(tt.input)
			if result != tt.expected {
				t.Errorf("SafeDeref(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// 辅助函数：创建字符串指针
func stringPtr(s string) *string {
	return &s
}

// BenchmarkIsNotEmpty 性能测试 IsNotEmpty 函数
func BenchmarkIsNotEmpty(b *testing.B) {
	testCases := []string{
		"hello world",
		"",
		"   ",
		"a",
		"this is a very long string for performance testing",
	}

	for _, tc := range testCases {
		b.Run("input_"+tc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				IsNotEmpty(tc)
			}
		})
	}
}

// BenchmarkStringSuffix8 性能测试 StringSuffix8 函数
func BenchmarkStringSuffix8(b *testing.B) {
	testCases := []string{
		"short",
		"12345678",
		"this is a very long string for performance testing",
		"",
	}

	for _, tc := range testCases {
		b.Run("input_len_"+string(rune(len(tc))), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				StringSuffix8(tc)
			}
		})
	}
}

// BenchmarkSafeDeref 性能测试 SafeDeref 函数
func BenchmarkSafeDeref(b *testing.B) {
	str := "hello world"
	ptr := &str

	b.Run("non_nil_pointer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SafeDeref(ptr)
		}
	})

	b.Run("nil_pointer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SafeDeref(nil)
		}
	})
}

// TestStringSuffix8_ByteLength 测试字节长度处理
func TestStringSuffix8_ByteLength(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// 测试函数按字节长度工作的特性
		{"9字节ASCII", "123456789", "23456789"},
		{"16字节ASCII", "abcdefghijklmnop", "ijklmnop"},
		{"正好8字节", "12345678", "12345678"},
		{"少于8字节", "hello", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringSuffix8(tt.input)
			if result != tt.expected {
				t.Errorf("StringSuffix8(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
			// 验证结果字节长度不超过8
			if len(result) > 8 {
				t.Errorf("StringSuffix8(%q) result byte length %d > 8", tt.input, len(result))
			}
		})
	}
}

// TestStringSuffix8_UTF8Behavior 测试UTF-8字符的实际行为
func TestStringSuffix8_UTF8Behavior(t *testing.T) {
	// 这个测试展示当前函数如何处理UTF-8字符
	// 注意：这些测试反映当前实现的行为，可能包含被截断的UTF-8字符

	input := "hello世界"
	result := StringSuffix8(input)

	// 验证结果字节长度不超过8
	if len(result) > 8 {
		t.Errorf("Result byte length %d > 8", len(result))
	}

	// 验证结果是输入的后缀
	if !strings.HasSuffix(input, result) {
		t.Errorf("Result %q is not a suffix of input %q", result, input)
	}
}
