package utils

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

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

const (
	// 使用位运算常量，1024 = 1 << 10
	_KB = 1 << 10 // 1024
	_MB = 1 << 20 // 1048576
	_GB = 1 << 30 // 1073741824
	_TB = 1 << 40 // 1099511627776
	_PB = 1 << 50 // 1125899906842624
)

// 预定义单位数组，避免每次函数调用时重新创建
var units = [6]string{"B", "KB", "MB", "GB", "TB", "PB"}

// FormatBytes 将字节数转换为人类可读的带单位的字符串
// 用于将字节数格式化为易读的存储单位格式，支持B到PB的转换
//
// 参数:
//   - bytes: 字节数（int64类型）
//
// 返回:
//   - string: 格式化后的字符串，如 "1.23 KB", "456.78 MB", "2.34 GB" 等
func FormatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	// 处理负数
	if bytes < 0 {
		return "-" + FormatBytes(-bytes)
	}

	// 使用条件判断代替循环，提高性能
	switch {
	case bytes < _KB:
		return strconv.FormatInt(bytes, 10) + " B"
	case bytes < _MB:
		return formatWithUnit(bytes, _KB, 0)
	case bytes < _GB:
		return formatWithUnit(bytes, _MB, 1)
	case bytes < _TB:
		return formatWithUnit(bytes, _GB, 2)
	case bytes < _PB:
		return formatWithUnit(bytes, _TB, 3)
	default:
		return formatWithUnit(bytes, _PB, 4)
	}
}

// formatWithUnit 格式化字节数为指定单位
// 用于将字节数按指定除数转换为对应单位的格式化字符串
//
// 参数:
//   - bytes: 字节数（int64类型）
//   - divisor: 除数，用于计算单位
//   - unitIndex: 单位索引，对应units数组中的位置
//
// 返回:
//   - string: 格式化后的字符串，保留两位小数
func formatWithUnit(bytes, divisor int64, unitIndex int) string {
	// 计算整数部分和小数部分
	quotient := bytes / divisor
	remainder := bytes % divisor

	// 计算两位小数（乘以100后除以divisor再取整）
	decimal := (remainder * 100) / divisor

	// 构建结果字符串
	if decimal == 0 {
		return strconv.FormatInt(quotient, 10) + " " + units[unitIndex+1]
	}

	// 格式化小数部分，确保两位数显示
	var decimalStr string
	if decimal < 10 {
		decimalStr = "0" + strconv.FormatInt(decimal, 10)
	} else {
		decimalStr = strconv.FormatInt(decimal, 10)
	}

	return strconv.FormatInt(quotient, 10) + "." + decimalStr + " " + units[unitIndex+1]
}
