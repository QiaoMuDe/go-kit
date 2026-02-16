package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// Move 通用移动函数，将文件或目录移动到目标位置
// 支持移动普通文件、目录、符号链接和特殊文件
// 优先使用 os.Rename（同文件系统内），失败时降级使用复制+删除（支持跨文件系统）
//
// 参数:
//   - src: 源路径 (支持文件、目录、符号链接、特殊文件)
//   - dst: 目标路径（支持文件、目录，自动创建父目录）
//
// 返回:
//   - error: 移动失败时返回错误，如果目标已存在则返回错误
func Move(src, dst string) error {
	return MoveEx(src, dst, false)
}

// MoveEx 通用移动函数 (可控制是否覆盖)，将文件或目录移动到目标位置
// 支持移动普通文件、目录、符号链接和特殊文件
// 优先使用 os.Rename（同文件系统内），失败时降级使用复制+删除（支持跨文件系统）
//
// 参数:
//   - src: 源路径 (支持文件、目录、符号链接、特殊文件)
//   - dst: 目标路径（支持文件、目录，自动创建父目录）
//   - overwrite: 是否允许覆盖已存在的目标文件/目录
//
// 返回:
//   - error: 移动失败时返回错误
//
// 智能路径处理:
//   - 如果 dst 是已存在的目录，会自动追加源文件名/目录名
//   - 例如: Move("a.txt", "existingDir") → 移动到 existingDir/a.txt
//   - 例如: Move("dirA", "existingDir") → 移动到 existingDir/dirA/
//
// 移动策略:
//  1. 优先使用 os.Rename（原子操作，同文件系统内高效）
//  2. rename 失败时降级使用 CopyEx + os.RemoveAll（支持跨文件系统）
func MoveEx(src, dst string, overwrite bool) error {
	// 获取绝对路径
	srcAbs, err := filepath.Abs(filepath.Clean(src))
	if err != nil {
		return fmt.Errorf("failed to get absolute path for source '%s': %w", src, err)
	}
	dstAbs, err := filepath.Abs(filepath.Clean(dst))
	if err != nil {
		return fmt.Errorf("failed to get absolute path for destination '%s': %w", dst, err)
	}

	// 基础路径验证：检查路径非空、源目标不相同、防止子目录循环
	// 注意：这里始终检查子目录，因为移动时这个问题和复制时一样严重
	if err := validatePaths(srcAbs, dstAbs, true); err != nil {
		return err
	}

	// 智能路径处理：如果目标是已存在的目录，自动追加源文件名/目录名
	dstAbs = resolveDestinationPathAbs(srcAbs, dstAbs)

	// 策略1：优先尝试 os.Rename（同文件系统内）
	// os.Rename 是原子操作，且只改变文件系统的元数据，不复制数据
	renameErr := tryRename(srcAbs, dstAbs, overwrite)
	if renameErr == nil {
		// rename 成功，直接返回
		return nil
	}

	// 如果 rename 失败是因为目标已存在且不允许覆盖，直接返回错误
	if !overwrite {
		if _, err := os.Lstat(dstAbs); err == nil {
			return renameErr
		}
	}

	// 策略2：rename 失败，降级使用复制+删除（跨文件系统场景）
	// 先执行复制操作
	if err := CopyEx(src, dst, overwrite); err != nil {
		return fmt.Errorf("failed to copy '%s' to '%s': %w", src, dst, err)
	}

	// 复制成功，删除源文件/目录
	if err := os.RemoveAll(srcAbs); err != nil {
		return fmt.Errorf("copy succeeded but failed to remove source '%s': %w", srcAbs, err)
	}

	return nil
}

// tryRename 尝试使用 os.Rename 移动文件/目录
// 适用于同文件系统内的快速移动（原子操作）
//
// 参数:
//   - srcAbs: 源绝对路径
//   - dstAbs: 目标绝对路径
//   - overwrite: 是否允许覆盖
//
// 返回:
//   - error: rename 失败时返回错误，成功返回 nil
func tryRename(srcAbs, dstAbs string, overwrite bool) error {
	// 检查目标是否存在
	_, err := os.Lstat(dstAbs)
	if err == nil {
		// 目标存在
		if !overwrite {
			return fmt.Errorf("destination '%s' already exists", dstAbs)
		}
		// 允许覆盖，删除目标
		if err := os.RemoveAll(dstAbs); err != nil {
			return fmt.Errorf("failed to remove destination '%s' for overwrite: %w", dstAbs, err)
		}
	}

	// 尝试 rename
	if err := os.Rename(srcAbs, dstAbs); err != nil {
		return fmt.Errorf("rename failed: %w", err)
	}

	return nil
}
