package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	if src == dst {
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
	tmp := dst + ".tmp"
	out, err := os.OpenFile(tmp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fi.Mode())
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

	// 数据拷贝
	if _, err := io.Copy(out, in); err != nil {
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
