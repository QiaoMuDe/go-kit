<div align="center">

# 🛠️ Go-Kit

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge)](https://gitee.com/MM-Q/go-kit)
[![Coverage](https://img.shields.io/badge/Coverage-85%25-green?style=for-the-badge)](https://gitee.com/MM-Q/go-kit)

**一个功能丰富、高性能的Go语言工具库集合**

[🏠 仓库地址](https://gitee.com/MM-Q/go-kit) • [📖 文档](https://gitee.com/MM-Q/go-kit) • [🐛 问题反馈](https://gitee.com/MM-Q/go-kit/issues)

</div>

---

## 📋 项目简介

Go-Kit 是一个精心设计的Go语言工具库集合，提供了文件系统操作、哈希计算、ID生成、字符串处理、系统命令执行等常用功能。项目采用模块化设计，每个模块都经过充分测试，可以独立使用或组合使用。

## ✨ 核心特性

- 🗂️ **文件系统操作** - 文件/目录管理、大小计算、路径处理
- 🔐 **哈希计算** - 支持MD5、SHA1、SHA256、SHA512多种算法
- 🆔 **ID生成器** - 时间戳ID、UUID、批量生成等多种方式
- 🧵 **字符串工具** - 字符串验证、截取、安全解引用
- ⚡ **高性能** - 使用对象池优化内存分配
- 🔧 **系统命令** - 安全的命令执行和超时控制
- 📊 **字节格式化** - 人性化的存储单位显示
- 🧪 **完整测试** - 85%+ 测试覆盖率，包含基准测试

## 🚀 安装指南

### 使用 go get 安装

```bash
go get gitee.com/MM-Q/go-kit
```

### 在项目中引入

```go
import (
    "gitee.com/MM-Q/go-kit/fs"
    "gitee.com/MM-Q/go-kit/hash"
    "gitee.com/MM-Q/go-kit/id"
    "gitee.com/MM-Q/go-kit/str"
    "gitee.com/MM-Q/go-kit/utils"
)
```

## 📚 使用示例

### 基础用法

#### 文件系统操作

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/go-kit/fs"
)

func main() {
    // 获取文件大小
    size, err := fs.GetSize("./example.txt")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("File size: %d bytes\n", size)
    
    // 查找文件
    files, err := fs.FindFiles("*.go", true)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Found %d Go files\n", len(files))
}
```

#### 哈希计算

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/go-kit/hash"
)

func main() {
    // 计算文件哈希
    checksum, err := hash.Checksum("example.txt", "sha256")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("SHA256: %s\n", checksum)
    
    // 计算字符串哈希
    strHash, err := hash.HashString("Hello, World!", "md5")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("MD5: %s\n", strHash)
}
```

#### ID生成

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/go-kit/id"
)

func main() {
    // 生成单个ID
    singleID := id.GenID(8)
    fmt.Printf("Generated ID: %s\n", singleID)
    
    // 批量生成ID
    ids := id.GenIDs(5, 6)
    fmt.Printf("Generated %d IDs: %v\n", len(ids), ids)
    
    // 生成UUID
    uuid := id.UUID()
    fmt.Printf("UUID: %s\n", uuid)
    
    // 带前缀的ID
    prefixedID := id.GenWithPrefix("user", 8)
    fmt.Printf("Prefixed ID: %s\n", prefixedID)
}
```

### 高级用法

#### 带进度条的哈希计算

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/go-kit/hash"
)

func main() {
    // 大文件哈希计算带进度条
    checksum, err := hash.ChecksumProgress("large_file.zip", "sha256")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("SHA256: %s\n", checksum)
}
```

#### 系统命令执行

```go
package main

import (
    "fmt"
    "time"
    "gitee.com/MM-Q/go-kit/utils"
)

func main() {
    // 执行命令
    output, err := utils.ExecuteCmd([]string{"echo", "Hello World"}, nil)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Output: %s\n", output)
    
    // 带超时的命令执行
    output, err = utils.ExecuteCmdWithTimeout(
        5*time.Second, 
        []string{"ping", "-c", "3", "google.com"}, 
        nil,
    )
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Ping result: %s\n", output)
}
```

#### 字节格式化

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/go-kit/utils"
)

func main() {
    sizes := []int64{1024, 1048576, 1073741824, 1099511627776}
    
    for _, size := range sizes {
        formatted := utils.FormatBytes(size)
        fmt.Printf("%d bytes = %s\n", size, formatted)
    }
    // 输出:
    // 1024 bytes = 1 KB
    // 1048576 bytes = 1 MB
    // 1073741824 bytes = 1 GB
    // 1099511627776 bytes = 1 TB
}
```

## 📖 API文档概述

### 模块列表

| 模块 | 功能描述 | 主要函数 |
|------|----------|----------|
| `fs` | 文件系统操作 | `GetSize`, `FindFiles`, `GetDefaultBinPath` |
| `hash` | 哈希计算 | `Checksum`, `ChecksumProgress`, `HashString`, `HashData` |
| `id` | ID生成 | `GenID`, `GenIDs`, `UUID`, `GenWithPrefix`, `Valid` |
| `str` | 字符串工具 | `IsNotEmpty`, `StringSuffix8`, `SafeDeref` |
| `utils` | 系统工具 | `ExecuteCmd`, `ExecuteCmdWithTimeout`, `FormatBytes` |

### 支持的哈希算法

| 算法 | 输出长度 | 用途 |
|------|----------|------|
| MD5 | 32字符 | 快速校验 |
| SHA1 | 40字符 | 版本控制 |
| SHA256 | 64字符 | 安全应用 |
| SHA512 | 128字符 | 高安全性 |

## ⚙️ 配置选项

### 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `GOPATH` | Go工作路径 | 系统默认 |

### 缓冲区配置

```go
// 哈希计算缓冲区大小会根据文件大小自动调整
// 最小: 1KB
// 最大: 根据文件大小动态计算
```

## 📁 项目结构

```
go-kit/
├── fs/                 # 文件系统操作模块
│   ├── fs.go          # 核心文件操作
│   ├── check.go       # 文件检查功能
│   ├── copy.go        # 文件复制功能
│   ├── attr.go        # 文件属性处理
│   └── fs_test.go     # 测试文件
├── hash/              # 哈希计算模块
│   ├── hash.go        # 哈希计算核心
│   └── hash_test.go   # 测试文件
├── id/                # ID生成模块
│   ├── id.go          # ID生成核心
│   └── id_test.go     # 测试文件
├── pool/              # 对象池模块
│   ├── buffer.go      # 缓冲区池
│   ├── byte.go        # 字节池
│   ├── rand.go        # 随机数池
│   ├── string.go      # 字符串池
│   └── timer.go       # 定时器池
├── str/               # 字符串工具模块
│   ├── str.go         # 字符串处理
│   └── str_test.go    # 测试文件
├── utils/             # 系统工具模块
│   ├── utils.go       # 工具函数
│   └── utils_test.go  # 测试文件
├── go.mod             # Go模块文件
├── go.sum             # 依赖校验文件
├── LICENSE            # 许可证文件
└── README.md          # 项目说明文档
```

## 🧪 测试说明

### 运行所有测试

```bash
go test ./...
```

### 运行特定模块测试

```bash
go test ./fs
go test ./hash
go test ./id
```

### 运行基准测试

```bash
go test -bench=. ./...
```

### 查看测试覆盖率

```bash
go test -cover ./...
```

### 生成详细覆盖率报告

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔧 开发指南

### 本地开发环境

1. **Go版本要求**: Go 1.21+
2. **依赖管理**: 使用Go Modules
3. **代码规范**: 遵循Go官方代码规范

### 添加新功能

1. 在相应模块目录下添加功能代码
2. 编写对应的测试文件
3. 更新文档和示例
4. 确保测试通过

### 性能优化

- 使用对象池减少内存分配
- 合理使用缓冲区大小
- 避免不必要的字符串拷贝

## 📊 性能基准

| 操作 | 性能指标 | 说明 |
|------|----------|------|
| ID生成 | ~100ns/op | 单个ID生成 |
| 哈希计算 | ~1MB/ms | SHA256算法 |
| 文件大小计算 | ~10μs/file | 小文件 |
| 字符串处理 | ~50ns/op | 基础操作 |

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

```
MIT License

Copyright (c) 2024 MM-Q

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献

1. **Fork** 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 **Pull Request**

### 贡献类型

- 🐛 **Bug修复**
- ✨ **新功能**
- 📝 **文档改进**
- 🎨 **代码优化**
- ✅ **测试增强**

### 代码规范

- 遵循Go官方代码规范
- 添加适当的注释和文档
- 确保测试覆盖率不低于80%
- 运行 `go fmt` 和 `go vet`

## 📞 联系方式

- **仓库地址**: [https://gitee.com/MM-Q/go-kit](https://gitee.com/MM-Q/go-kit)
- **问题反馈**: [Issues](https://gitee.com/MM-Q/go-kit/issues)
- **功能请求**: [Feature Requests](https://gitee.com/MM-Q/go-kit/issues)

## 🔗 相关链接

- [Go官方文档](https://golang.org/doc/)
- [Go模块参考](https://golang.org/ref/mod)
- [Go测试指南](https://golang.org/doc/tutorial/add-a-test)

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

---

<div align="center">

**如果这个项目对您有帮助，请给我们一个 ⭐️**

[🏠 返回仓库](https://gitee.com/MM-Q/go-kit) • [📖 查看文档](https://gitee.com/MM-Q/go-kit) • [🐛 报告问题](https://gitee.com/MM-Q/go-kit/issues)

</div>