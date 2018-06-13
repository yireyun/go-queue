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
	valCnt = uint32(16)
	values = make([]int, valCnt)
	valPut = make([]interface{}, valCnt)
)

func testQueueHigh(grp, cnt int, sched bool,
	printf func(format string, args ...interface{})) (int, int) {
	var wg sync.WaitGroup
	var Qt = newQtSum(grp)
	wg.Add(grp)
	q := sq.NewQueue(1024 * 1024)
	for i := 0; i < grp; i++ {
		go func(g int) {
			for j := 0; j < cnt; j++ {
				putCnt := uint32(0)
				puts, _ := q.Puts(valPut)
				putCnt += puts
				//var miss int32
				for putCnt < valCnt {
					//Qt.Go[g].putMiss++
					//atomic.AddInt32(&miss, 1)
					time.Sleep(time.Microsecond)
					puts, _ = q.Puts(valPut[putCnt:])
					putCnt += puts
					//if miss > 10000 {
					//panic(fmt.Sprintf("Put Fail "+
					//"putCnt:%12v, putMis:%12v, "+
					//"getCnt:%12v, getMis:%12v\n",
					//Qt.PutCnt(), Qt.PutMiss(), Qt.GetCnt(), Qt.GetMiss()))
					//}
				}
				//Qt.Go[g].putCnt++
			}
			wg.Done()
		}(i)
	}

	wg.Add(grp)
	for i := 0; i < grp; i++ {
		go func(g int) {
			var values = make([]interface{}, valCnt)
			for j := 0; j < cnt; j++ {
				//var miss int32
				getCnt := uint32(0)
				gets, _ := q.Gets(values) //该语句注释掉将导致运行结果不正确
				getCnt += gets
				for getCnt < valCnt {
					//Qt.Go[g].getMiss++
					//atomic.AddInt32(&miss, 1)
					time.Sleep(time.Microsecond * 100)
					gets, _ = q.Gets(values[getCnt:])
					getCnt += gets
					//if miss > 10000 {
					//panic(fmt.Sprintf("Get Miss "+
					//"putCnt:%12v, putMis:%12v, "+
					//"getCnt:%12v, getMis:%12v\n",
					//Qt.PutCnt(), Qt.PutMiss(),
					//Qt.GetCnt(), Qt.GetMiss()))
					//}
					//printf("Get.Fail\n")
				}
				//Qt.Go[g].getCnt++
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	return 0, int(Qt.PutMiss()) + int(Qt.GetMiss())
}

func TestQueueHigh() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	const (
		isPrintf = false
	)

	printf := func(format string, a ...interface{}) {
		if isPrintf {
			fmt.Printf(format, a...)
		}
	}

	pproF, _ := os.Create("pprof") // 创建记录文件
	pprof.StartCPUProfile(pproF)   // 开始cpu profile，结果写到文件f中
	defer pprof.StopCPUProfile()

	var remain, miss, Sum int
	var Use time.Duration
	for i := 1; i <= runtime.NumCPU()*4; i++ {
		//	for i := 1; i <= 2; i++ {
		cnt := 10000
		switch i {
		case 0, 1, 2:
			cnt = 10000 * 150
		case 3, 4, 5:
			cnt = 10000 * 100
		case 6, 7, 8:
			cnt = 10000 * 50
		case 9, 10, 11:
			cnt = 10000 * 40
		case 12, 13, 14:
			cnt = 10000 * 30
		case 15, 16:
			cnt = 10000 * 20
		}

		sum := i * cnt
		start := time.Now()
		if remain, miss = testQueueHigh(i, cnt, true, printf); remain > 0 {
			fmt.Printf("队列还剩下%d个数据\n", remain)
		}

		end := time.Now()
		use := end.Sub(start)
		op := use / time.Duration(sum)
		fmt.Printf("%v, Grp: %3d, Times: %10d, miss:%6v, use: %12v, %8v/op  %6v/op\n",
			runtime.Version(), i, sum, miss, use, op, uint32(op)/valCnt)
		Use += use
		Sum += sum
	}
	op := Use / time.Duration(Sum)
	fmt.Printf("%v %v, Grp: %3v, Times: %10d, miss:%6v, use: %12v, %8v/op  %6v/op\n",
		runtime.Version(), runtime.GOARCH, "Sum", Sum, 0, Use, op, uint32(op)/valCnt)
}

func main() {
	SetBatch(2)
	TestQueueHigh()
	fmt.Println()
	SetBatch(4)
	TestQueueHigh()
	fmt.Println()
	SetBatch(8)
	TestQueueHigh()
	fmt.Println()
	SetBatch(16)
	TestQueueHigh()
	fmt.Println()
	SetBatch(32)
	TestQueueHigh()
	fmt.Println()
	SetBatch(64)
	TestQueueHigh()
}

func SetBatch(cnt uint32) {
	fmt.Printf("----块尺寸-%v----\n", cnt)
	valCnt = cnt
	values = make([]int, valCnt)
	valPut = make([]interface{}, valCnt)
	for i := range values {
		values[i] = i
		valPut[i] = &values[i]
	}
}
