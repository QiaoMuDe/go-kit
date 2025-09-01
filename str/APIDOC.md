# Package str

```go
import "gitee.com/MM-Q/go-kit/str"
```

## Functions

### func IsNotEmpty

```go
func IsNotEmpty(s string) bool
```

IsNotEmpty 检查字符串是否不为空 用于验证字符串在去除首尾空格后是否包含有效内容

**参数:**
- `s`: 待检查的字符串

**返回:**
- `bool`: 字符串不为空返回true，否则返回false

### func SafeDeref

```go
func SafeDeref(s *string) string
```

SafeDeref 安全地解引用字符串指针 用于安全地获取字符串指针的值，避免空指针异常

**参数:**
- `s`: 字符串指针

**返回:**
- `string`: 解引用后的字符串，指针为nil时返回空字符串

### func StringSuffix8

```go
func StringSuffix8(s string) string
```

StringSuffix8 从给定字符串中获取最后8个字符。 如果字符串长度小于等于8，则返回原字符串。

**参数:**
- `s`: 输入字符串

**返回:**
- `string`: 字符串的最后8个字符，或原字符串（如果长度不足8），或空字符串（如果输入为空）

