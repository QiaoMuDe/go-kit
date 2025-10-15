package utils

import (
	"testing"
)

// TestNeedsEsc 测试 needsEsc 函数
func TestNeedsEsc(t *testing.T) {
	tests := []struct {
		input    byte
		expected bool
		name     string
	}{
		{0x00, true, "null character"},
		{0x01, true, "start of heading"},
		{0x1F, true, "unit separator"},
		{'"', true, "quotation mark"},
		{'\\', true, "backslash"},
		{' ', false, "space"},
		{'A', false, "uppercase A"},
		{'z', false, "lowercase z"},
		{0x20, false, "space hex"},
		{0x7F, false, "delete character"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := needsEsc(tt.input)
			if result != tt.expected {
				t.Errorf("needsEsc(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestQuoteBytes 测试 QuoteBytes 函数
func TestQuoteBytes(t *testing.T) {
	tests := []struct {
		input    []byte
		expected []byte
		name     string
	}{
		{[]byte(""), []byte(""), "empty string"},
		{[]byte("hello"), []byte("hello"), "no escape chars"},
		{[]byte("\""), []byte("\\\""), "quotation mark"},
		{[]byte("\\"), []byte("\\\\"), "backslash"},
		{[]byte("\b"), []byte("\\b"), "backspace"},
		{[]byte("\f"), []byte("\\f"), "form feed"},
		{[]byte("\n"), []byte("\\n"), "line feed"},
		{[]byte("\r"), []byte("\\r"), "carriage return"},
		{[]byte("\t"), []byte("\\t"), "tab"},
		{[]byte("hello\nworld"), []byte("hello\\nworld"), "newline in string"},
		{[]byte("foo\"bar\\baz"), []byte("foo\\\"bar\\\\baz"), "quotes and backslashes"},
		{[]byte{0x00}, []byte("\\u0000"), "null character"},
		{[]byte{0x01}, []byte("\\u0001"), "start of heading"},
		{[]byte{0x1F}, []byte("\\u001f"), "unit separator"},
		{[]byte("hello\x01world"), []byte("hello\\u0001world"), "control character in string"},
		{[]byte("混合test\t\n\"\\测试"), []byte("混合test\\t\\n\\\"\\\\测试"), "mixed unicode and escape chars"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := QuoteBytes(tt.input)
			if string(result) != string(tt.expected) {
				t.Errorf("QuoteBytes(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestQuoteString 测试 QuoteString 函数
func TestQuoteString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		name     string
	}{
		{"", "", "empty string"},
		{"hello", "hello", "no escape chars"},
		{"\"", "\\\"", "quotation mark"},
		{"\\", "\\\\", "backslash"},
		{"\b", "\\b", "backspace"},
		{"\f", "\\f", "form feed"},
		{"\n", "\\n", "line feed"},
		{"\r", "\\r", "carriage return"},
		{"\t", "\\t", "tab"},
		{"hello\nworld", "hello\\nworld", "newline in string"},
		{"foo\"bar\\baz", "foo\\\"bar\\\\baz", "quotes and backslashes"},
		{"\x00", "\\u0000", "null character"},
		{"\x01", "\\u0001", "start of heading"},
		{"\x1F", "\\u001f", "unit separator"},
		{"hello\x01world", "hello\\u0001world", "control character in string"},
		{"混合test\t\n\"\\测试", "混合test\\t\\n\\\"\\\\测试", "mixed unicode and escape chars"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := QuoteString(tt.input)
			if result != tt.expected {
				t.Errorf("QuoteString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// BenchmarkQuoteBytes 基准测试 QuoteBytes 函数
func BenchmarkQuoteBytes(b *testing.B) {
	input := []byte("hello\nworld\"foo\\bar\b\f\r\t\x01")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QuoteBytes(input)
	}
}

// BenchmarkQuoteString 基准测试 QuoteString 函数
func BenchmarkQuoteString(b *testing.B) {
	input := "hello\nworld\"foo\\bar\b\f\r\t\x01"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QuoteString(input)
	}
}

// TestQuoteBytesNoEscape 测试没有转义字符的情况
func TestQuoteBytesNoEscape(t *testing.T) {
	input := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := QuoteBytes(input)

	// 当没有转义字符时，应该返回原始切片
	if string(result) != string(input) {
		t.Errorf("QuoteBytes without escape chars should return original slice")
	}
}
