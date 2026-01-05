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
//   - showOutput: 是否显示命令输出
//   - showFormat: 是否显示格式化执行输出
//
// 返回：
//   - *EasySSH: 新创建的 EasySSH 实例
func New(hostsFile string, timeout time.Duration, showOutput, showFormat bool) *EasySSH {
	return &EasySSH{
		HostsFile:  hostsFile,
		Timeout:    timeout,
		ShowOutput: showOutput,
		ShowFormat: showFormat,
	}
}

// NewDef 创建 EasySSH 实例，使用默认设置
//
// 默认设置：
//   - 超时时间：3秒
//   - 显示命令输出：true
//   - 显示格式化执行输出：true
//
// 参数：
//   - hostsFile: 主机清单文件路径
//
// 返回：
//   - *EasySSH: 新创建的 EasySSH 实例
func NewDef(hostsFile string) *EasySSH {
	return &EasySSH{
		HostsFile:  hostsFile,
		Timeout:    3 * time.Second,
		ShowOutput: true,
		ShowFormat: true,
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
		if e.ShowFormat {
			fmt.Printf("==> 跳过 %s: 主机清单为空\n", description)
		}
		return nil
	}

	if e.ShowFormat {
		fmt.Printf("==> %s (%d hosts)\n", description, len(hosts))
		fmt.Println("----------------------------------------")
	}

	successCount := 0
	for _, host := range hosts {
		hostLabel := fmt.Sprintf("%s:%d", host.Host, host.Port)
		result := e.execOnHost(host, cmd)

		if result.Success {
			if e.ShowFormat {
				fmt.Printf("%-20s : [ ✓ ok ]\n", hostLabel)
			}
			if handleResult != nil {
				handleResult(hostLabel, result)
			}
			successCount++
		} else {
			if e.ShowFormat {
				fmt.Printf("%-20s : [ ✗ failed ]\n", hostLabel)
			}
			if result.Output != "" && e.ShowOutput {
				fmt.Printf("    %s\n", strings.TrimSpace(result.Output))
			}
		}
	}

	if e.ShowFormat {
		fmt.Println("----------------------------------------")
		fmt.Printf("==> 成功: %d/%d | 失败: %d/%d\n\n", successCount, len(hosts), len(hosts)-successCount, len(hosts))
	}
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
		if e.ShowOutput && result.Success {
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

// PingHosts 测试所有主机的连通性并打印结果
//
// 返回:
//   - error: 如果解析主机文件失败，返回错误
func (e *EasySSH) PingHosts() error {
	// 临时保存原始设置
	originalShowFormat := e.ShowFormat

	// 确保显示格式化输出
	e.ShowFormat = true
	defer func() {
		e.ShowFormat = originalShowFormat
	}()

	_, err := e.pingHosts()
	return err
}

// PingHostsRaw 测试所有主机的连通性，返回原始结果
//
// 返回：
//   - []PingResult: 每台主机的连通性测试结果
//   - error: 如果解析主机文件失败，返回错误
func (e *EasySSH) PingHostsRaw() ([]PingResult, error) {
	return e.pingHosts()
}

// pingHosts 测试所有主机的连通性
//
// 返回：
//   - []PingResult: 每台主机的连通性测试结果
//   - error: 如果解析主机文件失败，返回错误
func (e *EasySSH) pingHosts() ([]PingResult, error) {
	hosts, err := e.LoadHosts()
	if err != nil {
		return nil, fmt.Errorf("解析主机清单失败: %w", err)
	}

	if len(hosts) == 0 {
		if e.ShowFormat {
			fmt.Println("==> 跳过 PING: 主机清单为空")
		}
		return []PingResult{}, nil
	}

	if e.ShowFormat {
		fmt.Printf("==> PING (%d hosts)\n", len(hosts))
		fmt.Println("----------------------------------------")
	}

	results := make([]PingResult, 0, len(hosts))
	successCount := 0

	for _, host := range hosts {
		hostLabel := fmt.Sprintf("%s:%d", host.Host, host.Port)

		// 测试 TCP 连通性
		startTime := time.Now()
		result := e.pingSingleHost(host.Host, host.Port)
		latency := time.Since(startTime)

		if result.Connected {
			if e.ShowFormat {
				fmt.Printf("%-20s : [ ✓ ok (%.2fms) ]\n", hostLabel, float64(latency.Nanoseconds())/1e6)
			}
			result.Latency = latency
			successCount++
		} else {
			if e.ShowFormat {
				fmt.Printf("%-20s : [ ✗ failed ]\n", hostLabel)
			}
			if e.ShowOutput && result.Err != nil {
				fmt.Printf("    %v\n", result.Err)
			}
		}

		results = append(results, result)
	}

	if e.ShowFormat {
		fmt.Println("----------------------------------------")
		fmt.Printf("==> 成功: %d/%d | 失败: %d/%d\n\n", successCount, len(hosts), len(hosts)-successCount, len(hosts))
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
