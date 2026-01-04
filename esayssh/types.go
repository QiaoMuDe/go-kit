package esayssh

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
