<div align="center">

# 🛠️ Go-Kit

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge)](https://gitee.com/MM-Q/go-kit)

**一个功能丰富的Go语言工具包，提供常用的实用工具和组件**

[📖 文档](https://gitee.com/MM-Q/go-kit) • [🚀 快速开始](#-快速开始) • [📚 模块说明](#-模块说明) • [🤝 贡献](#-贡献指南)

</div>

---

## 📋 项目简介

Go-Kit 是一个精心设计的Go语言高性能工具包，集成了开发过程中常用的实用工具和组件。该项目采用分层架构设计，核心模块零外部依赖，旨在提高开发效率，减少重复代码，为Go开发者提供一套完整的基础工具集。

### 架构特点

- **分层设计**：`pool` 作为第0层基础模块，`fs`/`hash`/`id` 依赖 `pool`；其他模块独立
- **零外部依赖**：核心模块（pool/fs/hash/id/str/utils/fuzzy）仅使用Go标准库
- **高性能**：对象池复用减少GC压力，动态缓冲区计算，原子性文件操作
- **跨平台**：支持 Windows/Linux/macOS，跨平台属性检查适配

## ✨ 核心特性

- 🏊 **对象池** - 5种高性能对象池（Byte/Buffer/String/Rand/Timer），减少GC压力
- 🗂️ **文件系统工具** - 文件操作、路径处理、原子性复制、目录遍历、跨平台属性检查
- 🔐 **哈希工具** - MD5/SHA1/SHA256/SHA512多算法支持，文件/内存/流式哈希
- 🆔 **ID生成器** - UUID、时间戳+随机数、类UUID格式、纯随机字符串
- 📝 **字符串工具** - 安全解引用、字符串构建、截取、填充、掩码、模板替换
- 🖥️ **终端工具** - 输入读取、确认框、密码输入、基础菜单、结构化菜单
- 🔍 **模糊匹配** - 智能模糊搜索，支持评分排序，适用于文件名/代码符号搜索
- 🔌 **SSH工具** - 多主机批量操作、连通性测试、命令执行
- ⚙️ **通用工具** - 字节格式化、JSON转义等辅助函数
- 🎯 **零外部依赖** - 核心模块（pool/fs/hash/id/str/utils/fuzzy）零依赖

## 🚀 快速开始

### 安装

使用 `go get` 命令从Gitee安装：

```bash
go get gitee.com/MM-Q/go-kit
```

## 📚 模块说明

### 核心模块

| 模块 | 层级 | 描述 | 主要功能 | 外部依赖 |
|------|------|------|----------|----------|
| `pool` | 基础层 | 对象池 | Byte/Buf/Str/Rand/Timer池、动态缓冲区计算 | 无 |
| `fs` | 功能层 | 文件系统工具 | 文件操作、原子复制、目录遍历、跨平台属性 | 无 |
| `hash` | 功能层 | 哈希工具 | MD5/SHA1/SHA256/SHA512、文件/流式哈希 | progressbar |
| `id` | 功能层 | ID生成器 | UUID、时间戳ID、类UUID、随机字符串 | 无 |
| `str` | 独立 | 字符串工具 | 解引用、构建、截取、填充、掩码、模板 | 无 |
| `utils` | 独立 | 通用工具 | 字节格式化、JSON转义 | 无 |
| `term` | 独立 | 终端工具 | 输入读取、确认框、密码、菜单 | golang.org/x/term |
| `fuzzy` | 独立 | 模糊匹配 | 智能模糊搜索、评分排序、高亮匹配 | 无 |
| `easyssh` | 独立 | SSH工具 | 多主机操作、连通性测试、批量执行 | golang.org/x/crypto/ssh |

### 模块详细说明

#### 🏊 pool - 对象池（基础层）

`pool` 是整个工具库的基础模块，其他模块依赖它来实现高性能操作。

- **BytePool**：字节切片对象池，复用 `[]byte` 减少内存分配
- **BufPool**：`bytes.Buffer` 对象池，用于字符串构建
- **StrPool**：`strings.Builder` 对象池，高效的字符串拼接
- **RandPool**：随机数生成器池，复用 `rand.Rand` 实例
- **TimerPool**：定时器池，复用 `time.Timer` 减少GC压力
- **动态缓冲区计算**：`CalculateBufferSize` 根据文件大小智能选择1KB-2MB缓冲区

#### 📁 fs - 文件系统工具

提供完整的文件系统操作功能：

- **路径检查**：`Exists`/`IsFile`/`IsDir` 快速判断路径类型
- **原子性复制**：`Copy`/`CopyEx` 使用临时文件+`os.Rename` 保证原子性
- **目录遍历**：`Collect` 支持通配符的目录遍历
- **跨平台支持**：`attr_unix.go`/`attr_windows.go` 实现跨平台属性检查
- **安全特性**：原子操作、备份恢复、覆盖控制

#### 🔐 hash - 哈希工具

提供多种哈希算法的便捷封装：

- **常用算法**：MD5、SHA1、SHA256、SHA512
- **多输入类型**：文件哈希、内存数据哈希、流式哈希（`io.Reader`）
- **统一接口**：所有算法使用相同的API，方便切换
- **进度显示**：支持进度条显示（依赖 `progressbar` 库）

#### 🆔 id - ID生成器

提供多种ID生成方案：

- **GenID**：时间戳+随机数组合ID
- **UUID**：标准UUID v4格式（8-4-4-4-12共36位）
- **GenMaskedID**：隐藏时间戳ID
- **纯随机字符串**：基于 `crypto/rand` 保证安全性
- **依赖**：使用 `pool` 模块的随机数池和字符串构建器池

#### 📝 str - 字符串工具

丰富的字符串处理功能：

- **安全解引用**：`SafeDeref` 安全解引用字符串指针
- **字符串构建**：`BuildStr`/`BuildStrCap` 高效构建字符串
- **截取**：`Prefix`/`Suffix`/`Truncate` 字符串截取
- **填充**：`PadLeft`/`PadRight` 字符串填充
- **掩码**：`Mask` 字符串掩码（如手机号脱敏）
- **模板替换**：`Template` 简单的模板变量替换
- **Base64**：`Base64Encode`/`Base64Decode` 编解码

#### 🖥️ term - 终端工具

提供完整的终端交互功能：

- **输入读取**：字符串、整数、浮点数输入，支持默认值
- **确认框**：`Confirm` 支持默认值的确认对话框
- **密码输入**：`ReadPassword` 安全密码输入（无回显）
- **基础菜单**：`BasicMenu` 简单数字选项菜单，自动编号
- **结构化菜单**：`Menu` 支持自定义键、默认值、退出选项的复杂菜单
- **菜单样式**：可自定义前缀、分隔符、缩进等样式
- **循环菜单**：支持循环显示和交互处理

#### 🔍 fuzzy - 模糊匹配

智能模糊字符串匹配库：

- **评分系统**：首字符匹配(+10)、驼峰匹配(+20)、分隔符后匹配(+20)、相邻匹配(+5)
- **惩罚机制**：未匹配前导字符(-5)、未匹配字符惩罚
- **多种接口**：`Find`/`FindNoSort`/`FindFrom`/`FindFromNoSort`
- **高亮支持**：返回匹配字符索引，支持结果高亮显示
- **Unicode支持**：完整支持多语言字符（大小写不敏感）
- **命令行补全**：`Complete`/`CompletePrefix`/`CompleteExact` 专为命令行标志补全设计
  - 优先前缀匹配，其次模糊匹配
  - 精确匹配(1000分) > 前缀匹配(100-200分) > 模糊匹配(0-99分)
  - 适用于 CLI 工具自动补全场景

#### 🔌 easyssh - SSH工具

简化的多主机SSH连接和命令执行：

- **主机配置**：从配置文件解析主机列表（支持3字段和4字段格式）
- **连通性测试**：`PingHosts` 测试多主机SSH连接状态
- **批量执行**：`Exec` 在多主机上执行相同命令
- **回调支持**：`ExecWithCallback` 支持自定义回调处理输出
- **超时控制**：支持设置命令执行超时

#### ⚙️ utils - 通用工具

其他实用函数：

- **字节格式化**：`FormatBytes` 格式化字节为 B/PB 可读格式
- **JSON转义**：`QuoteBytes`/`QuoteString` JSON字符串转义

## 🧪 测试说明

### 运行所有测试

```bash
go test ./...
```

### 运行特定模块测试

```bash
go test ./pool
go test ./fs
go test ./hash
go test ./id
go test ./str
go test ./term
go test ./fuzzy
go test ./utils
go test ./easyssh
```

### 生成测试覆盖率报告

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 基准测试

```bash
go test -bench=. ./...
```

## 🔧 开发环境

### 要求

- Go 1.25 或更高版本
- Git

### 本地开发

1. 克隆仓库：
```bash
git clone https://gitee.com/MM-Q/go-kit.git
# 或者从GitHub镜像克隆
git clone https://github.com/MM-Q/go-kit.git
```

2. 进入项目目录：
```bash
cd go-kit
```

3. 安装依赖：
```bash
go mod tidy
```

4. 运行测试：
```bash
go test ./...
```

## 🤝 贡献指南

我们欢迎所有形式的贡献！请遵循以下步骤：

1. **Fork** 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 **Pull Request**

### 代码规范

- 遵循 Go 官方代码规范
- 添加适当的注释和文档
- 确保所有测试通过
- 保持测试覆盖率在 80% 以上

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 联系方式

- 🐛 问题反馈: [Gitee Issues](https://gitee.com/MM-Q/go-kit/issues)
- 💬 讨论: [Gitee 评论区](https://gitee.com/MM-Q/go-kit)
- 🔗 GitHub镜像: [GitHub Repository](https://github.com/MM-Q/go-kit)

## 🔗 相关链接

- 📖 [在线文档](https://pkg.go.dev/gitee.com/MM-Q/go-kit)
- 🏠 [项目主页](https://gitee.com/MM-Q/go-kit)
- 🪞 [GitHub镜像](https://github.com/MM-Q/go-kit)
- 📊 [Go Report Card](https://goreportcard.com/report/gitee.com/MM-Q/go-kit)

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给它一个星标！**

Made with ❤️ by [MMQ](https://gitee.com/MM-Q)

</div>
