// esQueue_test
package queue

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	q := NewQueue(8)
	ok, quantity := q.Put(&value)
	if !ok {
		t.Error("TestStack Get.Fail")
		return
	} else {
		t.Logf("TestStack Put value:%d[%v], quantity:%v\n", &value, value, quantity)
	}

	val, ok, quantity := q.Get()
	if !ok {
		t.Error("TestStack Get.Fail")
		return
	} else {
		t.Logf("TestStack Get value:%d[%v], quantity:%v\n", val, *(val.(*int)), quantity)
	}
	if q := q.Quantity(); q != 0 {
		t.Errorf("Quantity Error: [%v] <>[%v]", q, 0)
	}
}

func TestQueuePutGet(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	const (
		isPrintf = false
	)

	cnt := 10000
	sum := 0
	start := time.Now()
	var putD, getD time.Duration
	for i := 0; i <= runtime.NumCPU()*4; i++ {
		sum += i * cnt
		put, get := testQueuePutGet(t, i, cnt)
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

func TestQueueGeneral(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	const (
		isPrintf = false
	)

	var miss, Sum int
	var Use time.Duration
	for i := 1; i <= runtime.NumCPU()*4; i++ {
		cnt := 10000 * 10
		if i > 9 {
			cnt = 10000 * 1
		}
		sum := i * cnt
		start := time.Now()
		miss = testQueueGeneral(t, i, cnt)
		end := time.Now()
		use := end.Sub(start)
		op := use / time.Duration(sum)
		fmt.Printf("%v, Grp: %3d, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
			runtime.Version(), i, sum, miss, use, op)
		Use += use
		Sum += sum
	}
	op := Use / time.Duration(Sum)
	fmt.Printf("%v, Grp: %3v, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
		runtime.Version(), "Sum", Sum, 0, Use, op)
}

func TestQueuePutGoGet(t *testing.T) {
	var Sum, miss int
	var Use time.Duration
	for i := 1; i <= runtime.NumCPU()*4; i++ {
		//	for i := 2; i <= 2; i++ {
		cnt := 10000 * 100
		if i > 9 {
			cnt = 10000 * 10
		}
		sum := i * cnt
		start := time.Now()
		miss = testQueuePutGoGet(t, i, cnt)

		end := time.Now()
		use := end.Sub(start)
		op := use / time.Duration(sum)
		fmt.Printf("%v, Grp: %3d, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
			runtime.Version(), i, sum, miss, use, op)
		Use += use
		Sum += sum
	}
	op := Use / time.Duration(Sum)
	fmt.Printf("%v, Grp: %3v, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
		runtime.Version(), "Sum", Sum, 0, Use, op)
}

func TestQueuePutDoGet(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var miss, Sum int
	var Use time.Duration
	for i := 1; i <= runtime.NumCPU()*4; i++ {
		//	for i := 2; i <= 2; i++ {
		cnt := 10000 * 100
		if i > 9 {
			cnt = 10000 * 10
		}
		sum := i * cnt
		start := time.Now()
		miss = testQueuePutDoGet(t, i, cnt)
		end := time.Now()
		use := end.Sub(start)
		op := use / time.Duration(sum)
		fmt.Printf("%v, Grp: %3d, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
			runtime.Version(), i, sum, miss, use, op)
		Use += use
		Sum += sum
	}
	op := Use / time.Duration(Sum)
	fmt.Printf("%v, Grp: %3v, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
		runtime.Version(), "Sum", Sum, 0, Use, op)
}

func testQueuePutGet(t *testing.T, grp, cnt int) (
	put time.Duration, get time.Duration) {
	var wg sync.WaitGroup
	var id int32
	wg.Add(grp)
	q := NewQueue(1024 * 1024)
	start := time.Now()
	for i := 0; i < grp; i++ {
		go func(g int) {
			defer wg.Done()
			for j := 0; j < cnt; j++ {
				t := fmt.Sprintf("Node.%d.%d.%d", g, j, atomic.AddInt32(&id, 1))
				q.Put(t)
			}
		}(i)
	}
	wg.Wait()
	end := time.Now()
	put = end.Sub(start)

	wg.Add(grp)
	start = time.Now()
	for i := 0; i < grp; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < cnt; {
				_, ok, _ := q.Get()
				if !ok {
					runtime.Gosched()
				} else {
					j++
				}
			}
		}()
	}
	wg.Wait()
	end = time.Now()
	get = end.Sub(start)
	if q := q.Quantity(); q != 0 {
		t.Errorf("Grp:%v, Quantity Error: [%v] <>[%v]", grp, q, 0)
	}
	return put, get
}

func testQueueGeneral(t *testing.T, grp, cnt int) int {

	var wg sync.WaitGroup
	var idPut, idGet int32
	var miss = 0

	wg.Add(grp)
	q := NewQueue(1024 * 1024)
	for i := 0; i < grp; i++ {
		go func(g int) {
			defer wg.Done()
			for j := 0; j < cnt; j++ {
				t := fmt.Sprintf("Node.%d.%d.%d", g, j, atomic.AddInt32(&idPut, 1))
				q.Put(t)
			}
		}(i)
	}

	wg.Add(grp)
	for i := 0; i < grp; i++ {
		go func(g int) {
			defer wg.Done()
			ok := false
			for j := 0; j < cnt; j++ {
				_, ok, _ = q.Get() //该语句注释掉将导致运行结果不正确
				for !ok {
					miss++
					time.Sleep(time.Microsecond * 50)
					_, ok, _ = q.Get()
				}
				atomic.AddInt32(&idGet, 1)
			}
		}(i)
	}
	wg.Wait()
	if q := q.Quantity(); q != 0 {
		t.Errorf("Grp:%v, Quantity Error: [%v] <>[%v]", grp, q, 0)
	}
	return miss
}

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
	value int = 1
)

func testQueuePutGoGet(t *testing.T, grp, cnt int) int {
	var wg sync.WaitGroup
	//var Qt = newQtSum(grp)
	wg.Add(grp)
	q := NewQueue(1024 * 1024)
	for i := 0; i < grp; i++ {
		go func(g int) {
			ok := false
			for j := 0; j < cnt; j++ {
				ok, _ = q.Put(&value)
				//var miss int32
				for !ok {
					//Qt.Go[g].getMiss++
					//miss++
					//time.Sleep(time.Microsecond)
					ok, _ = q.Put(&value)
					//if miss > 10000 {
					//	panic(fmt.Sprintf("Put Fail PutId:%12v, GetId:%12v, "+
					//		"putCnt:%12v, putMis:%12v, "+
					//		"getCnt:%12v, getMis:%12v\n",
					//		q.eqPut, q.eqGet, Qt.PutCnt(), Qt.PutMiss(), Qt.GetCnt(), Qt.GetMiss()))
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
			ok := false
			for j := 0; j < cnt; j++ {
				//var miss int32
				_, ok, _ = q.Get() //该语句注释掉将导致运行结果不正确
				for !ok {
					//Qt.Go[g].putMiss++
					//miss++
					//time.Sleep(time.Microsecond * 100)
					_, ok, _ = q.Get()
					//if miss > 10000 {
					//	panic(fmt.Sprintf("Get Miss PutId:%12v, GetId:%12v, "+
					//		"putCnt:%12v, putMis:%12v, "+
					//		"getCnt:%12v, getMis:%12v\n",
					//		q.eqPut, q.eqGet, Qt.PutCnt(), Qt.PutMiss(),
					//		Qt.GetCnt(), Qt.GetMiss()))
					//}
					//printf("Get.Fail\n")
				}
				//Qt.Go[g].getCnt++
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return 0 //int(Qt.PutMiss()) + int(Qt.GetMiss())
}

func testQueuePutDoGet(t *testing.T, grp, cnt int) int {
	var wg sync.WaitGroup
	//var Qt = newQtSum(grp)
	wg.Add(grp)
	q := NewQueue(1024 * 1024)
	for i := 0; i < grp; i++ {
		go func(g int) {
			ok := false
			for j := 0; j < cnt; j++ {
				ok, _ = q.Put(&value)
				//var missPut int32
				for !ok {
					//Qt.Go[g].getMiss++
					//missPut++
					//time.Sleep(time.Microsecond)
					ok, _ = q.Put(&value)
					//if missPut > 10000 {
					//	panic(fmt.Sprintf("Put Fail PutId:%12v, GetId:%12v, "+
					//		"putCnt:%12v, putMis:%12v, "+
					//		"getCnt:%12v, getMis:%12v\n",
					//		q.eqPut, q.eqGet, Qt.PutCnt(), Qt.PutMiss(), Qt.GetCnt(), Qt.GetMiss()))
					//}
				}
				//Qt.Go[g].putCnt++

				//var missGet int32
				_, ok, _ = q.Get() //该语句注释掉将导致运行结果不正确
				for !ok {
					//Qt.Go[g].putMiss++
					//missGet++
					//time.Sleep(time.Microsecond * 100)
					_, ok, _ = q.Get()
					//if missGet > 10000 {
					//	panic(fmt.Sprintf("Get Miss PutId:%12v, GetId:%12v, "+
					//		"putCnt:%12v, putMis:%12v, "+
					//		"getCnt:%12v, getMis:%12v\n",
					//		q.eqPut, q.eqGet, Qt.PutCnt(), Qt.PutMiss(),
					//		Qt.GetCnt(), Qt.GetMiss()))
					//}
					//printf("Get.Fail\n")
				}
				//Qt.Go[g].getCnt++
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return 0 //int(Qt.PutMiss()) + int(Qt.GetMiss())
}

func testQueuePutGetOrder(t *testing.T, grp, cnt int) (
	residue int) {
	var wg sync.WaitGroup
	var idPut, idGet int32
	wg.Add(grp)
	q := NewQueue(1024 * 1024)
	for i := 0; i < grp; i++ {
		go func(g int) {
			defer wg.Done()
			for j := 0; j < cnt; j++ {
				v := atomic.AddInt32(&idPut, 1)
				q.Put(v)
			}
		}(i)
	}
	wg.Wait()
	wg.Add(grp)
	for i := 0; i < grp; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < cnt; {
				val, ok, _ := q.Get()
				if !ok {
					fmt.Printf("Get.Fail\n")
					runtime.Gosched()
				} else {
					j++
					idGet++
					if idGet != val.(int32) {
						t.Logf("Get.Err %d <> %d\n", idGet, val)
					}
				}
			}
		}()
	}
	wg.Wait()
	return
}

func TestQueuePutGetOrder(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	grp := 1
	cnt := 100

	testQueuePutGetOrder(t, grp, cnt)
	t.Logf("Grp: %d, Times: %d", grp, cnt)
}
