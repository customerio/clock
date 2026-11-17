package clock

import "testing"

func TestRealClock(t *testing.T) {
	clock := New()
	checkLess(t, clock.Now().Sub(now()), delay)
	checkRoughly(t, durationOf(func() { clock.Sleep(delay) }), delay, threshold)
	checkRoughly(t, durationOf(func() { <-clock.After(delay) }), delay, threshold)
	ticker := clock.Ticker(delay).C
	checkRoughly(t, durationOf(func() { <-ticker; <-ticker }), 2*delay, threshold)
	checkRoughly(t, durationOf(func() { <-clock.Tick(delay); <-clock.Tick(delay) }), 2*delay, threshold)
}
