package fs // import "gitee.com/MM-Q/go-kit/fs"


FUNCTIONS

func Collect(targetPath string, recursive bool) ([]string, error)
    Collect 收集指定路径下的所有文件 用于收集文件或目录中的文件，支持通配符匹配和递归遍历

    参数:
      - targetPath: 目标路径，支持通配符(*?[]{})
      - recursive: 是否递归遍历目录

    返回:
      - []string: 收集到的文件路径切片
      - error: 收集失败时返回错误

func CopyDir(src, dst string) error
    CopyDir 复制目录及其所有内容（默认覆盖已存在的文件） 用于递归复制整个目录，保持文件权限和目录结构

    参数:
      - src: 源目录路径
      - dst: 目标目录路径

    返回:
      - error: 复制失败时返回错误

func CopyDirWithOverwrite(src, dst string, overwrite bool) error
    CopyDirWithOverwrite 复制目录及其所有内容（可控制是否覆盖） 用于递归复制整个目录，保持文件权限和目录结构

    参数:
      - src: 源目录路径
      - dst: 目标目录路径
      - overwrite: 是否允许覆盖已存在的文件，false时如果目标目录或文件存在则返回错误

    返回:
      - error: 复制失败时返回错误

func CopyFile(src, dst string) error
    CopyFile 复制文件并继承权限（默认覆盖已存在的目标文件） 用于安全地复制文件，保持原文件的权限信息，失败时自动清理

    参数:
      - src: 源文件路径
      - dst: 目标文件路径

    返回:
      - error: 复制失败时返回错误

func CopyFileWithOverwrite(src, dst string, overwrite bool) error
    CopyFileWithOverwrite 复制文件并继承权限（可控制是否覆盖） 用于安全地复制文件，保持原文件的权限信息，失败时自动清理

    参数:
      - src: 源文件路径
      - dst: 目标文件路径
      - overwrite: 是否允许覆盖已存在的目标文件，false时如果目标文件存在则返回错误

    返回:
      - error: 复制失败时返回错误

func Exists(path string) bool
    Exists 检查指定路径的文件或目录是否存在 用于验证文件系统中指定路径是否存在，权限错误等异常情况视为不存在

    参数:
      - path: 要检查的路径

    返回:
      - bool: 文件或目录存在返回true，否则返回false

func GetDefaultBinPath() string
    GetDefaultBinPath 返回默认bin路径 用于获取Go程序的默认bin路径，采用多级回退策略确保总能返回有效路径

    返回:
      - string: 默认bin路径，优先级为GOPATH/bin > 用户主目录/go/bin > 当前工作目录/bin

func GetExecutablePath() string
    GetExecutablePath 获取程序的绝对安装路径 用于获取当前可执行文件的绝对路径，提供多级降级策略确保总能返回路径

    返回:
      - string: 程序的绝对路径，失败时降级为相对路径

func GetSize(path string) (int64, error)
    GetSize 获取文件或目录的大小 用于计算文件或目录的总字节数，目录会递归计算所有普通文件的大小

    参数:
      - path: 文件或目录路径

    返回:
      - int64: 文件或目录的总大小(字节)
      - error: 路径不存在或访问失败时返回错误

func GetUserHomeDir() string
    GetUserHomeDir 获取用户家目录 用于获取用户家目录路径，提供多级降级策略确保总能返回有效路径

    返回:
      - string: 用户家目录路径，失败时依次降级为工作目录或当前目录

func IsDir(path string) bool
    IsDir 检查指定路径是否为目录 用于验证指定路径是否为目录

    参数:
      - path: 要检查的路径

    返回:
      - bool: 是目录返回true，否则返回false

func IsFile(path string) bool
    IsFile 检查指定路径是否为文件 用于验证指定路径是否为普通文件

    参数:
      - path: 要检查的路径

    返回:
      - bool: 是文件返回true，否则返回false

func IsHidden(path string) bool
    IsHidden 判断文件或目录是否为隐藏 用于跨平台检查文件或目录的隐藏属性

    参数:
      - path: 文件或目录路径

    返回:
      - bool: 文件为隐藏返回true，否则返回false

func IsOctPerm(permission string) bool
    IsOctPerm 检查输入的权限是否是合法的4位八进制数 用于验证权限字符串格式是否符合八进制权限规范

    参数:
      - permission: 输入的权限字符串，例如 "0755" 或 "0644"

    返回:
      - bool: 权限格式合法返回true，否则返回false

func IsReadOnly(path string) bool
    IsReadOnly 判断文件或目录是否为只读 用于跨平台检查文件或目录的只读属性

    参数:
      - path: 文件或目录路径

    返回:
      - bool: 文件为只读返回true，否则返回false

func OctStrToMode(octalStr string) (os.FileMode, error)
    OctStrToMode 将4位八进制字符串权限转换为 os.FileMode 类型 用于将八进制权限字符串转换为Go标准库的文件权限类型

    参数:
      - octalStr: 4位八进制字符串，例如 "0755" 或 "0644"

    返回:
      - os.FileMode: 转换后的文件权限
      - error: 输入不合法时返回错误

    示例:

        mode, err := OctStrToMode("0755")
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(mode) // 输出: -rwxr-xr-x

