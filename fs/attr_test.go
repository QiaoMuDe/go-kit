package fs

import (
	"os"
	"path/filepath"
	"runtime"
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

// TestIsDriveRoot 测试盘符根目录检测（Windows 特有）
func TestIsDriveRoot(t *testing.T) {
	// 仅在 Windows 上测试
	if runtime.GOOS != "windows" {
		t.Skip("IsDriveRoot 仅在 Windows 上测试")
	}

	tests := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		{
			name:     "标准盘符 D:",
			path:     "D:",
			expected: true,
			desc:     "标准盘符格式",
		},
		{
			name:     "带反斜杠 D:\\",
			path:     "D:\\",
			expected: true,
			desc:     "带反斜杠的盘符",
		},
		{
			name:     "带斜杠 D:/",
			path:     "D:/",
			expected: true,
			desc:     "带斜杠的盘符",
		},
		{
			name:     "小写 d:",
			path:     "d:",
			expected: true,
			desc:     "小写盘符",
		},
		{
			name:     "小写带反斜杠 d:\\",
			path:     "d:\\",
			expected: true,
			desc:     "小写带反斜杠",
		},
		{
			name:     "带空格 D: ",
			path:     "D: ",
			expected: true,
			desc:     "带尾部空格的盘符",
		},
		{
			name:     "带前导空格  D:",
			path:     "  D:",
			expected: true,
			desc:     "带前导空格的盘符",
		},
		{
			name:     "非盘符路径",
			path:     "D:\\folder",
			expected: false,
			desc:     "包含子目录的路径",
		},
		{
			name:     "普通路径",
			path:     "C:\\Windows",
			expected: false,
			desc:     "普通文件路径",
		},
		{
			name:     "空字符串",
			path:     "",
			expected: false,
			desc:     "空字符串",
		},
		{
			name:     "单字符",
			path:     "D",
			expected: false,
			desc:     "单字符",
		},
		{
			name:     "无冒号",
			path:     "DX",
			expected: false,
			desc:     "无冒号",
		},
		{
			name:     "数字盘符",
			path:     "1:",
			expected: false,
			desc:     "数字不是有效盘符",
		},
		{
			name:     "太长路径",
			path:     "D:\\folder\\file",
			expected: false,
			desc:     "超过3个字符的路径",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDriveRoot(tt.path)
			if result != tt.expected {
				t.Errorf("IsDriveRoot(%q) = %v, 期望 %v - %s", tt.path, result, tt.expected, tt.desc)
			}
		})
	}
}

// TestGetFileOwner 测试获取文件所有者
func TestGetFileOwner(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "owner_test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		expectEmpty bool
		desc        string
	}{
		{
			name:        "存在的文件",
			path:        testFile,
			expectEmpty: false,
			desc:        "存在的文件应该返回所有者信息",
		},
		{
			name:        "不存在的文件",
			path:        filepath.Join(tempDir, "nonexistent.txt"),
			expectEmpty: true,
			desc:        "不存在的文件应该返回 ?",
		},
		{
			name:        "空路径",
			path:        "",
			expectEmpty: true,
			desc:        "空路径应该返回 ?",
		},
		{
			name:        "目录",
			path:        tempDir,
			expectEmpty: false,
			desc:        "目录也应该返回所有者信息",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, group := GetFileOwner(tt.path)

			if tt.expectEmpty {
				if owner != "?" || group != "?" {
					t.Errorf("GetFileOwner(%q) = (%q, %q), 期望 (?, ?) - %s",
						tt.path, owner, group, tt.desc)
				}
			} else {
				// 在 Unix 系统上应该返回有效的用户名
				// 在 Windows 上可能返回 ? 或实际用户名
				if runtime.GOOS != "windows" {
					if owner == "" || owner == "?" {
						t.Errorf("GetFileOwner(%q) 返回空所有者 - %s", tt.path, tt.desc)
					}
				}
				t.Logf("GetFileOwner(%q) = (owner: %q, group: %q)", tt.path, owner, group)
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
			owner, group := GetFileOwner(tt.path)

			// 对于无效路径，函数应该返回false或?
			t.Logf("路径 %q: IsHidden=%v, IsReadOnly=%v, Owner=%q, Group=%q - %s",
				tt.path, hidden, readonly, owner, group, tt.desc)
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

// BenchmarkIsDriveRoot 基准测试盘符根目录检测
func BenchmarkIsDriveRoot(b *testing.B) {
	if runtime.GOOS != "windows" {
		b.Skip("仅在 Windows 上测试")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsDriveRoot("C:\\")
		IsDriveRoot("D:")
		IsDriveRoot("E:/")
	}
}

// BenchmarkGetFileOwner 基准测试获取文件所有者
func BenchmarkGetFileOwner(b *testing.B) {
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "owner_benchmark.txt")

	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		b.Fatalf("创建基准测试文件失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetFileOwner(testFile)
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
				_, _ = GetFileOwner(hiddenFile)
			}
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}
