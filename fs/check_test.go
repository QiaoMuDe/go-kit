package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		setup    func() string
		expected bool
		desc     string
	}{
		{
			name: "存在的文件",
			setup: func() string {
				file := filepath.Join(tempDir, "existing_file.txt")
				f, err := os.Create(file)
				if err != nil {
					t.Fatalf("创建测试文件失败: %v", err)
				}
				_ = f.Close()
				return file
			},
			expected: true,
			desc:     "检查存在的文件应该返回true",
		},
		{
			name: "存在的目录",
			setup: func() string {
				dir := filepath.Join(tempDir, "existing_dir")
				if err := os.Mkdir(dir, 0755); err != nil {
					t.Fatalf("创建测试目录失败: %v", err)
				}
				return dir
			},
			expected: true,
			desc:     "检查存在的目录应该返回true",
		},
		{
			name: "不存在的文件",
			setup: func() string {
				return filepath.Join(tempDir, "non_existing_file.txt")
			},
			expected: false,
			desc:     "检查不存在的文件应该返回false",
		},
		{
			name: "空路径",
			setup: func() string {
				return ""
			},
			expected: false,
			desc:     "空路径应该返回false",
		},
		{
			name: "相对路径",
			setup: func() string {
				// 创建相对路径的文件
				relFile := "test_relative.txt"
				f, err := os.Create(relFile)
				if err != nil {
					t.Fatalf("创建相对路径测试文件失败: %v", err)
				}
				_ = f.Close()
				// 清理函数
				t.Cleanup(func() {
					_ = os.Remove(relFile)
				})
				return relFile
			},
			expected: true,
			desc:     "相对路径的存在文件应该返回true",
		},
		{
			name: "包含特殊字符的路径",
			setup: func() string {
				specialFile := filepath.Join(tempDir, "测试文件 with spaces & symbols!.txt")
				f, err := os.Create(specialFile)
				if err != nil {
					t.Fatalf("创建特殊字符测试文件失败: %v", err)
				}
				_ = f.Close()
				return specialFile
			},
			expected: true,
			desc:     "包含特殊字符的路径应该正确处理",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			result := Exists(path)

			if result != tt.expected {
				t.Errorf("Exists(%q) = %v, 期望 %v - %s", path, result, tt.expected, tt.desc)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
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
				file := filepath.Join(tempDir, "test_file.txt")
				f, err := os.Create(file)
				if err != nil {
					t.Fatalf("创建测试文件失败: %v", err)
				}
				_, _ = f.WriteString("test content")
				_ = f.Close()
				return file
			},
			expected: true,
			desc:     "普通文件应该返回true",
		},
		{
			name: "目录",
			setup: func() string {
				dir := filepath.Join(tempDir, "test_dir")
				if err := os.Mkdir(dir, 0755); err != nil {
					t.Fatalf("创建测试目录失败: %v", err)
				}
				return dir
			},
			expected: false,
			desc:     "目录应该返回false",
		},
		{
			name: "不存在的文件",
			setup: func() string {
				return filepath.Join(tempDir, "non_existing.txt")
			},
			expected: false,
			desc:     "不存在的文件应该返回false",
		},
		{
			name: "空文件",
			setup: func() string {
				file := filepath.Join(tempDir, "empty_file.txt")
				f, err := os.Create(file)
				if err != nil {
					t.Fatalf("创建空文件失败: %v", err)
				}
				_ = f.Close()
				return file
			},
			expected: true,
			desc:     "空文件也应该返回true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			result := IsFile(path)

			if result != tt.expected {
				t.Errorf("IsFile(%q) = %v, 期望 %v - %s", path, result, tt.expected, tt.desc)
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		setup    func() string
		expected bool
		desc     string
	}{
		{
			name: "普通目录",
			setup: func() string {
				dir := filepath.Join(tempDir, "test_directory")
				if err := os.Mkdir(dir, 0755); err != nil {
					t.Fatalf("创建测试目录失败: %v", err)
				}
				return dir
			},
			expected: true,
			desc:     "普通目录应该返回true",
		},
		{
			name: "文件",
			setup: func() string {
				file := filepath.Join(tempDir, "test_file.txt")
				f, err := os.Create(file)
				if err != nil {
					t.Fatalf("创建测试文件失败: %v", err)
				}
				_ = f.Close()
				return file
			},
			expected: false,
			desc:     "文件应该返回false",
		},
		{
			name: "不存在的目录",
			setup: func() string {
				return filepath.Join(tempDir, "non_existing_dir")
			},
			expected: false,
			desc:     "不存在的目录应该返回false",
		},
		{
			name: "嵌套目录",
			setup: func() string {
				nestedDir := filepath.Join(tempDir, "parent", "child")
				if err := os.MkdirAll(nestedDir, 0755); err != nil {
					t.Fatalf("创建嵌套目录失败: %v", err)
				}
				return nestedDir
			},
			expected: true,
			desc:     "嵌套目录应该返回true",
		},
		{
			name: "根目录",
			setup: func() string {
				return tempDir
			},
			expected: true,
			desc:     "根目录应该返回true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			result := IsDir(path)

			if result != tt.expected {
				t.Errorf("IsDir(%q) = %v, 期望 %v - %s", path, result, tt.expected, tt.desc)
			}
		})
	}
}

// 边界条件测试
func TestCheckFunctionsBoundaryConditions(t *testing.T) {
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
			name: "非常长的路径",
			path: string(make([]byte, 1000)),
			desc: "非常长的路径应该被正确处理",
		},
		{
			name: "包含null字符的路径",
			path: "test\x00file",
			desc: "包含null字符的路径应该被正确处理",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 这些调用不应该panic
			exists := Exists(tt.path)
			isFile := IsFile(tt.path)
			isDir := IsDir(tt.path)

			// 对于无效路径，所有函数都应该返回false
			if exists || isFile || isDir {
				t.Logf("路径 %q: Exists=%v, IsFile=%v, IsDir=%v - %s",
					tt.path, exists, isFile, isDir, tt.desc)
			}
		})
	}
}

// 性能测试
func BenchmarkExists(b *testing.B) {
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "benchmark_file.txt")

	// 创建测试文件
	f, err := os.Create(testFile)
	if err != nil {
		b.Fatalf("创建基准测试文件失败: %v", err)
	}
	_ = f.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Exists(testFile)
	}
}

func BenchmarkIsFile(b *testing.B) {
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "benchmark_file.txt")

	f, err := os.Create(testFile)
	if err != nil {
		b.Fatalf("创建基准测试文件失败: %v", err)
	}
	_ = f.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsFile(testFile)
	}
}

func BenchmarkIsDir(b *testing.B) {
	tempDir := b.TempDir()
	testDir := filepath.Join(tempDir, "benchmark_dir")

	if err := os.Mkdir(testDir, 0755); err != nil {
		b.Fatalf("创建基准测试目录失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsDir(testDir)
	}
}

// 并发测试
func TestCheckFunctionsConcurrency(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "concurrent_test.txt")
	testDir := filepath.Join(tempDir, "concurrent_dir")

	// 创建测试文件和目录
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("创建并发测试文件失败: %v", err)
	}
	_ = f.Close()

	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("创建并发测试目录失败: %v", err)
	}

	// 并发测试
	const numGoroutines = 100
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer func() { done <- true }()

			// 多次调用检查函数
			for j := 0; j < 10; j++ {
				if !Exists(testFile) {
					t.Errorf("并发测试中Exists返回了错误结果")
				}
				if !IsFile(testFile) {
					t.Errorf("并发测试中IsFile返回了错误结果")
				}
				if !IsDir(testDir) {
					t.Errorf("并发测试中IsDir返回了错误结果")
				}
			}
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}
