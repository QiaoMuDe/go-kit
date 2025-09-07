package str

import (
	"strings"
	"testing"
)

func TestSafeDeref(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: "",
		},
		{
			name:     "empty string pointer",
			input:    stringPtr(""),
			expected: "",
		},
		{
			name:     "normal string pointer",
			input:    stringPtr("hello"),
			expected: "hello",
		},
		{
			name:     "string with spaces",
			input:    stringPtr("  hello world  "),
			expected: "  hello world  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeDeref(tt.input)
			if result != tt.expected {
				t.Errorf("SafeDeref() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestBuildStr(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(*strings.Builder)
		expected string
	}{
		{
			name: "empty builder",
			fn: func(buf *strings.Builder) {
				// do nothing
			},
			expected: "",
		},
		{
			name: "single string",
			fn: func(buf *strings.Builder) {
				buf.WriteString("hello")
			},
			expected: "hello",
		},
		{
			name: "multiple operations",
			fn: func(buf *strings.Builder) {
				buf.WriteString("Hello")
				buf.WriteByte(' ')
				buf.WriteString("World")
			},
			expected: "Hello World",
		},
		{
			name: "complex string building",
			fn: func(buf *strings.Builder) {
				buf.WriteString("Name: ")
				buf.WriteString("Alice")
				buf.WriteString(", Age: ")
				buf.WriteString("25")
			},
			expected: "Name: Alice, Age: 25",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildStr(tt.fn)
			if result != tt.expected {
				t.Errorf("BuildStr() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestBuildStrCap(t *testing.T) {
	tests := []struct {
		name     string
		cap      int
		fn       func(*strings.Builder)
		expected string
	}{
		{
			name: "zero capacity",
			cap:  0,
			fn: func(buf *strings.Builder) {
				buf.WriteString("test")
			},
			expected: "test",
		},
		{
			name: "negative capacity",
			cap:  -1,
			fn: func(buf *strings.Builder) {
				buf.WriteString("test")
			},
			expected: "test",
		},
		{
			name: "positive capacity",
			cap:  64,
			fn: func(buf *strings.Builder) {
				buf.WriteString("Hello")
				buf.WriteByte(' ')
				buf.WriteString("World")
			},
			expected: "Hello World",
		},
		{
			name: "large capacity",
			cap:  1024,
			fn: func(buf *strings.Builder) {
				buf.WriteString("test")
			},
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildStrCap(tt.cap, tt.fn)
			if result != tt.expected {
				t.Errorf("BuildStrCap() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "non-empty string",
			input:    "hello",
			expected: false,
		},
		{
			name:     "string with spaces",
			input:    "   ",
			expected: false,
		},
		{
			name:     "single character",
			input:    "a",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmpty(tt.input)
			if result != tt.expected {
				t.Errorf("IsEmpty(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		n        int
		expected string
	}{
		{
			name:     "zero length",
			input:    "hello",
			n:        0,
			expected: "",
		},
		{
			name:     "negative length",
			input:    "hello",
			n:        -1,
			expected: "",
		},
		{
			name:     "length equals string length",
			input:    "hello",
			n:        5,
			expected: "hello",
		},
		{
			name:     "length greater than string length",
			input:    "hello",
			n:        10,
			expected: "hello",
		},
		{
			name:     "normal case",
			input:    "hello world",
			n:        5,
			expected: "hello",
		},
		{
			name:     "empty string",
			input:    "",
			n:        3,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Prefix(tt.input, tt.n)
			if result != tt.expected {
				t.Errorf("Prefix(%q, %d) = %q, want %q", tt.input, tt.n, result, tt.expected)
			}
		})
	}
}

func TestSuffix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		n        int
		expected string
	}{
		{
			name:     "zero length",
			input:    "hello",
			n:        0,
			expected: "",
		},
		{
			name:     "negative length",
			input:    "hello",
			n:        -1,
			expected: "",
		},
		{
			name:     "length equals string length",
			input:    "hello",
			n:        5,
			expected: "hello",
		},
		{
			name:     "length greater than string length",
			input:    "hello",
			n:        10,
			expected: "hello",
		},
		{
			name:     "normal case",
			input:    "hello world",
			n:        5,
			expected: "world",
		},
		{
			name:     "empty string",
			input:    "",
			n:        3,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Suffix(tt.input, tt.n)
			if result != tt.expected {
				t.Errorf("Suffix(%q, %d) = %q, want %q", tt.input, tt.n, result, tt.expected)
			}
		})
	}
}

func TestStringSuffix8(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "string shorter than 8",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "string exactly 8 characters",
			input:    "12345678",
			expected: "12345678",
		},
		{
			name:     "string longer than 8",
			input:    "hello world test",
			expected: "rld test",
		},
		{
			name:     "string with special characters",
			input:    "hello@#$%^&*()world",
			expected: "*()world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringSuffix8(tt.input)
			if result != tt.expected {
				t.Errorf("StringSuffix8(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "zero max length",
			input:    "hello",
			maxLen:   0,
			expected: "",
		},
		{
			name:     "negative max length",
			input:    "hello",
			maxLen:   -1,
			expected: "",
		},
		{
			name:     "max length equals string length",
			input:    "hello",
			maxLen:   5,
			expected: "hello",
		},
		{
			name:     "max length greater than string length",
			input:    "hello",
			maxLen:   10,
			expected: "hello",
		},
		{
			name:     "normal truncation",
			input:    "hello world",
			maxLen:   5,
			expected: "hello",
		},
		{
			name:     "empty string",
			input:    "",
			maxLen:   3,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Truncate(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "string with only spaces",
			input:    "   ",
			expected: false,
		},
		{
			name:     "string with tabs and spaces",
			input:    "\t  \n  ",
			expected: false,
		},
		{
			name:     "non-empty string",
			input:    "hello",
			expected: true,
		},
		{
			name:     "string with content and spaces",
			input:    "  hello  ",
			expected: true,
		},
		{
			name:     "single character",
			input:    "a",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotEmpty(tt.input)
			if result != tt.expected {
				t.Errorf("IsNotEmpty(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIfBlank(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultVal string
		expected   string
	}{
		{
			name:       "empty string",
			input:      "",
			defaultVal: "default",
			expected:   "default",
		},
		{
			name:       "string with only spaces",
			input:      "   ",
			defaultVal: "default",
			expected:   "default",
		},
		{
			name:       "string with tabs and newlines",
			input:      "\t\n  ",
			defaultVal: "default",
			expected:   "default",
		},
		{
			name:       "non-blank string",
			input:      "hello",
			defaultVal: "default",
			expected:   "hello",
		},
		{
			name:       "string with content and spaces",
			input:      "  hello  ",
			defaultVal: "default",
			expected:   "  hello  ",
		},
		{
			name:       "empty default value",
			input:      "",
			defaultVal: "",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IfBlank(tt.input, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("IfBlank(%q, %q) = %q, want %q", tt.input, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

func TestIfEmpty(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultVal string
		expected   string
	}{
		{
			name:       "empty string",
			input:      "",
			defaultVal: "default",
			expected:   "default",
		},
		{
			name:       "string with only spaces",
			input:      "   ",
			defaultVal: "default",
			expected:   "   ",
		},
		{
			name:       "non-empty string",
			input:      "hello",
			defaultVal: "default",
			expected:   "hello",
		},
		{
			name:       "empty default value",
			input:      "",
			defaultVal: "",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IfEmpty(tt.input, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("IfEmpty(%q, %q) = %q, want %q", tt.input, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		count    int
		expected string
	}{
		{
			name:     "zero count",
			input:    "hello",
			count:    0,
			expected: "",
		},
		{
			name:     "negative count",
			input:    "hello",
			count:    -1,
			expected: "",
		},
		{
			name:     "single repeat",
			input:    "hello",
			count:    1,
			expected: "hello",
		},
		{
			name:     "multiple repeats",
			input:    "ab",
			count:    3,
			expected: "ababab",
		},
		{
			name:     "empty string repeat",
			input:    "",
			count:    5,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Repeat(tt.input, tt.count)
			if result != tt.expected {
				t.Errorf("Repeat(%q, %d) = %q, want %q", tt.input, tt.count, result, tt.expected)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		pad      rune
		expected string
	}{
		{
			name:     "no padding needed",
			input:    "hello",
			length:   5,
			pad:      ' ',
			expected: "hello",
		},
		{
			name:     "string longer than target",
			input:    "hello world",
			length:   5,
			pad:      ' ',
			expected: "hello world",
		},
		{
			name:     "normal padding",
			input:    "123",
			length:   5,
			pad:      '0',
			expected: "00123",
		},
		{
			name:     "empty string padding",
			input:    "",
			length:   3,
			pad:      '*',
			expected: "***",
		},
		{
			name:     "zero length",
			input:    "hello",
			length:   0,
			pad:      ' ',
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PadLeft(tt.input, tt.length, tt.pad)
			if result != tt.expected {
				t.Errorf("PadLeft(%q, %d, %q) = %q, want %q", tt.input, tt.length, tt.pad, result, tt.expected)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		pad      rune
		expected string
	}{
		{
			name:     "no padding needed",
			input:    "hello",
			length:   5,
			pad:      ' ',
			expected: "hello",
		},
		{
			name:     "string longer than target",
			input:    "hello world",
			length:   5,
			pad:      ' ',
			expected: "hello world",
		},
		{
			name:     "normal padding",
			input:    "123",
			length:   5,
			pad:      '0',
			expected: "12300",
		},
		{
			name:     "empty string padding",
			input:    "",
			length:   3,
			pad:      '*',
			expected: "***",
		},
		{
			name:     "zero length",
			input:    "hello",
			length:   0,
			pad:      ' ',
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PadRight(tt.input, tt.length, tt.pad)
			if result != tt.expected {
				t.Errorf("PadRight(%q, %d, %q) = %q, want %q", tt.input, tt.length, tt.pad, result, tt.expected)
			}
		})
	}
}

func TestSafeIndex(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected int
	}{
		{
			name:     "empty source string",
			s:        "",
			substr:   "hello",
			expected: -1,
		},
		{
			name:     "empty substring",
			s:        "hello",
			substr:   "",
			expected: -1,
		},
		{
			name:     "both empty",
			s:        "",
			substr:   "",
			expected: -1,
		},
		{
			name:     "substring found at beginning",
			s:        "hello world",
			substr:   "hello",
			expected: 0,
		},
		{
			name:     "substring found in middle",
			s:        "hello world",
			substr:   "lo wo",
			expected: 3,
		},
		{
			name:     "substring not found",
			s:        "hello world",
			substr:   "xyz",
			expected: -1,
		},
		{
			name:     "substring at end",
			s:        "hello world",
			substr:   "world",
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeIndex(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("SafeIndex(%q, %q) = %d, want %d", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestEllipsis(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "zero max length",
			input:    "hello",
			maxLen:   0,
			expected: "",
		},
		{
			name:     "negative max length",
			input:    "hello",
			maxLen:   -1,
			expected: "",
		},
		{
			name:     "string shorter than max length",
			input:    "hello",
			maxLen:   10,
			expected: "hello",
		},
		{
			name:     "string equals max length",
			input:    "hello",
			maxLen:   5,
			expected: "hello",
		},
		{
			name:     "normal ellipsis",
			input:    "hello world",
			maxLen:   8,
			expected: "hello...",
		},
		{
			name:     "max length 1",
			input:    "hello",
			maxLen:   1,
			expected: ".",
		},
		{
			name:     "max length 2",
			input:    "hello",
			maxLen:   2,
			expected: "..",
		},
		{
			name:     "max length 3",
			input:    "hello",
			maxLen:   3,
			expected: "...",
		},
		{
			name:     "empty string",
			input:    "",
			maxLen:   5,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Ellipsis(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("Ellipsis(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestTemplate(t *testing.T) {
	tests := []struct {
		name     string
		tmpl     string
		data     map[string]string
		expected string
	}{
		{
			name:     "empty template",
			tmpl:     "",
			data:     map[string]string{"name": "Alice"},
			expected: "",
		},
		{
			name:     "nil data",
			tmpl:     "Hello {{name}}",
			data:     nil,
			expected: "Hello {{name}}",
		},
		{
			name:     "empty data",
			tmpl:     "Hello {{name}}",
			data:     map[string]string{},
			expected: "Hello {{name}}",
		},
		{
			name:     "single replacement",
			tmpl:     "Hello {{name}}",
			data:     map[string]string{"name": "Alice"},
			expected: "Hello Alice",
		},
		{
			name:     "multiple replacements",
			tmpl:     "Hello {{name}}, you are {{age}} years old",
			data:     map[string]string{"name": "Alice", "age": "25"},
			expected: "Hello Alice, you are 25 years old",
		},
		{
			name:     "no placeholders",
			tmpl:     "Hello World",
			data:     map[string]string{"name": "Alice"},
			expected: "Hello World",
		},
		{
			name:     "empty key in data",
			tmpl:     "Hello {{name}} and {{}}",
			data:     map[string]string{"name": "Alice", "": "Bob"},
			expected: "Hello Alice and {{}}",
		},
		{
			name:     "placeholder not in data",
			tmpl:     "Hello {{name}} and {{unknown}}",
			data:     map[string]string{"name": "Alice"},
			expected: "Hello Alice and {{unknown}}",
		},
		{
			name:     "repeated placeholders",
			tmpl:     "{{name}} says hello to {{name}}",
			data:     map[string]string{"name": "Alice"},
			expected: "Alice says hello to Alice",
		},
		{
			name:     "empty replacement value",
			tmpl:     "Hello {{name}}!",
			data:     map[string]string{"name": ""},
			expected: "Hello !",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Template(tt.tmpl, tt.data)
			if result != tt.expected {
				t.Errorf("Template(%q, %v) = %q, want %q", tt.tmpl, tt.data, result, tt.expected)
			}
		})
	}
}

func TestToBase64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple string",
			input:    "hello",
			expected: "aGVsbG8=",
		},
		{
			name:     "string with spaces",
			input:    "hello world",
			expected: "aGVsbG8gd29ybGQ=",
		},
		{
			name:     "special characters",
			input:    "hello@#$%",
			expected: "aGVsbG9AIyQl",
		},
		{
			name:     "unicode characters",
			input:    "你好",
			expected: "5L2g5aW9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToBase64(tt.input)
			if result != tt.expected {
				t.Errorf("ToBase64(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFromBase64(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			expectError: false,
		},
		{
			name:        "valid base64",
			input:       "aGVsbG8=",
			expected:    "hello",
			expectError: false,
		},
		{
			name:        "valid base64 with spaces",
			input:       "aGVsbG8gd29ybGQ=",
			expected:    "hello world",
			expectError: false,
		},
		{
			name:        "invalid base64",
			input:       "invalid!!!",
			expected:    "",
			expectError: true,
		},
		{
			name:        "unicode base64",
			input:       "5L2g5aW9",
			expected:    "你好",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromBase64(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("FromBase64(%q) expected error, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("FromBase64(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("FromBase64(%q) = %q, want %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{
			name:     "empty parts",
			parts:    []string{},
			expected: "",
		},
		{
			name:     "single part",
			parts:    []string{"hello"},
			expected: "hello",
		},
		{
			name:     "multiple parts",
			parts:    []string{"hello", "world", "test"},
			expected: "helloworldtest",
		},
		{
			name:     "parts with empty strings",
			parts:    []string{"hello", "", "world"},
			expected: "helloworld",
		},
		{
			name:     "all empty strings",
			parts:    []string{"", "", ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Join(tt.parts...)
			if result != tt.expected {
				t.Errorf("Join(%v) = %q, want %q", tt.parts, result, tt.expected)
			}
		})
	}
}

func TestJoinNonEmpty(t *testing.T) {
	tests := []struct {
		name     string
		sep      string
		parts    []string
		expected string
	}{
		{
			name:     "empty parts",
			sep:      ",",
			parts:    []string{},
			expected: "",
		},
		{
			name:     "single non-empty part",
			sep:      ",",
			parts:    []string{"hello"},
			expected: "hello",
		},
		{
			name:     "multiple non-empty parts",
			sep:      ",",
			parts:    []string{"hello", "world", "test"},
			expected: "hello,world,test",
		},
		{
			name:     "parts with empty strings",
			sep:      ",",
			parts:    []string{"hello", "", "world", ""},
			expected: "hello,world",
		},
		{
			name:     "all empty strings",
			sep:      ",",
			parts:    []string{"", "", ""},
			expected: "",
		},
		{
			name:     "different separator",
			sep:      " | ",
			parts:    []string{"a", "b", "c"},
			expected: "a | b | c",
		},
		{
			name:     "empty separator",
			sep:      "",
			parts:    []string{"a", "b", "c"},
			expected: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinNonEmpty(tt.sep, tt.parts...)
			if result != tt.expected {
				t.Errorf("JoinNonEmpty(%q, %v) = %q, want %q", tt.sep, tt.parts, result, tt.expected)
			}
		})
	}
}

func TestMask(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		start    int
		end      int
		maskChar rune
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			start:    0,
			end:      3,
			maskChar: '*',
			expected: "",
		},
		{
			name:     "negative start",
			input:    "hello",
			start:    -1,
			end:      3,
			maskChar: '*',
			expected: "hello",
		},
		{
			name:     "negative end",
			input:    "hello",
			start:    1,
			end:      -1,
			maskChar: '*',
			expected: "hello",
		},
		{
			name:     "end <= start",
			input:    "hello",
			start:    3,
			end:      3,
			maskChar: '*',
			expected: "hello",
		},
		{
			name:     "start >= string length",
			input:    "hello",
			start:    10,
			end:      15,
			maskChar: '*',
			expected: "hello",
		},
		{
			name:     "normal masking",
			input:    "13812345678",
			start:    3,
			end:      7,
			maskChar: '*',
			expected: "138****5678",
		},
		{
			name:     "end exceeds string length",
			input:    "hello",
			start:    2,
			end:      10,
			maskChar: '*',
			expected: "he***",
		},
		{
			name:     "zero mask character",
			input:    "hello",
			start:    1,
			end:      4,
			maskChar: 0,
			expected: "h***o",
		},
		{
			name:     "unicode characters",
			input:    "你好世界",
			start:    1,
			end:      3,
			maskChar: '●',
			expected: "你●●界",
		},
		{
			name:     "mask entire string",
			input:    "hello",
			start:    0,
			end:      5,
			maskChar: '#',
			expected: "#####",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Mask(tt.input, tt.start, tt.end, tt.maskChar)
			if result != tt.expected {
				t.Errorf("Mask(%q, %d, %d, %q) = %q, want %q", tt.input, tt.start, tt.end, tt.maskChar, result, tt.expected)
			}
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// Benchmark tests
func BenchmarkBuildStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildStr(func(buf *strings.Builder) {
			buf.WriteString("Hello")
			buf.WriteByte(' ')
			buf.WriteString("World")
		})
	}
}

func BenchmarkJoin(b *testing.B) {
	parts := []string{"hello", "world", "test", "benchmark"}
	for i := 0; i < b.N; i++ {
		Join(parts...)
	}
}

func BenchmarkTemplate(b *testing.B) {
	tmpl := "Hello {{name}}, you are {{age}} years old"
	data := map[string]string{"name": "Alice", "age": "25"}
	for i := 0; i < b.N; i++ {
		Template(tmpl, data)
	}
}

func BenchmarkMask(b *testing.B) {
	phone := "13812345678"
	for i := 0; i < b.N; i++ {
		Mask(phone, 3, 7, '*')
	}
}
