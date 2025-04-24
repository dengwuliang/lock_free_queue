package array_queue

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestChannelQueuePutGet(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU()) // 让 Go 使用所有 CPU 核心
	const (
		isPrintf = false
	)

	cnt := 10000
	sum := 0
	start := time.Now()
	var putD, getD time.Duration
	for i := 0; i <= runtime.NumCPU()*4; i++ {
		sum += i * cnt
		put, get := testChannelQueuePutGet(t, i, cnt)
		putD += put
		getD += get
	}
	end := time.Now()
	use := end.Sub(start)
	op := use / time.Duration(sum)
	t.Logf("Grp: %d, Times: %d, use: %v, %v/op", runtime.NumCPU()*4, sum, use, op)
	t.Logf("Put: %d, use: %v, %v/op", sum, putD, putD/time.Duration(sum))
	t.Logf("Get: %d, use: %v, %v/op", sum, getD, getD/time.Duration(sum))
}

func testChannelQueuePutGet(t *testing.T, grp, cnt int) (put time.Duration, get time.Duration) {
	var wg sync.WaitGroup
	var id int32
	wg.Add(grp)
	// 创建一个大小为 1024 * 1024 的队列，使用 Go 的 channel 实现
	ch := make(chan string, 1024*1024)
	start := time.Now()

	// 启动多个生产者 goroutine
	for i := 0; i < grp; i++ {
		go func(g int) {
			defer wg.Done()
			for j := 0; j < cnt; j++ {
				// 使用 atomic 来保证每个 ID 唯一
				val := fmt.Sprintf("Node.%d.%d.%d", g, j, atomic.AddInt32(&id, 1))
				// 放入 channel，如果 channel 满了则会阻塞，直到有空间
				ch <- val
			}
		}(i)
	}
	wg.Wait()
	end := time.Now()
	put = end.Sub(start)

	wg.Add(grp)
	start = time.Now()

	// 启动多个消费者 goroutine
	for i := 0; i < grp; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < cnt; {
				select {
				case <-ch:
					j++ // 成功获取元素
				default:
					runtime.Gosched() // 如果没有元素，调用 `Gosched` 让出时间片，避免空转
				}
			}
		}()
	}
	wg.Wait()
	end = time.Now()
	get = end.Sub(start)

	// 检查 channel 是否为空，确保所有元素都被消费
	if len(ch) != 0 {
		t.Errorf("Grp:%v, Channel not empty. Remaining items: %v", grp, len(ch))
	}

	return put, get
}
