package sortedAggr

import (
	"container/list"
	"sync"
)

type Latest struct {
	capacity int
	data     *list.List
	mutex    sync.Mutex
}

// NewLatest returns a queue of strings
func NewLatest(capacity int) *Latest {
	result := &Latest{}
	result.data = list.New()
	result.capacity = capacity
	return result
}

func (l *Latest) Add(elem string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.data.Len() >= l.capacity {
		l.data.Remove(l.data.Front())
	}
	l.data.PushBack(elem)
}

// Get returns a copy of the latest data
func (l *Latest) Get() []string {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	size := l.data.Len()
	result := make([]string, 0, size)
	node := l.data.Back()

	for node != nil {
		//in reverse order
		result = append(result, node.Value.(string))
		node = node.Prev()
	}
	return result
}
