// Package buffer 提供缓冲区管理功能，通过对象池优化内存使用。
//
// 该包实现了基于 sync.Pool 的缓冲区对象池，用于减少频繁的内存分配和回收。
// 通过复用缓冲区，可以显著提升文件读写、网络I/O等操作的性能。
//
// 主要功能：
//   - 缓冲区对象池管理
//   - 动态大小缓冲区获取
//   - 自动内存回收控制
//   - 防止内存泄漏的大小限制
//
// 性能优化：
//   - 使用 sync.Pool 减少 GC 压力
//   - 支持不同大小的缓冲区需求
//   - 自动限制大缓冲区回收
//
// 使用示例：
//
//	// 获取缓冲区
//	buf := buffer.Get(64 * 1024)
//
//	// 使用缓冲区进行文件操作
//	_, err := io.CopyBuffer(dst, src, buf)
//
//	// 归还缓冲区到对象池
//	buffer.Put(buf)
package buffer

import "sync"

// 全局默认对象池实例
var defaultPool = NewBytePool(32*1024, 1024*1024)

// Get 从默认对象池获取缓冲区
//
// 参数:
//   - size: 缓冲区大小
//
// 返回值:
//   - []byte: 获取到的缓冲区
func Get(size int) []byte {
	return defaultPool.Get(size)
}

// Put 将缓冲区归还到默认对象池
//
// 参数:
//   - buffer: 要归还的缓冲区
//
// 说明:
//   - 该函数将缓冲区归还到对象池，以便后续复用。
//   - 只有容量不超过1MB的缓冲区才会被归还，以避免对象池占用过多内存。
func Put(buffer []byte) {
	defaultPool.Put(buffer)
}

// BytePool 字节切片对象池，支持自定义配置
type BytePool struct {
	pool    sync.Pool // 缓冲区对象池
	maxSize int       // 最大回收缓冲区大小
}

// NewBytePool 创建新的字节切片对象池
//
// 参数:
//   - defaultSize: 默认缓冲区大小
//   - maxSize: 最大回收缓冲区大小，超过此大小的缓冲区不会被回收
//
// 返回值:
//   - *BytePool: 字节切片对象池实例
func NewBytePool(defaultSize, maxSize int) *BytePool {
	if defaultSize <= 0 {
		defaultSize = 32 * 1024 // 默认32KB
	}
	if maxSize <= 0 {
		maxSize = 1024 * 1024 // 默认1MB
	}

	return &BytePool{
		maxSize: maxSize,
		pool: sync.Pool{
			New: func() any {
				buf := make([]byte, defaultSize)
				return &buf // 返回指针避免装箱
			},
		},
	}
}

// Get 获取指定大小的缓冲区
//
// 参数:
//   - size: 需要的缓冲区大小
//
// 返回:
//   - []byte: 缓冲区切片
func (bp *BytePool) Get(size int) []byte {
	bufPtr, ok := bp.pool.Get().(*[]byte)
	if !ok {
		// 类型断言失败，创建新的
		return make([]byte, size)
	}

	// 获取缓冲区
	buffer := *bufPtr

	// 缓冲区容量不足，创建新的
	if cap(buffer) < size {
		return make([]byte, size)
	}

	return buffer[:size]
}

// Put 归还缓冲区到对象池
//
// 参数:
//   - buffer: 要归还的缓冲区
func (bp *BytePool) Put(buffer []byte) {
	if buffer == nil || cap(buffer) > bp.maxSize {
		return // 不回收nil或超过最大大小的缓冲区
	}

	// 清空缓冲区内容
	buffer = buffer[:0]

	bp.pool.Put(&buffer) // 传入指针避免装箱分配
}
