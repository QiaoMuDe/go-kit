package syncx

import "sync"

// FileLocks 是一个文件级锁管理器，可创建多个独立实例。
type FileLocks struct{ m sync.Map }

// Lock 对指定key加锁并返回解锁函数
// 用于为特定key创建互斥锁，首次访问时创建锁对象，后续复用
//
// 参数:
//   - key: 锁的键，通常是文件路径
//
// 返回:
//   - unlock: 解锁函数
func (fl *FileLocks) Lock(key string) (unlock func()) {
	actual, _ := fl.m.LoadOrStore(key, new(sync.Mutex))
	mu := actual.(*sync.Mutex)
	mu.Lock()
	return func() { mu.Unlock() }
}
