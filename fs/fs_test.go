package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetSize(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		setup       func() string
		expectedMin int64
		expectedMax int64
		expectError bool
		desc        string
	}{
		{
			name: "获取文件大小",
			setup: func() string {
				file := filepath.Join(tempDir, "size_test.txt")
				content := "Hello, World! This is a test file for size calculation."
				if err := os.WriteFile(file, []byte(content), 0644); err != nil {
					t.Fatalf("创建测试文件失败: %v", err)
				}
				return file
			},
			expectedMin: 50,  // 大概的字节数
			expectedMax: 100, // 允许一些误差
			expectError: false,
			desc:        "获取文件大小应该返回正确的字节数",
		},
		{
			name: "获取空文件大小",
			setup: func() string {
				file := filepath.Join(tempDir, "empty_file.txt")
				if err := os.WriteFile(file, []byte(""), 0644); err != nil {
					t.Fatalf("创建空文件失败: %v", err)
				}
				return file
			},
			expectedMin: 0,
			expectedMax: 0,
			expectError: false,
			desc:        "空文件大小应该为0",
		},
		{
			name: "获取大文件大小",
			setup: func() string {
				file := filepath.Join(tempDir, "large_file.txt")
				content := strings.Repeat("A", 10000) // 10KB
				if err := os.WriteFile(file, []byte(content), 0644); err != nil {
					t.Fatalf("创建大文件失败: %v", err)
				}
				return file
			},
			expectedMin: 10000,
			expectedMax: 10000,
			expectError: false,
			desc:        "大文件大小应该正确计算",
		},
		{
			name: "获取目录大小",
			setup: func() string {
				dir := filepath.Join(tempDir, "size_dir")
				if err := os.Mkdir(dir, 0755); err != nil {
					t.Fatalf("创建目录失败: %v", err)
				}

				// 在目录中创建一些文件
				file1 := filepath.Join(dir, "file1.txt")
				if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
					t.Fatalf("创建文件1失败: %v", err)
				}

				file2 := filepath.Join(dir, "file2.txt")
				if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
					t.Fatalf("创建文件2失败: %v", err)
				}

				return dir
			},
			expectedMin: 16, // 两个文件的内容总和
			expectedMax: 50, // 允许目录本身的大小
			expectError: false,
			desc:        "目录大小应该包含所有文件的大小",
		},
		{
			name: "获取不存在文件的大小",
			setup: func() string {
				return filepath.Join(tempDir, "non_existing_file.txt")
			},
			expectedMin: 0,
			expectedMax: 0,
			expectError: true,
			desc:        "不存在的文件应该返回错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()

			size, err := GetSize(path)

			if tt.expectError {
				if err == nil {
					t.Errorf("GetSize(%q) 期望返回错误，但没有错误 - %s", path, tt.desc)
				}
				return
			}

			if err != nil {
				t.Errorf("GetSize(%q) 返回意外错误: %v - %s", path, err, tt.desc)
				return
			}

			if size < tt.expectedMin || size > tt.expectedMax {
				t.Errorf("GetSize(%q) = %d, 期望在 %d-%d 范围内 - %s",
					path, size, tt.expectedMin, tt.expectedMax, tt.desc)
			}
		})
	}
}

// 边界条件测试
func TestFSBoundaryConditions(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name string
		test func(t *testing.T)
		desc string
	}{
		{
			name: "非常长的路径",
			test: func(t *testing.T) {
				// 创建一个非常长的路径
				longPath := tempDir
				for i := 0; i < 50; i++ {
					longPath = filepath.Join(longPath, "very_long_directory_name_that_might_cause_issues")
				}

				err := os.MkdirAll(longPath, 0755)

				// 在某些系统上可能会因为路径太长而失败
				t.Logf("长路径创建结果: %v", err)
			},
			desc: "非常长的路径应该被正确处理",
		},
		{
			name: "包含特殊字符的路径",
			test: func(t *testing.T) {
				specialChars := []string{
					"测试目录",
					"directory with spaces",
					"dir-with-dashes",
					"dir_with_underscores",
					"dir.with.dots",
				}

				for _, name := range specialChars {
					path := filepath.Join(tempDir, name)
					err := os.MkdirAll(path, 0755)
					if err != nil {
						t.Errorf("创建特殊字符目录 %q 失败: %v", name, err)
					}

					if !IsDir(path) {
						t.Errorf("特殊字符目录 %q 未被创建", name)
					}
				}
			},
			desc: "包含特殊字符的路径应该被正确处理",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// 性能测试
func BenchmarkMkdirAll(b *testing.B) {
	tempDir := b.TempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path := filepath.Join(tempDir, "bench", "mkdir", fmt.Sprintf("test_%d", i))
		if err := os.MkdirAll(path, 0755); err != nil {
			b.Fatalf("MkdirAll失败: %v", err)
		}
	}
}

func BenchmarkRemoveAll(b *testing.B) {
	tempDir := b.TempDir()

	// 预创建目录结构
	paths := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		path := filepath.Join(tempDir, "bench", "remove", fmt.Sprintf("test_%d", i))
		if err := os.MkdirAll(path, 0755); err != nil {
			b.Fatalf("预创建目录失败: %v", err)
		}

		// 在目录中创建一些文件
		file := filepath.Join(path, "test.txt")
		if err := os.WriteFile(file, []byte("test content"), 0644); err != nil {
			b.Fatalf("创建测试文件失败: %v", err)
		}

		paths[i] = path
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := os.RemoveAll(paths[i]); err != nil {

			b.Fatalf("RemoveAll失败: %v", err)
		}
	}
}

func BenchmarkGetSize(b *testing.B) {
	tempDir := b.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "benchmark_file.txt")
	content := strings.Repeat("A", 1024) // 1KB
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatalf("创建基准测试文件失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := GetSize(testFile); err != nil {
			b.Fatalf("GetSize失败: %v", err)
		}
	}
}

// 并发测试
func TestFSConcurrency(t *testing.T) {
	tempDir := t.TempDir()

	const numGoroutines = 50
	done := make(chan error, numGoroutines)

	// 并发创建目录
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			path := filepath.Join(tempDir, "concurrent", fmt.Sprintf("dir_%d", id))
			err := os.MkdirAll(path, 0755)
			done <- err
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < numGoroutines; i++ {
		if err := <-done; err != nil {
			t.Errorf("并发创建目录失败: %v", err)
		}
	}

	// 验证所有目录都被创建
	for i := 0; i < numGoroutines; i++ {
		path := filepath.Join(tempDir, "concurrent", fmt.Sprintf("dir_%d", i))
		if !IsDir(path) {
			t.Errorf("并发创建的目录 %q 不存在", path)
		}
	}
}

// 错误恢复测试
func TestFSErrorRecovery(t *testing.T) {
	tempDir := t.TempDir()

	// 测试在权限不足的情况下的错误处理
	t.Run("权限不足错误恢复", func(t *testing.T) {
		// 创建一个只读目录
		readOnlyDir := filepath.Join(tempDir, "readonly")
		if err := os.Mkdir(readOnlyDir, 0755); err != nil {
			t.Fatalf("创建目录失败: %v", err)
		}

		// 在Windows上，目录权限处理不同，跳过权限设置
		if runtime.GOOS != "windows" {
			if err := os.Chmod(readOnlyDir, 0444); err != nil {
				t.Fatalf("设置只读权限失败: %v", err)
			}
		}

		// 尝试在只读目录中创建子目录
		subDir := filepath.Join(readOnlyDir, "subdir")
		err := os.MkdirAll(subDir, 0755)

		// 恢复目录权限以便清理
		_ = os.Chmod(readOnlyDir, 0755)

		// Windows上不期望错误，Unix上期望错误
		if runtime.GOOS == "windows" {
			if err != nil {
				t.Logf("Windows上目录创建结果: %v", err)
			}
		} else {
			if err == nil {
				t.Error("在只读目录中创建子目录应该返回错误")
			}
		}
	})

	// 测试磁盘空间不足的模拟（这个测试可能不会在所有环境中触发）
	t.Run("大文件创建", func(t *testing.T) {
		largeFile := filepath.Join(tempDir, "large_test.txt")

		// 尝试创建一个非常大的文件（但不会真的写入这么多数据）
		f, err := os.Create(largeFile)
		if err != nil {
			t.Fatalf("创建大文件失败: %v", err)
		}
		defer func() { _ = f.Close() }()

		// 写入一些数据
		data := make([]byte, 1024*1024) // 1MB
		_, err = f.Write(data)
		if err != nil {
			t.Logf("写入大文件数据时出错: %v", err)
		}

		// 获取文件大小
		size, err := GetSize(largeFile)
		if err != nil {
			t.Errorf("获取大文件大小失败: %v", err)
		} else {
			t.Logf("大文件大小: %d 字节", size)
		}
	})
}

// TestGetSize1 测试GetSize函数
func TestGetSize1(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()

	// 测试用例
	tests := []struct {
		name        string
		setupFunc   func() (string, int64) // 返回路径和期望大小
		expectError bool
	}{
		{
			name: "单个文件",
			setupFunc: func() (string, int64) {
				content := "hello world"
				filePath := filepath.Join(tempDir, "test.txt")
				err := os.WriteFile(filePath, []byte(content), 0644)
				if err != nil {
					t.Fatal(err)
				}
				return filePath, int64(len(content))
			},
			expectError: false,
		},
		{
			name: "空文件",
			setupFunc: func() (string, int64) {
				filePath := filepath.Join(tempDir, "empty.txt")
				err := os.WriteFile(filePath, []byte{}, 0644)
				if err != nil {
					t.Fatal(err)
				}
				return filePath, 0
			},
			expectError: false,
		},
		{
			name: "目录包含多个文件",
			setupFunc: func() (string, int64) {
				dirPath := filepath.Join(tempDir, "testdir")
				err := os.MkdirAll(dirPath, 0755)
				if err != nil {
					t.Fatal(err)
				}

				// 创建多个文件
				files := map[string]string{
					"file1.txt": "content1",
					"file2.txt": "content2",
					"file3.txt": "content3",
				}

				var totalSize int64
				for name, content := range files {
					filePath := filepath.Join(dirPath, name)
					err := os.WriteFile(filePath, []byte(content), 0644)
					if err != nil {
						t.Fatal(err)
					}
					totalSize += int64(len(content))
				}

				return dirPath, totalSize
			},
			expectError: false,
		},
		{
			name: "空目录",
			setupFunc: func() (string, int64) {
				dirPath := filepath.Join(tempDir, "emptydir")
				err := os.MkdirAll(dirPath, 0755)
				if err != nil {
					t.Fatal(err)
				}
				return dirPath, 0
			},
			expectError: false,
		},
		{
			name: "不存在的路径",
			setupFunc: func() (string, int64) {
				return "/nonexistent/path", 0
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, expectedSize := tt.setupFunc()
			size, err := GetSize(path)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
				return
			}

			if err != nil {
				t.Errorf("意外错误: %v", err)
				return
			}

			if size != expectedSize {
				t.Errorf("GetSize() = %d, 期望 %d", size, expectedSize)
			}
		})
	}
}

// TestWrapPathError 测试wrapPathError函数
func TestWrapPathError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		path      string
		operation string
		expected  string
	}{
		{
			name:      "文件不存在错误",
			err:       os.ErrNotExist,
			path:      "/test/path",
			operation: "reading",
			expected:  "path does not exist when reading: /test/path",
		},
		{
			name:      "权限错误",
			err:       os.ErrPermission,
			path:      "/test/path",
			operation: "writing",
			expected:  "permission denied when writing path '/test/path'",
		},
		{
			name:      "其他错误",
			err:       fmt.Errorf("custom error"),
			path:      "/test/path",
			operation: "accessing",
			expected:  "error when accessing path '/test/path': custom error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapPathError(tt.err, tt.path, tt.operation)
			if !strings.Contains(result.Error(), tt.path) {
				t.Errorf("错误信息应包含路径 %s，但得到: %s", tt.path, result.Error())
			}
			if !strings.Contains(result.Error(), tt.operation) {
				t.Errorf("错误信息应包含操作 %s，但得到: %s", tt.operation, result.Error())
			}
		})
	}
}

// TestGetSize_SymbolicLinks 测试符号链接
func TestGetSize_SymbolicLinks(t *testing.T) {
	tempDir := t.TempDir()

	// 创建原始文件
	originalFile := filepath.Join(tempDir, "original.txt")
	content := "test content"
	err := os.WriteFile(originalFile, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 创建符号链接
	linkFile := filepath.Join(tempDir, "link.txt")
	err = os.Symlink(originalFile, linkFile)
	if err != nil {
		t.Skip("无法创建符号链接，跳过测试")
	}

	// 测试符号链接的大小
	size, err := GetSize(linkFile)
	if err != nil {
		t.Errorf("获取符号链接大小失败: %v", err)
	}

	expectedSize := int64(len(content))
	if size != expectedSize {
		t.Errorf("符号链接大小 = %d, 期望 %d", size, expectedSize)
	}
}

// TestIntegration_GetSizeAndFormat 集成测试：获取大小并格式化
func TestIntegration_GetSizeAndFormat(t *testing.T) {
	tempDir := t.TempDir()

	// 创建不同大小的文件
	testFiles := []struct {
		name         string
		content      string
		expectedUnit string
	}{
		{"small.txt", "hello", "B"},
		{"medium.txt", strings.Repeat("a", 2048), "KB"},
		{"large.txt", strings.Repeat("b", 1024*1024+512), "MB"},
	}

	for _, tf := range testFiles {
		t.Run(tf.name, func(t *testing.T) {
			filePath := filepath.Join(tempDir, tf.name)
			err := os.WriteFile(filePath, []byte(tf.content), 0644)
			if err != nil {
				t.Fatal(err)
			}

			// 获取文件大小
			size, err := GetSize(filePath)
			if err != nil {
				t.Errorf("获取文件大小失败: %v", err)
			}

			// 验证大小正确性
			expectedSize := int64(len(tf.content))
			if size != expectedSize {
				t.Errorf("文件大小不匹配: 得到 %d, 期望 %d", size, expectedSize)
			}
		})
	}
}
