package producer

import (
	"math/rand"
	"time"
)

// Dispatcher makes sure the users are spawned and killed accordnly to the TrafficPattern pattern
type Dispatcher struct {
	running    bool
	online     []User
	traffic    TrafficPattern
	usersQueen UserGenerator
	//lock         sync.Mutex //for Stop and Ticker in the same time
}

// NewDispatcher is making sure that [traffic.Next()] users are running and in
// total generating [eventsPerUserMin] events.
// The Traffic.next()] is called each [intervalSize].
func NewDispatcher(traffic TrafficPattern, queen UserGenerator) *Dispatcher {
	d := &Dispatcher{}
	d.traffic = traffic
	d.usersQueen = queen
	return d
}

// Start ticking and adjusting the online users based on the traffic pattern
func (d *Dispatcher) Start(intervalSize time.Duration) {
	if d.running {
		return
	}

	ticker := time.NewTicker(intervalSize)
	d.running = true

	//create first users
	d.tick()

	go func() {
	exit:
		for range ticker.C {
			if d.running == false {
				break exit
			}
			//must block to avoid concurrent d.online updates
			d.tick()
		}
		ticker.Stop()
		d.online = nil
	}()
}

// Stop pauses the dispatcher, the online users will still SEND events, but their numbers will not fluctuate any more
func (d *Dispatcher) Stop() {
	d.running = false
}

func (d *Dispatcher) tick() {
	usersShould := d.traffic.GetValueAndAdvance()
	usersNow := len(d.online)

	if usersNow == usersShould {
		return
	}

	diff := usersShould - usersNow
	if diff > 0 {
		//add more users
		for i := 0; i < diff; i++ {
			user := d.usersQueen()
			d.online = append(d.online, user)
		}
		return
	}

	//remove random users
	for len(d.online) > usersShould {
		killIndex := rand.Intn(len(d.online))
		user := d.online[killIndex]
		d.online = append(d.online[:killIndex], d.online[killIndex+1:]...)
		user.Kill()
	}
}
