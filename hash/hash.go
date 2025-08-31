package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// 字节单位定义
const (
	Byte = 1 << (10 * iota) // 1 字节
	KB                      // 千字节 (1024 B)
	MB                      // 兆字节 (1024 KB)
	GB                      // 吉字节 (1024 MB)
	TB                      // 太字节 (1024 GB)
)

// 支持的哈希算法列表
var supportedAlgorithms = map[string]func() hash.Hash{
	"md5":    md5.New,
	"sha1":   sha1.New,
	"sha256": sha256.New,
	"sha512": sha512.New,
}

// IsAlgorithmSupported 检查给定的哈希算法名称是否受支持。
// 匹配时会忽略算法名称的大小写。
//
// 参数:
//   - algorithm: 要检查的哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）。
//
// 返回:
//   - bool: 如果算法受支持则返回 true，否则返回 false。
func IsAlgorithmSupported(algorithm string) bool {
	// 如果算法名称为空，则返回 false
	if algorithm == "" {
		return false
	}

	_, ok := supportedAlgorithms[strings.ToLower(algorithm)]
	return ok
}

// getHashAlgorithm 根据算法名称获取对应的哈希函数构造器。
// 匹配时会忽略算法名称的大小写。
//
// 参数:
//   - algorithm: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）。
//
// 返回:
//   - func() hash.Hash: 对应的哈希函数构造器。
//   - error: 如果不支持该算法，则返回错误。
func getHashAlgorithm(algorithm string) (func() hash.Hash, error) {
	// 如果算法名称为空，则返回错误
	if algorithm == "" {
		return nil, fmt.Errorf("hash algorithm name cannot be empty")
	}

	algoFunc, ok := supportedAlgorithms[strings.ToLower(algorithm)]
	if !ok {
		return nil, fmt.Errorf("unsupported hash algorithm: %s", algorithm)
	}
	return algoFunc, nil
}

// Checksum 计算文件哈希值
//
// 参数:
//   - filePath: 文件路径
//   - algorithm: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）
//
// 返回:
//   - string: 文件的十六进制哈希值
//   - error: 错误信息，如果计算失败
//
// 注意:
//   - 根据文件大小动态分配缓冲区以提高性能
//   - 支持任何实现hash.Hash接口的哈希算法
//   - 使用io.CopyBuffer进行高效的文件读取和哈希计算
func Checksum(filePath string, algorithm string) (string, error) {
	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("file does not exist or is inaccessible: %v", err)
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("close file failed: %v\n", err)
		}
	}()

	// 获取哈希函数构造器
	hashFunc, err := getHashAlgorithm(algorithm)
	if err != nil {
		return "", err
	}
	// 创建哈希对象
	h := hashFunc()

	// 根据文件大小动态分配缓冲区
	fileSize := fileInfo.Size()
	bufferSize := CalculateBufferSize(fileSize)
	buffer := make([]byte, bufferSize)

	// 使用 io.CopyBuffer 进行高效复制并计算哈希
	if _, err := io.CopyBuffer(h, file, buffer); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// 返回哈希值的十六进制表示
	return hex.EncodeToString(h.Sum(nil)), nil
}

// ChecksumProgress 计算文件哈希值(带进度条)
//
// 参数:
//   - filePath: 文件路径
//   - algorithm: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）
//
// 返回:
//   - string: 文件的十六进制哈希值
//   - error: 错误信息，如果计算失败
//
// 注意:
//   - 根据文件大小动态分配缓冲区以提高性能
//   - 支持任何实现hash.Hash接口的哈希算法
//   - 使用io.CopyBuffer进行高效的文件读取和哈希计算
func ChecksumProgress(filePath string, algorithm string) (string, error) {
	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("file does not exist or is inaccessible: %v", err)
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("close file failed: %v\n", err)
		}
	}()

	// 获取哈希函数构造器
	hashFunc, err := getHashAlgorithm(algorithm)
	if err != nil {
		return "", err
	}
	// 创建哈希对象
	h := hashFunc()

	// 根据文件大小动态分配缓冲区
	fileSize := fileInfo.Size()
	bufferSize := CalculateBufferSize(fileSize)
	buffer := make([]byte, bufferSize)

	// 创建进度条
	bar := progressbar.NewOptions64(
		fileSize,                          // 总进度
		progressbar.OptionClearOnFinish(), // 完成后清除进度条
		progressbar.OptionSetDescription(file.Name()+" 计算中"), // 设置进度条描述
	)
	defer func() {
		// 完成进度条
		if err := bar.Finish(); err != nil {
			fmt.Printf("finish progress bar failed: %v\n", err)
		}

		// 关闭进度条
		if err := bar.Close(); err != nil {
			fmt.Printf("close progress bar failed: %v\n", err)
		}
	}()

	// 创建多路写入器
	multiWriter := io.MultiWriter(h, bar)

	// 使用 io.CopyBuffer 进行高效复制并计算哈希
	if _, err := io.CopyBuffer(multiWriter, file, buffer); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// 获取哈希值的十六进制表示
	hashStr := hex.EncodeToString(h.Sum(nil))

	// 返回哈希值的十六进制表示
	return hashStr, nil
}

// CalculateBufferSize 根据文件大小动态计算最佳缓冲区大小。
// 采用分层策略，平衡内存使用和I/O性能。
//
// 参数:
//   - fileSize: 文件大小（字节）
//
// 返回:
//   - int: 计算出的最佳缓冲区大小（字节）
//
// 缓冲区分配策略:
//   - ≤ 4KB: 使用文件实际大小，避免内存浪费
//   - 4KB - 32KB: 使用 8KB 缓冲区
//   - 32KB - 128KB: 使用 32KB 缓冲区
//   - 128KB - 512KB: 使用 64KB 缓冲区
//   - 512KB - 1MB: 使用 128KB 缓冲区
//   - 1MB - 4MB: 使用 256KB 缓冲区
//   - 4MB - 16MB: 使用 512KB 缓冲区
//   - 16MB - 64MB: 使用 1MB 缓冲区
//   - > 64MB: 使用 2MB 缓冲区
//
// 设计原则:
//   - 极小文件: 最小化内存占用
//   - 小文件: 适度缓冲，节省内存
//   - 大文件: 增大缓冲区，提升I/O吞吐量
//   - 超大文件: 限制最大缓冲区，避免过度内存消耗
func CalculateBufferSize(fileSize int64) int {
	switch {
	case fileSize <= 4*KB: // 极小文件直接使用文件大小作为缓冲区
		return int(fileSize)
	case fileSize < 32*KB: // 小于 32KB 的文件使用 8KB 缓冲区
		return int(8 * KB)
	case fileSize < 128*KB: // 32KB-128KB 使用 32KB 缓冲区
		return int(32 * KB)
	case fileSize < 512*KB: // 128KB-512KB 使用 64KB 缓冲区
		return int(64 * KB)
	case fileSize < 1*MB: // 512KB-1MB 使用 128KB 缓冲区
		return int(128 * KB)
	case fileSize < 4*MB: // 1MB-4MB 使用 256KB 缓冲区
		return int(256 * KB)
	case fileSize < 16*MB: // 4MB-16MB 使用 512KB 缓冲区
		return int(512 * KB)
	case fileSize < 64*MB: // 16MB-64MB 使用 1MB 缓冲区
		return int(1 * MB)
	default: // 大于 64MB 的文件使用 2MB 缓冲区
		return int(2 * MB)
	}
}
