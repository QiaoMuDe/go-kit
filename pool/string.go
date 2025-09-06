package pool

import (
	"strings"
	"sync"
)

// 全局默认字符串构建器池实例
//
// 说明:
//   - 该实例用于在全局范围内管理字符串构建器对象，避免频繁创建和销毁对象导致的性能问题。
//   - 初始容量为256，最大回收大小为32KB。
var defaultStringPool = NewStringPool(256, 32*1024)

// GetString 从默认字符串池获取默认容量的字符串构建器
//
// 返回值:
//   - *strings.Builder: 容量至少为默认大小的字符串构建器
func GetString() *strings.Builder {
	return defaultStringPool.Get()
}

// GetStringWithCapacity 从默认字符串池获取指定容量的字符串构建器
//
// 参数:
//   - capacity: 字符串构建器初始容量大小
//
// 返回值:
//   - *strings.Builder: 容量至少为capacity的字符串构建器
func GetStringWithCapacity(capacity int) *strings.Builder {
	return defaultStringPool.GetWithCapacity(capacity)
}

// PutString 将字符串构建器归还到默认字符串池
//
// 参数:
//   - builder: 要归还的字符串构建器
//
// 说明:
//   - 该函数将字符串构建器归还到对象池，以便后续复用。
//   - 只有容量不超过64KB的构建器才会被归还，以避免对象池占用过多内存。
func PutString(builder *strings.Builder) {
	defaultStringPool.Put(builder)
}

// SetStringMaxSize 动态调整默认字符串池的最大回收大小
//
// 参数:
//   - maxSize: 新的最大回收大小
func SetStringMaxSize(maxSize int) {
	defaultStringPool.SetMaxSize(maxSize)
}

// GetStringMaxSize 获取默认字符串池的当前最大回收大小
//
// 返回值:
//   - int: 当前最大回收大小
func GetStringMaxSize() int {
	return defaultStringPool.GetMaxSize()
}

// WarmString 预热默认字符串池
//
// 参数:
//   - count: 预分配的字符串构建器数量
//   - capacity: 每个字符串构建器的容量
func WarmString(count int, capacity int) {
	defaultStringPool.Warm(count, capacity)
}

// DrainString 清空默认字符串池
func DrainString() {
	defaultStringPool.Drain()
}

// WithString 使用默认容量的字符串构建器执行函数，自动管理获取和归还
//
// 参数:
//   - fn: 使用字符串构建器的函数
//
// 返回值:
//   - string: 函数执行后构建的字符串结果
//
// 使用示例:
//
//	result := pool.WithString(func(buf *strings.Builder) {
//	    buf.WriteString("Hello")
//	    buf.WriteByte(' ')
//	    buf.WriteString("World")
//	})
func WithString(fn func(*strings.Builder)) string {
	return defaultStringPool.WithString(fn)
}

// WithStringCapacity 使用指定容量的字符串构建器执行函数，自动管理获取和归还
//
// 参数:
//   - capacity: 字符串构建器初始容量大小
//   - fn: 使用字符串构建器的函数
//
// 返回值:
//   - string: 函数执行后构建的字符串结果
//
// 使用示例:
//
//	result := pool.WithStringCapacity(64, func(buf *strings.Builder) {
//	    buf.WriteString("Hello")
//	    buf.WriteByte(' ')
//	    buf.WriteString("World")
//	})
func WithStringCapacity(capacity int, fn func(*strings.Builder)) string {
	return defaultStringPool.WithStringCapacity(capacity, fn)
}

// StringPool 字符串构建器对象池，支持自定义配置
type StringPool struct {
	pool        sync.Pool // 字符串构建器对象池
	maxSize     int       // 最大回收构建器大小
	defaultSize int       // 默认构建器大小
}

// NewStringPool 创建新的字符串构建器对象池
//
// 参数:
//   - defaultSize: 默认字符串构建器容量大小
//   - maxSize: 最大回收构建器大小，超过此大小的构建器不会被回收
//
// 返回值:
//   - *StringPool: 字符串构建器对象池实例
func NewStringPool(defaultSize, maxSize int) *StringPool {
	if defaultSize <= 0 {
		defaultSize = 256 // 默认256字节
	}
	if maxSize <= 0 {
		maxSize = 32 * 1024 // 默认32KB
	}

	return &StringPool{
		maxSize:     maxSize,
		defaultSize: defaultSize,
		pool: sync.Pool{
			New: func() any {
				builder := &strings.Builder{}
				builder.Grow(defaultSize) // 预分配容量
				return builder
			},
		},
	}
}

// Get 获取默认容量的字符串构建器
//
// 返回:
//   - *strings.Builder: 容量至少为默认大小的字符串构建器
//
// 说明:
//   - 返回的字符串构建器已经重置为空状态，可以直接使用
//   - 底层容量可能大于默认大小，来自对象池的复用构建器
func (sp *StringPool) Get() *strings.Builder {
	return sp.GetWithCapacity(sp.defaultSize)
}

// GetWithCapacity 获取指定容量的字符串构建器
//
// 参数:
//   - capacity: 需要的字符串构建器容量大小
//
// 返回:
//   - *strings.Builder: 容量至少为capacity的字符串构建器
//
// 说明:
//   - 返回的字符串构建器已经重置为空状态，可以直接使用
//   - 底层容量可能大于capacity，来自对象池的复用构建器
func (sp *StringPool) GetWithCapacity(capacity int) *strings.Builder {
	builder, ok := sp.pool.Get().(*strings.Builder)
	if !ok {
		// 类型断言失败，创建新的
		builder = &strings.Builder{}
		builder.Grow(capacity) // 预分配容量
		builder.Reset()
		return builder
	}

	// 如果当前容量不足，扩容到所需大小
	if builder.Cap() < capacity {
		builder.Grow(capacity - builder.Cap())
	}

	// 重置构建器状态
	builder.Reset()

	return builder
}

// Put 归还字符串构建器到对象池
//
// 参数:
//   - builder: 要归还的字符串构建器
//
// 说明:
//   - nil构建器不会被回收
//   - 容量不超过maxSize的构建器直接重置后归还
//   - 容量超过maxSize的构建器会创建一个新的小容量构建器进行归还（智能缩容）
func (sp *StringPool) Put(builder *strings.Builder) {
	// 不回收nil构建器
	if builder == nil {
		return
	}

	// 如果容量不超过最大回收大小，直接重置后归还
	if builder.Cap() <= sp.maxSize {
		builder.Reset()
		sp.pool.Put(builder)
		return
	}

	// 对于容量超过最大回收大小的构建器，创建一个新的小容量构建器进行归还
	// 这样可以避免大容量构建器占用过多内存，同时保持对象池的复用性
	newBuilder := &strings.Builder{}
	newBuilder.Grow(sp.maxSize) // 预分配容量为maxSize
	newBuilder.Reset()
	sp.pool.Put(newBuilder)
}

// SetMaxSize 动态调整最大回收构建器大小
//
// 参数:
//   - maxSize: 新的最大回收大小
//
// 说明:
//   - 运行时动态调整配置
//   - 如果新的maxSize小于当前值，建议调用Drain()清空对象池
func (sp *StringPool) SetMaxSize(maxSize int) {
	if maxSize <= 0 {
		maxSize = 32 * 1024 // 默认32KB
	}
	sp.maxSize = maxSize
}

// GetMaxSize 获取当前最大回收构建器大小
//
// 返回:
//   - int: 当前最大回收大小
func (sp *StringPool) GetMaxSize() int {
	return sp.maxSize
}

// Warm 预热对象池
//
// 参数:
//   - count: 预分配的字符串构建器数量
//   - capacity: 每个字符串构建器的容量
//
// 说明:
//   - 在应用启动时调用，预分配指定数量的字符串构建器
//   - 减少冷启动时的内存分配延迟
//   - 提升初期性能表现
func (sp *StringPool) Warm(count int, capacity int) {
	if count <= 0 || capacity <= 0 {
		return
	}

	// 预分配指定数量的字符串构建器
	builders := make([]*strings.Builder, count)
	for i := 0; i < count; i++ {
		builder := &strings.Builder{}
		builder.Grow(capacity)
		builders[i] = builder
	}

	// 立即归还到对象池进行预热
	for _, builder := range builders {
		sp.Put(builder)
	}
}

// Drain 清空对象池中的所有字符串构建器
//
// 说明:
//   - 清空当前对象池中的所有字符串构建器
//   - 重新创建sync.Pool，释放可能占用的大量内存
//   - 适用于内存紧张或需要重置对象池状态的场景
func (sp *StringPool) Drain() {
	// 创建新的sync.Pool替换旧的
	sp.pool = sync.Pool{
		New: func() any {
			builder := &strings.Builder{}
			builder.Grow(sp.defaultSize) // 预分配容量
			return builder
		},
	}
}

// WithString 使用默认容量的字符串构建器执行函数，自动管理获取和归还
//
// 参数:
//   - fn: 使用字符串构建器的函数
//
// 返回值:
//   - string: 函数执行后构建的字符串结果
//
// 说明:
//   - 自动从对象池获取默认容量的字符串构建器
//   - 执行用户提供的函数
//   - 获取构建的字符串结果
//   - 自动归还字符串构建器到对象池
//   - 即使函数发生panic也会正确归还资源
func (sp *StringPool) WithString(fn func(*strings.Builder)) string {
	builder := sp.Get()
	defer sp.Put(builder)

	fn(builder)
	return builder.String()
}

// WithStringCapacity 使用指定容量的字符串构建器执行函数，自动管理获取和归还
//
// 参数:
//   - capacity: 字符串构建器初始容量大小
//   - fn: 使用字符串构建器的函数
//
// 返回值:
//   - string: 函数执行后构建的字符串结果
//
// 说明:
//   - 自动从对象池获取指定容量的字符串构建器
//   - 执行用户提供的函数
//   - 获取构建的字符串结果
//   - 自动归还字符串构建器到对象池
//   - 即使函数发生panic也会正确归还资源
func (sp *StringPool) WithStringCapacity(capacity int, fn func(*strings.Builder)) string {
	builder := sp.GetWithCapacity(capacity)
	defer sp.Put(builder)

	fn(builder)
	return builder.String()
}
