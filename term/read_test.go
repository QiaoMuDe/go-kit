package term

import (
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		prompt   string
		expected string
	}{
		{
			name:     "普通输入",
			input:    "hello world\n",
			prompt:   "请输入: ",
			expected: "hello world",
		},
		{
			name:     "带空格输入",
			input:    "  test  \n",
			prompt:   "请输入: ",
			expected: "test",
		},
		{
			name:     "空输入",
			input:    "\n",
			prompt:   "请输入: ",
			expected: "",
		},
		{
			name:     "只有空格",
			input:    "   \n",
			prompt:   "请输入: ",
			expected: "",
		},
		{
			name:     "中文输入",
			input:    "你好世界\n",
			prompt:   "请输入: ",
			expected: "你好世界",
		},
		{
			name:     "特殊字符",
			input:    "!@#$%^&*()\n",
			prompt:   "请输入: ",
			expected: "!@#$%^&*()",
		},
		{
			name:     "带空格的中文",
			input:    "  你好 世界  \n",
			prompt:   "请输入: ",
			expected: "你好 世界",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result := Read(reader, tt.prompt)

			if result != tt.expected {
				t.Errorf("Read() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestReadWithDef(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		prompt   string
		def      string
		expected string
	}{
		{
			name:     "有输入时使用输入值",
			input:    "user input\n",
			prompt:   "请输入: ",
			def:      "default",
			expected: "user input",
		},
		{
			name:     "无输入时使用默认值",
			input:    "\n",
			prompt:   "请输入: ",
			def:      "default",
			expected: "default",
		},
		{
			name:     "只有空格时使用默认值",
			input:    "   \n",
			prompt:   "请输入: ",
			def:      "default",
			expected: "default",
		},
		{
			name:     "空默认值",
			input:    "\n",
			prompt:   "请输入: ",
			def:      "",
			expected: "",
		},
		{
			name:     "输入值和默认值相同",
			input:    "same\n",
			prompt:   "请输入: ",
			def:      "same",
			expected: "same",
		},
		{
			name:     "中文默认值",
			input:    "\n",
			prompt:   "请输入: ",
			def:      "默认值",
			expected: "默认值",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result := ReadWithDef(reader, tt.prompt, tt.def)

			if result != tt.expected {
				t.Errorf("ReadWithDef() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestConfirm(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		prompt   string
		defVal   bool
		expected bool
	}{
		{
			name:     "输入 yes 返回 true",
			input:    "yes\n",
			prompt:   "确认吗?",
			defVal:   false,
			expected: true,
		},
		{
			name:     "输入 y 返回 true",
			input:    "y\n",
			prompt:   "确认吗?",
			defVal:   false,
			expected: true,
		},
		{
			name:     "输入 YES 返回 true",
			input:    "YES\n",
			prompt:   "确认吗?",
			defVal:   false,
			expected: true,
		},
		{
			name:     "输入 no 返回 false",
			input:    "no\n",
			prompt:   "确认吗?",
			defVal:   true,
			expected: false,
		},
		{
			name:     "输入 n 返回 false",
			input:    "n\n",
			prompt:   "确认吗?",
			defVal:   true,
			expected: false,
		},
		{
			name:     "输入 NO 返回 false",
			input:    "NO\n",
			prompt:   "确认吗?",
			defVal:   true,
			expected: false,
		},
		{
			name:     "空输入返回默认值 true",
			input:    "\n",
			prompt:   "确认吗?",
			defVal:   true,
			expected: true,
		},
		{
			name:     "空输入返回默认值 false",
			input:    "\n",
			prompt:   "确认吗?",
			defVal:   false,
			expected: false,
		},
		{
			name:     "无效输入返回默认值 true",
			input:    "maybe\n",
			prompt:   "确认吗?",
			defVal:   true,
			expected: true,
		},
		{
			name:     "无效输入返回默认值 false",
			input:    "maybe\n",
			prompt:   "确认吗?",
			defVal:   false,
			expected: false,
		},
		{
			name:     "带空格的 yes",
			input:    "  yes  \n",
			prompt:   "确认吗?",
			defVal:   false,
			expected: true,
		},
		{
			name:     "带空格的 no",
			input:    "  no  \n",
			prompt:   "确认吗?",
			defVal:   true,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result := Confirm(reader, tt.prompt, tt.defVal)

			if result != tt.expected {
				t.Errorf("Confirm() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReadInt(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		prompt    string
		expected  int
		expectErr bool
	}{
		{
			name:      "正整数",
			input:     "123\n",
			prompt:    "请输入: ",
			expected:  123,
			expectErr: false,
		},
		{
			name:      "负整数",
			input:     "-456\n",
			prompt:    "请输入: ",
			expected:  -456,
			expectErr: false,
		},
		{
			name:      "零",
			input:     "0\n",
			prompt:    "请输入: ",
			expected:  0,
			expectErr: false,
		},
		{
			name:      "空输入",
			input:     "\n",
			prompt:    "请输入: ",
			expected:  0,
			expectErr: true,
		},
		{
			name:      "无效的整数",
			input:     "abc\n",
			prompt:    "请输入: ",
			expected:  0,
			expectErr: true,
		},
		{
			name:      "浮点数",
			input:     "12.5\n",
			prompt:    "请输入: ",
			expected:  0,
			expectErr: true,
		},
		{
			name:      "带空格的整数",
			input:     "  789  \n",
			prompt:    "请输入: ",
			expected:  789,
			expectErr: false,
		},
		{
			name:      "大整数",
			input:     "2147483647\n",
			prompt:    "请输入: ",
			expected:  2147483647,
			expectErr: false,
		},
		{
			name:      "负大整数",
			input:     "-2147483648\n",
			prompt:    "请输入: ",
			expected:  -2147483648,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result, err := ReadInt(reader, tt.prompt)

			if tt.expectErr {
				if err == nil {
					t.Errorf("ReadInt() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ReadInt() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("ReadInt() = %d, want %d", result, tt.expected)
				}
			}
		})
	}
}

func TestReadIntWithDef(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		prompt    string
		def       int
		expected  int
		expectErr bool
	}{
		{
			name:      "有输入时使用输入值",
			input:     "100\n",
			prompt:    "请输入: ",
			def:       18,
			expected:  100,
			expectErr: false,
		},
		{
			name:      "空输入时使用默认值",
			input:     "\n",
			prompt:    "请输入: ",
			def:       18,
			expected:  18,
			expectErr: false,
		},
		{
			name:      "只有空格时使用默认值",
			input:     "   \n",
			prompt:    "请输入: ",
			def:       18,
			expected:  18,
			expectErr: false,
		},
		{
			name:      "无效整数",
			input:     "abc\n",
			prompt:    "请输入: ",
			def:       18,
			expected:  0,
			expectErr: true,
		},
		{
			name:      "负数默认值",
			input:     "\n",
			prompt:    "请输入: ",
			def:       -1,
			expected:  -1,
			expectErr: false,
		},
		{
			name:      "零默认值",
			input:     "\n",
			prompt:    "请输入: ",
			def:       0,
			expected:  0,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result, err := ReadIntWithDef(reader, tt.prompt, tt.def)

			if tt.expectErr {
				if err == nil {
					t.Errorf("ReadIntWithDef() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ReadIntWithDef() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("ReadIntWithDef() = %d, want %d", result, tt.expected)
				}
			}
		})
	}
}

func TestReadFloat(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		prompt    string
		expected  float64
		expectErr bool
	}{
		{
			name:      "正浮点数",
			input:     "123.45\n",
			prompt:    "请输入: ",
			expected:  123.45,
			expectErr: false,
		},
		{
			name:      "负浮点数",
			input:     "-456.78\n",
			prompt:    "请输入: ",
			expected:  -456.78,
			expectErr: false,
		},
		{
			name:      "零",
			input:     "0\n",
			prompt:    "请输入: ",
			expected:  0,
			expectErr: false,
		},
		{
			name:      "整数",
			input:     "100\n",
			prompt:    "请输入: ",
			expected:  100,
			expectErr: false,
		},
		{
			name:      "科学计数法",
			input:     "1.23e2\n",
			prompt:    "请输入: ",
			expected:  123,
			expectErr: false,
		},
		{
			name:      "空输入",
			input:     "\n",
			prompt:    "请输入: ",
			expected:  0,
			expectErr: true,
		},
		{
			name:      "无效的数字",
			input:     "abc\n",
			prompt:    "请输入: ",
			expected:  0,
			expectErr: true,
		},
		{
			name:      "带空格的浮点数",
			input:     "  789.12  \n",
			prompt:    "请输入: ",
			expected:  789.12,
			expectErr: false,
		},
		{
			name:      "小数点前导零",
			input:     "0.5\n",
			prompt:    "请输入: ",
			expected:  0.5,
			expectErr: false,
		},
		{
			name:      "大浮点数",
			input:     "1.7976931348623157e+308\n",
			prompt:    "请输入: ",
			expected:  1.7976931348623157e+308,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result, err := ReadFloat(reader, tt.prompt)

			if tt.expectErr {
				if err == nil {
					t.Errorf("ReadFloat() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ReadFloat() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("ReadFloat() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestReadFloatWithDef(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		prompt    string
		def       float64
		expected  float64
		expectErr bool
	}{
		{
			name:      "有输入时使用输入值",
			input:     "99.99\n",
			prompt:    "请输入: ",
			def:       0.0,
			expected:  99.99,
			expectErr: false,
		},
		{
			name:      "空输入时使用默认值",
			input:     "\n",
			prompt:    "请输入: ",
			def:       1.0,
			expected:  1.0,
			expectErr: false,
		},
		{
			name:      "只有空格时使用默认值",
			input:     "   \n",
			prompt:    "请输入: ",
			def:       1.0,
			expected:  1.0,
			expectErr: false,
		},
		{
			name:      "无效数字",
			input:     "abc\n",
			prompt:    "请输入: ",
			def:       1.0,
			expected:  0,
			expectErr: true,
		},
		{
			name:      "负数默认值",
			input:     "\n",
			prompt:    "请输入: ",
			def:       -1.5,
			expected:  -1.5,
			expectErr: false,
		},
		{
			name:      "零默认值",
			input:     "\n",
			prompt:    "请输入: ",
			def:       0.0,
			expected:  0.0,
			expectErr: false,
		},
		{
			name:      "科学计数法默认值",
			input:     "\n",
			prompt:    "请输入: ",
			def:       1.23e2,
			expected:  123,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result, err := ReadFloatWithDef(reader, tt.prompt, tt.def)

			if tt.expectErr {
				if err == nil {
					t.Errorf("ReadFloatWithDef() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ReadFloatWithDef() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("ReadFloatWithDef() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}
