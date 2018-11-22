package sortedAggr

import (
	"container/heap"
	"sync"
)

type Elem struct {
	Key   string
	Value float32
}

type innerMaxHeap []*Elem

func (h innerMaxHeap) Len() int           { return len(h) }
func (h innerMaxHeap) Less(i, j int) bool { return h[i].Value > h[j].Value }
func (h innerMaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *innerMaxHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*Elem))
}

func (h *innerMaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

//linear search for a specific key
func (h *innerMaxHeap) getIndex(key string) int {
	for i, elem := range *h {
		if elem.Key == key {
			return i
		}
	}
	return -1
}

type AggHeap struct {
	data  *innerMaxHeap
	mutex sync.RWMutex
}

func NewAggregateHeap() *AggHeap {
	result := &AggHeap{}
	heap.Init(result.data)

	return result
}

func (ag *AggHeap) Increment(key string, value float32) {
	ag.mutex.RLock()
	indexOf := ag.data.getIndex(key)
	ag.mutex.RUnlock()

	ag.mutex.Lock()
	defer ag.mutex.Unlock()

	if indexOf == -1 {
		elem := &Elem{key, value}
		heap.Push(ag.data, elem)
		return
	}

	elem := (*ag.data)[indexOf]
	elem.Value += value
	heap.Fix(ag.data, indexOf)
}
