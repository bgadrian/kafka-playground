package sortedAggr

import (
	"fmt"
	"sort"
	"sync"
)

type AggArray struct {
	byKey  map[string]*Elem
	sorted []*Elem
	mutex  sync.Mutex
}

func NewAggArray() *AggArray {
	result := &AggArray{}
	result.byKey = make(map[string]*Elem, 100)
	result.sorted = make([]*Elem, 0, 100)
	return result
}

func (ag *AggArray) Increment(key string, value float32) {
	ag.mutex.Lock()
	defer ag.mutex.Unlock()

	elem, exists := ag.byKey[key]
	if exists {
		elem.Value += value
	} else {
		elem = &Elem{key, value}
		ag.byKey[key] = elem
		ag.sorted = append(ag.sorted, elem)
	}

	//TODO replace with a more efficient sorting alg on partial sorts
	sort.Slice(ag.sorted, func(i, j int) bool {
		return ag.sorted[i].Value > ag.sorted[j].Value
	})
}

func (ag *AggArray) TopList() []string {
	elems := ag.Top(10)
	result := make([]string, len(elems))
	for i, elem := range elems {
		result[i] = fmt.Sprintf("%s: %.2f", elem.Key, elem.Value)
	}
	return result
}

func (ag *AggArray) Top(n int) []*Elem {
	ag.mutex.Lock()
	defer ag.mutex.Unlock()

	if n > len(ag.sorted) {
		n = len(ag.sorted)
	}

	result := make([]*Elem, n)
	for i := 0; i < n; i++ {
		result[i] = ag.sorted[i]
	}
	return result
}
