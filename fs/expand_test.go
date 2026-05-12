package fs

import (
	"os"
	"path/filepath"
	"testing"
)

// setupExpandTestDir 创建测试目录结构
func setupExpandTestDir(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir()

	// 创建测试文件和目录
	files := []string{
		"main.go",
		"utils.go",
		"README.md",
		"src/app.go",
		"src/helper.go",
		"config.yaml",
	}

	for _, file := range files {
		path := filepath.Join(tempDir, file)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("创建目录失败: %v", err)
		}
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("创建文件失败: %v", err)
		}
	}

	// 创建空目录
	if err := os.MkdirAll(filepath.Join(tempDir, "emptydir"), 0755); err != nil {
		t.Fatalf("创建空目录失败: %v", err)
	}

	return tempDir
}

func TestExpand(t *testing.T) {
	testDir := setupExpandTestDir(t)

	tests := []struct {
		name     string
		patterns []string
		wantMin  int // 最小匹配数量
		wantMax  int // 最大匹配数量
		wantErr  bool
	}{
		{
			name:     "空切片",
			patterns: []string{},
			wantMin:  0,
			wantMax:  0,
			wantErr:  false,
		},
		{
			name:     "单个通配符",
			patterns: []string{filepath.Join(testDir, "*.go")},
			wantMin:  2, // main.go, utils.go
			wantMax:  2,
			wantErr:  false,
		},
		{
			name:     "多个通配符",
			patterns: []string{filepath.Join(testDir, "*.go"), filepath.Join(testDir, "*.md")},
			wantMin:  3, // main.go, utils.go, README.md
			wantMax:  3,
			wantErr:  false,
		},
		{
			name:     "具体路径",
			patterns: []string{filepath.Join(testDir, "config.yaml")},
			wantMin:  1,
			wantMax:  1,
			wantErr:  false,
		},
		{
			name:     "无匹配保留原模式",
			patterns: []string{filepath.Join(testDir, "*.notexist")},
			wantMin:  1,
			wantMax:  1,
			wantErr:  false,
		},
		{
			name:     "子目录通配符",
			patterns: []string{filepath.Join(testDir, "src", "*.go")},
			wantMin:  2, // app.go, helper.go
			wantMax:  2,
			wantErr:  false,
		},
		{
			name: "混合模式",
			patterns: []string{
				filepath.Join(testDir, "*.go"),
				filepath.Join(testDir, "src", "*.go"),
			},
			wantMin: 4, // main.go, utils.go, app.go, helper.go
			wantMax: 4,
			wantErr: false,
		},
		{
			name: "去重测试",
			patterns: []string{
				filepath.Join(testDir, "main.go"),
				filepath.Join(testDir, "main.go"),
			},
			wantMin: 1,
			wantMax: 1,
			wantErr: false,
		},
		{
			name: "路径清洗",
			patterns: []string{
				filepath.Join(testDir, "./main.go"),
			},
			wantMin: 1,
			wantMax: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Expand(tt.patterns)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) < tt.wantMin || len(got) > tt.wantMax {
				t.Errorf("Expand() returned %d items, want between %d and %d", len(got), tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestExpandFiles(t *testing.T) {
	testDir := setupExpandTestDir(t)

	tests := []struct {
		name     string
		patterns []string
		wantMin  int
		wantMax  int
		wantErr  bool
	}{
		{
			name:     "过滤目录",
			patterns: []string{filepath.Join(testDir, "*")},
			wantMin:  4, // main.go, utils.go, README.md, config.yaml（不包含 emptydir）
			wantMax:  4,
			wantErr:  false,
		},
		{
			name:     "只返回文件",
			patterns: []string{filepath.Join(testDir, "src", "*")},
			wantMin:  2, // app.go, helper.go（不包含 src 目录本身）
			wantMax:  2,
			wantErr:  false,
		},
		{
			name:     "空目录通配符",
			patterns: []string{filepath.Join(testDir, "emptydir", "*")},
			wantMin:  1, // 无匹配时保留原模式
			wantMax:  1,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpandFiles(tt.patterns)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) < tt.wantMin || len(got) > tt.wantMax {
				t.Errorf("ExpandFiles() returned %d items, want between %d and %d", len(got), tt.wantMin, tt.wantMax)
			}
			// 验证都是文件
			for _, path := range got {
				if isDir(path) {
					t.Errorf("ExpandFiles() returned directory: %s", path)
				}
			}
		})
	}
}

func TestExpandPattern(t *testing.T) {
	testDir := setupExpandTestDir(t)

	tests := []struct {
		name    string
		pattern string
		wantMin int
		wantMax int
		wantErr bool
	}{
		{
			name:    "空模式",
			pattern: "",
			wantMin: 0,
			wantMax: 0,
			wantErr: true,
		},
		{
			name:    "通配符",
			pattern: filepath.Join(testDir, "*.go"),
			wantMin: 2,
			wantMax: 2,
			wantErr: false,
		},
		{
			name:    "具体文件",
			pattern: filepath.Join(testDir, "main.go"),
			wantMin: 1,
			wantMax: 1,
			wantErr: false,
		},
		{
			name:    "无匹配",
			pattern: filepath.Join(testDir, "*.notexist"),
			wantMin: 1,
			wantMax: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpandPattern(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) < tt.wantMin || len(got) > tt.wantMax {
				t.Errorf("ExpandPattern() returned %d items, want between %d and %d", len(got), tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestExpandPattern_InvalidPattern(t *testing.T) {
	// 测试无效的模式（包含非法通配符语法）
	// 注意：filepath.Glob 对大多数无效模式不会返回错误
	// 这里测试空模式
	_, err := ExpandPattern("")
	if err == nil {
		t.Error("ExpandPattern(\"\") should return error")
	}
}

func TestExpand_DuplicatePatterns(t *testing.T) {
	testDir := setupExpandTestDir(t)

	// 测试重复模式的去重
	patterns := []string{
		filepath.Join(testDir, "*.go"),
		filepath.Join(testDir, "*.go"),
	}

	got, err := Expand(patterns)
	if err != nil {
		t.Fatalf("Expand() error = %v", err)
	}

	// 即使有重复模式，结果也不应该重复
	seen := make(map[string]bool)
	for _, path := range got {
		if seen[path] {
			t.Errorf("Duplicate path in result: %s", path)
		}
		seen[path] = true
	}
}

func TestExpand_CleanPath(t *testing.T) {
	testDir := setupExpandTestDir(t)

	// 测试路径清洗
	patterns := []string{
		filepath.Join(testDir, "./main.go"),
		filepath.Join(testDir, ".", "main.go"),
	}

	got, err := Expand(patterns)
	if err != nil {
		t.Fatalf("Expand() error = %v", err)
	}

	// 清洗后应该只有一个结果
	if len(got) != 1 {
		t.Errorf("Expand() with duplicate paths after cleaning returned %d items, want 1", len(got))
	}
}
