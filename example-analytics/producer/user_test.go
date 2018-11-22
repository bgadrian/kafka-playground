package producer

import (
	"testing"
	"time"
)

func TestNewSimpleUser(t *testing.T) {
	gen := NewEmptyGenerator()
	collector := make(chan Event, 1000)

	tickMs := 10
	ticksPerMinute := 60000 / tickMs
	eventsPerMinute := 1000
	eventsPerTick := eventsPerMinute / ticksPerMinute

	ticksToWait := 15
	testTimeMs := ticksToWait * tickMs
	testShouldCount := ticksToWait * eventsPerTick

	u := NewSimpleUser(ticksPerMinute, gen, collector,
		tickMs, tickMs)
	gotCount := 0

	go func() {
		for range collector {
			gotCount++
		}
	}()
	time.Sleep(time.Millisecond * time.Duration(testTimeMs))
	close(collector)
	if testShouldCount != gotCount {
		t.Errorf("wrong count events generated, got %d wanted %d", gotCount, testShouldCount)
	}

	u.Kill()

	//if the user tries to write it will panic, because
	//we closed the collector, this way we test the Kill
	time.Sleep(time.Millisecond * time.Duration(testTimeMs))

}
