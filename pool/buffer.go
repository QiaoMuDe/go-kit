package pool

import (
	"bytes"
	"sync"
)

// 全局默认缓冲区池实例
var defBufPool = NewBufferPool(256, 32*1024)

// GetBuffer 从默认缓冲区池获取默认容量的字节缓冲区
//
// 返回值:
//   - *bytes.Buffer: 容量至少为默认容量的字节缓冲区
func GetBuffer() *bytes.Buffer { return defBufPool.Get() }

// GetBufferWithCapacity 从默认缓冲区池获取指定容量的字节缓冲区
//
// 参数:
//   - capacity: 缓冲区初始容量
//
// 返回值:
//   - *bytes.Buffer: 容量至少为capacity的字节缓冲区
func GetBufferWithCapacity(capacity int) *bytes.Buffer {
	return defBufPool.GetWithCapacity(capacity)
}

// PutBuffer 将字节缓冲区归还到默认缓冲区池
//
// 参数:
//   - buffer: 要归还的字节缓冲区
//
// 说明:
//   - 该函数将字节缓冲区归还到对象池，以便后续复用。
func PutBuffer(buffer *bytes.Buffer) { defBufPool.Put(buffer) }

// WithBuffer 使用默认容量的字节缓冲区执行函数，自动管理获取和归还
//
// 参数:
//   - fn: 使用字节缓冲区的函数
//
// 返回值:
//   - []byte: 函数执行后缓冲区的字节数据副本
//
// 使用示例:
//
//	data := pool.WithBuffer(func(buf *bytes.Buffer) {
//	    buf.WriteString("Hello")
//	    buf.WriteByte(' ')
//	    buf.WriteString("World")
//	})
func WithBuffer(fn func(*bytes.Buffer)) []byte { return defBufPool.WithBuffer(fn) }

// WithBufferCapacity 使用指定容量的字节缓冲区执行函数，自动管理获取和归还
//
// 参数:
//   - capacity: 字节缓冲区初始容量
//   - fn: 使用字节缓冲区的函数
//
// 返回值:
//   - []byte: 函数执行后缓冲区的字节数据副本
//
// 使用示例:
//
//	data := pool.WithBufferCapacity(1024, func(buf *bytes.Buffer) {
//	    buf.WriteString("Hello")
//	    buf.WriteByte(' ')
//	    buf.WriteString("World")
//	})
func WithBufferCapacity(capacity int, fn func(*bytes.Buffer)) []byte {
	return defBufPool.WithBufferCapacity(capacity, fn)
}

// BufferPool 字节缓冲区对象池，支持自定义配置
type BufferPool struct {
	pool            sync.Pool // 字节缓冲区对象池
	maxCapacity     int       // 最大回收缓冲区容量
	defaultCapacity int       // 默认缓冲区容量
}

// NewBufferPool 创建新的字节缓冲区对象池
//
// 参数:
//   - defaultCapacity: 默认字节缓冲区容量
//   - maxCapacity: 最大回收缓冲区容量，超过此容量的缓冲区不会被回收
//
// 返回值:
//   - *BufferPool: 字节缓冲区对象池实例
func NewBufferPool(defaultCapacity, maxCapacity int) *BufferPool {
	if defaultCapacity <= 0 {
		defaultCapacity = 256 // 默认256字节
	}
	if maxCapacity <= 0 {
		maxCapacity = 32 * 1024 // 默认32KB
	}

	return &BufferPool{
		maxCapacity:     maxCapacity,
		defaultCapacity: defaultCapacity,
		pool: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

// Get 获取默认容量的字节缓冲区
//
// 返回:
//   - *bytes.Buffer: 容量至少为默认容量的字节缓冲区
//
// 说明:
//   - 返回的字节缓冲区已经重置为空状态，可以直接使用
//   - 底层容量可能大于默认容量，来自对象池的复用缓冲区
func (bp *BufferPool) Get() *bytes.Buffer { return bp.GetWithCapacity(bp.defaultCapacity) }

// GetWithCapacity 获取指定容量的字节缓冲区
//
// 参数:
//   - capacity: 需要的字节缓冲区容量
//
// 返回:
//   - *bytes.Buffer: 容量至少为capacity的字节缓冲区
//
// 说明:
//   - 返回的字节缓冲区已经重置为空状态，可以直接使用
//   - 底层容量可能大于capacity，来自对象池的复用缓冲区
//   - 如果capacity <= 0, 返回默认容量的缓冲区
func (bp *BufferPool) GetWithCapacity(capacity int) *bytes.Buffer {
	if capacity <= 0 {
		capacity = bp.defaultCapacity
	}

	buffer, ok := bp.pool.Get().(*bytes.Buffer)
	if !ok {
		// 一旦触发说明代码契约被破坏了,直接panic比静默继续更安全
		panic("bufferPool: unexpected type")
	}

	// 如果当前容量不足，扩容到所需容量
	if buffer.Cap() < capacity {
		buffer.Grow(capacity)
	}

	// 重置缓冲区状态
	buffer.Reset()

	return buffer
}

// Put 归还字节缓冲区到对象池
//
// 参数:
//   - buffer: 要归还的字节缓冲区
func (bp *BufferPool) Put(buf *bytes.Buffer) {
	if buf == nil || buf.Cap() > bp.maxCapacity {
		return // 直接扔
	}
	buf.Reset()
	bp.pool.Put(buf)
}

// WithBuffer 使用默认容量的字节缓冲区执行函数，自动管理获取和归还
//
// 参数:
//   - fn: 使用字节缓冲区的函数
//
// 返回值:
//   - []byte: 函数执行后缓冲区的字节数据副本
//
// 说明:
//   - 自动从对象池获取默认容量的字节缓冲区
//   - 执行用户提供的函数
//   - 获取缓冲区字节数据的副本
//   - 自动归还字节缓冲区到对象池
//   - 即使函数发生panic也会正确归还资源
func (bp *BufferPool) WithBuffer(fn func(*bytes.Buffer)) []byte {
	buf := bp.Get()
	defer bp.Put(buf)

	fn(buf)
	return append([]byte(nil), buf.Bytes()...) // 一次性拷贝
}

// WithBufferCapacity 使用指定容量的字节缓冲区执行函数，自动管理获取和归还
//
// 参数:
//   - capacity: 字节缓冲区初始容量
//   - fn: 使用字节缓冲区的函数
//
// 返回值:
//   - []byte: 函数执行后缓冲区的字节数据副本
//
// 说明:
//   - 自动从对象池获取指定容量的字节缓冲区
//   - 执行用户提供的函数
//   - 获取缓冲区字节数据的副本
//   - 自动归还字节缓冲区到对象池
//   - 即使函数发生panic也会正确归还资源
func (bp *BufferPool) WithBufferCapacity(capacity int, fn func(*bytes.Buffer)) []byte {
	buf := bp.GetWithCapacity(capacity)
	defer bp.Put(buf)
	fn(buf)
	return append([]byte(nil), buf.Bytes()...) // 一次性拷贝
}
