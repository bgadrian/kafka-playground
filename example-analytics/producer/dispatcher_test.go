package producer

import (
	"testing"
	"time"
)

func TestNewDispatcher(t *testing.T) {
	pattern := []float64{0, 0.5, 1}
	should := []int{0, 50, 100} //peak * each cycle factor
	traffic := NewTrafficPattern(100, pattern)

	maxTicks := len(pattern) * 2
	tickSize := time.Millisecond * time.Duration(20)

	disp := NewDispatcher(traffic, NewDummyUser)

	disp.Start(tickSize)

	//wait for tick 0
	time.Sleep(time.Millisecond * 3)

	for tick := 0; tick < maxTicks; tick++ {
		got := len(disp.online)
		wanted := should[tick%len(should)]
		if got != wanted {
			t.Errorf("for tick %d got %d should %d",
				tick, got, wanted)
		}
		//we have to wait for the Tick to finish
		time.Sleep(tickSize + time.Millisecond)
	}
}
