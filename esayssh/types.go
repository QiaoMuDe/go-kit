package esayssh

import "time"

// HostConfig 存储单台主机的 SSH 配置
type HostConfig struct {
	Host     string // 主机地址
	Port     int    // 端口，默认22
	Username string // 用户名
	Password string // 密码
}

// RemoteExecResult 远程命令执行结果结构体
type RemoteExecResult struct {
	Success bool   // 执行是否成功
	Output  string // 命令输出内容（标准输出+标准错误）
	Err     error  // 执行过程中的错误信息
}

// PingResult Ping 结果结构体
type PingResult struct {
	Host      string        // 主机地址
	Port      int           // 端口
	Connected bool          // 是否连接成功
	Latency   time.Duration // 连接延迟
	Err       error         // 错误信息
}

// EasySSH SSH管理器（基础版）
type EasySSH struct {
	HostsFile string        // 主机配置文件路径
	Timeout   time.Duration // 连接超时时间
	Verbose   bool          // 是否打印详细输出
	hosts     []HostConfig  // 缓存的主机列表
}
