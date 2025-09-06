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

// 全局默认对象池实例, 默认容量为256，最大容量为32KB
var defaultPool = NewBytePool(256, 32*1024)

// GetByte 从默认字节池获取默认容量的缓冲区
//
// 返回值:
//   - []byte: 长度为默认容量, 容量至少为默认容量的缓冲区
func GetByte() []byte {
	return defaultPool.Get()
}

// GetByteWithCapacity 从默认字节池获取指定容量的缓冲区
//
// 参数:
//   - capacity: 缓冲区容量
//
// 返回值:
//   - []byte: 长度为capacity, 容量至少为capacity的缓冲区
func GetByteWithCapacity(capacity int) []byte {
	return defaultPool.GetByteWithCapacity(capacity)
}

// PutByte 将缓冲区归还到默认字节池
//
// 参数:
//   - buffer: 要归还的缓冲区
//
// 说明:
//   - 该函数将缓冲区归还到对象池, 以便后续复用。
func PutByte(buffer []byte) {
	defaultPool.Put(buffer)
}

// GetEmptyByte 从默认字节池获取空缓冲区
//
// 参数:
//   - capacity: 指定容量要求
//
// 返回值:
//   - []byte: 长度为0但容量至少为capacity的缓冲区切片
func GetEmptyByte(capacity int) []byte {
	return defaultPool.GetEmpty(capacity)
}

// WarmByte 预热默认字节池
//
// 参数:
//   - count: 预分配的缓冲区数量
//   - capacity: 每个缓冲区的容量
func WarmByte(count int, capacity int) {
	defaultPool.Warm(count, capacity)
}

// DrainByte 清空默认字节池
func DrainByte() {
	defaultPool.Drain()
}

// WithByte 使用默认容量的字节切片执行函数, 自动管理获取和归还
//
// 参数:
//   - fn: 使用字节切片的函数
//
// 返回值:
//   - []byte: 函数执行后字节切片的数据副本
//
// 使用示例:
//
//	data := pool.WithByte(func(buf []byte) {
//	    // buf长度为默认容量(256字节), 可以直接使用
//	    n := copy(buf, []byte("Hello World"))
//	    // 注意: 只使用了前n个字节
//	})
//	// data包含完整的buf内容(256字节)
func WithByte(fn func([]byte)) []byte {
	return defaultPool.WithByte(fn)
}

// WithByteCapacity 使用指定容量的字节切片执行函数, 自动管理获取和归还
//
// 参数:
//   - capacity: 字节切片初始容量
//   - fn: 使用字节切片的函数
//
// 返回值:
//   - []byte: 函数执行后字节切片的数据副本
//
// 使用示例:
//
//	data := pool.WithByteCapacity(1024, func(buf []byte) {
//	    // buf长度为1024字节, 可以直接使用
//	    n := copy(buf, []byte("Hello World"))
//	    // 可以继续使用buf的其他部分
//	    buf[n] = '\n' // 添加换行符
//	})
//	// data包含完整的buf内容(1024字节)
func WithByteCapacity(capacity int, fn func([]byte)) []byte {
	return defaultPool.WithByteCapacity(capacity, fn)
}

// WithEmptyByte 使用空字节切片执行函数, 自动管理获取和归还
//
// 参数:
//   - capacity: 字节切片初始容量
//   - fn: 使用字节切片的函数, 通过append等操作构建数据
//
// 返回值:
//   - []byte: 函数执行后字节切片的数据副本
//
// 使用示例:
//
//	data := pool.WithEmptyByte(1024, func(buf []byte) []byte {
//	    buf = append(buf, []byte("Hello")...)
//	    buf = append(buf, ' ')
//	    buf = append(buf, []byte("World")...)
//	    return buf
//	})
func WithEmptyByte(capacity int, fn func([]byte) []byte) []byte {
	return defaultPool.WithEmptyByte(capacity, fn)
}

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
				buf := make([]byte, 0, defaultCapacity)
				return &buf // 返回指针避免装箱
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

	bufPtr, ok := bp.pool.Get().(*[]byte)
	if !ok {
		// 类型断言失败, 创建新的
		return make([]byte, capacity)
	}

	buffer := *bufPtr

	// 缓冲区容量不足, 扩容
	if cap(buffer) < capacity {
		// 创建新的更大容量的缓冲区
		return make([]byte, capacity)
	}

	// 清空缓冲区内容并设置长度
	return buffer[:capacity]
}

// Put 归还缓冲区到对象池
//
// 参数:
//   - buffer: 要归还的缓冲区
func (bp *BytePool) Put(buffer []byte) {
	if buffer == nil {
		return // 不回收nil
	}

	// 容量小于等于最大回收容量, 归还到对象池
	if cap(buffer) <= bp.maxCapacity {
		// 清空缓冲区内容
		buffer = buffer[:0]
		bp.pool.Put(&buffer) // 传入指针避免装箱分配
		return
	}

	/* 容量大于最大回收容量, 智能缩容 */

	// 创建小容量缓冲区, 避免池变空
	newBuffer := make([]byte, 0, bp.maxCapacity)
	bp.pool.Put(&newBuffer) // 传入指针避免装箱分配
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

	bufPtr, ok := bp.pool.Get().(*[]byte)
	if !ok {
		// 类型断言失败, 创建新的
		return make([]byte, 0, capacity)
	}

	buffer := *bufPtr

	// 缓冲区容量不足, 创建新的
	if cap(buffer) < capacity {
		return make([]byte, 0, capacity)
	}

	return buffer[:0] // 返回长度为0但保持容量的切片
}

// Warm 预热对象池
//
// 参数:
//   - count: 预分配的缓冲区数量
//   - capacity: 每个缓冲区的容量
//
// 说明:
//   - 在应用启动时调用, 预分配指定数量的缓冲区
//   - 减少冷启动时的内存分配延迟
//   - 提升初期性能表现
func (bp *BytePool) Warm(count int, capacity int) {
	if count <= 0 || capacity <= 0 {
		return
	}

	// 预分配指定数量的缓冲区
	buffers := make([][]byte, count)
	for i := 0; i < count; i++ {
		buffers[i] = make([]byte, 0, capacity) // 长度0, 容量capacity
	}

	// 立即归还到对象池进行预热
	for _, buf := range buffers {
		bp.Put(buf)
	}
}

// Drain 清空对象池中的所有缓冲区
//
// 说明:
//   - 清空当前对象池中的所有缓冲区
//   - 重新创建sync.Pool, 释放可能占用的大量内存
//   - 适用于内存紧张或需要重置对象池状态的场景
func (bp *BytePool) Drain() {
	// 创建新的sync.Pool替换旧的
	bp.pool = sync.Pool{
		New: func() any {
			buf := make([]byte, 0, bp.defaultCapacity)
			return &buf // 返回指针避免装箱
		},
	}
}

// WithByte 使用默认容量的字节切片执行函数, 自动管理获取和归还
//
// 参数:
//   - fn: 使用字节切片的函数
//
// 返回值:
//   - []byte: 函数执行后字节切片的数据副本
//
// 说明:
//   - 自动从对象池获取默认容量的字节切片
//   - 执行用户提供的函数
//   - 获取字节切片数据的副本
//   - 自动归还字节切片到对象池
//   - 即使函数发生panic也会正确归还资源
func (bp *BytePool) WithByte(fn func([]byte)) []byte {
	buffer := bp.Get()
	defer bp.Put(buffer)

	fn(buffer)
	// 返回数据的副本, 避免在归还后访问
	result := make([]byte, len(buffer))
	copy(result, buffer)
	return result
}

// WithByteCapacity 使用指定容量的字节切片执行函数, 自动管理获取和归还
//
// 参数:
//   - capacity: 字节切片初始容量
//   - fn: 使用字节切片的函数
//
// 返回值:
//   - []byte: 函数执行后字节切片的数据副本
//
// 说明:
//   - 自动从对象池获取指定容量的字节切片
//   - 执行用户提供的函数
//   - 获取字节切片数据的副本
//   - 自动归还字节切片到对象池
//   - 即使函数发生panic也会正确归还资源
func (bp *BytePool) WithByteCapacity(capacity int, fn func([]byte)) []byte {
	buffer := bp.GetByteWithCapacity(capacity)
	defer bp.Put(buffer)

	fn(buffer)
	// 返回数据的副本, 避免在归还后访问
	result := make([]byte, len(buffer))
	copy(result, buffer)
	return result
}

// WithEmptyByte 使用空字节切片执行函数, 自动管理获取和归还
//
// 参数:
//   - capacity: 字节切片初始容量
//   - fn: 使用字节切片的函数, 通过append等操作构建数据
//
// 返回值:
//   - []byte: 函数执行后字节切片的数据副本
//
// 说明:
//   - 自动从对象池获取空字节切片（长度为0）
//   - 执行用户提供的函数, 函数需要返回构建后的切片
//   - 获取字节切片数据的副本
//   - 自动归还字节切片到对象池
//   - 即使函数发生panic也会正确归还资源
func (bp *BytePool) WithEmptyByte(capacity int, fn func([]byte) []byte) []byte {
	buffer := bp.GetEmpty(capacity)
	defer bp.Put(buffer)

	result := fn(buffer)
	// 返回数据的副本, 避免在归还后访问
	finalResult := make([]byte, len(result))
	copy(finalResult, result)
	return finalResult
}
