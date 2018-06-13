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
				for d := range valPut {
					ok := false
					ok, _ = q.Put(valPut[d])
					var miss int32
					for !ok {
						//Qt.Go[g].putMiss++
						miss++
						time.Sleep(time.Microsecond * time.Duration(miss))
						ok, _ = q.Put(valPut[d])
						if miss > 1000 {
							panic(fmt.Sprintf("Put Fail "+
								"putCnt:%12v, putMis:%12v, "+
								"getCnt:%12v, getMis:%12v\n%s\n",
								Qt.PutCnt(), Qt.PutMiss(), Qt.GetCnt(), Qt.GetMiss(), q))
						}
					}
					//Qt.Go[g].putCnt++
				}
			}
			wg.Done()
		}(i)
	}

	wg.Add(1)
	go func(g int) {
		var values = make([]interface{}, valCnt)
		for j := 0; j < cnt*grp; j++ {
			var miss int32 = 0
			getCnt := uint32(0)
			gets, _ := q.Gets(values) //该语句注释掉将导致运行结果不正确
			getCnt += gets
			for getCnt < valCnt {
				miss++
				time.Sleep(time.Microsecond * time.Duration(miss))
				gets, _ = q.Gets(values[getCnt:])
				getCnt += gets
				if miss > 100 {
					panic(fmt.Sprintf("Get Miss "+
						"putCnt:%12v, putMis:%12v, "+
						"getCnt:%12v, getMis:%12v, gets:%12v\n%s\n",
						Qt.PutCnt(), Qt.PutMiss(),
						Qt.GetCnt(), Qt.GetMiss(), gets, q))
				}
				//printf("Get.Fail\n")
			}
			//Qt.Go[g].getCnt++
		}
		wg.Done()
	}(0)

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
	for i := 1; i <= runtime.NumCPU()*2; i++ {
		//	for i := 1; i <= 2; i++ {
		cnt := 10000
		switch i {
		case 0, 1, 2:
			cnt = 10000 * 150
		case 3, 4, 5:
			cnt = 10000 * 80
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
		sum = sum * int(valCnt)
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
	SetBatch(2)
	TestQueueHigh()
	SetBatch(4)
	TestQueueHigh()
	SetBatch(8)
	TestQueueHigh()
	SetBatch(16)
	TestQueueHigh()
	SetBatch(32)
	TestQueueHigh()
	SetBatch(64)
	TestQueueHigh()
}

func Sleep(d int) {
	var n1 = 0
	for i := 0; i < d; i++ {
		n1 *= 2
	}
	if n1 < 0 {
		n1 = 0
	}
}

//12M晶振 1毫秒延迟
func Delay(z int) {
	var n = 0
	for x := z * 2; x > 0; x-- {
		//		for y := 0; y > 0; y++ /* y > 0; y--*/ {

		//		}
		//		n *= 2
	}
	if n < 0 {
		n = 0
	}
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
