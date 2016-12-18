// esQueue
package esQueue

import (
	"runtime"
	"sync/atomic"
)

type esCache struct {
	value interface{}
	mark  bool
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
	q.cache = make([]esCache, q.capaciity)
	return q
}

func (q *EsQueue) Capaciity() uint32 {
	return q.capaciity
}

func (q *EsQueue) Quantity() uint32 {
	var putPos, getPos uint32
	var quantity uint32
	getPos = q.getPos
	putPos = q.putPos

	if putPos >= getPos {
		quantity = putPos - getPos
	} else {
		quantity = q.capMod + putPos - getPos
	}

	return quantity
}

// put queue functions
func (q *EsQueue) Put(val interface{}) (ok bool, quantity uint32) {
	var putPos, putPosNew, getPos, posCnt uint32
	var cache *esCache
	capMod := q.capMod
	for {
		getPos = q.getPos
		putPos = q.putPos

		if putPos >= getPos {
			posCnt = putPos - getPos
		} else {
			posCnt = capMod + putPos - getPos
		}

		if posCnt >= capMod {
			runtime.Gosched()
			return false, posCnt
		}

		putPosNew = putPos + 1
		if atomic.CompareAndSwapUint32(&q.putPos, putPos, putPosNew) {
			break
		} else {
			runtime.Gosched()
		}
	}

	cache = &q.cache[putPosNew&capMod]

	for {
		if !cache.mark {
			cache.value = val
			cache.mark = true
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
	for {
		putPos = q.putPos
		getPos = q.getPos

		if putPos >= getPos {
			posCnt = putPos - getPos
		} else {
			posCnt = capMod + putPos - getPos
		}

		if posCnt < 1 {
			runtime.Gosched()
			return nil, false, posCnt
		}

		getPosNew = getPos + 1
		if atomic.CompareAndSwapUint32(&q.getPos, getPos, getPosNew) {
			break
		} else {
			runtime.Gosched()
		}
	}

	cache = &q.cache[getPosNew&capMod]

	for {
		if cache.mark {
			val = cache.value
			cache.mark = false
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
