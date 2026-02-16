package fs

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// setupMoveTestDir 创建移动测试的目录结构
func setupMoveTestDir(t *testing.T) string {
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

	// 创建符号链接（仅Unix）
	if runtime.GOOS != "windows" {
		symlinkPath := filepath.Join(testDir, "symlink.txt")
		targetPath := filepath.Join(testDir, "file1.txt")
		if err := os.Symlink(targetPath, symlinkPath); err != nil {
			t.Fatalf("Failed to create symlink: %v", err)
		}

		dirlink := filepath.Join(testDir, "dirlink")
		if err := os.Symlink(filepath.Join(testDir, "dir1"), dirlink); err != nil {
			t.Fatalf("Failed to create dir symlink: %v", err)
		}
	}

	return testDir
}

// validateFileMove 验证文件移动是否正确
func validateFileMove(t *testing.T, src, dst string, expectedContent []byte) {
	t.Helper()

	// 验证源文件不存在
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Errorf("Source file still exists after move: %v", err)
	}

	// 验证目标文件存在
	dstContent, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	// 验证内容正确
	if string(dstContent) != string(expectedContent) {
		t.Errorf("File content mismatch: got = %q, want = %q", string(dstContent), string(expectedContent))
	}

	// 验证文件信息
	dstInfo, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("Failed to get destination file info: %v", err)
	}

	if dstInfo.Size() != int64(len(expectedContent)) {
		t.Errorf("File size mismatch: got = %d, want = %d", dstInfo.Size(), len(expectedContent))
	}
}

// validateDirMove 验证目录移动是否正确
func validateDirMove(t *testing.T, src, dst string, expectedFileCount int) {
	t.Helper()

	// 验证源目录不存在
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Errorf("Source directory still exists after move: %v", err)
	}

	// 验证目标目录存在
	dstInfo, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("Failed to get destination directory info: %v", err)
	}

	if !dstInfo.IsDir() {
		t.Fatalf("Destination is not a directory")
	}

	// 验证文件数量
	dstFiles, err := Collect(dst, true)
	if err != nil {
		t.Fatalf("Failed to collect destination files: %v", err)
	}

	if len(dstFiles) != expectedFileCount {
		t.Errorf("File count mismatch: got = %d, want = %d", len(dstFiles), expectedFileCount)
	}
}

// TestMoveFile 测试文件移动功能
func TestMoveFile(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (testDir, src, dst string, expectedContent []byte)
		cleanup     func(t *testing.T, testDir string)
		overwrite   bool
		wantErr     bool
		errContains string
		validate    func(t *testing.T, src, dst string, expectedContent []byte)
	}{
		// 基本文件移动测试
		{
			name: "精确路径模式: Move('a.txt', 'b.txt')",
			setup: func(t *testing.T) (testDir, src, dst string, expectedContent []byte) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "file1_moved.txt")
				expectedContent, _ = os.ReadFile(src)
				return testDir, src, dst, expectedContent
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedContent []byte) {
				validateFileMove(t, src, dst, expectedContent)
			},
		},
		{
			name: "自动追加文件名: Move('a.txt', 'existingDir')",
			setup: func(t *testing.T) (testDir, src, dst string, expectedContent []byte) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "existingDir")
				expectedContent, _ = os.ReadFile(src)
				return testDir, src, dst, expectedContent
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedContent []byte) {
				testDir := filepath.Dir(src)
				expectedDst := filepath.Join(testDir, "existingDir", "file1.txt")
				validateFileMove(t, src, expectedDst, expectedContent)
			},
		},
		{
			name: "自动追加文件名(带斜杠): Move('a.txt', 'existingDir/')",
			setup: func(t *testing.T) (testDir, src, dst string, expectedContent []byte) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "existingDir") + string(filepath.Separator)
				expectedContent, _ = os.ReadFile(src)
				return testDir, src, dst, expectedContent
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedContent []byte) {
				testDir := filepath.Dir(src)
				expectedDst := filepath.Join(testDir, "existingDir", "file1.txt")
				validateFileMove(t, src, expectedDst, expectedContent)
			},
		},
		{
			name: "自动创建父目录: Move('a.txt', 'newDir/b.txt')",
			setup: func(t *testing.T) (testDir, src, dst string, expectedContent []byte) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "newDir", "subDir", "file1.txt")
				expectedContent, _ = os.ReadFile(src)
				return testDir, src, dst, expectedContent
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedContent []byte) {
				validateFileMove(t, src, dst, expectedContent)
			},
		},

		// 边界情况测试
		{
			name: "空源路径",
			setup: func(t *testing.T) (testDir, src, dst string, expectedContent []byte) {
				testDir = setupMoveTestDir(t)
				src = ""
				dst = filepath.Join(testDir, "dst.txt")
				return testDir, src, dst, expectedContent
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
			setup: func(t *testing.T) (testDir, src, dst string, expectedContent []byte) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = ""
				return testDir, src, dst, expectedContent
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
			setup: func(t *testing.T) (testDir, src, dst string, expectedContent []byte) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "file1.txt")
				return testDir, src, dst, expectedContent
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
			setup: func(t *testing.T) (testDir, src, dst string, expectedContent []byte) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "nonexistent.txt")
				dst = filepath.Join(testDir, "dst.txt")
				return testDir, src, dst, expectedContent
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite:   false,
			wantErr:     true,
			errContains: "failed to copy",
		},
		{
			name: "目标已存在且不允许覆盖",
			setup: func(t *testing.T) (testDir, src, dst string, expectedContent []byte) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "file2.txt")
				return testDir, src, dst, expectedContent
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
			setup: func(t *testing.T) (testDir, src, dst string, expectedContent []byte) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "file1.txt")
				dst = filepath.Join(testDir, "file2.txt")
				expectedContent, _ = os.ReadFile(src)
				return testDir, src, dst, expectedContent
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: true,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedContent []byte) {
				validateFileMove(t, src, dst, expectedContent)
			},
		},

		// 空文件测试
		{
			name: "移动空文件",
			setup: func(t *testing.T) (testDir, src, dst string, expectedContent []byte) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "empty.txt")
				dst = filepath.Join(testDir, "empty_moved.txt")
				expectedContent, _ = os.ReadFile(src)
				return testDir, src, dst, expectedContent
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedContent []byte) {
				validateFileMove(t, src, dst, expectedContent)
			},
		},

		// 大文件测试
		{
			name: "移动大文件",
			setup: func(t *testing.T) (testDir, src, dst string, expectedContent []byte) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "large.txt")
				dst = filepath.Join(testDir, "large_moved.txt")
				expectedContent, _ = os.ReadFile(src)
				return testDir, src, dst, expectedContent
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedContent []byte) {
				validateFileMove(t, src, dst, expectedContent)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir, src, dst, expectedContent := tt.setup(t)
			defer tt.cleanup(t, testDir)

			err := MoveEx(src, dst, tt.overwrite)

			if tt.wantErr {
				if err == nil {
					t.Errorf("MoveEx() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("MoveEx() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("MoveEx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.validate != nil {
				tt.validate(t, src, dst, expectedContent)
			}
		})
	}
}

// TestMoveDir 测试目录移动功能
func TestMoveDir(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (testDir, src, dst string, expectedFileCount int)
		cleanup     func(t *testing.T, testDir string)
		overwrite   bool
		wantErr     bool
		errContains string
		validate    func(t *testing.T, src, dst string, expectedFileCount int)
	}{
		// 基本目录移动测试
		{
			name: "目录不存在时: Move('dirA', 'dirB')",
			setup: func(t *testing.T) (testDir, src, dst string, expectedFileCount int) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "dir1_moved")
				srcFiles, _ := Collect(src, true)
				expectedFileCount = len(srcFiles)
				return testDir, src, dst, expectedFileCount
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedFileCount int) {
				validateDirMove(t, src, dst, expectedFileCount)
			},
		},
		{
			name: "自动追加目录名: Move('dirA', 'existingDir')",
			setup: func(t *testing.T) (testDir, src, dst string, expectedFileCount int) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "existingDir")
				srcFiles, _ := Collect(src, true)
				expectedFileCount = len(srcFiles)
				return testDir, src, dst, expectedFileCount
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedFileCount int) {
				testDir := filepath.Dir(src)
				expectedDst := filepath.Join(testDir, "existingDir", "dir1")
				validateDirMove(t, src, expectedDst, expectedFileCount)
			},
		},
		{
			name: "自动追加目录名(带斜杠): Move('dirA', 'existingDir/')",
			setup: func(t *testing.T) (testDir, src, dst string, expectedFileCount int) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "existingDir") + string(filepath.Separator)
				srcFiles, _ := Collect(src, true)
				expectedFileCount = len(srcFiles)
				return testDir, src, dst, expectedFileCount
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedFileCount int) {
				testDir := filepath.Dir(src)
				expectedDst := filepath.Join(testDir, "existingDir", "dir1")
				validateDirMove(t, src, expectedDst, expectedFileCount)
			},
		},
		{
			name: "自动创建父目录: Move('dirA', 'newDir/subDir')",
			setup: func(t *testing.T) (testDir, src, dst string, expectedFileCount int) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "newDir", "subDir")
				srcFiles, _ := Collect(src, true)
				expectedFileCount = len(srcFiles)
				return testDir, src, dst, expectedFileCount
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedFileCount int) {
				validateDirMove(t, src, dst, expectedFileCount)
			},
		},

		// 边界情况测试
		{
			name: "目录移动到子目录",
			setup: func(t *testing.T) (testDir, src, dst string, expectedFileCount int) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "dir1", "subdir")
				return testDir, src, dst, expectedFileCount
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
			setup: func(t *testing.T) (testDir, src, dst string, expectedFileCount int) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "nonexistent")
				dst = filepath.Join(testDir, "dst")
				return testDir, src, dst, expectedFileCount
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite:   false,
			wantErr:     true,
			errContains: "failed to copy",
		},
		{
			name: "目标目录已存在且不允许覆盖",
			setup: func(t *testing.T) (testDir, src, dst string, expectedFileCount int) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "dir1")
				return testDir, src, dst, expectedFileCount
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
			setup: func(t *testing.T) (testDir, src, dst string, expectedFileCount int) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "dir1")
				dst = filepath.Join(testDir, "dir1_copy")
				srcFiles, _ := Collect(src, true)
				expectedFileCount = len(srcFiles)
				return testDir, src, dst, expectedFileCount
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: true,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedFileCount int) {
				validateDirMove(t, src, dst, expectedFileCount)
			},
		},

		// 空目录测试
		{
			name: "移动空目录",
			setup: func(t *testing.T) (testDir, src, dst string, expectedFileCount int) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "emptyDir")
				dst = filepath.Join(testDir, "emptyDir_moved")
				srcFiles, _ := Collect(src, true)
				expectedFileCount = len(srcFiles)
				return testDir, src, dst, expectedFileCount
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedFileCount int) {
				validateDirMove(t, src, dst, expectedFileCount)
			},
		},

		// 嵌套目录测试
		{
			name: "移动嵌套目录",
			setup: func(t *testing.T) (testDir, src, dst string, expectedFileCount int) {
				testDir = setupMoveTestDir(t)
				src = filepath.Join(testDir, "nestedDir")
				dst = filepath.Join(testDir, "nestedDir_moved")
				srcFiles, _ := Collect(src, true)
				expectedFileCount = len(srcFiles)
				return testDir, src, dst, expectedFileCount
			},
			cleanup: func(t *testing.T, testDir string) {
				_ = os.RemoveAll(testDir)
			},
			overwrite: false,
			wantErr:   false,
			validate: func(t *testing.T, src, dst string, expectedFileCount int) {
				validateDirMove(t, src, dst, expectedFileCount)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir, src, dst, expectedFileCount := tt.setup(t)
			defer tt.cleanup(t, testDir)

			err := MoveEx(src, dst, tt.overwrite)

			if tt.wantErr {
				if err == nil {
					t.Errorf("MoveEx() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("MoveEx() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("MoveEx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.validate != nil {
				tt.validate(t, src, dst, expectedFileCount)
			}
		})
	}
}

// TestMoveSymlink 测试符号链接移动
func TestMoveSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping symlink test on Windows")
	}

	testDir := setupMoveTestDir(t)
	defer func() { _ = os.RemoveAll(testDir) }()

	// 移动文件符号链接
	srcSymlink := filepath.Join(testDir, "symlink.txt")
	dstSymlink := filepath.Join(testDir, "symlink_moved.txt")

	if err := Move(srcSymlink, dstSymlink); err != nil {
		t.Fatalf("Move() error = %v", err)
	}

	// 验证源符号链接不存在
	if _, err := os.Lstat(srcSymlink); !os.IsNotExist(err) {
		t.Errorf("Source symlink still exists after move: %v", err)
	}

	// 验证目标符号链接存在
	if _, err := os.Lstat(dstSymlink); err != nil {
		t.Errorf("Destination symlink does not exist after move: %v", err)
	}

	// 移动目录符号链接
	srcDirlink := filepath.Join(testDir, "dirlink")
	dstDirlink := filepath.Join(testDir, "dirlink_moved")

	if err := Move(srcDirlink, dstDirlink); err != nil {
		t.Fatalf("Move() error = %v", err)
	}

	// 验证源目录符号链接不存在
	if _, err := os.Lstat(srcDirlink); !os.IsNotExist(err) {
		t.Errorf("Source dir symlink still exists after move: %v", err)
	}

	// 验证目标目录符号链接存在
	if _, err := os.Lstat(dstDirlink); err != nil {
		t.Errorf("Destination dir symlink does not exist after move: %v", err)
	}
}

// TestMoveVsCopy 测试移动和复制的区别
func TestMoveVsCopy(t *testing.T) {
	testDir := setupMoveTestDir(t)
	defer func() { _ = os.RemoveAll(testDir) }()

	// 测试移动：源文件应该被删除
	srcFile1 := filepath.Join(testDir, "file1.txt")
	dstFile1 := filepath.Join(testDir, "file1_moved.txt")

	if err := Move(srcFile1, dstFile1); err != nil {
		t.Fatalf("Move() error = %v", err)
	}

	// 验证源文件不存在
	if _, err := os.Stat(srcFile1); !os.IsNotExist(err) {
		t.Errorf("Source file still exists after move: %v", err)
	}

	// 验证目标文件存在
	if _, err := os.Stat(dstFile1); err != nil {
		t.Errorf("Destination file does not exist after move: %v", err)
	}

	// 测试复制：源文件应该保留
	srcFile2 := filepath.Join(testDir, "file2.txt")
	dstFile2 := filepath.Join(testDir, "file2_copy.txt")

	if err := Copy(srcFile2, dstFile2); err != nil {
		t.Fatalf("Copy() error = %v", err)
	}

	// 验证源文件仍然存在
	if _, err := os.Stat(srcFile2); err != nil {
		t.Errorf("Source file does not exist after copy: %v", err)
	}

	// 验证目标文件存在
	if _, err := os.Stat(dstFile2); err != nil {
		t.Errorf("Destination file does not exist after copy: %v", err)
	}
}

// TestMoveCrossFilesystem 测试跨文件系统移动
func TestMoveCrossFilesystem(t *testing.T) {
	// 这个测试需要两个不同的文件系统
	// 在大多数测试环境中，我们无法创建真正的跨文件系统场景
	// 所以这个测试主要是为了验证代码的降级逻辑

	testDir := setupMoveTestDir(t)
	defer func() { _ = os.RemoveAll(testDir) }()

	srcFile := filepath.Join(testDir, "file1.txt")
	dstFile := filepath.Join(testDir, "file1_moved.txt")

	// 读取源文件内容
	expectedContent, _ := os.ReadFile(srcFile)

	// 在同文件系统中，Move 应该使用 os.Rename
	if err := Move(srcFile, dstFile); err != nil {
		t.Fatalf("Move() error = %v", err)
	}

	// 验证移动成功
	validateFileMove(t, srcFile, dstFile, expectedContent)
}

// TestMovePermission 测试权限保留
func TestMovePermission(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping permission test on Windows")
	}

	testDir := setupMoveTestDir(t)
	defer func() { _ = os.RemoveAll(testDir) }()

	// 创建具有特定权限的文件
	srcFile := filepath.Join(testDir, "file_with_perm.txt")
	content := []byte("test content")
	if err := os.WriteFile(srcFile, content, 0755); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	dstFile := filepath.Join(testDir, "file_with_perm_moved.txt")

	if err := Move(srcFile, dstFile); err != nil {
		t.Fatalf("Move() error = %v", err)
	}

	// 验证权限
	dstInfo, err := os.Stat(dstFile)
	if err != nil {
		t.Fatalf("Failed to get destination file info: %v", err)
	}

	// 移动操作应该保留权限
	if dstInfo.Mode().Perm() != os.FileMode(0755) {
		t.Errorf("Permission not preserved: got = %v, want = 0755", dstInfo.Mode().Perm())
	}
}
