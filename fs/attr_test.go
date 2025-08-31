package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsHidden(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		setup    func() string
		expected bool
		desc     string
	}{
		{
			name: "普通文件",
			setup: func() string {
				file := filepath.Join(tempDir, "normal_file.txt")
				if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
					t.Fatalf("创建普通文件失败: %v", err)
				}
				return file
			},
			expected: false,
			desc:     "普通文件应该不是隐藏文件",
		},
		{
			name: "点开头的隐藏文件",
			setup: func() string {
				file := filepath.Join(tempDir, ".hidden_file")
				if err := os.WriteFile(file, []byte("hidden content"), 0644); err != nil {
					t.Fatalf("创建隐藏文件失败: %v", err)
				}
				return file
			},
			expected: true,
			desc:     "以点开头的文件应该是隐藏文件",
		},
		{
			name: "点开头的隐藏目录",
			setup: func() string {
				dir := filepath.Join(tempDir, ".hidden_dir")
				if err := os.Mkdir(dir, 0755); err != nil {
					t.Fatalf("创建隐藏目录失败: %v", err)
				}
				return dir
			},
			expected: true,
			desc:     "以点开头的目录应该是隐藏目录",
		},
		{
			name: "只有一个点的文件名",
			setup: func() string {
				// 这种情况通常不会被认为是隐藏文件
				return "."
			},
			expected: false,
			desc:     "单个点不应该被认为是隐藏文件",
		},
		{
			name: "不存在的文件",
			setup: func() string {
				return filepath.Join(tempDir, "non_existing_file")
			},
			expected: false,
			desc:     "不存在的文件应该返回false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			result := IsHidden(path)

			if result != tt.expected {
				t.Errorf("IsHidden(%q) = %v, 期望 %v - %s", path, result, tt.expected, tt.desc)
			}
		})
	}
}

func TestIsReadOnly(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		setup    func() string
		expected bool
		desc     string
	}{
		{
			name: "普通可写文件",
			setup: func() string {
				file := filepath.Join(tempDir, "writable_file.txt")
				if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
					t.Fatalf("创建可写文件失败: %v", err)
				}
				return file
			},
			expected: false,
			desc:     "普通可写文件应该不是只读",
		},
		{
			name: "只读文件",
			setup: func() string {
				file := filepath.Join(tempDir, "readonly_file.txt")
				if err := os.WriteFile(file, []byte("readonly content"), 0444); err != nil {
					t.Fatalf("创建只读文件失败: %v", err)
				}
				return file
			},
			expected: true,
			desc:     "只读文件应该返回true",
		},
		{
			name: "不存在的文件",
			setup: func() string {
				return filepath.Join(tempDir, "non_existing_readonly_file")
			},
			expected: false,
			desc:     "不存在的文件应该返回false",
		},
		{
			name: "目录",
			setup: func() string {
				dir := filepath.Join(tempDir, "test_directory")
				if err := os.Mkdir(dir, 0755); err != nil {
					t.Fatalf("创建测试目录失败: %v", err)
				}
				return dir
			},
			expected: false,
			desc:     "普通目录应该不是只读",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			result := IsReadOnly(path)

			if result != tt.expected {
				t.Errorf("IsReadOnly(%q) = %v, 期望 %v - %s", path, result, tt.expected, tt.desc)
			}
		})
	}
}

// 边界条件测试
func TestAttrBoundaryConditions(t *testing.T) {
	tests := []struct {
		name string
		path string
		desc string
	}{
		{
			name: "空字符串路径",
			path: "",
			desc: "空字符串路径应该被正确处理",
		},
		{
			name: "只有空格的路径",
			path: "   ",
			desc: "只包含空格的路径应该被正确处理",
		},
		{
			name: "包含特殊字符的路径",
			path: "test\x00file",
			desc: "包含null字符的路径应该被正确处理",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 这些调用不应该panic
			hidden := IsHidden(tt.path)
			readonly := IsReadOnly(tt.path)

			// 对于无效路径，函数应该返回false
			t.Logf("路径 %q: IsHidden=%v, IsReadOnly=%v - %s",
				tt.path, hidden, readonly, tt.desc)
		})
	}
}

// 性能测试
func BenchmarkIsHidden(b *testing.B) {
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, ".hidden_benchmark_file")

	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		b.Fatalf("创建基准测试文件失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsHidden(testFile)
	}
}

func BenchmarkIsReadOnly(b *testing.B) {
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "readonly_benchmark_file")

	if err := os.WriteFile(testFile, []byte("content"), 0444); err != nil {
		b.Fatalf("创建基准测试文件失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsReadOnly(testFile)
	}
}

// 并发测试
func TestAttrConcurrency(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	hiddenFile := filepath.Join(tempDir, ".hidden_concurrent_test")
	readonlyFile := filepath.Join(tempDir, "readonly_concurrent_test")

	if err := os.WriteFile(hiddenFile, []byte("hidden content"), 0644); err != nil {
		t.Fatalf("创建隐藏测试文件失败: %v", err)
	}

	if err := os.WriteFile(readonlyFile, []byte("readonly content"), 0444); err != nil {
		t.Fatalf("创建只读测试文件失败: %v", err)
	}

	const numGoroutines = 100
	done := make(chan bool, numGoroutines)

	// 并发测试
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer func() { done <- true }()

			// 多次调用属性检查函数
			for j := 0; j < 10; j++ {
				if !IsHidden(hiddenFile) {
					t.Errorf("并发测试中IsHidden返回了错误结果")
				}
				if !IsReadOnly(readonlyFile) {
					t.Errorf("并发测试中IsReadOnly返回了错误结果")
				}
			}
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}
