package clock

import "testing"

func TestWorkNow(t *testing.T) {
	Work = New()
	checkLess(t, Now().Sub(now()), delay)
}

func TestWorkSleep(t *testing.T) {
	checkRoughly(t, durationOf(func() { Sleep(delay) }), delay, threshold)
}

func TestWorkAfter(t *testing.T) {
	checkRoughly(t, durationOf(func() { <-After(delay) }), delay, threshold)
}

func TestWorkTicker(t *testing.T) {
	ticker := Ticker(delay).C
	checkRoughly(t, durationOf(func() { <-ticker; <-ticker }), 2*delay, threshold)
}

func TestWorkTick(t *testing.T) {
	checkRoughly(t, durationOf(func() { <-Tick(delay); <-Tick(delay) }), 2*delay, threshold)
}
