package fs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupCopyTestDir 创建复制测试的目录结构
func setupCopyTestDir(t *testing.T) string {
	t.Helper()

	testDir := t.TempDir()

	// 创建文件
	files := map[string]string{
		"file1.txt":                  "Hello, World!",
		"file2.txt":                  "Go is awesome",
		"empty.txt":                  "",
		"large.txt":                  strings.Repeat("This is a large file for testing. ", 10000),
		"dir1/file1.txt":             "File in dir1",
		"dir1/file2.txt":             "Another file in dir1",
		"dir2/file1.txt":             "File in dir2",
		"nestedDir/subdir/file1.txt": "Nested file",
	}

	for filePath, content := range files {
		fullPath := filepath.Join(testDir, filePath)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// 创建空目录
	emptyDir := filepath.Join(testDir, "emptyDir")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("Failed to create empty dir: %v", err)
	}

	// 创建已存在的目录（用于测试自动追加）
	existingDir := filepath.Join(testDir, "existingDir")
	if err := os.MkdirAll(existingDir, 0755); err != nil {
		t.Fatalf("Failed to create existing dir: %v", err)
	}

	return testDir
}

// validateFileCopy 验证文件复制是否正确
func validateFileCopy(t *testing.T, src, dst string) {
	t.Helper()

	// 验证源文件仍然存在
	if _, err := os.Stat(src); err != nil {
		t.Errorf("Source file does not exist after copy: %v", err)
	}

	// 验证目标文件存在
	dstContent, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	// 验证内容正确
	srcContent, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("Failed to read source file: %v", err)
	}

	if string(dstContent) != string(srcContent) {
		t.Errorf("File content mismatch: got = %q, want = %q", string(dstContent), string(srcContent))
	}

	// 验证文件信息
	dstInfo, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("Failed to get destination file info: %v", err)
	}

	if dstInfo.Size() != int64(len(srcContent)) {
		t.Errorf("File size mismatch: got = %d, want = %d", dstInfo.Size(), len(srcContent))
	}
}

// validateDirCopy 验证目录复制是否正确
func validateDirCopy(t *testing.T, src, dst string) {
	t.Helper()

	// 验证源目录仍然存在
	if _, err := os.Stat(src); err != nil {
		t.Errorf("Source directory does not exist after copy: %v", err)
	}

	// 验证目标目录存在
	dstFiles, err := Collect(dst, true)
	if err != nil {
		t.Fatalf("Failed to collect destination files: %v", err)
	}

	// 验证文件数量
	srcFiles, err := Collect(src, true)
	if err != nil {
		t.Fatalf("Failed to collect source files: %v", err)
	}

	if len(dstFiles) != len(srcFiles) {
		t.Errorf("File count mismatch: got = %d, want = %d", len(dstFiles), len(srcFiles))
	}
}

// TestCopyFile 测试文件复制功能
func TestCopyFile(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (testDir, src, dst string)
		cleanup     func(t *testing.T, testDir string)
		overwrite   bool
		wantErr     bool
		errContains string
		validate    func(t *testing.T, src, dst string)
	}{
		// 基本文件复制测试
		{
			name: "精确路径模式: Copy('a.txt', 'b.txt')",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "file1_copy.txt")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				validateFileCopy(t, src, dst)
			},
		},
		{
			name: "自动追加文件名: Copy('a.txt', 'existingDir')",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "existingDir")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				testDir := filepath.Dir(src)
				expectedDst := filepath.Join(testDir, "existingDir", "file1.txt")
				validateFileCopy(t, src, expectedDst)
			},
		},
		{
			name: "自动追加文件名(带斜杠): Copy('a.txt', 'existingDir/')",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "existingDir") + string(filepath.Separator)
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				testDir := filepath.Dir(src)
				expectedDst := filepath.Join(testDir, "existingDir", "file1.txt")
				validateFileCopy(t, src, expectedDst)
			},
		},
		{
			name: "自动创建父目录: Copy('a.txt', 'newDir/b.txt')",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "newDir", "subDir", "file1.txt")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				validateFileCopy(t, src, dst)
			},
		},

		// 边界情况测试
		{
			name: "空源路径",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = ""
				dst = filepath.Join(testDir, "dst.txt")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite:   false,
			wantErr:     true,
			errContains: "cannot be empty",
		},
		{
			name: "空目标路径",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = ""
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite:   false,
			wantErr:     true,
			errContains: "cannot be empty",
		},
		{
			name: "源路径和目标路径相同",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "file1.txt")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite:   false,
			wantErr:     true,
			errContains: "cannot be the same",
		},
		{
			name: "源文件不存在",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "nonexistent.txt")
				dst = filepath.Join(testDir, "dst.txt")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite:   false,
			wantErr:     true,
			errContains: "failed to get source info",
		},
		{
			name: "目标已存在且不允许覆盖",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "file2.txt")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite:   false,
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name: "目标已存在但允许覆盖",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "file2.txt")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: true,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				validateFileCopy(t, src, dst)
			},
		},

		// 空文件测试
		{
			name: "复制空文件",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "empty.txt")
				dst = filepath.Join(testDir, "empty_copy.txt")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				validateFileCopy(t, src, dst)
			},
		},

		// 大文件测试
		{
			name: "复制大文件",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "large.txt")
				dst = filepath.Join(testDir, "large_copy.txt")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				validateFileCopy(t, src, dst)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir, src, dst := tt.setup(t)
			defer tt.cleanup(t, testDir)

			err := CopyEx(src, dst, tt.overwrite)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CopyEx() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CopyEx() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("CopyEx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.validate != nil {
				tt.validate(t, src, dst)
			}
		})
	}
}

// TestCopyDir 测试目录复制功能
func TestCopyDir(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (testDir, src, dst string)
		cleanup     func(t *testing.T, testDir string)
		overwrite   bool
		wantErr     bool
		errContains string
		validate    func(t *testing.T, src, dst string)
	}{
		// 基本目录复制测试
		{
			name: "目录不存在时: Copy('dirA', 'dirB')",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "dir1_copy")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				validateDirCopy(t, src, dst)
			},
		},
		{
			name: "自动追加目录名: Copy('dirA', 'existingDir')",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "existingDir")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				testDir := filepath.Dir(src)
				expectedDst := filepath.Join(testDir, "existingDir", "dir1")
				validateDirCopy(t, src, expectedDst)
			},
		},
		{
			name: "自动追加目录名(带斜杠): Copy('dirA', 'existingDir/')",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "existingDir") + string(filepath.Separator)
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				testDir := filepath.Dir(src)
				expectedDst := filepath.Join(testDir, "existingDir", "dir1")
				validateDirCopy(t, src, expectedDst)
			},
		},
		{
			name: "自动创建父目录: Copy('dirA', 'newDir/subDir')",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "newDir", "subDir")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				validateDirCopy(t, src, dst)
			},
		},

		// 边界情况测试
		{
			name: "目录复制到子目录",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "dir1", "subdir")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite:   false,
			wantErr:     true,
			errContains: "cannot copy directory",
		},
		{
			name: "源目录不存在",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "nonexistent")
				dst = filepath.Join(testDir, "dst")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite:   false,
			wantErr:     true,
			errContains: "failed to get source info",
		},
		{
			name: "目标目录已存在且不允许覆盖",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "dir1")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite:   false,
			wantErr:     true,
			errContains: "cannot be the same",
		},
		{
			name: "目标目录已存在但允许覆盖",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "dir1_copy")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: true,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				validateDirCopy(t, src, dst)
			},
		},

		// 空目录测试
		{
			name: "复制空目录",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "emptyDir")
				dst = filepath.Join(testDir, "emptyDir_copy")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				validateDirCopy(t, src, dst)
			},
		},

		// 嵌套目录测试
		{
			name: "复制嵌套目录",
			setup: func(t *testing.T) (testDir, src, dst string) {
				testDir = setupCopyTestDir(t)
				src = filepath.Join(testDir, "nestedDir")
				dst = filepath.Join(testDir, "nestedDir_copy")
				return testDir, src, dst
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string) {
				validateDirCopy(t, src, dst)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir, src, dst := tt.setup(t)
			defer tt.cleanup(t, testDir)

			err := CopyEx(src, dst, tt.overwrite)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CopyEx() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CopyEx() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("CopyEx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.validate != nil {
				tt.validate(t, src, dst)
			}
		})
	}
}
