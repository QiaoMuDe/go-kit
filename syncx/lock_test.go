package syncx

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestFileLocks(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "基本功能测试_单个key加锁解锁",
			testFunc: func(t *testing.T) {
				fl := &FileLocks{}
				key := "test-file.txt"

				// 获取锁
				unlock := fl.Lock(key)
				if unlock == nil {
					t.Error("期望返回非空的解锁函数")
				}

				// 解锁
				unlock()
			},
		},
		{
			name: "基本功能测试_多个不同key并发加锁",
			testFunc: func(t *testing.T) {
				fl := &FileLocks{}
				var wg sync.WaitGroup

				keys := []string{"file1.txt", "file2.txt", "file3.txt"}
				wg.Add(len(keys))

				// 并发获取不同key的锁
				for _, key := range keys {
					go func(k string) {
						defer wg.Done()
						unlock := fl.Lock(k)
						time.Sleep(10 * time.Millisecond) // 模拟工作
						unlock()
					}(key)
				}

				wg.Wait()
			},
		},
		{
			name: "并发测试_同一key多个goroutine竞争",
			testFunc: func(t *testing.T) {
				fl := &FileLocks{}
				key := "shared-file.txt"
				var counter int
				var wg sync.WaitGroup

				goroutineCount := 10
				wg.Add(goroutineCount)

				// 多个 goroutine 竞争同一个key的锁
				for i := 0; i < goroutineCount; i++ {
					go func(id int) {
						defer wg.Done()

						unlock := fl.Lock(key)
						// 临界区操作
						temp := counter
						time.Sleep(time.Millisecond) // 增加竞态条件发生概率
						counter = temp + 1
						unlock()
					}(i)
				}

				wg.Wait()

				// 验证计数器值正确
				if counter != goroutineCount {
					t.Errorf("期望计数器值为 %d，实际为 %d", goroutineCount, counter)
				}
			},
		},
		{
			name: "锁复用测试_同一key多次加锁",
			testFunc: func(t *testing.T) {
				fl := &FileLocks{}
				key := "reuse-file.txt"

				// 多次对同一key加锁解锁
				for i := 0; i < 5; i++ {
					unlock := fl.Lock(key)
					if unlock == nil {
						t.Errorf("第 %d 次加锁失败", i+1)
					}
					unlock()
				}
			},
		},
		{
			name: "边界测试_空字符串key",
			testFunc: func(t *testing.T) {
				fl := &FileLocks{}
				key := ""

				// 空字符串key应该也能正常工作
				unlock := fl.Lock(key)
				if unlock == nil {
					t.Error("空字符串key加锁失败")
				}
				unlock()
			},
		},
		{
			name: "边界测试_特殊字符key",
			testFunc: func(t *testing.T) {
				fl := &FileLocks{}
				specialKeys := []string{
					"/path/to/file.txt",
					"file with spaces.txt",
					"file-with-dashes.txt",
					"file_with_underscores.txt",
					"文件名中文.txt",
					"file@#$%^&*().txt",
				}

				for _, key := range specialKeys {
					unlock := fl.Lock(key)
					if unlock == nil {
						t.Errorf("特殊字符key '%s' 加锁失败", key)
					}
					unlock()
				}
			},
		},
		{
			name: "内存泄漏测试_大量不同key",
			testFunc: func(t *testing.T) {
				fl := &FileLocks{}

				// 创建大量不同的key
				for i := 0; i < 1000; i++ {
					key := fmt.Sprintf("file-%d.txt", i)
					unlock := fl.Lock(key)
					unlock()
				}

				// 验证FileLocks仍然可以正常工作
				unlock := fl.Lock("test-after-many.txt")
				if unlock == nil {
					t.Error("大量key操作后，FileLocks无法正常工作")
				}
				unlock()
			},
		},
		{
			name: "并发安全测试_混合操作",
			testFunc: func(t *testing.T) {
				fl := &FileLocks{}
				var wg sync.WaitGroup
				var results sync.Map

				// 并发执行混合操作
				operations := 100
				wg.Add(operations)

				for i := 0; i < operations; i++ {
					go func(id int) {
						defer wg.Done()

						// 使用不同的key模式
						key := fmt.Sprintf("file-%d.txt", id%10) // 10个不同的key

						unlock := fl.Lock(key)

						// 记录操作
						if _, loaded := results.LoadOrStore(key, 1); loaded {
							// 如果key已存在，增加计数
							if val, ok := results.Load(key); ok {
								results.Store(key, val.(int)+1)
							}
						}

						time.Sleep(time.Microsecond) // 模拟短暂工作
						unlock()
					}(i)
				}

				wg.Wait()

				// 验证所有操作都完成了
				totalOps := 0
				results.Range(func(key, value interface{}) bool {
					totalOps += value.(int)
					return true
				})

				if totalOps != operations {
					t.Errorf("期望总操作数 %d，实际 %d", operations, totalOps)
				}
			},
		},
		{
			name: "解锁函数测试_正确使用解锁",
			testFunc: func(t *testing.T) {
				fl := &FileLocks{}
				key := "unlock-test.txt"

				// 测试正常的加锁解锁流程
				unlock1 := fl.Lock(key)
				unlock1() // 正常解锁

				// 再次加锁应该成功
				unlock2 := fl.Lock(key)
				unlock2() // 正常解锁

				// 验证可以继续使用
				unlock3 := fl.Lock(key)
				if unlock3 == nil {
					t.Error("期望能够再次获取锁")
				}
				unlock3()
			},
		},
		{
			name: "零值测试_未初始化的FileLocks",
			testFunc: func(t *testing.T) {
				var fl FileLocks // 零值
				key := "zero-value-test.txt"

				// 零值的FileLocks应该也能正常工作
				unlock := fl.Lock(key)
				if unlock == nil {
					t.Error("零值FileLocks加锁失败")
				}
				unlock()
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
func BenchmarkFileLocks_Lock(b *testing.B) {
	fl := &FileLocks{}
	key := "benchmark-file.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		unlock := fl.Lock(key)
		unlock()
	}
}

// 基准测试：不同key的性能
func BenchmarkFileLocks_DifferentKeys(b *testing.B) {
	fl := &FileLocks{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("file-%d.txt", i%100) // 100个不同的key
		unlock := fl.Lock(key)
		unlock()
	}
}

// 基准测试：并发性能
func BenchmarkFileLocks_Concurrent(b *testing.B) {
	fl := &FileLocks{}
	key := "concurrent-file.txt"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			unlock := fl.Lock(key)
			unlock()
		}
	})
}

// 基准测试：对比直接使用sync.Mutex的性能
func BenchmarkFileLocks_vs_Mutex(b *testing.B) {
	b.Run("FileLocks", func(b *testing.B) {
		fl := &FileLocks{}
		key := "comparison-file.txt"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			unlock := fl.Lock(key)
			unlock()
		}
	})

	b.Run("DirectMutex", func(b *testing.B) {
		var mu sync.Mutex
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mu.Lock()
			// 模拟最小的临界区操作
			_ = i
			mu.Unlock()
		}
	})
}

// 示例测试
func ExampleFileLocks_Lock() {
	fl := &FileLocks{}

	// 对文件加锁
	unlock := fl.Lock("/path/to/file.txt")
	defer unlock()

	// 执行需要同步的文件操作
	println("正在安全地操作文件...")

	// unlock() 会在defer中自动调用
}

// 压力测试：长时间运行的并发测试
func TestFileLocks_Stress(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过压力测试")
	}

	fl := &FileLocks{}
	var wg sync.WaitGroup
	var operations int64

	// 运行时间
	duration := 2 * time.Second
	done := make(chan struct{})

	// 启动定时器
	go func() {
		time.Sleep(duration)
		close(done)
	}()

	// 启动多个工作 goroutine
	workerCount := 20
	keyCount := 10
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func(id int) {
			defer wg.Done()

			for {
				select {
				case <-done:
					return
				default:
					// 随机选择key
					key := fmt.Sprintf("stress-file-%d.txt", id%keyCount)
					unlock := fl.Lock(key)
					operations++
					// 模拟短暂的工作
					time.Sleep(time.Microsecond * 10)
					unlock()
				}
			}
		}(i)
	}

	wg.Wait()

	t.Logf("压力测试完成，总操作数: %d", operations)

	if operations == 0 {
		t.Error("压力测试期间没有成功的操作")
	}
}

// 辅助函数测试：验证 FileLocks 的类型和方法签名
func TestFileLocks_TypeSignature(t *testing.T) {
	// 验证 FileLocks 结构体
	var fl FileLocks

	// 验证 Lock 方法签名
	var lockMethod func(string) func() = fl.Lock
	_ = lockMethod

	// 验证返回的解锁函数
	unlock := fl.Lock("signature-test.txt")
	var unlockFunc func() = unlock
	_ = unlockFunc

	unlock()
}
