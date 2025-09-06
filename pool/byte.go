// Package pool 提供高性能对象池管理功能, 通过复用对象优化内存使用。
//
// 该包实现了基于 sync.Pool 的多种对象池, 用于减少频繁的内存分配和回收。
// 通过复用对象, 可以显著提升应用程序的性能, 特别是在高并发场景下。
//
// 主要功能：
//   - 字节切片对象池管理
//   - 动态容量对象获取
//   - 自动内存回收控制
//   - 防止内存泄漏的容量限制
//   - 支持多种对象类型的池化
//
// 性能优化：
//   - 使用 sync.Pool 减少 GC 压力
//   - 支持不同容量的对象需求
//   - 自动限制大对象回收
//   - 预热机制提升冷启动性能
//
// 使用示例：
//
//	// 获取默认容量的字节缓冲区
//	buf := pool.GetByte()
//
//	// 获取指定容量的字节缓冲区
//	largeBuf := pool.GetByteWithCapacity(64 * 1024)
//
//	// 使用缓冲区进行文件操作
//	_, err := io.CopyBuffer(dst, src, largeBuf)
//
//	// 归还缓冲区到对象池
//	pool.PutByte(largeBuf)
package pool

import "sync"

// 全局默认对象池实例, 默认容量为256, 最大容量为32KB
var defaultPool = NewBytePool(256, 32*1024)

// GetByte 从默认字节池获取默认容量的缓冲区
//
// 返回值:
//   - []byte: 长度为默认容量, 容量至少为默认容量的缓冲区
func GetByte() []byte { return defaultPool.Get() }

// GetByteWithCapacity 从默认字节池获取指定容量的缓冲区
//
// 参数:
//   - capacity: 缓冲区容量
//
// 返回值:
//   - []byte: 长度为capacity, 容量至少为capacity的缓冲区
func GetByteWithCapacity(capacity int) []byte { return defaultPool.GetByteWithCapacity(capacity) }

// PutByte 将缓冲区归还到默认字节池
//
// 参数:
//   - buffer: 要归还的缓冲区
//
// 说明:
//   - 该函数将缓冲区归还到对象池, 以便后续复用。
func PutByte(buffer []byte) { defaultPool.Put(buffer) }

// GetEmptyByte 从默认字节池获取空缓冲区
//
// 参数:
//   - capacity: 指定容量要求
//
// 返回值:
//   - []byte: 长度为0但容量至少为capacity的缓冲区切片
func GetEmptyByte(capacity int) []byte { return defaultPool.GetEmpty(capacity) }

// BytePool 字节切片对象池, 支持自定义配置
type BytePool struct {
	pool            sync.Pool // 缓冲区对象池
	maxCapacity     int       // 最大回收缓冲区容量
	defaultCapacity int       // 默认缓冲区容量
}

// NewBytePool 创建新的字节切片对象池
//
// 参数:
//   - defaultCapacity: 默认缓冲区容量
//   - maxCapacity: 最大回收缓冲区容量, 超过此容量的缓冲区不会被回收
//
// 返回值:
//   - *BytePool: 字节切片对象池实例
func NewBytePool(defaultCapacity, maxCapacity int) *BytePool {
	if defaultCapacity <= 0 {
		defaultCapacity = 256 // 默认256字节
	}
	if maxCapacity <= 0 {
		maxCapacity = 32 * 1024 // 默认32KB
	}

	return &BytePool{
		maxCapacity:     maxCapacity,
		defaultCapacity: defaultCapacity,
		pool: sync.Pool{
			New: func() any {
				return make([]byte, 0, defaultCapacity)
			},
		},
	}
}

// Get 获取默认容量的缓冲区
//
// 返回:
//   - []byte: 长度为默认容量, 容量至少为默认容量的缓冲区切片
//
// 说明:
//   - 返回的缓冲区长度等于默认容量, 可以直接使用
//   - 底层容量可能大于默认容量, 来自对象池的复用缓冲区
func (bp *BytePool) Get() []byte {
	return bp.GetByteWithCapacity(bp.defaultCapacity)
}

// GetByteWithCapacity 获取指定容量的缓冲区
//
// 参数:
//   - capacity: 需要的缓冲区容量
//
// 返回:
//   - []byte: 长度为capacity, 容量至少为capacity的缓冲区切片
//
// 说明:
//   - 返回的缓冲区长度等于请求的capacity, 可以直接使用
//   - 底层容量可能大于capacity, 来自对象池的复用缓冲区
//   - 如果capacity <= 0, 使用默认容量
func (bp *BytePool) GetByteWithCapacity(capacity int) []byte {
	if capacity <= 0 {
		capacity = bp.defaultCapacity
	}

	buf, ok := bp.pool.Get().([]byte)
	if !ok {
		// 类型断言失败, 创建新的
		return make([]byte, capacity)
	}

	// 容量足够，返回
	if cap(buf) >= capacity {
		return buf[:capacity] // 返回长度为capacity的切片
	}

	// 容量不足，创建新的
	bp.pool.Put(buf) // 归还旧的
	return make([]byte, capacity)
}

// Put 归还缓冲区到对象池
//
// 参数:
//   - buffer: 要归还的缓冲区
func (bp *BytePool) Put(buffer []byte) {
	if buffer == nil || cap(buffer) > bp.maxCapacity {
		return // 不回收空指针或容量超过最大回收容量
	}

	bp.pool.Put(buffer[:0]) // 归还空缓冲区
}

// GetEmpty 获取指定容量的空缓冲区
//
// 参数:
//   - capacity: 指定容量要求
//
// 返回:
//   - []byte: 长度为0但容量至少为capacity的缓冲区切片
//
// 说明:
//   - 适用于需要使用append操作逐步构建数据的场景
//   - 避免频繁的内存重新分配
//   - 如果capacity <= 0, 使用默认容量
func (bp *BytePool) GetEmpty(capacity int) []byte {
	if capacity <= 0 {
		capacity = bp.defaultCapacity
	}

	buf, ok := bp.pool.Get().([]byte)
	if !ok {
		// 类型断言失败, 创建新的
		return make([]byte, 0, capacity)
	}

	// 缓冲区容量足够, 返回空切片
	if cap(buf) >= capacity {
		return buf[:0]
	}

	// 缓冲区容量不足, 创建新的
	bp.pool.Put(buf) // 归还旧的
	return make([]byte, 0, capacity)
}
