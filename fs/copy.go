package fs

import (
	"io"
	"os"
	"path/filepath"
)

// CopyFile 复制文件并继承权限
// 用于安全地复制文件，保持原文件的权限信息，失败时自动清理
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标文件路径
//
// 返回:
//   - error: 复制失败时返回错误
func CopyFile(src, dst string) (err error) {
	// 打开源文件
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	// 获取源文件元数据
	fi, err := in.Stat()
	if err != nil {
		return err
	}

	// 目标目录必须存在
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	// 创建临时文件（与目标同目录，保证 rename 原子）
	tmp := dst + ".tmp"
	out, err := os.OpenFile(tmp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fi.Mode())
	if err != nil {
		return err
	}

	// 捕获所有错误，统一清理
	success := false
	defer func() {
		_ = out.Close()
		if !success {
			_ = os.Remove(tmp) // 忽略清理错误
		}
	}()

	// 数据拷贝
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	// 强制刷盘，确保 rename 前数据已落盘
	if err := out.Sync(); err != nil {
		return err
	}

	// 复制文件权限
	if err := os.Chmod(tmp, fi.Mode()); err != nil {
		return err
	}

	// 原子重命名
	if err := os.Rename(tmp, dst); err != nil {
		return err
	}
	success = true
	return nil
}
