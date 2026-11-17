package clock

import (
	"testing"
	"time"
)

var (
	delay     = 200 * time.Millisecond
	threshold = 100 * time.Millisecond
)

func now() time.Time { return time.Now() }

func timeDiff(c Clock) time.Duration { return time.Now().Sub(c.Now()) }

func durationOf(fn func()) time.Duration {
	before := now()
	fn()
	return now().Sub(before)
}

func checkRoughly(t *testing.T, got, want, tol time.Duration) {
	t.Helper()
	if d := got - want; d < -tol || d > tol {
		t.Fatalf("got %v, want %v +/- %v", got, want, tol)
	}
}

func checkLess(t *testing.T, got, max time.Duration) {
	t.Helper()
	if got >= max {
		t.Fatalf("got %v, want < %v", got, max)
	}
}
