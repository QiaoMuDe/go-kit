# 命令行补全搜索器设计方案（简洁版）

## 设计原则

1. **参考 fuzzy 包**：一个核心函数 + 几个便捷函数
2. **返回类型一致**：使用 `[]Match` 切片
3. **单文件实现**：所有代码放在一个文件
4. **前缀优先**：命令行补全场景，前缀匹配权重最高

---

## 核心函数

```go
// Complete 核心补全函数
// 优先前缀匹配，其次模糊匹配
// 返回按匹配质量排序的结果
func Complete(pattern string, candidates []string) Matches

// CompletePrefix 仅前缀匹配
// 只返回以 pattern 开头的候选
func CompletePrefix(pattern string, candidates []string) Matches

// CompleteExact 仅精确匹配
// 只返回完全相等的候选
func CompleteExact(pattern string, candidates []string) Matches
```

---

## 评分规则

| 匹配类型 | 基础分数 | 说明 |
|----------|----------|------|
| 精确匹配 | 1000 | 完全相等 |
| 前缀匹配 | 100 + 长度奖励 | 开头匹配，越短分数越高 |
| 模糊匹配 | 0-99 | 复用 fuzzy 原有评分 |

---

## 实现要点

```go
// Complete 实现逻辑
func Complete(pattern string, candidates []string) Matches {
    var matches Matches
    
    for i, candidate := range candidates {
        score, matchedIndexes := matchComplete(pattern, candidate)
        if score > 0 {
            matches = append(matches, Match{
                Str:            candidate,
                Index:          i,
                Score:          score,
                MatchedIndexes: matchedIndexes,
            })
        }
    }
    
    // 按分数降序排序
    sort.Stable(matches)
    return matches
}

// matchComplete 判断匹配类型并评分
func matchComplete(pattern, candidate string) (int, []int) {
    // 1. 精确匹配
    if pattern == candidate {
        return 1000, getAllIndexes(candidate)
    }
    
    // 2. 前缀匹配（不区分大小写）
    if strings.HasPrefix(strings.ToLower(candidate), strings.ToLower(pattern)) {
        // 越短的前缀匹配分数越高（优先短候选）
        score := 100 + (100 - len(candidate))
        if score < 100 {
            score = 100
        }
        return score, getPrefixIndexes(len(pattern))
    }
    
    // 3. 模糊匹配
    if m := fuzzyMatch(pattern, candidate); m.Score > 0 {
        return m.Score / 10, m.MatchedIndexes // 模糊匹配分数降低
    }
    
    return 0, nil
}
```

---

## 使用示例

```go
flags := []string{
    "--verbose",
    "--version", 
    "--help",
    "-v",
    "-h",
}

// 模糊补全
matches := fuzzy.Complete("--v", flags)
// 结果: [--verbose, --version, -v]（前缀匹配优先）

// 仅前缀
matches := fuzzy.CompletePrefix("--v", flags)
// 结果: [--verbose, --version]

// 仅精确
matches := fuzzy.CompleteExact("--verbose", flags)
// 结果: [--verbose]
```

---

## 文件结构

```
fuzzy/
├── types.go      # 现有
├── fuzzy.go      # 现有
└── complete.go   # 新增：补全搜索器
```

---

## 与现有 fuzzy 的关系

```go
// complete.go 复用 fuzzy.go 的算法
// 不修改现有代码，只新增函数
```
