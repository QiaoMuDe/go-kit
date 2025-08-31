package fs

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestCollect(t *testing.T) {
	// 创建测试目录结构
	testDir := setupTestDir(t)
	defer func() { _ = os.RemoveAll(testDir) }()

	tests := []struct {
		name        string
		targetPath  string
		recursive   bool
		want        []string
		wantErr     bool
		errContains string
	}{
		// 基本功能测试
		{
			name:        "empty path should return error",
			targetPath:  "",
			recursive:   false,
			want:        nil,
			wantErr:     true,
			errContains: "target path cannot be empty",
		},
		{
			name:       "single file",
			targetPath: filepath.Join(testDir, "file1.txt"),
			recursive:  false,
			want:       []string{filepath.Join(testDir, "file1.txt")},
			wantErr:    false,
		},
		{
			name:        "non-existent file",
			targetPath:  filepath.Join(testDir, "nonexistent.txt"),
			recursive:   false,
			want:        nil,
			wantErr:     true,
			errContains: "failed to get path info",
		},

		// 目录遍历测试
		{
			name:       "directory non-recursive",
			targetPath: testDir,
			recursive:  false,
			want: []string{
				filepath.Join(testDir, "file1.txt"),
				filepath.Join(testDir, "file2.go"),
			},
			wantErr: false,
		},
		{
			name:       "directory recursive",
			targetPath: testDir,
			recursive:  true,
			want: []string{
				filepath.Join(testDir, "file1.txt"),
				filepath.Join(testDir, "file2.go"),
				filepath.Join(testDir, "subdir", "file3.txt"),
				filepath.Join(testDir, "subdir", "file4.go"),
				filepath.Join(testDir, "subdir", "nested", "file5.txt"),
			},
			wantErr: false,
		},

		// 通配符测试
		{
			name:       "glob pattern *.txt non-recursive",
			targetPath: filepath.Join(testDir, "*.txt"),
			recursive:  false,
			want: []string{
				filepath.Join(testDir, "file1.txt"),
			},
			wantErr: false,
		},
		{
			name:       "glob pattern *.go non-recursive",
			targetPath: filepath.Join(testDir, "*.go"),
			recursive:  false,
			want: []string{
				filepath.Join(testDir, "file2.go"),
			},
			wantErr: false,
		},
		{
			name:       "glob pattern subdir/* non-recursive",
			targetPath: filepath.Join(testDir, "subdir", "*"),
			recursive:  false,
			want: []string{
				filepath.Join(testDir, "subdir", "file3.txt"),
				filepath.Join(testDir, "subdir", "file4.go"),
				// 通配符匹配到 nested 目录，会收集该目录下的文件（但不递归子目录）
				filepath.Join(testDir, "subdir", "nested", "file5.txt"),
			},
			wantErr: false,
		},
		{
			name:       "glob pattern subdir/* recursive",
			targetPath: filepath.Join(testDir, "subdir", "*"),
			recursive:  true,
			want: []string{
				filepath.Join(testDir, "subdir", "file3.txt"),
				filepath.Join(testDir, "subdir", "file4.go"),
				filepath.Join(testDir, "subdir", "nested", "file5.txt"),
			},
			wantErr: false,
		},
		{
			name:        "invalid glob pattern",
			targetPath:  filepath.Join(testDir, "["),
			recursive:   false,
			want:        nil,
			wantErr:     true,
			errContains: "invalid path pattern",
		},
		{
			name:        "no matching files for glob",
			targetPath:  filepath.Join(testDir, "*.xyz"),
			recursive:   false,
			want:        nil,
			wantErr:     true,
			errContains: "no matching files found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Collect(tt.targetPath, tt.recursive)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Collect() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Collect() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Collect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 排序结果以便比较
			sort.Strings(got)
			sort.Strings(tt.want)

			if !equalStringSlices(got, tt.want) {
				t.Errorf("Collect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalkDir(t *testing.T) {
	// 创建测试目录结构
	testDir := setupTestDir(t)
	defer func() { _ = os.RemoveAll(testDir) }()

	tests := []struct {
		name      string
		dirPath   string
		recursive bool
		want      []string
		wantErr   bool
	}{
		{
			name:      "empty path should return error",
			dirPath:   "",
			recursive: false,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "non-existent directory",
			dirPath:   filepath.Join(testDir, "nonexistent"),
			recursive: false,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty directory non-recursive",
			dirPath:   filepath.Join(testDir, "empty"),
			recursive: false,
			want:      []string{},
			wantErr:   false,
		},
		{
			name:      "directory with files non-recursive",
			dirPath:   testDir,
			recursive: false,
			want: []string{
				filepath.Join(testDir, "file1.txt"),
				filepath.Join(testDir, "file2.go"),
			},
			wantErr: false,
		},
		{
			name:      "directory with files recursive",
			dirPath:   testDir,
			recursive: true,
			want: []string{
				filepath.Join(testDir, "file1.txt"),
				filepath.Join(testDir, "file2.go"),
				filepath.Join(testDir, "subdir", "file3.txt"),
				filepath.Join(testDir, "subdir", "file4.go"),
				filepath.Join(testDir, "subdir", "nested", "file5.txt"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := walkDir(tt.dirPath, tt.recursive)

			if tt.wantErr {
				if err == nil {
					t.Errorf("walkDir() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("walkDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 排序结果以便比较
			sort.Strings(got)
			sort.Strings(tt.want)

			if !equalStringSlices(got, tt.want) {
				t.Errorf("walkDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

// setupTestDir 创建测试目录结构
func setupTestDir(t *testing.T) string {
	testDir, err := os.MkdirTemp("", "fs_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// 创建文件和目录结构
	files := map[string]string{
		"file1.txt":               "content1",
		"file2.go":                "package main",
		"subdir/file3.txt":        "content3",
		"subdir/file4.go":         "package sub",
		"subdir/nested/file5.txt": "content5",
	}

	// 创建空目录
	emptyDir := filepath.Join(testDir, "empty")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("Failed to create empty dir: %v", err)
	}

	// 创建文件
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

	return testDir
}

// equalStringSlices 比较两个字符串切片是否相等
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// BenchmarkCollect 性能测试
func BenchmarkCollect(b *testing.B) {
	testDir := setupBenchmarkDir(b)
	defer func() { _ = os.RemoveAll(testDir) }()

	b.Run("SingleFile", func(b *testing.B) {
		filePath := filepath.Join(testDir, "file1.txt")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := Collect(filePath, false)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DirectoryNonRecursive", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := Collect(testDir, false)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DirectoryRecursive", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := Collect(testDir, true)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("GlobPattern", func(b *testing.B) {
		pattern := filepath.Join(testDir, "*.txt")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := Collect(pattern, false)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// setupBenchmarkDir 创建性能测试目录结构
func setupBenchmarkDir(b *testing.B) string {
	testDir, err := os.MkdirTemp("", "fs_bench")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	// 创建更多文件用于性能测试
	for i := 0; i < 100; i++ {
		fileName := filepath.Join(testDir, "file"+string(rune('0'+i%10))+".txt")
		if err := os.WriteFile(fileName, []byte("content"), 0644); err != nil {
			b.Fatalf("Failed to create file: %v", err)
		}
	}

	// 创建子目录
	subDir := filepath.Join(testDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		b.Fatalf("Failed to create subdir: %v", err)
	}

	for i := 0; i < 50; i++ {
		fileName := filepath.Join(subDir, "subfile"+string(rune('0'+i%10))+".go")
		if err := os.WriteFile(fileName, []byte("package main"), 0644); err != nil {
			b.Fatalf("Failed to create subfile: %v", err)
		}
	}

	return testDir
}
