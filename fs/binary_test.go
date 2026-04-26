package fs

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIsBinaryFile_Text 测试文本文件检测
func TestIsBinaryFile_Text(t *testing.T) {
	// 创建临时文本文件
	tmpDir := t.TempDir()
	textFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("Hello, World!\nThis is a text file.\n中文测试。")
	if err := os.WriteFile(textFile, content, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	file, err := os.Open(textFile)
	if err != nil {
		t.Fatalf("打开文件失败: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	isBinary, err := IsBinaryFile(file)
	if err != nil {
		t.Errorf("IsBinaryFile() 返回错误: %v", err)
	}
	if isBinary {
		t.Errorf("IsBinaryFile() = true, 期望 false (文本文件)")
	}
}

// TestIsBinaryFile_Binary 测试二进制文件检测
func TestIsBinaryFile_Binary(t *testing.T) {
	// 创建临时二进制文件（包含空字符）
	tmpDir := t.TempDir()
	binaryFile := filepath.Join(tmpDir, "test.bin")
	content := []byte{0x00, 0x01, 0x02, 0x03, 0x00, 0xFF, 0xFE}
	if err := os.WriteFile(binaryFile, content, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	file, err := os.Open(binaryFile)
	if err != nil {
		t.Fatalf("打开文件失败: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	isBinary, err := IsBinaryFile(file)
	if err != nil {
		t.Errorf("IsBinaryFile() 返回错误: %v", err)
	}
	if !isBinary {
		t.Errorf("IsBinaryFile() = false, 期望 true (二进制文件)")
	}
}

// TestIsBinaryFile_Empty 测试空文件检测
func TestIsBinaryFile_Empty(t *testing.T) {
	// 创建空文件
	tmpDir := t.TempDir()
	emptyFile := filepath.Join(tmpDir, "empty.txt")
	if err := os.WriteFile(emptyFile, []byte{}, 0644); err != nil {
		t.Fatalf("创建空文件失败: %v", err)
	}

	file, err := os.Open(emptyFile)
	if err != nil {
		t.Fatalf("打开文件失败: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	isBinary, err := IsBinaryFile(file)
	if err != nil {
		t.Errorf("IsBinaryFile() 返回错误: %v", err)
	}
	if isBinary {
		t.Errorf("IsBinaryFile() = true, 期望 false (空文件应视为文本)")
	}
}

// TestIsBinaryFile_NonRegular 测试非普通文件（如目录）
func TestIsBinaryFile_NonRegular(t *testing.T) {
	// 创建目录
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "testdir")
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("创建目录失败: %v", err)
	}

	// 尝试打开目录
	file, err := os.Open(testDir)
	if err != nil {
		t.Fatalf("打开目录失败: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	isBinary, err := IsBinaryFile(file)
	if err != nil {
		t.Errorf("IsBinaryFile() 返回错误: %v", err)
	}
	if isBinary {
		t.Errorf("IsBinaryFile() = true, 期望 false (目录应视为文本)")
	}
}

// TestIsBinaryFile_ResetPointer 测试文件指针重置
func TestIsBinaryFile_ResetPointer(t *testing.T) {
	// 创建测试文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("Hello, World!")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	file, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("打开文件失败: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	// 调用检测函数
	_, err = IsBinaryFile(file)
	if err != nil {
		t.Fatalf("IsBinaryFile() 返回错误: %v", err)
	}

	// 检测后应该能重新读取文件内容
	buf := make([]byte, 5)
	n, err := file.Read(buf)
	if err != nil {
		t.Errorf("检测后读取文件失败: %v", err)
	}
	if string(buf[:n]) != "Hello" {
		t.Errorf("文件指针未重置到开头，读取到: %s", string(buf[:n]))
	}
}

// TestIsBinaryFilePath_Text 测试路径版文本文件检测
func TestIsBinaryFilePath_Text(t *testing.T) {
	tmpDir := t.TempDir()
	textFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("Hello, World!")
	if err := os.WriteFile(textFile, content, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	isBinary, err := IsBinaryFilePath(textFile)
	if err != nil {
		t.Errorf("IsBinaryFilePath() 返回错误: %v", err)
	}
	if isBinary {
		t.Errorf("IsBinaryFilePath() = true, 期望 false")
	}
}

// TestIsBinaryFilePath_Binary 测试路径版二进制文件检测
func TestIsBinaryFilePath_Binary(t *testing.T) {
	tmpDir := t.TempDir()
	binaryFile := filepath.Join(tmpDir, "test.bin")
	content := []byte{0x00, 0x01, 0x02}
	if err := os.WriteFile(binaryFile, content, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	isBinary, err := IsBinaryFilePath(binaryFile)
	if err != nil {
		t.Errorf("IsBinaryFilePath() 返回错误: %v", err)
	}
	if !isBinary {
		t.Errorf("IsBinaryFilePath() = false, 期望 true")
	}
}

// TestIsBinaryFilePath_NotExist 测试不存在的文件
func TestIsBinaryFilePath_NotExist(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistFile := filepath.Join(tmpDir, "nonexistent.txt")

	_, err := IsBinaryFilePath(nonExistFile)
	if err == nil {
		t.Errorf("IsBinaryFilePath() 应该返回错误，但没有")
	}
}

// TestIsBinary 测试简洁版函数
func TestIsBinary(t *testing.T) {
	tmpDir := t.TempDir()

	// 测试文本文件
	textFile := filepath.Join(tmpDir, "text.txt")
	if err := os.WriteFile(textFile, []byte("Hello"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	file1, err := os.Open(textFile)
	if err != nil {
		t.Fatalf("打开文件失败: %v", err)
	}
	defer func() {
		_ = file1.Close()
	}()

	if IsBinary(file1) {
		t.Errorf("IsBinary(文本文件) = true, 期望 false")
	}

	// 测试二进制文件
	binaryFile := filepath.Join(tmpDir, "binary.bin")
	if err := os.WriteFile(binaryFile, []byte{0x00, 0x01}, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	file2, err := os.Open(binaryFile)
	if err != nil {
		t.Fatalf("打开文件失败: %v", err)
	}
	defer func() {
		_ = file2.Close()
	}()

	if !IsBinary(file2) {
		t.Errorf("IsBinary(二进制文件) = false, 期望 true")
	}
}

// TestIsBinaryPath 测试路径简洁版函数
func TestIsBinaryPath(t *testing.T) {
	tmpDir := t.TempDir()

	// 测试文本文件
	textFile := filepath.Join(tmpDir, "text.txt")
	if err := os.WriteFile(textFile, []byte("Hello"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	if IsBinaryPath(textFile) {
		t.Errorf("IsBinaryPath(文本文件) = true, 期望 false")
	}

	// 测试二进制文件
	binaryFile := filepath.Join(tmpDir, "binary.bin")
	if err := os.WriteFile(binaryFile, []byte{0x00, 0x01}, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	if !IsBinaryPath(binaryFile) {
		t.Errorf("IsBinaryPath(二进制文件) = false, 期望 true")
	}

	// 测试不存在的文件（应该返回 false，不 panic）
	nonExistFile := filepath.Join(tmpDir, "nonexistent.txt")
	if IsBinaryPath(nonExistFile) {
		t.Errorf("IsBinaryPath(不存在的文件) = true, 期望 false")
	}
}

// TestIsBinary_LargeFile 测试大文件（超过 8000 字节）
func TestIsBinary_LargeFile(t *testing.T) {
	tmpDir := t.TempDir()
	largeFile := filepath.Join(tmpDir, "large.bin")

	// 创建大文件，前 8000 字节是文本，后面是二进制
	content := make([]byte, 10000)
	for i := 0; i < 8000; i++ {
		content[i] = byte('A')
	}
	content[8001] = 0x00 // 在第 8001 字节放入空字符

	if err := os.WriteFile(largeFile, content, 0644); err != nil {
		t.Fatalf("创建大文件失败: %v", err)
	}

	// 应该检测为文本（因为只检查前 8000 字节）
	isBinary, err := IsBinaryFilePath(largeFile)
	if err != nil {
		t.Errorf("IsBinaryFilePath() 返回错误: %v", err)
	}
	if isBinary {
		t.Errorf("大文件（空字符在 8000 字节后）应该视为文本，但返回 true")
	}
}

// TestIsBinary_BinaryAtStart 测试开头就有空字符的文件
func TestIsBinary_BinaryAtStart(t *testing.T) {
	tmpDir := t.TempDir()
	binaryFile := filepath.Join(tmpDir, "binary.bin")

	// 文件开头就有空字符，后面是文本
	content := []byte{0x00, 0x01, byte('H'), byte('e'), byte('l'), byte('l'), byte('o')}
	if err := os.WriteFile(binaryFile, content, 0644); err != nil {
		t.Fatalf("创建文件失败: %v", err)
	}

	isBinary, err := IsBinaryFilePath(binaryFile)
	if err != nil {
		t.Errorf("IsBinaryFilePath() 返回错误: %v", err)
	}
	if !isBinary {
		t.Errorf("开头有空字符的文件应该视为二进制，但返回 false")
	}
}

// BenchmarkIsBinaryFile 基准测试
func BenchmarkIsBinaryFile(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("Hello, World! This is a test file for benchmarking."), 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		file, err := os.Open(testFile)
		if err != nil {
			b.Fatalf("打开文件失败: %v", err)
		}
		_, _ = IsBinaryFile(file)
		_ = file.Close()
	}
}

// BenchmarkIsBinaryFilePath 基准测试路径版
func BenchmarkIsBinaryFilePath(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("Hello, World! This is a test file for benchmarking."), 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = IsBinaryFilePath(testFile)
	}
}

// BenchmarkIsBinaryPath 基准测试简洁版
func BenchmarkIsBinaryPath(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("Hello, World! This is a test file for benchmarking."), 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsBinaryPath(testFile)
	}
}
