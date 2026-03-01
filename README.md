<div align="center">

# 🛠️ Go-Kit

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge)](https://gitee.com/MM-Q/go-kit)

**一个功能丰富的Go语言工具包，提供常用的实用工具和组件**

[📖 文档](https://gitee.com/MM-Q/go-kit) • [🚀 快速开始](#-快速开始) • [� 模块说明](#-模块说明) • [🤝 贡献](#-贡献指南)

</div>

---

## 📋 项目简介

Go-Kit 是一个精心设计的Go语言工具包，集成了开发过程中常用的实用工具和组件。该项目旨在提高开发效率，减少重复代码，为Go开发者提供一套完整的基础工具集。

## ✨ 核心特性

- 🗂️ **文件系统工具** - 提供文件操作、路径处理、复制移动等实用功能
- 🔐 **哈希工具** - 支持多种哈希算法的便捷封装
- 🆔 **ID生成器** - 提供UUID、雪花算法等ID生成方案
- 🏊 **对象池** - 高性能的对象池实现，优化内存使用
- 📝 **字符串工具** - 丰富的字符串处理和操作函数
- �️ **终端工具** - 完整的终端输入、菜单显示和用户交互功能
- �🔧 **通用工具** - 其他常用的辅助函数和工具
- 🔌 **SSH工具** - 简化的多主机SSH连接和命令执行工具

## 🚀 快速开始

### 安装

使用 `go get` 命令从Gitee安装：

```bash
go get gitee.com/MM-Q/go-kit
```

## 📚 模块说明

### 核心模块

| 模块 | 描述 | 主要功能 |
|------|------|----------|
| `fs` | 文件系统工具 | 文件操作、路径处理、复制移动、目录管理 |
| `hash` | 哈希工具 | MD5、SHA1、SHA256等哈希算法 |
| `id` | ID生成器 | UUID、雪花算法、随机ID |
| `pool` | 对象池 | 高性能对象复用池、字节缓冲池 |
| `str` | 字符串工具 | 字符串转换、格式化、验证、处理 |
| `term` | 终端工具 | 终端输入、菜单显示、用户交互 |
| `utils` | 通用工具 | 其他实用函数、JSON工具 |
| `esayssh` | SSH工具 | 多主机SSH连接、批量执行命令、连通性测试 |

### 模块详细说明

#### � fs - 文件系统工具

提供完整的文件系统操作功能，包括：

- **文件操作**：复制、移动、删除文件
- **目录操作**：创建、复制、移动、遍历目录
- **路径处理**：路径验证、绝对路径解析、智能路径处理
- **文件检查**：存在性检查、类型判断、权限检查、隐藏文件判断
- **文件信息**：获取文件大小、权限、修改时间等
- **特殊支持**：符号链接、特殊文件（设备文件、命名管道等）
- **安全特性**：原子性操作、备份恢复机制、覆盖控制

#### 🔐 hash - 哈希工具

提供多种哈希算法的便捷封装：

- **常用算法**：MD5、SHA1、SHA256、SHA512
- **便捷接口**：字符串哈希、文件哈希、字节切片哈希
- **统一接口**：所有算法使用相同的API，方便切换

#### 🆔 id - ID生成器

提供多种ID生成方案：

- **UUID**：标准UUID v4生成
- **雪花算法**：分布式唯一ID生成
- **随机ID**：基于随机数的ID生成
- **自定义格式**：支持自定义ID格式

#### 🏊 pool - 对象池

高性能的对象池实现：

- **字节缓冲池**：复用字节缓冲区，减少内存分配
- **自动计算**：根据文件大小自动选择合适的缓冲区大小
- **线程安全**：支持并发使用
- **性能优化**：减少GC压力，提高性能

#### 📝 str - 字符串工具

丰富的字符串处理功能：

- **转换**：大小写转换、类型转换、编码转换
- **格式化**：字符串格式化、填充、截断
- **验证**：格式验证、长度验证、内容验证
- **处理**：分割、连接、替换、去除空白
- **工具**：随机字符串生成、字符串比较、子串查找

#### �️ term - 终端工具

提供完整的终端交互功能，包括：

- **基础输入**：字符串、整数、浮点数输入，支持默认值
- **确认框**：支持默认值的确认对话框
- **密码输入**：安全的密码输入，不显示回显
- **基础菜单**：简单的数字选项菜单，自动编号
- **结构化菜单**：支持自定义键、默认值、退出选项的复杂菜单
- **菜单样式**：可自定义前缀、分隔符、缩进等样式
- **循环菜单**：支持循环显示和交互处理
- **输入验证**：严格的输入验证和错误处理
- **依赖注入**：支持自定义输入源，便于测试

#### � utils - 通用工具

其他实用函数：

- **JSON工具**：JSON格式化、压缩、验证
- **时间工具**：时间格式化、转换、计算
- **系统工具**：获取系统信息、环境变量
- **其他**：其他常用的辅助函数

#### 🔌 esayssh - SSH工具

简化的多主机SSH连接和命令执行：

- **多主机管理**：支持从文件读取主机列表
- **批量执行**：在多个主机上执行相同命令
- **连通性测试**：测试主机SSH连接状态
- **回调支持**：支持自定义回调处理输出
- **超时控制**：支持设置命令执行超时
- **格式化输出**：支持格式化输出命令结果

## 📁 项目结构

```
go-kit/
├── fs/                 # 文件系统工具
│   ├── copy.go         # 文件复制功能
│   ├── move.go         # 文件移动功能
│   ├── fs.go           # 核心文件操作
│   ├── check.go        # 文件检查功能
│   ├── attr.go         # 文件属性操作
│   ├── path.go         # 路径处理功能
│   ├── collect.go      # 文件收集功能
│   ├── APIDOC.md       # API文档
│   └── *_test.go       # 单元测试
├── hash/              # 哈希工具
│   ├── hash.go         # 哈希算法实现
│   └── hash_test.go    # 单元测试
├── id/                # ID生成器
│   ├── id.go           # ID生成实现
│   └── id_test.go      # 单元测试
├── pool/              # 对象池
│   ├── pool.go         # 对象池实现
│   └── pool_test.go    # 单元测试
├── str/               # 字符串工具
│   ├── str.go          # 字符串处理
│   └── str_test.go     # 单元测试
├── term/              # 终端工具
│   ├── read.go         # 基础输入功能
│   ├── menu.go         # 结构化菜单功能
│   ├── basic_menu.go   # 基础菜单功能
│   ├── APIDOC.md       # API文档
│   └── *_test.go       # 单元测试
├── utils/             # 通用工具
│   ├── utils.go        # 实用函数
│   ├── utils_test.go   # 单元测试
│   ├── json.go         # JSON 工具
│   └── json_test.go    # 单元测试
├── esayssh/           # SSH工具
│   ├── esayssh.go      # 核心SSH连接和命令执行
│   ├── ssh.go          # SSH协议实现
│   ├── types.go        # 数据结构定义
│   └── APIDOC.md       # API文档
├── go.mod             # Go模块文件
├── go.sum             # 依赖校验文件
├── LICENSE            # 许可证文件
└── README.md          # 项目说明
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
go test ./pool
go test ./str
go test ./term
go test ./utils
go test ./esayssh
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
