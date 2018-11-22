package producer

// TrafficPattern simulates each hour traffic amount for a normal week day
type TrafficPattern interface {
	GetValueAndAdvance() int
}

// count simulates a 24h pattern of online users, with a peak.
type count struct {
	currentCycleIndex int
	peak              float64
	ratios            []float64
}

// manual hard-coded Perlin noise, curve
var hourFactor = []float64{
	//00-06 hours, night is low
	0.1, 0.2, 0.2, 0.2, 0.2, 0.3,
	//07-12 small peak of users morning
	0.4, 0.5, 0.65, 0.60, 0.55, 0.4,
	//13-16 lunch time, low again
	0.5, 0.45, 0.4, 0.35, 0.45, 0.5,
	//17-24 peak evening
	0.65, 0.85, 1, 0.9, 0.6, 0.3,
}

// NewTrafficNormalDay creates a struct that simulates a normal day traffic
// starting from a peak given at the evenings, to a low of 10% of peak in the night
func NewTrafficNormalDay(peak int) TrafficPattern {
	if peak <= 1 {
		panic("peak must be > 1")
	}
	c := &count{}
	c.peak = float64(peak)
	c.ratios = hourFactor
	return c
}

// NewTrafficPattern returns a new custom pattern based on cycleRatios and peak.
func NewTrafficPattern(peak float64, cycleRatios []float64) TrafficPattern {
	if peak <= 1 {
		panic("peak must be > 1")
	}
	if len(cycleRatios) < 1 {
		panic("must have at least 1 ratio")
	}

	c := &count{}
	c.peak = peak
	//TODO copy the cycleratios, do not trust it will not modify
	c.ratios = cycleRatios
	return c

}

func (c *count) GetValueAndAdvance() int {
	if len(c.ratios) == 0 {
		panic("no ratios to choose from")
	}
	old := int(c.ratios[c.currentCycleIndex] * c.peak)

	c.currentCycleIndex++
	if c.currentCycleIndex >= len(c.ratios) {
		c.currentCycleIndex = 0
	}
	return old
}
