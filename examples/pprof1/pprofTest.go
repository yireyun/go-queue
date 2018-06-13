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

type QtObj struct {
	getMiss int32
	putMiss int32
	putCnt  int32
	getCnt  int32
}

type QtSum struct {
	Go []QtObj
}

func newQtSum(grp int) *QtSum {
	qt := new(QtSum)
	qt.Go = make([]QtObj, grp)
	return qt
}

func (q *QtSum) GetMiss() (num int32) {
	for i := range q.Go {
		num += q.Go[i].getMiss
	}
	return
}
func (q *QtSum) PutMiss() (num int32) {
	for i := range q.Go {
		num += q.Go[i].putMiss
	}
	return
}
func (q *QtSum) PutCnt() (num int32) {
	for i := range q.Go {
		num += q.Go[i].putCnt
	}
	return
}
func (q *QtSum) GetCnt() (num int32) {
	for i := range q.Go {
		num += q.Go[i].getCnt
	}
	return
}

var (
	value = 1
)

func testQueueHigh(grp, cnt int) int {
	var wg sync.WaitGroup
	var Qt = newQtSum(grp)
	wg.Add(grp)
	q := sq.NewQueue(1024 * 1024)
	for i := 0; i < grp; i++ {
		go func(g int) {
			ok := false
			for j := 0; j < cnt; j++ {
				ok, _ = q.Put(&value)
				var miss int32
				for !ok {
					Qt.Go[g].putMiss++
					atomic.AddInt32(&miss, 1)
					time.Sleep(time.Microsecond)
					ok, _ = q.Put(&value)
					if miss > 10000 {
						panic(fmt.Sprintf("Put Fail "+
							"putCnt:%12v, putMis:%12v, "+
							"getCnt:%12v, getMis:%12v ",
							Qt.PutCnt(), Qt.PutMiss(),
							Qt.GetCnt(), Qt.GetMiss()))
					}
				}
				Qt.Go[g].putCnt++
			}
			wg.Done()
		}(i)
	}
	wg.Add(grp)
	for i := 0; i < grp; i++ {
		go func(g int) {
			ok := false
			for j := 0; j < cnt; j++ {
				var miss int32
				_, ok, _ = q.Get() //该语句注释掉将导致运行结果不正确
				for !ok {
					Qt.Go[g].getMiss++
					atomic.AddInt32(&miss, 1)
					time.Sleep(time.Microsecond * 100)
					_, ok, _ = q.Get()
					if miss > 10000 {
						panic(fmt.Sprintf("Get Miss "+
							"putCnt:%12v, putMis:%12v, "+
							"getCnt:%12v, getMis:%12v\n",
							Qt.PutCnt(), Qt.PutMiss(),
							Qt.GetCnt(), Qt.GetMiss()))
					}
				}
				Qt.Go[g].getCnt++
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return int(Qt.PutMiss()) + int(Qt.GetMiss())
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
		fmt.Printf("%v, Grp: %3d, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
			runtime.Version(), i, sum, miss, use, op)
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
