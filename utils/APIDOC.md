# Package utils

```go
import "gitee.com/MM-Q/go-kit/utils"
```

## Functions

### func ExecuteCmd

```go
func ExecuteCmd(args []string, env []string) ([]byte, error)
```

ExecuteCmd 执行指定的系统命令，并可设置独立的环境变量。 此函数会等待命令执行完成，不设置超时。

**参数:**
- `args`: 命令行参数切片，其中 args[0] 为要执行的命令本身（如 "ls", "go"）， 后续元素为命令的参数（如 "-l", "main.go"）。
- `env`: 一个完整的环境变量切片，形如 "KEY=VALUE"。 如果传入 nil 或空切片，则命令将继承当前进程的环境变量。 如果传入非空切片，则命令的环境变量将仅限于此切片中定义的内容， 不会继承当前进程的任何环境变量。

**返回:**
- `[]byte`: 命令的标准输出和标准错误合并后的内容。
- `error`: 如果命令执行失败（如命令不存在、权限问题、命令返回非零退出码）， 或在执行过程中发生其他错误，则返回相应的错误信息。

### func ExecuteCmdWithTimeout

```go
func ExecuteCmdWithTimeout(timeout time.Duration, args []string, env []string) ([]byte, error)
```

ExecuteCmdWithTimeout 执行指定的系统命令，并设置超时时间及独立的环境变量。 此函数会等待命令执行完成，支持设置超时时间。

**参数:**
- `timeout`: 命令允许执行的最长时间。如果命令在此时间内未完成，将被终止并返回超时错误。 如果 timeout 为 0，则表示不设置超时。
- `args`: 命令行参数切片，其中 args[0] 为要执行的命令本身。
- `env`: 一个完整的环境变量切片，形如 "KEY=VALUE"。 如果传入 nil 或空切片，则命令将继承当前进程的环境变量。 如果传入非空切片，则命令的环境变量将仅限于此切片中定义的内容。

**返回:**
- `[]byte`: 命令的标准输出和标准错误合并后的内容。
- `error`: 如果命令执行失败、超时，或在执行过程中发生其他错误，则返回相应的错误信息。

### func FormatBytes

```go
func FormatBytes(bytes int64) string
```

FormatBytes 将字节数转换为人类可读的带单位的字符串 用于将字节数格式化为易读的存储单位格式，支持B到PB的转换

**参数:**
- `bytes`: 字节数（int64类型）

**返回:**
- `string`: 格式化后的字符串，如 "1.23 KB", "456.78 MB", "2.34 GB" 等

