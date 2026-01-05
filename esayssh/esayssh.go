package esayssh

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// New 创建 EasySSH 实例
//
// 参数：
//   - hostsFile: 主机清单文件路径
//   - timeout: 超时时间
//   - verbose: 是否启用详细模式
//
// 返回：
//   - *EasySSH: 新创建的 EasySSH 实例
func New(hostsFile string, timeout time.Duration, verbose bool) *EasySSH {
	return &EasySSH{
		HostsFile: hostsFile,
		Timeout:   timeout,
		Verbose:   verbose,
	}
}

// LoadHosts 加载主机配置文件
//
// 返回：
//   - []HostConfig: 解析得到的主机配置列表
//   - error: 解析错误，如果发生错误则返回非 nil 错误
func (e *EasySSH) LoadHosts() ([]HostConfig, error) {
	if e.hosts == nil {
		hosts, err := ParseHostsFile(e.HostsFile)
		if err != nil {
			return nil, err
		}
		e.hosts = hosts
	}
	return e.hosts, nil
}

// ReloadHosts 重新加载主机配置文件
func (e *EasySSH) ReloadHosts() error {
	hosts, err := ParseHostsFile(e.HostsFile)
	if err != nil {
		return err
	}
	e.hosts = hosts
	return nil
}

// execOnHost 在单台主机上执行命令（私有方法）
//
// 参数：
//   - host: 要执行的主机配置
//   - cmd: 要执行的命令
//
// 返回：
//   - RemoteExecResult: 执行结果
func (e *EasySSH) execOnHost(host HostConfig, cmd string) RemoteExecResult {
	result := ExecRemoteCmd(host, cmd, e.Timeout)
	result.Host = host.Host
	result.Port = host.Port
	return result
}

// Exec 在所有主机上执行命令并打印结果
//
// 参数：
//   - cmd: 要执行的命令
//   - description: 描述信息
//
// 返回：
//   - error: 执行错误，如果发生错误则返回非 nil 错误
func (e *EasySSH) Exec(cmd, description string) error {
	results, err := e.ExecRaw(cmd)
	if err != nil {
		return err
	}

	// 打印结果
	if len(results) == 0 {
		fmt.Printf("==> 跳过 %s: 主机清单为空\n", description)
		return nil
	}

	fmt.Printf("==> %s (%d hosts)\n", description, len(results))
	fmt.Println("----------------------------------------")

	successCount := 0
	for _, result := range results {
		hostLabel := fmt.Sprintf("%s:%d", result.Host, result.Port)

		if result.Success {
			fmt.Printf("%-20s : [ ✓ ok ]\n", hostLabel)
			if e.Verbose && result.Output != "" {
				fmt.Printf("    %s\n", strings.TrimSpace(result.Output))
			}
			successCount++
		} else {
			fmt.Printf("%-20s : [ ✗ failed ]\n", hostLabel)
			if result.Output != "" {
				fmt.Printf("    %s\n", strings.TrimSpace(result.Output))
			}
		}
	}

	fmt.Println("----------------------------------------")
	fmt.Printf("==> 成功: %d/%d | 失败: %d/%d\n\n", successCount, len(results), len(results)-successCount, len(results))
	return nil
}

// ExecRaw 在所有主机上执行命令并返回结构化结果
//
// 参数：
//   - cmd: 要执行的命令
//
// 返回：
//   - []RemoteExecResult: 每台主机的命令执行结果
//   - error: 如果解析主机文件失败，返回错误
func (e *EasySSH) ExecRaw(cmd string) ([]RemoteExecResult, error) {
	hosts, err := e.LoadHosts()
	if err != nil {
		return nil, fmt.Errorf("解析主机清单失败: %w", err)
	}

	if len(hosts) == 0 {
		return []RemoteExecResult{}, nil
	}

	results := make([]RemoteExecResult, 0, len(hosts))

	for _, host := range hosts {
		result := e.execOnHost(host, cmd)
		results = append(results, result)
	}

	return results, nil
}

// ExecWithCallback 在所有主机上执行命令，并使用回调函数处理结果
//
// 参数：
//   - cmd: 要执行的命令
//   - description: 描述信息
//   - processFunc: 处理结果函数，接收两个参数：hostLabel 和 output，分别表示服务器标签和输出结果
func (e *EasySSH) ExecWithCallback(cmd, description string, processFunc func(hostLabel, output string)) {
	results, err := e.ExecRaw(cmd)
	if err != nil {
		fmt.Printf("执行命令失败: %v\n", err)
		return
	}

	for _, result := range results {
		hostLabel := fmt.Sprintf("%s:%d", result.Host, result.Port)
		if result.Success {
			output := strings.TrimSpace(result.Output)
			processFunc(hostLabel, output)
		}
	}
}

// PingHosts 测试所有主机的连通性并打印结果
//
// 返回：
//   - error: 如果解析主机文件失败，返回错误
func (e *EasySSH) PingHosts() error {
	results, err := e.PingHostsRaw()
	if err != nil {
		return err
	}

	// 打印结果
	if len(results) == 0 {
		fmt.Println("==> 跳过 PING: 主机清单为空")
		return nil
	}

	fmt.Printf("==> PING (%d hosts)\n", len(results))
	fmt.Println("----------------------------------------")

	successCount := 0
	for _, result := range results {
		hostLabel := fmt.Sprintf("%s:%d", result.Host, result.Port)

		if result.Connected {
			fmt.Printf("%-20s : [ ✓ ok (%.2fms) ]\n", hostLabel, float64(result.Latency.Nanoseconds())/1e6)
			successCount++
		} else {
			fmt.Printf("%-20s : [ ✗ failed ]\n", hostLabel)
			if e.Verbose && result.Err != nil {
				fmt.Printf("    %v\n", result.Err)
			}
		}
	}

	fmt.Println("----------------------------------------")
	fmt.Printf("==> 成功: %d/%d | 失败: %d/%d\n\n", successCount, len(results), len(results)-successCount, len(results))
	return nil
}

// PingHostsRaw 测试所有主机的连通性并返回结构化结果
//
// 返回：
//   - []PingResult: 每台主机的连通性测试结果
//   - error: 如果解析主机文件失败，返回错误
func (e *EasySSH) PingHostsRaw() ([]PingResult, error) {
	hosts, err := e.LoadHosts()
	if err != nil {
		return nil, fmt.Errorf("解析主机清单失败: %w", err)
	}

	if len(hosts) == 0 {
		return []PingResult{}, nil
	}

	results := make([]PingResult, 0, len(hosts))

	for _, host := range hosts {
		// 测试 TCP 连通性
		startTime := time.Now()
		result := e.pingSingleHost(host.Host, host.Port)
		latency := time.Since(startTime)

		if result.Connected {
			result.Latency = latency
		}

		results = append(results, result)
	}

	return results, nil
}

// pingSingleHost 测试单个主机的连通性（私有方法）
func (e *EasySSH) pingSingleHost(host string, port int) PingResult {
	result := PingResult{
		Host: host,
		Port: port,
	}

	timeout := e.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second // 默认超时5秒
	}

	// 使用 net.DialTimeout 测试 TCP 连通性
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		result.Connected = false
		result.Err = err
		return result
	}
	defer func() { _ = conn.Close() }()

	result.Connected = true
	return result
}
