package id

import (
	"math/rand"
	"sync"
	"time"
)

// 随机数生成器池
// 用于复用随机数生成器，避免频繁创建和销毁
var pool = sync.Pool{
	New: func() interface{} {
		return rand.New(rand.NewSource(time.Now().UnixNano()))
	},
}

// getRand 获取随机数生成器
// 从池中获取随机数生成器
// 如果池为空，则创建一个新的随机数生成器
//
// 返回:
//   - 随机数生成器
func getRand() *rand.Rand {
	if r := pool.Get(); r != nil {
		if gen, ok := r.(*rand.Rand); ok {
			return gen
		}
	}
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// putRand 归还随机数生成器
// 将随机数生成器归还到池中，以便后续复用
//
// 参数:
//   - r: 要归还的随机数生成器
func putRand(r *rand.Rand) {
	if r != nil {
		pool.Put(r)
	}
}
