// esQueue
package queue

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

type esCache struct {
	putNo uint32
	getNo uint32
	value interface{}
}

// lock free queue
type EsQueue struct {
	capaciity uint32
	capMod    uint32
	putPos    uint32
	getPos    uint32
	cache     []esCache
}

func NewQueue(capaciity uint32) *EsQueue {
	q := new(EsQueue)
	q.capaciity = minQuantity(capaciity)
	q.capMod = q.capaciity - 1
	q.putPos = 0
	q.getPos = 0
	q.cache = make([]esCache, q.capaciity)
	for i := range q.cache {
		cache := &q.cache[i]
		cache.getNo = uint32(i)
		cache.putNo = uint32(i)
	}
	cache := &q.cache[0]
	cache.getNo = q.capaciity
	cache.putNo = q.capaciity
	return q
}

func (q *EsQueue) String() string {
	getPos := atomic.LoadUint32(&q.getPos)
	putPos := atomic.LoadUint32(&q.putPos)
	return fmt.Sprintf("Queue{capaciity: %v, capMod: %v, putPos: %v, getPos: %v}",
		q.capaciity, q.capMod, putPos, getPos)
}

func (q *EsQueue) Capaciity() uint32 {
	return q.capaciity
}

func (q *EsQueue) Quantity() uint32 {
	var putPos, getPos uint32
	var quantity uint32
	getPos = atomic.LoadUint32(&q.getPos)
	putPos = atomic.LoadUint32(&q.putPos)

	if putPos >= getPos {
		quantity = putPos - getPos
	} else {
		quantity = q.capMod + (putPos - getPos)
	}

	return quantity
}

// put queue functions
func (q *EsQueue) Put(val interface{}) (ok bool, quantity uint32) {
	var putPos, putPosNew, getPos, posCnt uint32
	var cache *esCache
	capMod := q.capMod

	getPos = atomic.LoadUint32(&q.getPos)
	putPos = atomic.LoadUint32(&q.putPos)

	if putPos >= getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = capMod + (putPos - getPos)
	}

	if posCnt >= capMod-1 {
		runtime.Gosched()
		return false, posCnt
	}

	putPosNew = putPos + 1
	if !atomic.CompareAndSwapUint32(&q.putPos, putPos, putPosNew) {
		runtime.Gosched()
		return false, posCnt
	}

	cache = &q.cache[putPosNew&capMod]

	for {
		getNo := atomic.LoadUint32(&cache.getNo)
		putNo := atomic.LoadUint32(&cache.putNo)
		if putPosNew == putNo && getNo == putNo {
			cache.value = val
			atomic.AddUint32(&cache.putNo, q.capaciity)
			return true, posCnt + 1
		} else {
			runtime.Gosched()
		}
	}
}

// get queue functions
func (q *EsQueue) Get() (val interface{}, ok bool, quantity uint32) {
	var putPos, getPos, getPosNew, posCnt uint32
	var cache *esCache
	capMod := q.capMod

	putPos = atomic.LoadUint32(&q.putPos)
	getPos = atomic.LoadUint32(&q.getPos)

	if putPos >= getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = capMod + (putPos - getPos)
	}

	if posCnt < 1 {
		runtime.Gosched()
		return nil, false, posCnt
	}

	getPosNew = getPos + 1
	if !atomic.CompareAndSwapUint32(&q.getPos, getPos, getPosNew) {
		runtime.Gosched()
		return nil, false, posCnt
	}

	cache = &q.cache[getPosNew&capMod]

	for {
		getNo := atomic.LoadUint32(&cache.getNo)
		putNo := atomic.LoadUint32(&cache.putNo)
		if getPosNew == getNo && getNo == putNo-q.capaciity {
			val = cache.value
			atomic.AddUint32(&cache.getNo, q.capaciity)
			return val, true, posCnt - 1
		} else {
			runtime.Gosched()
		}
	}
}

// round 到最近的2的倍数
func minQuantity(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
