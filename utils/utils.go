package utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// GetSize 获取文件或目录的大小
// 用于计算文件或目录的总字节数，目录会递归计算所有普通文件的大小
//
// 参数:
//   - path: 文件或目录路径
//
// 返回:
//   - int64: 文件或目录的总大小(字节)
//   - error: 路径不存在或访问失败时返回错误
func GetSize(path string) (int64, error) {
	// 获取路径信息
	info, err := os.Stat(path)
	if err != nil {
		// 判断错误类型，返回精准的错误信息
		if os.IsNotExist(err) {
			return 0, fmt.Errorf("路径不存在: %s", path)
		}
		if os.IsPermission(err) {
			return 0, fmt.Errorf("访问路径 '%s' 时权限不足: %w", path, err)
		}
		// 其他错误
		return 0, fmt.Errorf("获取路径 '%s' 信息失败: %w", path, err)
	}

	// 如果是普通文件，直接返回文件大小
	if info.Mode().IsRegular() {
		return info.Size(), nil
	}

	// 如果不是目录，提前返回 0(符号链接等特殊文件)
	if !info.IsDir() {
		return 0, nil
	}

	// 如果是目录，遍历计算总大小
	var totalSize int64
	walkDirErr := filepath.WalkDir(path, func(walkPath string, entry os.DirEntry, err error) error {
		if err != nil {
			// 判断错误类型
			if os.IsNotExist(err) {
				// 文件不存在，忽略并继续遍历
				return nil
			}
			if os.IsPermission(err) {
				// 权限错误，返回具体错误信息
				return fmt.Errorf("访问路径 '%s' 时权限不足: %w", walkPath, err)
			}
			// 其他错误，返回通用错误信息
			return fmt.Errorf("访问路径 '%s' 时出错: %w", walkPath, err)
		}

		// 只计算普通文件的大小
		if entry.Type().IsRegular() {
			if info, err := entry.Info(); err == nil {
				// 累加文件大小
				totalSize += info.Size()
			} else {
				// 判断获取文件信息时的错误类型
				if os.IsNotExist(err) {
					// 文件在遍历过程中被删除，忽略
					return nil
				}
				if os.IsPermission(err) {
					// 权限错误
					return fmt.Errorf("获取文件 '%s' 信息时权限不足: %w", walkPath, err)
				}
				// 其他错误
				return fmt.Errorf("获取文件 '%s' 信息时出错: %w", walkPath, err)
			}
		}

		return nil
	})

	if walkDirErr != nil {
		return 0, fmt.Errorf("遍历目录失败: %w", walkDirErr)
	}

	return totalSize, nil
}

// ExecuteCmd 执行指定的系统命令，并可设置独立的环境变量。
// 此函数会等待命令执行完成，不设置超时。
//
// 参数:
//   - args: 命令行参数切片，其中 args[0] 为要执行的命令本身（如 "ls", "go"），
//     后续元素为命令的参数（如 "-l", "main.go"）。
//   - env: 一个完整的环境变量切片，形如 "KEY=VALUE"。
//     如果传入 nil 或空切片，则命令将继承当前进程的环境变量。
//     如果传入非空切片，则命令的环境变量将仅限于此切片中定义的内容，
//     不会继承当前进程的任何环境变量。
//
// 返回:
//   - []byte: 命令的标准输出和标准错误合并后的内容。
//   - error: 如果命令执行失败（如命令不存在、权限问题、命令返回非零退出码），
//     或在执行过程中发生其他错误，则返回相应的错误信息。
func ExecuteCmd(args []string, env []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command")
	}
	cmd := exec.Command(args[0], args[1:]...)
	if len(env) > 0 {
		cmd.Env = env // 直接覆盖，不再继承系统环境
	}
	return cmd.CombinedOutput()
}

// ExecuteCmdWithTimeout 执行指定的系统命令，并设置超时时间及独立的环境变量。
// 此函数会等待命令执行完成，支持设置超时时间。
//
// 参数:
//   - timeout: 命令允许执行的最长时间。如果命令在此时间内未完成，将被终止并返回超时错误。
//     如果 timeout 为 0，则表示不设置超时。
//   - args: 命令行参数切片，其中 args[0] 为要执行的命令本身。
//   - env: 一个完整的环境变量切片，形如 "KEY=VALUE"。
//     如果传入 nil 或空切片，则命令将继承当前进程的环境变量。
//     如果传入非空切片，则命令的环境变量将仅限于此切片中定义的内容。
//
// 返回:
//   - []byte: 命令的标准输出和标准错误合并后的内容。
//   - error: 如果命令执行失败、超时，或在执行过程中发生其他错误，则返回相应的错误信息。
func ExecuteCmdWithTimeout(timeout time.Duration, args []string, env []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 创建命令
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	// 设置命令的额外环境变量
	if len(env) > 0 {
		cmd.Env = env // 直接覆盖，不再继承系统环境
	}

	// 执行命令并返回结果
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 检查是否为超时错误
		if ctx.Err() == context.DeadlineExceeded {
			return output, fmt.Errorf("命令超时 (超过 %v)", timeout)
		}
		// 其他错误类型
		return output, fmt.Errorf("执行命令失败: %v 错误: %v", args, err)
	}

	return output, nil
}
