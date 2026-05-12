# Expand 文件展开函数设计方案

## 概述

为命令行工具提供简洁的通配符展开功能，类似于 shell 的 glob 展开。

---

## 核心设计理念

1. **简洁实用**: 只做一件事 - 展开通配符
2. **兼容 shell 行为**: 行为类似于 bash 的 glob 展开
3. **优雅降级**: 没有匹配时保留原路径（允许后续检查）
4. **去重**: 自动去除重复路径

---

## 函数设计

### 1. Expand

```go
// Expand 展开文件路径列表（支持通配符）
// 将包含通配符的模式展开为具体路径，返回所有匹配的文件和目录
//
// 参数:
//   - patterns: 文件路径模式列表，支持 *、?、[]、** 等通配符
//
// 返回值:
//   - []string: 展开后的路径列表（去重、排序），包含文件和目录
//   - error: 模式语法错误时返回错误
//
// 行为说明:
//   - 通配符匹配成功: 返回匹配的所有路径（文件和目录）
//   - 通配符无匹配: 保留原模式（不报错，由调用者处理）
//   - 具体路径: 直接保留
//
// 示例:
//   paths, err := fs.Expand([]string{"*.go"})                    // [main.go utils.go]
//   paths, err := fs.Expand([]string{"**/*.go"})                 // 递归匹配所有 .go 文件
//   paths, err := fs.Expand([]string{"src/*"})                   // [src/main.go src/utils src/pkg]
//   paths, err := fs.Expand([]string{"config.yaml"})             // [config.yaml]（原样保留）
//   paths, err := fs.Expand([]string{"*.notexist"})              // [*.notexist]（无匹配时保留）
func Expand(patterns []string) ([]string, error)
```

---

### 2. ExpandFiles

```go
// ExpandFiles 展开文件路径列表，只返回文件（排除目录）
// 与 Expand 类似，但自动过滤掉目录路径
//
// 参数:
//   - patterns: 文件路径模式列表，支持 *、?、[]、** 等通配符
//
// 返回值:
//   - []string: 展开后的文件路径列表（去重、排序），只包含文件
//   - error: 模式语法错误时返回错误
//
// 示例:
//   files, err := fs.ExpandFiles([]string{"*.go"})               // [main.go utils.go]
//   files, err := fs.ExpandFiles([]string{"src/*"})              // [src/main.go]（排除 src/utils 目录）
//   files, err := fs.ExpandFiles([]string{"**/*.go"})            // 递归匹配所有 .go 文件
func ExpandFiles(patterns []string) ([]string, error)
```

---

### 3. ExpandPattern

```go
// ExpandPattern 展开单个路径模式
// 与 Expand 类似，但只接收单个模式，返回所有匹配的路径（文件和目录）
//
// 参数:
//   - pattern: 文件路径模式，支持 *、?、[]、** 等通配符
//
// 返回值:
//   - []string: 展开后的路径列表（排序），包含文件和目录
//   - error: 模式语法错误时返回错误
//
// 示例:
//   paths, err := fs.ExpandPattern("*.go")                 // [main.go utils.go]
//   paths, err := fs.ExpandPattern("src/*")                // [src/main.go src/utils src/pkg]
//   paths, err := fs.ExpandPattern("config.yaml")          // [config.yaml]（原样保留）
//   paths, err := fs.ExpandPattern("*.notexist")           // [*.notexist]（无匹配时保留）
func ExpandPattern(pattern string) ([]string, error)
```

---

### 4. Match

```go
// Match 检查路径是否匹配模式（支持 ** 双星号）
//
// 参数:
//   - pattern: 通配符模式，支持 ** 匹配任意层级目录
//   - path: 要检查的路径
//
// 返回值:
//   - bool: 是否匹配
//   - error: 模式语法错误时返回错误
//
// 示例:
//   matched, _ := fs.Match("**/*.go", "src/main.go")      // true
//   matched, _ := fs.Match("*.go", "src/main.go")         // false
func Match(pattern, path string) (bool, error)
```

---

## 使用场景

### 场景 1: 命令行参数展开（所有路径）
```go
// 命令行: mytool *.go src/*
patterns := os.Args[1:]
paths, err := fs.Expand(patterns)
// paths 可能包含: [main.go utils.go src/main.go src/utils(目录)]
```

### 场景 2: 只处理文件（排除目录）
```go
// 只获取文件，排除目录
files, err := fs.ExpandFiles([]string{"src/*"})
// 结果只包含文件，不包含 src/utils 这样的目录
```

### 场景 3: 单个模式展开
```go
// 展开单个模式
paths, err := fs.ExpandPattern("*.go")

// 检查是否无匹配
if len(paths) == 1 && paths[0] == "*.go" {
    // 无匹配，保留原模式
    log.Println("No .go files found")
}
```

### 场景 4: 配置文件中的通配符
```go
// 配置文件: sources = ["src/**/*.go", "tests/*_test.go"]
config := loadConfig()
files, err := fs.ExpandFiles(config.Sources)
```

---

## 实现要点

```go
func Expand(patterns []string) ([]string, error) {
    seen := make(map[string]bool)
    var result []string

    for _, pattern := range patterns {
        matches, err := expandPattern(pattern)
        if err != nil {
            return nil, err
        }

        for _, match := range matches {
            if !seen[match] {
                seen[match] = true
                result = append(result, match)
            }
        }
    }

    sort.Strings(result)
    return result, nil
}

func ExpandFiles(patterns []string) ([]string, error) {
    paths, err := Expand(patterns)
    if err != nil {
        return nil, err
    }

    var files []string
    for _, path := range paths {
        if !isDir(path) {
            files = append(files, path)
        }
    }

    return files, nil
}
```
