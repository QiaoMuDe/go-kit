package pool

import (
	"sync"
	"testing"
	"time"
)

func TestTimerPool_Get(t *testing.T) {
	timer := GetTimer(time.Second)
	if timer == nil {
		t.Fatal("GetTimer() returned nil")
	}

	// 验证timer是停止状态
	select {
	case <-timer.C:
		t.Error("Timer should be stopped when retrieved from pool")
	default:
		// 正确，timer应该是停止的
	}

	PutTimer(timer)
}

func TestTimerPool_Put(t *testing.T) {
	timer := GetTimer(time.Second)

	// 启动timer
	timer.Reset(10 * time.Millisecond)

	// 等待触发
	select {
	case <-timer.C:
		// timer触发了
	case <-time.After(50 * time.Millisecond):
		t.Error("Timer should have fired")
	}

	PutTimer(timer)

	// 再次获取应该是停止状态
	timer2 := GetTimer(time.Second)
	select {
	case <-timer2.C:
		t.Error("Timer should be stopped when retrieved from pool after put")
	default:
		// 正确
	}

	PutTimer(timer2)
}

func TestTimerPool_Reset(t *testing.T) {
	timer := GetTimer(time.Second)

	// 测试Reset方法
	timer.Reset(20 * time.Millisecond)

	// 立即重置为更长时间
	timer.Reset(100 * time.Millisecond)

	// 在20ms内不应该触发
	select {
	case <-timer.C:
		t.Error("Timer should not fire after reset")
	case <-time.After(30 * time.Millisecond):
		// 正确，timer被重置了
	}

	// 停止timer
	timer.Stop()
	PutTimer(timer)
}

func TestTimerPool_Stop(t *testing.T) {
	timer := GetTimer(time.Second)

	// 启动timer
	timer.Reset(50 * time.Millisecond)

	// 立即停止
	stopped := timer.Stop()
	if !stopped {
		t.Log("Timer was already stopped or fired")
	}

	// 验证不会触发
	select {
	case <-timer.C:
		t.Error("Timer should not fire after stop")
	case <-time.After(100 * time.Millisecond):
		// 正确
	}

	PutTimer(timer)
}

func TestTimerPool_Reuse(t *testing.T) {
	timer1 := GetTimer(time.Second)
	PutTimer(timer1)

	timer2 := GetTimer(time.Second)
	// 在单线程环境下可能复用同一个对象
	if timer1 == timer2 {
		t.Log("Reused the same timer object")
	}

	PutTimer(timer2)
}

func TestTimerPool_Concurrent(t *testing.T) {
	const numGoroutines = 50
	const numOperations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				timer := GetTimer(time.Second)
				if timer == nil {
					t.Errorf("GetTimer() returned nil in goroutine %d", id)
					return
				}

				// 设置一个短暂的超时
				duration := time.Duration(1+j%10) * time.Millisecond
				timer.Reset(duration)

				// 等待触发或超时
				select {
				case <-timer.C:
					// timer正常触发
				case <-time.After(duration + 10*time.Millisecond):
					t.Errorf("Timer did not fire in time in goroutine %d", id)
					timer.Stop()
				}

				PutTimer(timer)
			}
		}(i)
	}

	wg.Wait()
}

func TestTimerPool_MultipleResets(t *testing.T) {
	timer := GetTimer(time.Second)

	// 多次重置timer
	for i := 0; i < 10; i++ {
		timer.Reset(time.Duration(i+1) * time.Millisecond)

		// 立即重置，前一个应该被取消
		if i < 9 {
			continue
		}

		// 最后一次等待触发
		select {
		case <-timer.C:
			// 正确触发
		case <-time.After(50 * time.Millisecond):
			t.Error("Final timer should have fired")
		}
	}

	PutTimer(timer)
}

func TestTimerPool_EdgeCases(t *testing.T) {
	// 测试零时间
	timer := GetTimer(time.Second)
	timer.Reset(0)

	select {
	case <-timer.C:
		// 应该立即触发
	case <-time.After(10 * time.Millisecond):
		t.Error("Zero duration timer should fire immediately")
	}

	PutTimer(timer)

	// 测试负时间（应该被当作0处理）
	timer2 := GetTimer(time.Second)
	timer2.Reset(-time.Second)

	select {
	case <-timer2.C:
		// 应该立即触发
	case <-time.After(10 * time.Millisecond):
		t.Error("Negative duration timer should fire immediately")
	}

	PutTimer(timer2)
}

func TestTimerPool_LongDuration(t *testing.T) {
	timer := GetTimer(time.Second)

	// 设置一个很长的时间
	timer.Reset(time.Hour)

	// 立即停止
	stopped := timer.Stop()
	if !stopped {
		t.Error("Should be able to stop long duration timer")
	}

	// 验证不会触发
	select {
	case <-timer.C:
		t.Error("Stopped timer should not fire")
	case <-time.After(10 * time.Millisecond):
		// 正确
	}

	PutTimer(timer)
}

func TestTimerPool_DrainChannel(t *testing.T) {
	timer := GetTimer(time.Second)

	// 让timer触发
	timer.Reset(1 * time.Millisecond)

	// 等待触发
	select {
	case <-timer.C:
		// timer触发了
	case <-time.After(10 * time.Millisecond):
		t.Fatal("Timer should have fired")
	}

	// 现在channel应该是空的
	select {
	case <-timer.C:
		t.Error("Timer channel should be empty after firing")
	default:
		// 正确
	}

	PutTimer(timer)
}

func BenchmarkTimerPool_GetPut(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		timer := GetTimer(time.Microsecond)
		timer.Reset(time.Microsecond)
		timer.Stop()
		PutTimer(timer)
	}
}

func BenchmarkTimerPool_vs_New(b *testing.B) {
	b.Run("Pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			timer := GetTimer(time.Microsecond)
			timer.Reset(time.Microsecond)
			timer.Stop()
			PutTimer(timer)
		}
	})

	b.Run("New", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			timer := time.NewTimer(time.Microsecond)
			timer.Stop()
		}
	})
}

func BenchmarkTimerPool_FireAndReset(b *testing.B) {
	timer := GetTimer(time.Nanosecond)
	defer PutTimer(timer)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		timer.Reset(time.Nanosecond)
		<-timer.C
	}
}
