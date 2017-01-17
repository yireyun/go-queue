// esQueue_test
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	sq "github.com/yireyun/go-queue"
)

var (
	value = 1
)

func testQueueHigh(grp, cnt int) int {
	var wg sync.WaitGroup
	wg.Add(grp)
	q := sq.NewQueue(1024 * 1024)
	for i := 0; i < grp; i++ {
		go func(g int) {
			ok := false
			for j := 0; j < cnt; j++ {
				ok, _ = q.Put(&value)
				for !ok {
					time.Sleep(time.Microsecond)
					ok, _ = q.Put(&value)
				}
			}
			wg.Done()
		}(i)
	}
	wg.Add(grp)
	for i := 0; i < grp; i++ {
		go func(g int) {
			ok := false
			for j := 0; j < cnt; j++ {
				_, ok, _ = q.Get() //该语句注释掉将导致运行结果不正确
				for !ok {
					time.Sleep(time.Microsecond * 100)
					_, ok, _ = q.Get()
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return 0
}

func TestQueueHigh() {
	pproF, _ := os.Create("pprof") // 创建记录文件
	pprof.StartCPUProfile(pproF)   // 开始cpu profile，结果写到文件f中
	defer pprof.StopCPUProfile()

	var miss, Sum int
	var Use time.Duration
	for i := 1; i <= runtime.NumCPU()*4; i++ {
		cnt := 10000 * 1000
		if i > 9 {
			cnt = 10000 * 100
		}
		sum := i * cnt
		start := time.Now()
		miss = testQueueHigh(i, cnt)
		end := time.Now()
		use := end.Sub(start)
		op := use / time.Duration(sum)
		fmt.Printf("%v %v, Grp: %3d, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
			runtime.Version(), runtime.GOARCH, i, sum, miss, use, op)
		Use += use
		Sum += sum
	}
	op := Use / time.Duration(Sum)
	fmt.Printf("%v %v, Grp: %3v, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
		runtime.Version(), runtime.GOARCH, "Sum", Sum, 0, Use, op)
}

func main() {
	TestQueueHigh()
}
