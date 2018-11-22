package producer

import (
	"math/rand"
	"strconv"
)

// EventGenerator creates a batch of random events
type EventGenerator interface {
	NewEvents(count int) []Event
}

// NewEmptyGenerator is used for testing purposes.
// Its events names are random numbers and props: [fake=>empty]
func NewEmptyGenerator() EventGenerator {
	return &fakeGenerator{}
}

type fakeGenerator struct {
}

func (f *fakeGenerator) NewEvents(count int) []Event {
	batchID := rand.Int()
	result := make([]Event, 0, count)
	for i := 0; i < count; i++ {
		name := strconv.Itoa(batchID) + ":" + strconv.Itoa(i)
		result = append(result, NewSimpleEvent(name, map[string]string{
			"fake": "empty",
		}))
	}
	return result
}
