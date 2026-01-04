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
	return ExecRemoteCmd(host, cmd, e.Timeout)
}

// execAll 通用执行逻辑（私有方法）
func (e *EasySSH) execAll(cmd, description string, handleResult func(hostLabel string, result RemoteExecResult)) error {
	hosts, err := e.LoadHosts()
	if err != nil {
		return fmt.Errorf("解析主机清单失败: %w", err)
	}

	if len(hosts) == 0 {
		fmt.Printf("==> 跳过 %s: 主机清单为空\n", description)
		return nil
	}

	fmt.Printf("==> 开始 %s (共 %d 台主机)...\n", description, len(hosts))

	successCount := 0
	for i, host := range hosts {
		hostLabel := fmt.Sprintf("%s:%d", host.Host, host.Port)
		fmt.Printf("  [%d/%d] %s: ", i+1, len(hosts), hostLabel)

		result := e.execOnHost(host, cmd)

		if result.Success {
			fmt.Println("✓ ok")
			if handleResult != nil {
				handleResult(hostLabel, result)
			}
			successCount++
		} else {
			fmt.Printf("✗ failed: %v\n", result.Err)
			if result.Output != "" {
				fmt.Printf("    %s\n", strings.TrimSpace(result.Output))
			}
		}
	}

	fmt.Printf("==> %s完成: 成功 %d/%d\n\n", description, successCount, len(hosts))
	return nil
}

// Exec 在所有主机上执行命令
//
// 参数：
//   - cmd: 要执行的命令
//   - description: 描述信息
//
// 返回：
//   - error: 执行错误，如果发生错误则返回非 nil 错误
func (e *EasySSH) Exec(cmd, description string) error {
	return e.execAll(cmd, description, func(hostLabel string, result RemoteExecResult) {
		if e.Verbose && result.Success {
			output := strings.TrimSpace(result.Output)
			fmt.Printf("    %s\n", output)
		}
	})
}

// ExecWithCallback 在所有主机上执行命令，并使用回调函数处理结果
//
// 参数：
//   - cmd: 要执行的命令
//   - description: 描述信息
//   - processFunc: 处理结果函数，接收两个参数：hostLabel 和 output，分别表示服务器标签和输出结果
func (e *EasySSH) ExecWithCallback(cmd, description string, processFunc func(hostLabel, output string)) {
	_ = e.execAll(cmd, description, func(hostLabel string, result RemoteExecResult) {
		if result.Success {
			output := strings.TrimSpace(result.Output)
			processFunc(hostLabel, output)
		}
	})
}

// PingHosts 测试所有主机的连通性
//
// 参数：
//   - description: 描述信息
//
// 返回：
//   - []PingResult: 每台主机的连通性测试结果
//   - error: 如果解析主机文件失败，返回错误
func (e *EasySSH) PingHosts(description string) ([]PingResult, error) {
	hosts, err := e.LoadHosts()
	if err != nil {
		return nil, fmt.Errorf("解析主机清单失败: %w", err)
	}

	if len(hosts) == 0 {
		fmt.Printf("==> 跳过 %s: 主机清单为空\n", description)
		return []PingResult{}, nil
	}

	fmt.Printf("==> 开始 %s (共 %d 台主机)...\n", description, len(hosts))

	results := make([]PingResult, 0, len(hosts))
	successCount := 0

	for i, host := range hosts {
		hostLabel := fmt.Sprintf("%s:%d", host.Host, host.Port)
		fmt.Printf("  [%d/%d] %s: ", i+1, len(hosts), hostLabel)

		// 测试 TCP 连通性
		startTime := time.Now()
		result := e.pingSingleHost(host.Host, host.Port)
		latency := time.Since(startTime)

		if result.Connected {
			fmt.Printf("✓ ok (%.2fms)\n", float64(latency.Nanoseconds())/1e6)
			result.Latency = latency
			successCount++
		} else {
			fmt.Printf("✗ failed\n")
			if e.Verbose && result.Err != nil {
				fmt.Printf("    %v\n", result.Err)
			}
		}

		results = append(results, result)
	}

	fmt.Printf("==> %s完成: 成功 %d/%d\n\n", description, successCount, len(hosts))
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
