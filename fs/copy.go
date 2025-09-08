package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/go-kit/pool"
)

// Copy 通用复制函数，自动判断源路径类型并调用相应的复制函数
// 支持复制普通文件、目录、符号链接和特殊文件 (设备文件、命名管道等)
//
// 参数:
//   - src: 源路径 (支持文件、目录、符号链接、特殊文件)
//   - dst: 目标路径
//
// 返回:
//   - error: 复制失败时返回错误，如果目标已存在则返回错误
func Copy(src, dst string) error {
	return CopyEx(src, dst, false)
}

// CopyEx 通用复制函数 (可控制是否覆盖)，自动判断源路径类型并调用相应的复制函数
// 支持复制普通文件、目录、符号链接和特殊文件 (设备文件、命名管道等)
//
// 参数:
//   - src: 源路径 (支持文件、目录、符号链接、特殊文件)
//   - dst: 目标路径
//   - overwrite: 是否允许覆盖已存在的目标文件/目录
//
// 返回:
//   - error: 复制失败时返回错误
func CopyEx(src, dst string, overwrite bool) error {
	// 获取源路径信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source info '%s': %w", src, err)
	}

	// 根据源路径类型调用相应的复制函数
	if srcInfo.IsDir() {
		return copyDirInternal(src, dst, overwrite)
	} else {
		// 处理所有文件类型（普通文件、符号链接、特殊文件等）
		return copyFileRouter(src, dst, srcInfo.Mode(), overwrite)
	}
}

// validateCopyPaths 验证复制操作的源路径和目标路径
// 检查路径是否为空、获取绝对路径并验证源目标路径不相同
//
// 参数:
//   - src: 源路径
//   - dst: 目标路径
//   - checkSubdir: 是否检查目录复制到子目录的情况（仅对目录复制有效）
//
// 返回:
//   - error: 验证失败时返回错误
func validateCopyPaths(src, dst string, checkSubdir bool) error {
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

	// 检查是否尝试将目录复制到自己的子目录中（仅对目录复制）
	if checkSubdir {
		if strings.HasPrefix(dstAbs+string(filepath.Separator), srcAbs+string(filepath.Separator)) {
			return fmt.Errorf("cannot copy directory '%s' to its subdirectory '%s'", src, dst)
		}
	}

	return nil
}

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
	// 验证路径
	if err := validateCopyPaths(src, dst, false); err != nil {
		return err
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
		if out != nil {
			_ = out.Close()
			out = nil // 防止重复关闭
		}
		if !success {
			_ = os.Remove(tmp) // 清理临时文件
		}
	}()

	// 根据文件大小计算最佳缓冲区
	bufSize := pool.CalculateBufferSize(fi.Size())
	buf := pool.GetByteCap(bufSize)
	defer pool.PutByte(buf)

	// 使用缓冲区进行数据拷贝
	if _, err := io.CopyBuffer(out, in, buf); err != nil {
		return fmt.Errorf("failed to copy data from '%s' to '%s': %w", src, tmp, err)
	}

	// 强制刷盘，确保数据持久化
	if err := out.Sync(); err != nil {
		return fmt.Errorf("failed to sync temporary file '%s': %w", tmp, err)
	}

	// 在重命名前关闭文件句柄(Windows要求)
	if err := out.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file '%s': %w", tmp, err)
	}
	out = nil // 标记为已关闭

	// 如果允许覆盖且目标文件存在，先删除目标文件（确保跨平台兼容性）
	if overwrite {
		_ = os.Remove(dst) // 忽略错误，可能文件不存在
	}

	// 原子重命名
	if err := os.Rename(tmp, dst); err != nil {
		return fmt.Errorf("failed to rename temporary file '%s' to '%s': %w", tmp, dst, err)
	}

	success = true
	return nil
}

// copySymlink 复制符号链接
// 读取源符号链接的目标，然后在目标位置创建相同的符号链接
//
// 参数:
//   - src: 源符号链接路径
//   - dst: 目标符号链接路径
//   - overwrite: 是否允许覆盖已存在的目标
//
// 返回:
//   - error: 复制失败时返回错误
func copySymlink(src, dst string, overwrite bool) error {
	// 检查目标是否存在
	if !overwrite {
		if _, err := os.Lstat(dst); err == nil {
			return fmt.Errorf("destination symlink '%s' already exists", dst)
		}
	}

	// 读取符号链接的目标
	target, err := os.Readlink(src)
	if err != nil {
		return fmt.Errorf("failed to read symlink '%s': %w", src, err)
	}

	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return fmt.Errorf("failed to create destination directory '%s': %w", dstDir, err)
	}

	// 如果目标存在且允许覆盖，先删除
	if overwrite {
		_ = os.Remove(dst) // 忽略错误，可能不存在
	}

	// 创建符号链接
	if err := os.Symlink(target, dst); err != nil {
		return fmt.Errorf("failed to create symlink '%s' -> '%s': %w", dst, target, err)
	}

	return nil
}

// copySpecialFile 复制特殊文件（设备文件、命名管道、套接字等）
// 对于特殊文件，只创建一个具有相同权限模式的空文件
//
// 参数:
//   - src: 源特殊文件路径
//   - dst: 目标文件路径
//   - overwrite: 是否允许覆盖已存在的目标
//
// 返回:
//   - error: 复制失败时返回错误
func copySpecialFile(src, dst string, overwrite bool) error {
	// 检查目标是否存在
	if !overwrite {
		if _, err := os.Stat(dst); err == nil {
			return fmt.Errorf("destination file '%s' already exists", dst)
		}
	}

	// 获取源文件信息
	srcInfo, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("failed to get source file info '%s': %w", src, err)
	}

	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return fmt.Errorf("failed to create destination directory '%s': %w", dstDir, err)
	}

	// 如果目标存在且允许覆盖，先删除
	if overwrite {
		_ = os.Remove(dst) // 忽略错误，可能不存在
	}

	// 创建一个空文件，保持相同的权限模式
	file, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create special file '%s': %w", dst, err)
	}

	// 立即关闭文件
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close special file '%s': %w", dst, err)
	}

	return nil
}

// copyFileRouter 文件复制路由函数
// 根据文件类型调用相应的复制函数
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标文件路径
//   - fileType: 文件类型（从 os.DirEntry.Type() 获取）
//   - overwrite: 是否允许覆盖已存在的目标
//
// 返回:
//   - error: 复制失败时返回错误
func copyFileRouter(src, dst string, fileType os.FileMode, overwrite bool) error {
	switch {
	case fileType.IsRegular():
		// 普通文件
		return copyFileInternal(src, dst, overwrite)
	case fileType&os.ModeSymlink != 0:
		// 符号链接
		return copySymlink(src, dst, overwrite)
	default:
		// 其他特殊文件（设备文件、命名管道、套接字等）
		return copySpecialFile(src, dst, overwrite)
	}
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
	// 验证路径
	if err := validateCopyPaths(src, dst, true); err != nil {
		return err
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
	} else {
		// 如果允许覆盖且目标目录存在，先删除整个目录
		if err := os.RemoveAll(dst); err != nil {
			return fmt.Errorf("failed to remove existing destination directory '%s': %w", dst, err)
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

		// 处理所有文件类型（普通文件、符号链接、特殊文件等）
		return copyFileRouter(path, dstPath, entry.Type(), overwrite)
	})
}
