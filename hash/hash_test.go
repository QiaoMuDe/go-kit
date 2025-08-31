package hash

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsAlgorithmSupported(t *testing.T) {
	t.Run("Supported algorithms", func(t *testing.T) {
		supportedAlgos := []string{"md5", "sha1", "sha256", "sha512"}

		for _, algo := range supportedAlgos {
			if !IsAlgorithmSupported(algo) {
				t.Errorf("Algorithm %q should be supported", algo)
			}
		}
	})

	t.Run("Case insensitive", func(t *testing.T) {
		testCases := []string{"MD5", "Sha1", "SHA256", "sha512", "ShA1"}

		for _, algo := range testCases {
			if !IsAlgorithmSupported(algo) {
				t.Errorf("Algorithm %q should be supported (case insensitive)", algo)
			}
		}
	})

	t.Run("Unsupported algorithms", func(t *testing.T) {
		unsupportedAlgos := []string{"md4", "sha3", "blake2", "crc32", "unknown"}

		for _, algo := range unsupportedAlgos {
			if IsAlgorithmSupported(algo) {
				t.Errorf("Algorithm %q should not be supported", algo)
			}
		}
	})

	t.Run("Empty string", func(t *testing.T) {
		if IsAlgorithmSupported("") {
			t.Error("Empty string should not be supported")
		}
	})

	t.Run("Whitespace", func(t *testing.T) {
		if IsAlgorithmSupported(" ") || IsAlgorithmSupported("\t") || IsAlgorithmSupported("\n") {
			t.Error("Whitespace should not be supported")
		}
	})
}

func TestGetHashAlgorithm(t *testing.T) {
	t.Run("Valid algorithms", func(t *testing.T) {
		supportedAlgos := []string{"md5", "sha1", "sha256", "sha512"}

		for _, algo := range supportedAlgos {
			hashFunc, err := getHashAlgorithm(algo)
			if err != nil {
				t.Errorf("getHashAlgorithm(%q) should not return error: %v", algo, err)
			}
			if hashFunc == nil {
				t.Errorf("getHashAlgorithm(%q) should return a function", algo)
			}
		}
	})

	t.Run("Case insensitive", func(t *testing.T) {
		testCases := []string{"MD5", "Sha1", "SHA256", "sha512"}

		for _, algo := range testCases {
			hashFunc, err := getHashAlgorithm(algo)
			if err != nil {
				t.Errorf("getHashAlgorithm(%q) should not return error: %v", algo, err)
			}
			if hashFunc == nil {
				t.Errorf("getHashAlgorithm(%q) should return a function", algo)
			}
		}
	})

	t.Run("Invalid algorithms", func(t *testing.T) {
		invalidAlgos := []string{"md4", "sha3", "blake2", "unknown"}

		for _, algo := range invalidAlgos {
			hashFunc, err := getHashAlgorithm(algo)
			if err == nil {
				t.Errorf("getHashAlgorithm(%q) should return error", algo)
			}
			if hashFunc != nil {
				t.Errorf("getHashAlgorithm(%q) should return nil function", algo)
			}
		}
	})

	t.Run("Empty string", func(t *testing.T) {
		hashFunc, err := getHashAlgorithm("")
		if err == nil {
			t.Error("getHashAlgorithm(\"\") should return error")
		}
		if hashFunc != nil {
			t.Error("getHashAlgorithm(\"\") should return nil function")
		}
	})
}

func TestChecksum(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	t.Run("MD5 checksum", func(t *testing.T) {
		content := "hello world"
		filename := filepath.Join(tempDir, "test_md5.txt")
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		result, err := Checksum(filename, "md5")
		if err != nil {
			t.Fatalf("Checksum failed: %v", err)
		}

		// 验证MD5哈希长度和格式
		if len(result) != 32 {
			t.Errorf("MD5 checksum length should be 32, got %d", len(result))
		}

		// 验证只包含十六进制字符
		for _, char := range result {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
				t.Errorf("MD5 checksum contains invalid character: %c", char)
			}
		}

		// 验证一致性 - 同样内容应该产生相同哈希
		result2, err := Checksum(filename, "md5")
		if err != nil {
			t.Fatalf("Second checksum failed: %v", err)
		}
		if result != result2 {
			t.Errorf("MD5 checksum should be consistent: %q != %q", result, result2)
		}
	})

	t.Run("SHA1 checksum", func(t *testing.T) {
		content := "hello world"
		filename := filepath.Join(tempDir, "test_sha1.txt")
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		result, err := Checksum(filename, "sha1")
		if err != nil {
			t.Fatalf("Checksum failed: %v", err)
		}

		// 验证SHA1哈希长度和格式
		if len(result) != 40 {
			t.Errorf("SHA1 checksum length should be 40, got %d", len(result))
		}

		// 验证只包含十六进制字符
		for _, char := range result {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
				t.Errorf("SHA1 checksum contains invalid character: %c", char)
			}
		}
	})

	t.Run("SHA256 checksum", func(t *testing.T) {
		content := "hello world"
		filename := filepath.Join(tempDir, "test_sha256.txt")
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		result, err := Checksum(filename, "sha256")
		if err != nil {
			t.Fatalf("Checksum failed: %v", err)
		}

		// 验证SHA256哈希长度和格式
		if len(result) != 64 {
			t.Errorf("SHA256 checksum length should be 64, got %d", len(result))
		}

		// 验证只包含十六进制字符
		for _, char := range result {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
				t.Errorf("SHA256 checksum contains invalid character: %c", char)
			}
		}
	})

	t.Run("SHA512 checksum", func(t *testing.T) {
		content := "hello world"
		filename := filepath.Join(tempDir, "test_sha512.txt")
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		result, err := Checksum(filename, "sha512")
		if err != nil {
			t.Fatalf("Checksum failed: %v", err)
		}

		// 验证SHA512长度应该是128个十六进制字符
		if len(result) != 128 {
			t.Errorf("SHA512 checksum length should be 128, got %d", len(result))
		}

		// 验证只包含十六进制字符
		for _, char := range result {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
				t.Errorf("SHA512 checksum contains invalid character: %c", char)
			}
		}
	})

	t.Run("Small file", func(t *testing.T) {
		content := "a" // 单个字符，避免空文件问题
		filename := filepath.Join(tempDir, "small.txt")
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create small test file: %v", err)
		}

		result, err := Checksum(filename, "sha256")
		if err != nil {
			t.Fatalf("Checksum failed: %v", err)
		}

		// 验证SHA256哈希长度
		if len(result) != 64 {
			t.Errorf("SHA256 checksum should be 64 chars, got %d", len(result))
		}

		// 验证只包含十六进制字符
		for _, char := range result {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
				t.Errorf("SHA256 checksum contains invalid character: %c", char)
			}
		}
	})

	t.Run("Large file", func(t *testing.T) {
		// 创建大文件（1MB）
		content := make([]byte, 1024*1024)
		for i := range content {
			content[i] = byte(i % 256)
		}

		filename := filepath.Join(tempDir, "large.txt")
		err := os.WriteFile(filename, content, 0644)
		if err != nil {
			t.Fatalf("Failed to create large test file: %v", err)
		}

		result, err := Checksum(filename, "sha256")
		if err != nil {
			t.Fatalf("Checksum failed: %v", err)
		}

		// 验证哈希长度
		if len(result) != 64 {
			t.Errorf("SHA256 checksum length should be 64, got %d", len(result))
		}
	})

	t.Run("Binary file", func(t *testing.T) {
		// 创建二进制文件
		content := make([]byte, 1000)
		_, err := rand.Read(content)
		if err != nil {
			t.Fatalf("Failed to generate random content: %v", err)
		}

		filename := filepath.Join(tempDir, "binary.bin")
		err = os.WriteFile(filename, content, 0644)
		if err != nil {
			t.Fatalf("Failed to create binary test file: %v", err)
		}

		result, err := Checksum(filename, "md5")
		if err != nil {
			t.Fatalf("Checksum failed: %v", err)
		}

		// 验证MD5哈希长度和格式
		if len(result) != 32 {
			t.Errorf("MD5 checksum length should be 32, got %d", len(result))
		}

		for _, char := range result {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
				t.Errorf("MD5 checksum contains invalid character: %c", char)
			}
		}
	})

	t.Run("Case insensitive algorithm", func(t *testing.T) {
		content := "test content"
		filename := filepath.Join(tempDir, "case_test.txt")
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		algorithms := []string{"md5", "MD5", "Md5", "mD5"}
		var results []string

		for _, algo := range algorithms {
			result, err := Checksum(filename, algo)
			if err != nil {
				t.Fatalf("Checksum with %q failed: %v", algo, err)
			}
			results = append(results, result)
		}

		// 所有结果应该相同
		for i := 1; i < len(results); i++ {
			if results[i] != results[0] {
				t.Errorf("Case insensitive algorithm should produce same result: %q != %q", results[i], results[0])
			}
		}
	})

	t.Run("Non-existent file", func(t *testing.T) {
		filename := filepath.Join(tempDir, "nonexistent.txt")

		result, err := Checksum(filename, "sha256")
		if err == nil {
			t.Errorf("Checksum should fail for non-existent file, got result: %q", result)
		}
	})

	t.Run("Directory instead of file", func(t *testing.T) {
		dirname := filepath.Join(tempDir, "testdir")
		err := os.Mkdir(dirname, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		result, err := Checksum(dirname, "sha256")
		if err == nil {
			t.Errorf("Checksum should fail for directory, got result: %q", result)
		}
	})

	t.Run("Unsupported algorithm", func(t *testing.T) {
		content := "test content"
		filename := filepath.Join(tempDir, "unsupported_test.txt")
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		result, err := Checksum(filename, "unsupported")
		if err == nil {
			t.Errorf("Checksum should fail for unsupported algorithm, got result: %q", result)
		}
	})

	t.Run("Same content different files", func(t *testing.T) {
		content := "identical content"

		file1 := filepath.Join(tempDir, "file1.txt")
		file2 := filepath.Join(tempDir, "file2.txt")

		err := os.WriteFile(file1, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file1: %v", err)
		}

		err = os.WriteFile(file2, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file2: %v", err)
		}

		hash1, err := Checksum(file1, "sha256")
		if err != nil {
			t.Fatalf("Checksum file1 failed: %v", err)
		}

		hash2, err := Checksum(file2, "sha256")
		if err != nil {
			t.Fatalf("Checksum file2 failed: %v", err)
		}

		if hash1 != hash2 {
			t.Errorf("Files with same content should have same checksum: %q != %q", hash1, hash2)
		}
	})
}

func TestChecksumProgress(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	t.Run("Basic progress checksum", func(t *testing.T) {
		content := "hello world"
		filename := filepath.Join(tempDir, "progress_test.txt")
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		result, err := ChecksumProgress(filename, "sha256")
		if err != nil {
			t.Fatalf("ChecksumProgress failed: %v", err)
		}

		// 应该与普通Checksum结果相同
		expected, err := Checksum(filename, "sha256")
		if err != nil {
			t.Fatalf("Checksum failed: %v", err)
		}

		if result != expected {
			t.Errorf("ChecksumProgress result should match Checksum: %q != %q", result, expected)
		}
	})

	t.Run("Large file with progress", func(t *testing.T) {
		// 创建较大文件以测试进度条
		content := make([]byte, 100*1024) // 100KB
		for i := range content {
			content[i] = byte(i % 256)
		}

		filename := filepath.Join(tempDir, "large_progress.txt")
		err := os.WriteFile(filename, content, 0644)
		if err != nil {
			t.Fatalf("Failed to create large test file: %v", err)
		}

		result, err := ChecksumProgress(filename, "md5")
		if err != nil {
			t.Fatalf("ChecksumProgress failed: %v", err)
		}

		// 验证结果格式
		if len(result) != 32 {
			t.Errorf("MD5 checksum length should be 32, got %d", len(result))
		}
	})

	t.Run("Progress with non-existent file", func(t *testing.T) {
		filename := filepath.Join(tempDir, "nonexistent_progress.txt")

		result, err := ChecksumProgress(filename, "sha256")
		if err == nil {
			t.Errorf("ChecksumProgress should fail for non-existent file, got result: %q", result)
		}
	})
}

// 边界测试
func TestEdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("Very large file", func(t *testing.T) {
		// 创建大文件（5MB）
		content := make([]byte, 5*1024*1024)
		for i := range content {
			content[i] = byte(i % 256)
		}

		filename := filepath.Join(tempDir, "very_large.txt")
		err := os.WriteFile(filename, content, 0644)
		if err != nil {
			t.Fatalf("Failed to create very large test file: %v", err)
		}

		// 测试不同算法
		algorithms := []string{"md5", "sha1", "sha256", "sha512"}
		expectedLengths := []int{32, 40, 64, 128}

		for i, algo := range algorithms {
			result, err := Checksum(filename, algo)
			if err != nil {
				t.Fatalf("Checksum with %q failed: %v", algo, err)
			}

			if len(result) != expectedLengths[i] {
				t.Errorf("%s checksum length should be %d, got %d", strings.ToUpper(algo), expectedLengths[i], len(result))
			}
		}
	})

	t.Run("File with special characters in name", func(t *testing.T) {
		content := "special file content"
		// 注意：在Windows上某些字符可能不被支持
		filename := filepath.Join(tempDir, "special_file_测试.txt")
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create special file: %v", err)
		}

		result, err := Checksum(filename, "sha256")
		if err != nil {
			t.Fatalf("Checksum failed for special filename: %v", err)
		}

		if len(result) != 64 {
			t.Errorf("SHA256 checksum length should be 64, got %d", len(result))
		}
	})

	t.Run("Consistency across multiple calls", func(t *testing.T) {
		content := "consistency test content"
		filename := filepath.Join(tempDir, "consistency.txt")
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// 多次调用应该返回相同结果
		var results []string
		for i := 0; i < 10; i++ {
			result, err := Checksum(filename, "sha256")
			if err != nil {
				t.Fatalf("Checksum failed at iteration %d: %v", i, err)
			}
			results = append(results, result)
		}

		// 验证所有结果相同
		for i := 1; i < len(results); i++ {
			if results[i] != results[0] {
				t.Errorf("Checksum inconsistent at iteration %d: %q != %q", i, results[i], results[0])
			}
		}
	})

	t.Run("All supported algorithms on same file", func(t *testing.T) {
		content := "multi-algorithm test"
		filename := filepath.Join(tempDir, "multi_algo.txt")
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		algorithms := []string{"md5", "sha1", "sha256", "sha512"}
		expectedLengths := []int{32, 40, 64, 128}
		results := make(map[string]string)

		for i, algo := range algorithms {
			result, err := Checksum(filename, algo)
			if err != nil {
				t.Fatalf("Checksum with %q failed: %v", algo, err)
			}

			if len(result) != expectedLengths[i] {
				t.Errorf("%s checksum length should be %d, got %d", strings.ToUpper(algo), expectedLengths[i], len(result))
			}

			results[algo] = result
		}

		// 验证不同算法产生不同结果
		uniqueResults := make(map[string]bool)
		for _, result := range results {
			uniqueResults[result] = true
		}

		if len(uniqueResults) != len(algorithms) {
			t.Error("Different algorithms should produce different checksums")
		}
	})
}

// 性能基准测试
func BenchmarkChecksum(b *testing.B) {
	tempDir := b.TempDir()
	content := strings.Repeat("benchmark data ", 1000)
	filename := filepath.Join(tempDir, "benchmark.txt")
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to create benchmark file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Checksum(filename, "sha256")
		if err != nil {
			b.Fatalf("Checksum failed: %v", err)
		}
	}
}

func BenchmarkChecksumMD5(b *testing.B) {
	tempDir := b.TempDir()
	content := strings.Repeat("benchmark data ", 1000)
	filename := filepath.Join(tempDir, "benchmark_md5.txt")
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to create benchmark file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Checksum(filename, "md5")
		if err != nil {
			b.Fatalf("Checksum failed: %v", err)
		}
	}
}

func BenchmarkChecksumSHA512(b *testing.B) {
	tempDir := b.TempDir()
	content := strings.Repeat("benchmark data ", 1000)
	filename := filepath.Join(tempDir, "benchmark_sha512.txt")
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to create benchmark file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Checksum(filename, "sha512")
		if err != nil {
			b.Fatalf("Checksum failed: %v", err)
		}
	}
}

// 并发测试
func TestConcurrentChecksum(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("Concurrent same file", func(t *testing.T) {
		content := "concurrent test content"
		filename := filepath.Join(tempDir, "concurrent.txt")
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		const numGoroutines = 10
		results := make(chan string, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				result, err := Checksum(filename, "sha256")
				if err != nil {
					errors <- err
					return
				}
				results <- result
			}()
		}

		// 收集结果
		var checksums []string
		for i := 0; i < numGoroutines; i++ {
			select {
			case result := <-results:
				checksums = append(checksums, result)
			case err := <-errors:
				t.Fatalf("Concurrent checksum failed: %v", err)
			}
		}

		// 验证所有结果相同
		for i := 1; i < len(checksums); i++ {
			if checksums[i] != checksums[0] {
				t.Errorf("Concurrent checksum inconsistent: %q != %q", checksums[i], checksums[0])
			}
		}
	})

	t.Run("Concurrent different files", func(t *testing.T) {
		const numFiles = 5
		filenames := make([]string, numFiles)
		expectedResults := make([]string, numFiles)

		// 创建不同的测试文件
		for i := 0; i < numFiles; i++ {
			content := strings.Repeat("file content ", i+1)
			filename := filepath.Join(tempDir, "concurrent_"+string(rune('a'+i))+".txt")
			err := os.WriteFile(filename, []byte(content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file %d: %v", i, err)
			}
			filenames[i] = filename

			// 预先计算期望结果
			expected, err := Checksum(filename, "md5")
			if err != nil {
				t.Fatalf("Failed to compute expected result for file %d: %v", i, err)
			}
			expectedResults[i] = expected
		}

		// 并发计算所有文件的校验和
		results := make(chan string, numFiles)
		errors := make(chan error, numFiles)

		for i := 0; i < numFiles; i++ {
			go func(idx int) {
				result, err := Checksum(filenames[idx], "md5")
				if err != nil {
					errors <- err
					return
				}
				results <- result
			}(i)
		}

		// 收集结果
		var actualResults []string
		for i := 0; i < numFiles; i++ {
			select {
			case result := <-results:
				actualResults = append(actualResults, result)
			case err := <-errors:
				t.Fatalf("Concurrent checksum failed: %v", err)
			}
		}

		// 验证结果数量正确
		if len(actualResults) != numFiles {
			t.Errorf("Expected %d results, got %d", numFiles, len(actualResults))
		}
	})
}
