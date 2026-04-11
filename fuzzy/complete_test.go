package fuzzy

import (
	"testing"
)

func TestComplete(t *testing.T) {
	flags := []string{
		"--verbose",
		"--version",
		"--help",
		"-v",
		"-h",
		"--input-file",
		"--output",
	}

	tests := []struct {
		name     string
		pattern  string
		expected []string
	}{
		{
			name:     "前缀匹配",
			pattern:  "--v",
			expected: []string{"--verbose", "--version"},
		},
		{
			name:     "精确匹配优先",
			pattern:  "-v",
			expected: []string{"-v", "--verbose", "--version"},
		},
		{
			name:     "不区分大小写前缀",
			pattern:  "--V",
			expected: []string{"--verbose", "--version"},
		},
		{
			name:     "短横线匹配",
			pattern:  "-",
			expected: []string{"-v", "-h", "--help", "--output", "--verbose", "--version", "--input-file"},
		},
		{
			name:     "无匹配",
			pattern:  "xyz",
			expected: []string{},
		},
		{
			name:     "空模式",
			pattern:  "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := Complete(tt.pattern, flags)

			if len(tt.expected) == 0 {
				if len(matches) != 0 {
					t.Errorf("期望无匹配，但得到 %d 个结果", len(matches))
				}
				return
			}

			if len(matches) != len(tt.expected) {
				t.Errorf("期望 %d 个结果，得到 %d 个", len(tt.expected), len(matches))
				return
			}

			for i, expected := range tt.expected {
				if matches[i].Str != expected {
					t.Errorf("位置 %d: 期望 %q，得到 %q", i, expected, matches[i].Str)
				}
			}
		})
	}
}

func TestCompletePrefix(t *testing.T) {
	flags := []string{
		"--verbose",
		"--version",
		"--help",
		"-v",
		"-h",
	}

	tests := []struct {
		name     string
		pattern  string
		expected []string
	}{
		{
			name:     "双横线前缀",
			pattern:  "--v",
			expected: []string{"--verbose", "--version"},
		},
		{
			name:     "单横线前缀",
			pattern:  "-v",
			expected: []string{"-v"},
		},
		{
			name:     "不区分大小写",
			pattern:  "--V",
			expected: []string{"--verbose", "--version"},
		},
		{
			name:     "短候选优先",
			pattern:  "-",
			expected: []string{"-v", "-h", "--help", "--verbose", "--version"},
		},
		{
			name:     "模糊不匹配",
			pattern:  "vrb",
			expected: []string{},
		},
		{
			name:     "空模式",
			pattern:  "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := CompletePrefix(tt.pattern, flags)

			if len(tt.expected) == 0 {
				if len(matches) != 0 {
					t.Errorf("期望无匹配，但得到 %d 个结果", len(matches))
				}
				return
			}

			if len(matches) != len(tt.expected) {
				t.Errorf("期望 %d 个结果，得到 %d 个", len(tt.expected), len(matches))
				return
			}

			for i, expected := range tt.expected {
				if matches[i].Str != expected {
					t.Errorf("位置 %d: 期望 %q，得到 %q", i, expected, matches[i].Str)
				}
			}
		})
	}
}

func TestCompleteExact(t *testing.T) {
	flags := []string{
		"--verbose",
		"--version",
		"-v",
	}

	tests := []struct {
		name     string
		pattern  string
		expected string
		found    bool
	}{
		{
			name:     "精确匹配",
			pattern:  "-v",
			expected: "-v",
			found:    true,
		},
		{
			name:     "区分大小写",
			pattern:  "-V",
			expected: "",
			found:    false,
		},
		{
			name:     "前缀不匹配",
			pattern:  "--v",
			expected: "",
			found:    false,
		},
		{
			name:     "空模式",
			pattern:  "",
			expected: "",
			found:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := CompleteExact(tt.pattern, flags)

			if !tt.found {
				if len(matches) != 0 {
					t.Errorf("期望无匹配，但得到 %d 个结果", len(matches))
				}
				return
			}

			if len(matches) != 1 {
				t.Errorf("期望 1 个结果，得到 %d 个", len(matches))
				return
			}

			if matches[0].Str != tt.expected {
				t.Errorf("期望 %q，得到 %q", tt.expected, matches[0].Str)
			}
		})
	}
}

func TestCompleteScore(t *testing.T) {
	// 验证分数排序
	flags := []string{
		"--verbose",
		"--version",
		"-v",
	}

	matches := Complete("-v", flags)

	if len(matches) != 3 {
		t.Fatalf("期望 3 个结果，得到 %d 个", len(matches))
	}

	// 验证分数顺序：精确匹配 > 前缀匹配 > 模糊匹配
	if matches[0].Str != "-v" {
		t.Errorf("第一个应该是精确匹配 '-v'，得到 %q", matches[0].Str)
	}
	if matches[0].Score != 1000 {
		t.Errorf("精确匹配分数应该是 1000，得到 %d", matches[0].Score)
	}

	// -v 匹配 --verbose 是模糊匹配（不是前缀匹配，因为 --verbose 以 -- 开头）
	// 模糊匹配分数较低
	if matches[1].Score >= 100 {
		t.Errorf("模糊匹配分数应该小于 100，得到 %d", matches[1].Score)
	}
}

func TestCompleteEmptyCandidates(t *testing.T) {
	// 空候选列表
	matches := Complete("test", []string{})
	if matches != nil {
		t.Errorf("空候选列表应该返回 nil，得到 %v", matches)
	}

	matches = CompletePrefix("test", []string{})
	if matches != nil {
		t.Errorf("空候选列表应该返回 nil，得到 %v", matches)
	}

	matches = CompleteExact("test", []string{})
	if matches != nil {
		t.Errorf("空候选列表应该返回 nil，得到 %v", matches)
	}
}

func TestCompleteMatchedIndexes(t *testing.T) {
	flags := []string{"--verbose"}

	// 测试前缀匹配的索引记录
	matches := CompletePrefix("--v", flags)
	if len(matches) != 1 {
		t.Fatalf("期望 1 个结果，得到 %d 个", len(matches))
	}

	// 前缀 "--v" 应该匹配位置 0, 1, 2
	expectedIndexes := []int{0, 1, 2}
	if len(matches[0].MatchedIndexes) != len(expectedIndexes) {
		t.Errorf("期望 %d 个匹配索引，得到 %d 个",
			len(expectedIndexes), len(matches[0].MatchedIndexes))
	}

	for i, expected := range expectedIndexes {
		if matches[0].MatchedIndexes[i] != expected {
			t.Errorf("索引 %d: 期望 %d，得到 %d",
				i, expected, matches[0].MatchedIndexes[i])
		}
	}
}
