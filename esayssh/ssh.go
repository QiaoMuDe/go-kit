package esayssh

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

// ParseHostsFile 解析主机配置文件，返回 HostConfig 切片和错误信息
//
// 参数：
//   - filePath: 主机配置文件路径
//
// 返回：
//   - []HostConfig: 解析后的主机配置切片
//   - error: 如果解析过程中出错，返回具体的错误信息；否则返回 nil
func ParseHostsFile(filePath string) ([]HostConfig, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("关闭文件失败: %v\n", closeErr)
		}
	}()

	var hosts []HostConfig
	scanner := bufio.NewScanner(file)

	// 逐行读取解析
	for lineNum := 1; scanner.Scan(); lineNum++ {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和 # 开头的注释行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 按空格分割字段，自动忽略连续空格
		fields := strings.Fields(line)
		var cfg HostConfig

		switch len(fields) {
		case 3:
			// 三字段格式：主机地址 用户名 密码，端口默认22
			cfg = HostConfig{
				Host:     fields[0],
				Port:     22,
				Username: fields[1],
				Password: fields[2],
			}
		case 4:
			// 四字段格式：主机地址 端口 用户名 密码
			port, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid port: %w", lineNum, err)
			}
			if port <= 0 || port > 65535 {
				return nil, fmt.Errorf("line %d: port out of range (1-65535)", lineNum)
			}
			cfg = HostConfig{
				Host:     fields[0],
				Port:     port,
				Username: fields[2],
				Password: fields[3],
			}
		default:
			// 字段数量不合法
			return nil, fmt.Errorf("line %d: invalid field count (expected 3 or 4, got %d)", lineNum, len(fields))
		}

		hosts = append(hosts, cfg)
	}

	// 检查扫描过程中是否出错
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}

	return hosts, nil
}

// ExecRemoteCmd 远程执行命令的核心函数
//
// 参数：
//   - host: 主机信息结构体，包含连接信息（主机地址、端口、用户名、密码）
//   - cmd: 要执行的命令字符串
//
// 返回：
//   - RemoteExecResult: 命令执行结果结构体
func ExecRemoteCmd(host HostConfig, cmd string) RemoteExecResult {
	// 1. 校验入参合法性
	if err := validateHostConfig(host); err != nil {
		return RemoteExecResult{
			Success: false,
			Output:  "",
			Err:     err,
		}
	}
	if strings.TrimSpace(cmd) == "" {
		return RemoteExecResult{
			Success: false,
			Output:  "",
			Err:     errors.New("执行的命令不能为空"),
		}
	}

	// 2. 配置SSH客户端参数
	config := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password), // 密码认证
		},
		// 生产环境需替换为 ssh.FixedHostKey(hostKey) 进行主机密钥校验
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// 超时时间（可选），需要导入 "time" 包：Timeout: time.Second * 30
	}

	// 3. 建立SSH连接
	addr := fmt.Sprintf("%s:%d", host.Host, host.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return RemoteExecResult{
			Success: false,
			Output:  "",
			Err:     fmt.Errorf("SSH连接失败: %w", err),
		}
	}
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			// EOF是SSH连接关闭时的正常情况，不需要记录为错误
			if !errors.Is(closeErr, io.EOF) {
				fmt.Printf("关闭SSH客户端失败: %v\n", closeErr)
			}
		}
	}() // 延迟关闭客户端连接

	// 4. 创建SSH会话
	session, err := client.NewSession()
	if err != nil {
		return RemoteExecResult{
			Success: false,
			Output:  "",
			Err:     fmt.Errorf("创建SSH会话失败: %w", err),
		}
	}
	defer func() {
		if closeErr := session.Close(); closeErr != nil {
			// EOF是SSH会话关闭时的正常情况，不需要记录为错误
			if !errors.Is(closeErr, io.EOF) {
				fmt.Printf("关闭SSH会话失败: %v\n", closeErr)
			}
		}
	}() // 延迟关闭会话

	// 5. 执行远程命令并获取输出
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return RemoteExecResult{
			Success: false,
			Output:  string(output),
			Err:     fmt.Errorf("命令执行失败: %w", err),
		}
	}

	// 6. 执行成功返回结果
	return RemoteExecResult{
		Success: true,
		Output:  string(output),
		Err:     nil,
	}
}

// validateHostConfig 校验主机信息的合法性
//
// 参数：
//   - host: 主机信息结构体，包含连接信息（主机地址、端口、用户名、密码）
//
// 返回：
//   - error: 如果校验失败，返回具体的错误信息；否则返回 nil
func validateHostConfig(host HostConfig) error {
	if strings.TrimSpace(host.Host) == "" {
		return errors.New("主机地址不能为空")
	}
	if host.Port <= 0 || host.Port > 65535 {
		return errors.New("端口号必须在1-65535之间")
	}
	if strings.TrimSpace(host.Username) == "" {
		return errors.New("登录用户名不能为空")
	}
	if strings.TrimSpace(host.Password) == "" {
		return errors.New("登录密码不能为空")
	}
	return nil
}

// ExecOnAllHosts 在指定主机配置文件中的所有主机上执行命令
//
// 参数：
//   - hostsFilePath: 主机配置文件路径
//   - cmd: 要执行的命令字符串
//   - description: 操作描述（用于日志输出）
//   - processFunc: 处理单个主机结果的回调函数
//     参数: hostLabel 主机标签, output 命令输出
//
// 处理函数签名：func(hostLabel string, output string)
func ExecOnAllHosts(hostsFilePath, cmd, description string, processFunc func(hostLabel, output string)) {
	// 解析主机清单
	hosts, err := ParseHostsFile(hostsFilePath)
	if err != nil {
		fmt.Printf("解析主机清单失败: %v\n", err)
		return
	}

	if len(hosts) == 0 {
		fmt.Printf("主机清单为空，没有服务器需要%s\n", description)
		return
	}

	fmt.Printf("%s, 共 %d 台服务器...\n", description, len(hosts))

	// 遍历所有主机，执行命令
	successCount := 0
	for i, host := range hosts {
		hostLabel := fmt.Sprintf("服务器[%d] %s:%d", i+1, host.Host, host.Port)

		// 执行命令
		result := ExecRemoteCmd(host, cmd)

		if result.Success {
			// 调用处理函数处理成功的结果
			output := strings.TrimSpace(result.Output)
			processFunc(hostLabel, output)
			successCount++
		} else {
			fmt.Printf("%-25s 执行失败: %v\n", hostLabel, result.Err)
			if result.Output != "" {
				fmt.Printf(" 输出: %s\n", result.Output)
			}
		}
	}

	// 输出统计信息
	fmt.Println()
	fmt.Printf("%s完成: 成功 %d/%d 台服务器\n", description, successCount, len(hosts))
}

// ExecCmdOnHosts 在指定主机文件路径中的所有主机上执行命令
//
// 参数：
//   - hostsFilePath: 主机配置文件路径
//   - cmd: 要执行的命令字符串
//   - description: 操作描述（用于日志输出）
//   - verbose: 是否打印详细输出到终端
//
// 返回：
//   - error: 如果解析主机文件失败，返回错误
func ExecCmdOnHosts(hostsFilePath, cmd, description string, verbose bool) error {
	// 解析主机清单
	hosts, err := ParseHostsFile(hostsFilePath)
	if err != nil {
		return fmt.Errorf("解析主机清单失败: %w", err)
	}

	if len(hosts) == 0 {
		fmt.Printf("主机清单为空，没有服务器需要%s\n", description)
		return nil
	}

	fmt.Printf("%s, 共 %d 台服务器...\n", description, len(hosts))

	// 遍历所有主机，执行命令
	successCount := 0
	for i, host := range hosts {
		hostLabel := fmt.Sprintf("服务器[%d] %s:%d", i+1, host.Host, host.Port)

		// 执行命令
		result := ExecRemoteCmd(host, cmd)

		if result.Success {
			// 打印输出（如果启用详细模式）
			if verbose {
				output := strings.TrimSpace(result.Output)
				fmt.Printf("%-25s 执行成功\n", hostLabel)
				if output != "" {
					fmt.Printf(" 输出: %s\n", output)
				}
			}
			successCount++
		} else {
			fmt.Printf("%-25s 执行失败: %v\n", hostLabel, result.Err)
			if result.Output != "" {
				fmt.Printf(" 输出: %s\n", result.Output)
			}
		}
	}

	// 输出统计信息
	fmt.Println()
	fmt.Printf("%s完成: 成功 %d/%d 台服务器\n", description, successCount, len(hosts))

	return nil
}
