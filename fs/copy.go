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
func CopyEx(src, dst string, overwrite bool) (err error) {
	// 捕获 panic 并转换为错误
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("copy operation panicked: %v", r)
		}
	}()

	// 获取源路径信息（使用 Lstat 避免跟随符号链接）
	srcInfo, localErr := os.Lstat(src)
	if localErr != nil {
		err = fmt.Errorf("failed to get source info '%s': %w", src, localErr)
		return
	}

	// 根据源路径类型调用相应的复制函数
	if srcInfo.IsDir() {
		err = copyDir(src, dst, overwrite)
	} else {
		// 处理所有文件类型（普通文件、符号链接、特殊文件等）
		err = copyFileRouter(src, dst, srcInfo.Mode(), overwrite)
	}
	return
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

// handleBackupAndRestore 统一的备份恢复处理函数
// 用于在覆盖操作前创建备份，失败时恢复备份
//
// 参数:
//   - dst: 目标路径
//   - overwrite: 是否允许覆盖
//
// 返回:
//   - backupPath: 备份文件路径（如果创建了备份）
//   - error: 处理失败时返回错误
func handleBackupAndRestore(dst string, overwrite bool) (string, error) {
	// 检查目标是否存在
	if _, err := os.Lstat(dst); err != nil {
		// 目标不存在，无需备份
		return "", nil
	}

	// 目标存在
	if !overwrite {
		return "", fmt.Errorf("destination '%s' already exists", dst)
	}

	// 允许覆盖，创建备份
	backupPath := dst + ".backup." + fmt.Sprintf("%d", os.Getpid())
	if err := os.Rename(dst, backupPath); err != nil {
		return "", fmt.Errorf("failed to backup existing '%s': %w", dst, err)
	}

	return backupPath, nil
}

// restoreBackup 恢复备份文件
// 在操作失败时调用，尽力恢复原始文件
//
// 参数:
//   - dst: 目标路径
//   - backupPath: 备份文件路径
func restoreBackup(dst, backupPath string) {
	if backupPath != "" {
		_ = os.Rename(backupPath, dst) // 忽略恢复错误，尽力而为
	}
}

// cleanupBackup 清理备份文件/目录
// 在操作成功时调用，删除不再需要的备份
//
// 参数:
//   - backupPath: 备份文件或目录路径
func cleanupBackup(backupPath string) {
	if backupPath != "" {
		_ = os.RemoveAll(backupPath) // 忽略删除错误，支持文件和目录
	}
}

// copyFile 内部复制文件逻辑
// 用于安全地复制文件，保持原文件的权限信息，失败时自动清理
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标文件路径
//   - overwrite: 是否允许覆盖已存在的目标文件
//
// 返回:
//   - error: 复制失败时返回错误
func copyFile(src, dst string, overwrite bool) error {
	// 验证路径
	if err := validateCopyPaths(src, dst, false); err != nil {
		return err
	}

	// 安全覆盖机制：处理已存在的目标文件
	backupPath, err := handleBackupAndRestore(dst, overwrite)
	if err != nil {
		return err
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

	// 根据文件大小选择处理方式
	if fi.Size() == 0 {
		// 空文件：跳过数据复制，直接进行后续操作
	} else {
		// 非空文件：使用缓冲区进行数据拷贝
		bufSize := pool.CalculateBufferSize(fi.Size())
		buf := pool.GetByteCap(bufSize)
		defer pool.PutByte(buf)

		if _, err := io.CopyBuffer(out, in, buf); err != nil {
			return fmt.Errorf("failed to copy data from '%s' to '%s': %w", src, tmp, err)
		}

		// 强制刷盘，确保数据持久化（仅对非空文件）
		if err := out.Sync(); err != nil {
			return fmt.Errorf("failed to sync temporary file '%s': %w", tmp, err)
		}

	}

	// 在重命名前关闭文件句柄(Windows要求)
	if err := out.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file '%s': %w", tmp, err)
	}
	out = nil // 标记为已关闭

	// 原子重命名
	if err := os.Rename(tmp, dst); err != nil {
		// 复制失败，恢复备份文件
		restoreBackup(dst, backupPath)
		return fmt.Errorf("failed to rename temporary file '%s' to '%s': %w", tmp, dst, err)
	}

	// 复制成功，删除备份文件
	cleanupBackup(backupPath)

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
	// 安全覆盖机制：处理已存在的目标符号链接
	backupPath, err := handleBackupAndRestore(dst, overwrite)
	if err != nil {
		return err
	}

	// 读取符号链接的目标
	target, err := os.Readlink(src)
	if err != nil {
		restoreBackup(dst, backupPath)
		return fmt.Errorf("failed to read symlink '%s': %w", src, err)
	}

	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		restoreBackup(dst, backupPath)
		return fmt.Errorf("failed to create destination directory '%s': %w", dstDir, err)
	}

	// 创建符号链接
	if err := os.Symlink(target, dst); err != nil {
		restoreBackup(dst, backupPath)
		return fmt.Errorf("failed to create symlink '%s' -> '%s': %w", dst, target, err)
	}

	// 创建成功，删除备份
	cleanupBackup(backupPath)
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
	// 安全覆盖机制：处理已存在的目标文件
	backupPath, err := handleBackupAndRestore(dst, overwrite)
	if err != nil {
		return err
	}

	// 获取源文件信息
	srcInfo, err := os.Lstat(src)
	if err != nil {
		restoreBackup(dst, backupPath)
		return fmt.Errorf("failed to get source file info '%s': %w", src, err)
	}

	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		restoreBackup(dst, backupPath)
		return fmt.Errorf("failed to create destination directory '%s': %w", dstDir, err)
	}

	// 创建一个空文件，保持相同的权限模式
	file, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, srcInfo.Mode())
	if err != nil {
		restoreBackup(dst, backupPath)
		return fmt.Errorf("failed to create special file '%s': %w", dst, err)
	}

	// 立即关闭文件
	if err := file.Close(); err != nil {
		// 关闭失败，删除创建的文件并恢复备份
		_ = os.Remove(dst)
		restoreBackup(dst, backupPath)
		return fmt.Errorf("failed to close special file '%s': %w", dst, err)
	}

	// 创建成功，删除备份
	cleanupBackup(backupPath)
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
		return copyFile(src, dst, overwrite)

	case fileType&os.ModeSymlink != 0:
		// 符号链接
		return copySymlink(src, dst, overwrite)

	default:
		// 其他特殊文件（设备文件、命名管道、套接字等）
		return copySpecialFile(src, dst, overwrite)
	}
}

// copyDir 内部复制目录逻辑
// 用于递归复制整个目录，保持文件权限和目录结构
//
// 参数:
//   - src: 源目录路径
//   - dst: 目标目录路径
//   - overwrite: 是否允许覆盖已存在的目标文件
//
// 返回:
//   - error: 复制失败时返回错误
func copyDir(src, dst string, overwrite bool) error {
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

	// 安全覆盖机制：处理已存在的目标目录
	backupPath, err := handleBackupAndRestore(dst, overwrite)
	if err != nil {
		return err
	}

	// 创建目标目录，使用合适的权限（至少需要写权限以便后续操作）
	dirMode := srcInfo.Mode()
	if dirMode&0o200 == 0 {
		// 如果源目录没有写权限，临时添加写权限以便复制操作
		dirMode |= 0o200
	}
	if err := os.MkdirAll(dst, dirMode); err != nil {
		restoreBackup(dst, backupPath)
		return fmt.Errorf("failed to create destination directory '%s': %w", dst, err)
	}

	// 遍历源目录
	copyErr := filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
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

	// 处理复制结果
	if copyErr != nil {
		// 复制失败，清理已复制的内容并恢复备份
		_ = os.RemoveAll(dst) // 清理部分复制的内容
		restoreBackup(dst, backupPath)
		return copyErr
	}

	// 复制成功，恢复目录的原始权限（如果之前临时修改了权限）
	if srcInfo.Mode() != dirMode {
		_ = os.Chmod(dst, srcInfo.Mode()) // 忽略权限恢复错误
	}

	// 复制成功，删除备份目录
	cleanupBackup(backupPath)
	return nil
}
