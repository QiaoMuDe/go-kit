package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsHidden(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"normal file", "file.txt", false},
		{"hidden file", ".hidden", true},
		{"hidden dir", ".git", true},
		{"current dir", ".", true},
		{"parent dir", "..", true},
		{"normal dir", "folder", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsHidden(tt.path)
			if result != tt.expected {
				t.Errorf("IsHidden(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsReadOnly(t *testing.T) {
	// 创建临时文件进行测试
	tmpDir := t.TempDir()

	// 创建普通文件
	normalFile := filepath.Join(tmpDir, "normal.txt")
	if err := os.WriteFile(normalFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// 创建只读文件
	readOnlyFile := filepath.Join(tmpDir, "readonly.txt")
	if err := os.WriteFile(readOnlyFile, []byte("test"), 0444); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"normal file", normalFile, false},
		{"readonly file", readOnlyFile, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsReadOnly(tt.path)
			if result != tt.expected {
				t.Errorf("IsReadOnly(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}
