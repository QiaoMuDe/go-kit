# esayssh

```go
package esayssh // import "gitee.com/MM-Q/go-kit/esayssh"
```

## 主机清单格式

EasySSH 支持两种主机配置格式：

### 3字段格式（推荐）
```
# 格式：主机地址 用户名 密码
192.168.1.100 root mypassword
192.168.1.101 admin 123456
```

### 4字段格式
```
# 格式：主机地址 端口 用户名 密码
192.168.1.100 22 root mypassword
192.168.1.101 2222 admin 123456
```

### 注意事项
- 空行和以 `#` 开头的行将被忽略
- 3字段格式下，端口默认为 22
- 端口必须在 1-65535 范围内
- 字段之间使用空格分隔，支持多个连续空格

## TYPES

### type EasySSH struct

```go
type EasySSH struct {
	HostsFile string        // 主机配置文件路径
	Timeout   time.Duration // 连接超时时间
	ShowOutput bool         // 是否显示命令输出
	ShowFormat bool         // 是否显示格式化执行输出
	// Has unexported fields.
}
```

EasySSH SSH管理器

#### func New

```go
func New(hostsFile string, timeout time.Duration, showOutput, showFormat bool) *EasySSH
```

New 创建 EasySSH 实例

参数：
  - hostsFile: 主机清单文件路径
  - timeout: 超时时间
  - showOutput: 是否显示命令输出
  - showFormat: 是否显示格式化执行输出

返回：
  - *EasySSH: 新创建的 EasySSH 实例

#### func NewDef

```go
func NewDef(hostsFile string) *EasySSH
```

NewDef 创建 EasySSH 实例，使用默认设置

默认设置：
  - 超时时间：3秒
  - 显示命令输出：true
  - 显示格式化执行输出：true

参数：
  - hostsFile: 主机清单文件路径

返回：
  - *EasySSH: 新创建的 EasySSH 实例

#### func (*EasySSH) Exec

```go
func (e *EasySSH) Exec(cmd, description string) error
```

Exec 在所有主机上执行命令

参数：
  - cmd: 要执行的命令
  - description: 描述信息

返回：
  - error: 执行错误，如果发生错误则返回非 nil 错误

#### func (*EasySSH) ExecWithCallback

```go
func (e *EasySSH) ExecWithCallback(cmd, description string, processFunc func(hostLabel, output string))
```

ExecWithCallback 在所有主机上执行命令，并使用回调函数处理结果

参数：
  - cmd: 要执行的命令
  - description: 描述信息
  - processFunc: 处理结果函数，接收两个参数：hostLabel 和 output，分别表示服务器标签和输出结果

#### func (*EasySSH) LoadHosts

```go
func (e *EasySSH) LoadHosts() ([]HostConfig, error)
```

LoadHosts 加载主机配置文件

返回：
  - []HostConfig: 解析得到的主机配置列表
  - error: 解析错误，如果发生错误则返回非 nil 错误

#### func (*EasySSH) PingHosts

```go
func (e *EasySSH) PingHosts() error
```

PingHosts 测试所有主机的连通性

返回:
  - error: 如果解析主机文件失败，返回错误

#### func (*EasySSH) PingHostsRaw

```go
func (e *EasySSH) PingHostsRaw() ([]PingResult, error)
```

PingHostsRaw 测试所有主机的连通性，返回原始结果

返回：
  - []PingResult: 每台主机的连通性测试结果
  - error: 如果解析主机文件失败，返回错误

#### func (*EasySSH) ReloadHosts

```go
func (e *EasySSH) ReloadHosts() error
```

ReloadHosts 重新加载主机配置文件

### type HostConfig struct

```go
type HostConfig struct {
	Host     string // 主机地址
	Port     int    // 端口，默认22
	Username string // 用户名
	Password string // 密码
}
```

HostConfig 存储单台主机的 SSH 配置

#### func ParseHostsFile

```go
func ParseHostsFile(filePath string) ([]HostConfig, error)
```

ParseHostsFile 解析主机配置文件，返回 HostConfig 切片和错误信息

参数：
  - filePath: 主机配置文件路径

返回：
  - []HostConfig: 解析后的主机配置切片
  - error: 如果解析过程中出错，返回具体的错误信息；否则返回 nil

### type PingResult struct

```go
type PingResult struct {
	Host      string        // 主机地址
	Port      int           // 端口
	Connected bool          // 是否连接成功
	Latency   time.Duration // 连接延迟
	Err       error         // 错误信息
}
```

PingResult Ping 结果结构体

### type RemoteExecResult struct

```go
type RemoteExecResult struct {
	Success bool   // 执行是否成功
	Output  string // 命令输出内容（标准输出+标准错误）
	Err     error  // 执行过程中的错误信息
}
```

RemoteExecResult 远程命令执行结果结构体

#### func ExecRemoteCmd

```go
func ExecRemoteCmd(host HostConfig, cmd string, timeout time.Duration) RemoteExecResult
```

ExecRemoteCmd 远程执行命令的核心函数

参数：
  - host: 主机信息结构体，包含连接信息（主机地址、端口、用户名、密码）
  - cmd: 要执行的命令字符串
  - timeout: 连接超时时间（零值表示不设置超时）

返回：
  - RemoteExecResult: 命令执行结果结构体