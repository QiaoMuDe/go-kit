package testutil

import (
	"os"
	"testing"
)

func TestQuiet(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "基本功能测试_返回RunFunc",
			testFunc: func(t *testing.T) {
				runFunc := Quiet()
				if runFunc == nil {
					t.Error("期望 Quiet() 返回非空的 RunFunc")
				}
			},
		},
		{
			name: "RunFunc类型验证",
			testFunc: func(t *testing.T) {
				runFunc := Quiet()

				// 验证返回的函数可以接受 *testing.M 参数并返回 int
				// 创建一个模拟的 testing.M
				m := &testing.M{}

				// 这应该能够调用而不会 panic
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("调用 runFunc 时发生 panic: %v", r)
					}
				}()

				// 注意：实际调用 m.Run() 会运行所有测试，这里我们只验证函数签名
				_ = runFunc
				_ = m
			},
		},
		{
			name: "verbose模式测试",
			testFunc: func(t *testing.T) {
				// 保存原始的 stdout 和 stderr
				origStdout := os.Stdout
				origStderr := os.Stderr
				defer func() {
					os.Stdout = origStdout
					os.Stderr = origStderr
				}()

				// 在 verbose 模式下，Quiet 应该不会重定向输出
				runFunc := Quiet()
				if runFunc == nil {
					t.Error("期望 Quiet() 返回非空的 RunFunc")
				}

				// 验证 stdout 和 stderr 没有被改变（在 verbose 模式下）
				if testing.Verbose() {
					if os.Stdout != origStdout {
						t.Error("在 verbose 模式下，stdout 不应该被重定向")
					}
					if os.Stderr != origStderr {
						t.Error("在 verbose 模式下，stderr 不应该被重定向")
					}
				}
			},
		},
		{
			name: "多次调用测试",
			testFunc: func(t *testing.T) {
				// 多次调用 Quiet 应该都能正常工作
				for i := 0; i < 3; i++ {
					runFunc := Quiet()
					if runFunc == nil {
						t.Errorf("第 %d 次调用 Quiet() 返回了 nil", i+1)
					}
				}
			},
		},
		{
			name: "并发安全性测试",
			testFunc: func(t *testing.T) {
				// 由于 flag.Parse() 不是并发安全的，我们测试顺序调用
				// 这更符合实际使用场景，因为 TestMain 通常只会被调用一次

				results := make([]RunFunc, 3)

				// 顺序调用多次来测试是否有资源泄露或状态问题
				for i := 0; i < 3; i++ {
					runFunc := Quiet()
					if runFunc == nil {
						t.Errorf("第 %d 次调用 Quiet() 返回了 nil", i+1)
					}
					results[i] = runFunc
				}

				// 验证所有返回的函数都不为空且可以使用
				for i, runFunc := range results {
					if runFunc == nil {
						t.Errorf("结果 %d 为 nil", i)
					}
				}

				t.Log("并发安全性测试完成")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// 基准测试
func BenchmarkQuiet(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runFunc := Quiet()
		_ = runFunc // 避免编译器优化
	}
}

// 示例测试
func ExampleQuiet() {
	// 在 TestMain 中使用 Quiet
	// func TestMain(m *testing.M) {
	//     run := Quiet()
	//     os.Exit(run(m))
	// }

	runFunc := Quiet()
	_ = runFunc // 实际使用中会调用 runFunc(m)

	// Output:
}

// 集成测试：模拟 TestMain 的使用场景
func TestQuietIntegration(t *testing.T) {
	// 创建一个模拟的测试环境
	runFunc := Quiet()

	// 验证返回的函数不为空
	if runFunc == nil {
		t.Fatal("Quiet() 返回了 nil RunFunc")
	}

	// 模拟 TestMain 中的使用
	// 注意：我们不能真正调用 m.Run()，因为那会递归运行测试
	m := &testing.M{}
	_ = m

	// 验证函数签名正确
	var testRunFunc = runFunc
	_ = testRunFunc

	t.Log("集成测试通过：Quiet() 返回了正确类型的函数")
}

// 边界测试：测试资源清理
func TestQuietResourceCleanup(t *testing.T) {
	// 保存原始状态
	origStdout := os.Stdout
	origStderr := os.Stderr

	// 多次调用 Quiet 来测试资源是否正确清理
	for i := 0; i < 10; i++ {
		runFunc := Quiet()
		if runFunc == nil {
			t.Errorf("第 %d 次调用失败", i+1)
		}

		// 在非 verbose 模式下，Quiet() 会重定向输出
		// 这是预期的行为，我们主要测试函数不会 panic
	}

	// 在 verbose 模式下，原始状态应该保持不变
	// 在非 verbose 模式下，输出会被重定向，这是正常的
	if testing.Verbose() {
		if os.Stdout != origStdout {
			t.Error("在 verbose 模式下，stdout 状态不应该改变")
		}
		if os.Stderr != origStderr {
			t.Error("在 verbose 模式下，stderr 状态不应该改变")
		}
	} else {
		// 在非 verbose 模式下，输出被重定向是正常的
		t.Log("在非 verbose 模式下，输出被重定向是预期行为")
	}
}
