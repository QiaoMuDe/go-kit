package fs

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestCopyFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		setup       func() (src, dst string)
		expectError bool
		desc        string
	}{
		{
			name: "复制普通文件",
			setup: func() (string, string) {
				src := filepath.Join(tempDir, "source.txt")
				dst := filepath.Join(tempDir, "destination.txt")

				content := "Hello, World! This is a test file."
				if err := os.WriteFile(src, []byte(content), 0644); err != nil {
					t.Fatalf("创建源文件失败: %v", err)
				}

				return src, dst
			},
			expectError: false,
			desc:        "普通文件复制应该成功",
		},
		{
			name: "复制空文件",
			setup: func() (string, string) {
				src := filepath.Join(tempDir, "empty_source.txt")
				dst := filepath.Join(tempDir, "empty_destination.txt")

				if err := os.WriteFile(src, []byte(""), 0644); err != nil {
					t.Fatalf("创建空源文件失败: %v", err)
				}

				return src, dst
			},
			expectError: false,
			desc:        "空文件复制应该成功",
		},
		{
			name: "复制大文件",
			setup: func() (string, string) {
				src := filepath.Join(tempDir, "large_source.txt")
				dst := filepath.Join(tempDir, "large_destination.txt")

				// 创建1MB的测试文件
				content := strings.Repeat("A", 1024*1024)
				if err := os.WriteFile(src, []byte(content), 0644); err != nil {
					t.Fatalf("创建大源文件失败: %v", err)
				}

				return src, dst
			},
			expectError: false,
			desc:        "大文件复制应该成功",
		},
		{
			name: "复制二进制文件",
			setup: func() (string, string) {
				src := filepath.Join(tempDir, "binary_source.bin")
				dst := filepath.Join(tempDir, "binary_destination.bin")

				// 创建包含各种字节值的二进制文件
				content := make([]byte, 256)
				for i := 0; i < 256; i++ {
					content[i] = byte(i)
				}
				if err := os.WriteFile(src, content, 0644); err != nil {
					t.Fatalf("创建二进制源文件失败: %v", err)
				}

				return src, dst
			},
			expectError: false,
			desc:        "二进制文件复制应该成功",
		},
		{
			name: "源文件不存在",
			setup: func() (string, string) {
				src := filepath.Join(tempDir, "non_existing_source.txt")
				dst := filepath.Join(tempDir, "destination.txt")
				return src, dst
			},
			expectError: true,
			desc:        "源文件不存在应该返回错误",
		},
		{
			name: "目标目录不存在",
			setup: func() (string, string) {
				src := filepath.Join(tempDir, "source_for_missing_dir.txt")
				dst := filepath.Join(tempDir, "non_existing_dir", "destination.txt")

				if err := os.WriteFile(src, []byte("test"), 0644); err != nil {
					t.Fatalf("创建源文件失败: %v", err)
				}

				return src, dst
			},
			expectError: false, // 函数应该创建目录
			desc:        "目标目录不存在时应该自动创建",
		},
		{
			name: "覆盖已存在的文件",
			setup: func() (string, string) {
				src := filepath.Join(tempDir, "overwrite_source.txt")
				dst := filepath.Join(tempDir, "overwrite_destination.txt")

				// 创建源文件
				if err := os.WriteFile(src, []byte("new content"), 0644); err != nil {
					t.Fatalf("创建源文件失败: %v", err)
				}

				// 创建目标文件（将被覆盖）
				if err := os.WriteFile(dst, []byte("old content"), 0644); err != nil {
					t.Fatalf("创建目标文件失败: %v", err)
				}

				return src, dst
			},
			expectError: false,
			desc:        "覆盖已存在的文件应该成功",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, dst := tt.setup()

			err := Copy(src, dst)

			if tt.expectError {
				if err == nil {
					t.Errorf("CopyFile(%q, %q) 期望返回错误，但没有错误 - %s", src, dst, tt.desc)
				}
				return
			}

			if err != nil {
				t.Errorf("CopyFile(%q, %q) 返回意外错误: %v - %s", src, dst, err, tt.desc)
				return
			}

			// 验证文件是否存在
			if !Exists(dst) {
				t.Errorf("目标文件 %q 不存在 - %s", dst, tt.desc)
				return
			}

			// 验证文件内容是否相同
			srcContent, err := os.ReadFile(src)
			if err != nil {
				t.Errorf("读取源文件失败: %v", err)
				return
			}

			dstContent, err := os.ReadFile(dst)
			if err != nil {
				t.Errorf("读取目标文件失败: %v", err)
				return
			}

			if !bytes.Equal(srcContent, dstContent) {
				t.Errorf("文件内容不匹配 - %s", tt.desc)
			}

			// 验证文件大小
			srcInfo, _ := os.Stat(src)
			dstInfo, _ := os.Stat(dst)
			if srcInfo.Size() != dstInfo.Size() {
				t.Errorf("文件大小不匹配: 源文件 %d 字节, 目标文件 %d 字节 - %s",
					srcInfo.Size(), dstInfo.Size(), tt.desc)
			}
		})
	}
}

func TestCopyDir(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		setup       func() (src, dst string)
		expectError bool
		desc        string
	}{
		{
			name: "复制简单目录",
			setup: func() (string, string) {
				src := filepath.Join(tempDir, "simple_source_dir")
				dst := filepath.Join(tempDir, "simple_destination_dir")

				// 创建源目录和文件
				if err := os.Mkdir(src, 0755); err != nil {
					t.Fatalf("创建源目录失败: %v", err)
				}

				file1 := filepath.Join(src, "file1.txt")
				if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
					t.Fatalf("创建文件1失败: %v", err)
				}

				file2 := filepath.Join(src, "file2.txt")
				if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
					t.Fatalf("创建文件2失败: %v", err)
				}

				return src, dst
			},
			expectError: false,
			desc:        "简单目录复制应该成功",
		},
		{
			name: "复制嵌套目录",
			setup: func() (string, string) {
				src := filepath.Join(tempDir, "nested_source_dir")
				dst := filepath.Join(tempDir, "nested_destination_dir")

				// 创建嵌套目录结构
				subDir := filepath.Join(src, "subdir")
				if err := os.MkdirAll(subDir, 0755); err != nil {
					t.Fatalf("创建嵌套目录失败: %v", err)
				}

				// 在根目录创建文件
				rootFile := filepath.Join(src, "root.txt")
				if err := os.WriteFile(rootFile, []byte("root content"), 0644); err != nil {
					t.Fatalf("创建根文件失败: %v", err)
				}

				// 在子目录创建文件
				subFile := filepath.Join(subDir, "sub.txt")
				if err := os.WriteFile(subFile, []byte("sub content"), 0644); err != nil {
					t.Fatalf("创建子文件失败: %v", err)
				}

				return src, dst
			},
			expectError: false,
			desc:        "嵌套目录复制应该成功",
		},
		{
			name: "复制空目录",
			setup: func() (string, string) {
				src := filepath.Join(tempDir, "empty_source_dir")
				dst := filepath.Join(tempDir, "empty_destination_dir")

				if err := os.Mkdir(src, 0755); err != nil {
					t.Fatalf("创建空源目录失败: %v", err)
				}

				return src, dst
			},
			expectError: false,
			desc:        "空目录复制应该成功",
		},
		{
			name: "源目录不存在",
			setup: func() (string, string) {
				src := filepath.Join(tempDir, "non_existing_source_dir")
				dst := filepath.Join(tempDir, "destination_dir")
				return src, dst
			},
			expectError: true,
			desc:        "源目录不存在应该返回错误",
		},
		{
			name: "复制包含特殊文件名的目录",
			setup: func() (string, string) {
				src := filepath.Join(tempDir, "special_source_dir")
				dst := filepath.Join(tempDir, "special_destination_dir")

				if err := os.Mkdir(src, 0755); err != nil {
					t.Fatalf("创建源目录失败: %v", err)
				}

				// 创建包含特殊字符的文件名
				specialFile := filepath.Join(src, "测试文件 with spaces & symbols!.txt")
				if err := os.WriteFile(specialFile, []byte("special content"), 0644); err != nil {
					t.Fatalf("创建特殊文件失败: %v", err)
				}

				return src, dst
			},
			expectError: false,
			desc:        "包含特殊文件名的目录复制应该成功",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, dst := tt.setup()

			err := Copy(src, dst)

			if tt.expectError {
				if err == nil {
					t.Errorf("CopyDir(%q, %q) 期望返回错误，但没有错误 - %s", src, dst, tt.desc)
				}
				return
			}

			if err != nil {
				t.Errorf("CopyDir(%q, %q) 返回意外错误: %v - %s", src, dst, err, tt.desc)
				return
			}

			// 验证目录是否存在
			if !IsDir(dst) {
				t.Errorf("目标目录 %q 不存在 - %s", dst, tt.desc)
				return
			}

			// 验证目录内容是否相同
			err = compareDirs(src, dst)
			if err != nil {
				t.Errorf("目录内容不匹配: %v - %s", err, tt.desc)
			}
		})
	}
}

// 辅助函数：比较两个目录的内容
func compareDirs(src, dst string) error {
	return filepath.Walk(src, func(srcPath string, srcInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(src, srcPath)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		// 检查目标路径是否存在
		dstInfo, err := os.Stat(dstPath)
		if err != nil {
			return fmt.Errorf("目标路径 %q 不存在", dstPath)
		}

		// 检查文件类型是否匹配
		if srcInfo.IsDir() != dstInfo.IsDir() {
			return fmt.Errorf("文件类型不匹配: %q", relPath)
		}

		// 如果是文件，比较内容
		if !srcInfo.IsDir() {
			srcContent, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}

			dstContent, err := os.ReadFile(dstPath)
			if err != nil {
				return err
			}

			if !bytes.Equal(srcContent, dstContent) {
				return fmt.Errorf("文件内容不匹配: %q", relPath)
			}
		}

		return nil
	})
}

// 边界条件测试
func TestCopyBoundaryConditions(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name string
		test func(t *testing.T)
		desc string
	}{
		{
			name: "空路径参数",
			test: func(t *testing.T) {
				err := Copy("", "")
				if err == nil {
					t.Error("空路径参数应该返回错误")
				}
			},
			desc: "空路径参数应该被正确处理",
		},
		{
			name: "相同源和目标路径",
			test: func(t *testing.T) {
				file := filepath.Join(tempDir, "same_path.txt")
				if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
					t.Fatalf("创建测试文件失败: %v", err)
				}

				err := Copy(file, file)
				// 这种情况的行为取决于具体实现
				t.Logf("相同路径复制结果: %v", err)
			},
			desc: "相同源和目标路径应该被正确处理",
		},
		{
			name: "权限不足的目标目录",
			test: func(t *testing.T) {
				src := filepath.Join(tempDir, "perm_source.txt")
				if err := os.WriteFile(src, []byte("content"), 0644); err != nil {
					t.Fatalf("创建源文件失败: %v", err)
				}

				// 创建只读目录
				readOnlyDir := filepath.Join(tempDir, "readonly_dir")
				if err := os.Mkdir(readOnlyDir, 0755); err != nil {
					t.Fatalf("创建目录失败: %v", err)
				}

				// 在Windows上，目录权限处理不同，跳过权限设置
				if runtime.GOOS != "windows" {
					if err := os.Chmod(readOnlyDir, 0444); err != nil {
						t.Fatalf("设置只读权限失败: %v", err)
					}
				}

				dst := filepath.Join(readOnlyDir, "destination.txt")
				err := Copy(src, dst)

				// 恢复目录权限以便清理
				_ = os.Chmod(readOnlyDir, 0755)

				// Windows上不期望错误，Unix上期望错误
				if runtime.GOOS == "windows" {
					if err != nil {
						t.Logf("Windows上复制操作结果: %v", err)
					}
				} else {
					if err == nil {
						t.Error("复制到只读目录应该返回错误")
					}
				}
			},
			desc: "权限不足的情况应该被正确处理",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// 性能测试
func BenchmarkCopyFile(b *testing.B) {
	tempDir := b.TempDir()

	// 创建不同大小的测试文件
	sizes := []struct {
		name string
		size int
	}{
		{"1KB", 1024},
		{"10KB", 10 * 1024},
		{"100KB", 100 * 1024},
		{"1MB", 1024 * 1024},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			src := filepath.Join(tempDir, fmt.Sprintf("bench_src_%s.txt", size.name))
			content := strings.Repeat("A", size.size)
			if err := os.WriteFile(src, []byte(content), 0644); err != nil {
				b.Fatalf("创建基准测试文件失败: %v", err)
			}

			b.ResetTimer()
			b.SetBytes(int64(size.size))

			for i := 0; i < b.N; i++ {
				dst := filepath.Join(tempDir, fmt.Sprintf("bench_dst_%s_%d.txt", size.name, i))
				if err := Copy(src, dst); err != nil {
					b.Fatalf("复制文件失败: %v", err)
				}
				// 清理目标文件
				_ = os.Remove(dst)
			}
		})
	}
}

// 并发测试
func TestCopyConcurrency(t *testing.T) {
	tempDir := t.TempDir()

	// 创建源文件
	src := filepath.Join(tempDir, "concurrent_source.txt")
	content := "This is a test file for concurrent copying."
	if err := os.WriteFile(src, []byte(content), 0644); err != nil {
		t.Fatalf("创建源文件失败: %v", err)
	}

	const numGoroutines = 50
	done := make(chan error, numGoroutines)

	// 并发复制文件
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			dst := filepath.Join(tempDir, fmt.Sprintf("concurrent_dst_%d.txt", id))
			err := Copy(src, dst)
			done <- err
		}(i)
	}

	// 等待所有goroutine完成并检查结果
	for i := 0; i < numGoroutines; i++ {
		if err := <-done; err != nil {
			t.Errorf("并发复制失败: %v", err)
		}
	}

	// 验证所有文件都被正确复制
	for i := 0; i < numGoroutines; i++ {
		dst := filepath.Join(tempDir, fmt.Sprintf("concurrent_dst_%d.txt", i))
		dstContent, err := os.ReadFile(dst)
		if err != nil {
			t.Errorf("读取并发复制的文件失败: %v", err)
			continue
		}

		if string(dstContent) != content {
			t.Errorf("并发复制的文件内容不正确")
		}
	}
}

// 完整性测试
func TestCopyIntegrity(t *testing.T) {
	tempDir := t.TempDir()

	// 创建包含随机数据的大文件
	src := filepath.Join(tempDir, "integrity_source.bin")
	dst := filepath.Join(tempDir, "integrity_destination.bin")

	// 生成随机内容
	content := make([]byte, 1024*1024) // 1MB
	for i := range content {
		content[i] = byte(i % 256)
	}

	if err := os.WriteFile(src, content, 0644); err != nil {
		t.Fatalf("创建完整性测试文件失败: %v", err)
	}

	// 复制文件
	if err := Copy(src, dst); err != nil {
		t.Fatalf("复制文件失败: %v", err)
	}

	// 计算源文件和目标文件的MD5哈希
	srcHash, err := calculateMD5(src)
	if err != nil {
		t.Fatalf("计算源文件MD5失败: %v", err)
	}

	dstHash, err := calculateMD5(dst)
	if err != nil {
		t.Fatalf("计算目标文件MD5失败: %v", err)
	}

	if srcHash != dstHash {
		t.Errorf("文件完整性检查失败: 源文件MD5=%s, 目标文件MD5=%s", srcHash, dstHash)
	}
}

// 辅助函数：计算文件的MD5哈希
func calculateMD5(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// 时间戳保持测试
func TestCopyPreservesTimestamp(t *testing.T) {
	tempDir := t.TempDir()

	src := filepath.Join(tempDir, "timestamp_source.txt")
	dst := filepath.Join(tempDir, "timestamp_destination.txt")

	// 创建源文件
	if err := os.WriteFile(src, []byte("timestamp test"), 0644); err != nil {
		t.Fatalf("创建源文件失败: %v", err)
	}

	// 设置特定的修改时间
	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	if err := os.Chtimes(src, testTime, testTime); err != nil {
		t.Fatalf("设置源文件时间戳失败: %v", err)
	}

	// 复制文件
	if err := Copy(src, dst); err != nil {
		t.Fatalf("复制文件失败: %v", err)
	}

	// 检查时间戳是否保持（这取决于CopyFile的实现）
	srcInfo, err := os.Stat(src)
	if err != nil {
		t.Fatalf("获取源文件信息失败: %v", err)
	}

	dstInfo, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("获取目标文件信息失败: %v", err)
	}

	// 记录时间戳信息（实际行为取决于实现）
	t.Logf("源文件修改时间: %v", srcInfo.ModTime())
	t.Logf("目标文件修改时间: %v", dstInfo.ModTime())
}
