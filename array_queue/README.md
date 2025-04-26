# 无锁队列 (Lock-Free Queue)

这是一个基于Go语言实现的高性能无锁队列库，设计用于高并发环境下的生产者-消费者场景。

## 特性

- **无锁设计**：使用原子操作而非互斥锁，大幅提高并发性能
- **高性能**：针对多核处理器优化，支持高吞吐量
- **线程安全**：安全用于多线程/多协程环境
- **简单API**：提供简洁直观的接口
- **批量操作**：支持批量入队和出队操作
- **无GC压力**：高效内存使用，减少垃圾收集压力

## 安装

```bash
go get github.com/dengwuliang/lock_free_queue
```

## 使用示例

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/dengwuliang/lock_free_queue/array_queue"
)

func main() {
    // 创建一个容量为1024的队列
    q := array_queue.NewQueue(1024)
    
    // 入队操作
    value := 42
    ok, quantity := q.Put(&value)
    if ok {
        fmt.Printf("成功入队，当前队列元素数量: %d\n", quantity)
    }
    
    // 出队操作
    val, ok, quantity := q.Get()
    if ok {
        fmt.Printf("成功出队值: %d, 剩余队列元素数量: %d\n", *(val.(*int)), quantity)
    }
}
```

### 批量操作

```go
package main

import (
    "fmt"
    "github.com/dengwuliang/lock_free_queue/array_queue"
)

func main() {
    q := array_queue.NewQueue(1024)
    
    // 批量入队
    values := make([]interface{}, 10)
    for i := 0; i < 10; i++ {
        val := i
        values[i] = &val
    }
    
    puts, quantity := q.Puts(values)
    fmt.Printf("成功入队 %d 个元素，当前队列元素数量: %d\n", puts, quantity)
    
    // 批量出队
    results := make([]interface{}, 10)
    gets, quantity := q.Gets(results)
    fmt.Printf("成功出队 %d 个元素，剩余队列元素数量: %d\n", gets, quantity)
    
    // 处理结果
    for i := 0; i < int(gets); i++ {
        fmt.Printf("值: %d\n", *(results[i].(*int)))
    }
}
```

## API 参考

### 队列创建

- `NewQueue(capacity uint32) *EsQueue`: 创建一个新的无锁队列，容量会被调整为大于等于给定容量的最小2的幂

### 队列操作

- `Put(val interface{}) (ok bool, quantity uint32)`: 入队单个元素
- `Puts(values []interface{}) (puts, quantity uint32)`: 批量入队多个元素
- `Get() (val interface{}, ok bool, quantity uint32)`: 出队单个元素
- `Gets(values []interface{}) (gets, quantity uint32)`: 批量出队多个元素
- `Quantity() uint32`: 获取当前队列中的元素数量
- `Capacity() uint32`: 获取队列的容量

## 性能

该无锁队列的性能随着CPU核心数的增加而接近线性扩展。在多生产者多消费者的场景下表现出色，特别适合高并发环境。

## 实现原理

该实现基于数组环形缓冲区设计，使用原子操作避免了传统锁的性能瓶颈。通过精心设计的读写位置管理，确保了在高并发环境下的数据一致性和高吞吐量。

## 许可证

MIT License
