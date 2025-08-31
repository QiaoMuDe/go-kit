package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/go-kit/pool"
)

// copyFileInternal 内部复制文件逻辑
// 用于安全地复制文件，保持原文件的权限信息，失败时自动清理
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标文件路径
//   - overwrite: 是否允许覆盖已存在的目标文件
//
// 返回:
//   - error: 复制失败时返回错误
func copyFileInternal(src, dst string, overwrite bool) error {
	// 参数验证
	if src == "" || dst == "" {
		return fmt.Errorf("source and destination paths cannot be empty")
	}

	// 清理路径并检查是否相同
	srcAbs, err := filepath.Abs(filepath.Clean(src))
	if err != nil {
		return fmt.Errorf("failed to get absolute path for source '%s': %w", src, err)
	}
	dstAbs, err := filepath.Abs(filepath.Clean(dst))
	if err != nil {
		return fmt.Errorf("failed to get absolute path for destination '%s': %w", dst, err)
	}
	if srcAbs == dstAbs {
		return fmt.Errorf("source and destination paths cannot be the same")
	}

	// 检查目标文件是否存在
	if !overwrite {
		if _, err := os.Stat(dst); err == nil {
			return fmt.Errorf("destination file '%s' already exists", dst)
		}
	}

	// 打开源文件
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file '%s': %w", src, err)
	}
	defer func() { _ = in.Close() }()

	// 获取源文件元数据
	fi, err := in.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info '%s': %w", src, err)
	}

	// 检查是否为普通文件
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("source '%s' is not a regular file", src)
	}

	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return fmt.Errorf("failed to create destination directory '%s': %w", dstDir, err)
	}

	// 创建临时文件（与目标同目录，保证 rename 原子性）
	// 使用更安全的临时文件名，避免冲突
	tmp := dst + ".tmp." + fmt.Sprintf("%d", os.Getpid())
	out, err := os.OpenFile(tmp, os.O_RDWR|os.O_CREATE|os.O_EXCL, fi.Mode())
	if err != nil {
		return fmt.Errorf("failed to create temporary file '%s': %w", tmp, err)
	}

	// 统一清理资源
	success := false
	defer func() {
		_ = out.Close()
		if !success {
			_ = os.Remove(tmp) // 清理临时文件
		}
	}()

	// 根据文件大小计算最佳缓冲区
	bufSize := pool.CalculateBufferSize(fi.Size())
	buf := pool.GetByte(bufSize)
	defer pool.PutByte(buf)

	// 使用缓冲区进行数据拷贝
	if _, err := io.CopyBuffer(out, in, buf); err != nil {
		return fmt.Errorf("failed to copy data from '%s' to '%s': %w", src, tmp, err)
	}

	// 强制刷盘，确保数据持久化
	if err := out.Sync(); err != nil {
		return fmt.Errorf("failed to sync temporary file '%s': %w", tmp, err)
	}

	// 原子重命名
	if err := os.Rename(tmp, dst); err != nil {
		return fmt.Errorf("failed to rename temporary file '%s' to '%s': %w", tmp, dst, err)
	}

	success = true
	return nil
}

// CopyFile 复制文件并继承权限（默认覆盖已存在的目标文件）
// 用于安全地复制文件，保持原文件的权限信息，失败时自动清理
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标文件路径
//
// 返回:
//   - error: 复制失败时返回错误
func CopyFile(src, dst string) error {
	return copyFileInternal(src, dst, true)
}

// CopyFileWithOverwrite 复制文件并继承权限（可控制是否覆盖）
// 用于安全地复制文件，保持原文件的权限信息，失败时自动清理
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标文件路径
//   - overwrite: 是否允许覆盖已存在的目标文件，false时如果目标文件存在则返回错误
//
// 返回:
//   - error: 复制失败时返回错误
func CopyFileWithOverwrite(src, dst string, overwrite bool) error {
	return copyFileInternal(src, dst, overwrite)
}

// copyDirInternal 内部复制目录逻辑
// 用于递归复制整个目录，保持文件权限和目录结构
//
// 参数:
//   - src: 源目录路径
//   - dst: 目标目录路径
//   - overwrite: 是否允许覆盖已存在的目标文件
//
// 返回:
//   - error: 复制失败时返回错误
func copyDirInternal(src, dst string, overwrite bool) error {
	// 参数验证
	if src == "" || dst == "" {
		return fmt.Errorf("source and destination paths cannot be empty")
	}

	// 清理路径并检查是否相同
	srcAbs, err := filepath.Abs(filepath.Clean(src))
	if err != nil {
		return fmt.Errorf("failed to get absolute path for source '%s': %w", src, err)
	}
	dstAbs, err := filepath.Abs(filepath.Clean(dst))
	if err != nil {
		return fmt.Errorf("failed to get absolute path for destination '%s': %w", dst, err)
	}
	if srcAbs == dstAbs {
		return fmt.Errorf("source and destination paths cannot be the same")
	}

	// 检查是否尝试将目录复制到自己的子目录中
	if strings.HasPrefix(dstAbs+string(filepath.Separator), srcAbs+string(filepath.Separator)) {
		return fmt.Errorf("cannot copy directory '%s' to its subdirectory '%s'", src, dst)
	}

	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source directory info '%s': %w", src, err)
	}

	// 检查源路径是否为目录
	if !srcInfo.IsDir() {
		return fmt.Errorf("source '%s' is not a directory", src)
	}

	// 检查目标目录是否存在
	if !overwrite {
		if _, err := os.Stat(dst); err == nil {
			return fmt.Errorf("destination directory '%s' already exists", dst)
		}
	}

	// 创建目标目录
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination directory '%s': %w", dst, err)
	}

	// 遍历源目录
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path '%s': %w", path, err)
		}

		// 计算相对路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for '%s': %w", path, err)
		}

		// 跳过根目录本身
		if relPath == "." {
			return nil
		}

		// 构建目标路径
		dstPath := filepath.Join(dst, relPath)

		// 处理目录
		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("failed to get directory info '%s': %w", path, err)
			}
			if err := os.MkdirAll(dstPath, info.Mode()); err != nil {
				return fmt.Errorf("failed to create directory '%s': %w", dstPath, err)
			}
			return nil
		}

		// 处理普通文件
		if entry.Type().IsRegular() {
			return copyFileInternal(path, dstPath, overwrite)
		}

		// 跳过其他类型的文件（符号链接、设备文件等）
		return nil
	})
}

// CopyDir 复制目录及其所有内容（默认覆盖已存在的文件）
// 用于递归复制整个目录，保持文件权限和目录结构
//
// 参数:
//   - src: 源目录路径
//   - dst: 目标目录路径
//
// 返回:
//   - error: 复制失败时返回错误
func CopyDir(src, dst string) error {
	return copyDirInternal(src, dst, true)
}

// CopyDirWithOverwrite 复制目录及其所有内容（可控制是否覆盖）
// 用于递归复制整个目录，保持文件权限和目录结构
//
// 参数:
//   - src: 源目录路径
//   - dst: 目标目录路径
//   - overwrite: 是否允许覆盖已存在的文件，false时如果目标目录或文件存在则返回错误
//
// 返回:
//   - error: 复制失败时返回错误
func CopyDirWithOverwrite(src, dst string, overwrite bool) error {
	return copyDirInternal(src, dst, overwrite)
}
