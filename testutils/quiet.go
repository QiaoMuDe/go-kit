package testutil

import (
	"flag"
	"os"
	"testing"
)

// RunFunc 测试运行函数类型
type RunFunc func(m *testing.M) int

// Quiet 创建静默测试运行函数
// 用于在非verbose模式下抑制测试输出，verbose模式下正常输出
//
// 返回:
//   - RunFunc: 测试运行函数，用于TestMain中执行测试
//
// 示例:
//
//	func TestMain(m *testing.M) {
//	    run := testutil.QuietUnlessVerbose()
//	    os.Exit(run(m))
//	}
func Quiet() RunFunc {
	flag.Parse() // 让 -test.v 等参数先被解析

	var (
		restoreStdout, restoreStderr func()
	)

	var nullFile *os.File

	if !testing.Verbose() {
		var err error
		nullFile, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0o666)
		if err != nil {
			panic("testutil: open /dev/null: " + err.Error())
		}

		origOut, origErr := os.Stdout, os.Stderr
		os.Stdout, os.Stderr = nullFile, nullFile

		restoreStdout = func() { os.Stdout = origOut }
		restoreStderr = func() { os.Stderr = origErr }
	}

	return func(m *testing.M) int {
		defer func() {
			if restoreStdout != nil {
				restoreStdout()
			}
			if restoreStderr != nil {
				restoreStderr()
			}
			if nullFile != nil {
				nullFile.Close()
			}
		}()

		return m.Run()
	}
}
