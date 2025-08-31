package pool_test

import (
	"fmt"
	"strconv"

	"gitee.com/MM-Q/go-kit/pool"
)

// ExampleGetBuffer 演示字节缓冲区对象池的基本使用
func ExampleGetBuffer() {
	// 获取字节缓冲区
	buffer := pool.GetBuffer(100)

	// 写入数据
	buffer.WriteString("Hello, ")
	buffer.WriteString("World!")
	buffer.WriteString(" Number: ")
	buffer.WriteString(strconv.Itoa(42))

	// 获取结果
	result := buffer.String()
	fmt.Println(result)

	// 归还到对象池
	pool.PutBuffer(buffer)

	// Output: Hello, World! Number: 42
}

// ExampleBufferPool 演示自定义缓冲区池的使用
func ExampleBufferPool() {
	// 创建自定义缓冲区池
	bufferPool := pool.NewBufferPool(512, 32*1024)

	// 预热对象池
	bufferPool.Warm(10, 1024)

	// 获取字节缓冲区
	buffer := bufferPool.Get(200)

	// 构建二进制数据
	buffer.WriteByte(0x48) // 'H'
	buffer.WriteByte(0x65) // 'e'
	buffer.WriteByte(0x6C) // 'l'
	buffer.WriteByte(0x6C) // 'l'
	buffer.WriteByte(0x6F) // 'o'

	// 写入字符串
	buffer.WriteString(", Buffer Pool!")

	// 获取结果
	result := buffer.String()
	fmt.Println(result)

	// 获取字节数据
	data := buffer.Bytes()
	fmt.Printf("Length: %d, Capacity: %d\n", len(data), buffer.Cap())

	// 归还到对象池
	bufferPool.Put(buffer)

	// Output: Hello, Buffer Pool!
	// Length: 19, Capacity: 512
}