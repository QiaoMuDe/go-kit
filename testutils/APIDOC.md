package testutil // import "gitee.com/MM-Q/go-kit/testutils"


TYPES

type RunFunc func(m *testing.M) int
    RunFunc 测试运行函数类型

func Quiet() RunFunc
    Quiet 创建静默测试运行函数 用于在非verbose模式下抑制测试输出，verbose模式下正常输出

    返回:
      - RunFunc: 测试运行函数，用于TestMain中执行测试

    示例:

        func TestMain(m *testing.M) {
            run := testutil.QuietUnlessVerbose()
            os.Exit(run(m))
        }

