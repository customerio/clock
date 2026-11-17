package clock

import (
	"testing"
	"time"
)

var fixedTime = time.Unix(1415926535, 0)

func TestMockReturnsTime(t *testing.T) {
	clock := NewMock()
	if clock.IsFrozen() {
		t.Fatal("new mock should not be frozen")
	}
	checkLess(t, timeDiff(clock), time.Millisecond)
}

func TestMockSetsTime(t *testing.T) {
	clock := NewMock().Set(fixedTime)
	checkLess(t, clock.Now().Sub(fixedTime), time.Millisecond)
	time.Sleep(delay)
	checkRoughly(t, clock.Now().Sub(fixedTime), delay, threshold)
}

func TestMockAddNotFrozen(t *testing.T) {
	checkRoughly(t, timeDiff(NewMock().Add(time.Hour)), -time.Hour, threshold)
}

func TestMockAddFrozen(t *testing.T) {
	clock := NewMock().Freeze()
	clock.Add(time.Hour)
	checkRoughly(t, timeDiff(clock), -time.Hour, threshold)
}

func TestMockFreezes(t *testing.T) {
	clock := NewMock().Add(time.Hour).Freeze()
	if !clock.IsFrozen() {
		t.Fatal("expected frozen")
	}
	at := clock.Now()
	time.Sleep(delay)
	if !clock.Now().Equal(at) {
		t.Fatalf("frozen clock moved: %v -> %v", at, clock.Now())
	}
}

func TestMockFreezesAtPassedInTime(t *testing.T) {
	clock := NewMock().Freeze().Set(fixedTime)
	time.Sleep(delay)
	if !clock.Now().Equal(fixedTime) {
		t.Fatalf("got %v, want %v", clock.Now(), fixedTime)
	}
}

func TestMockUnfreezes(t *testing.T) {
	clock := NewMock().Freeze()
	old := clock.Now()
	time.Sleep(delay)
	clock.Unfreeze()
	if clock.IsFrozen() {
		t.Fatal("expected unfrozen")
	}
	checkRoughly(t, old.Sub(clock.Now()), 0, threshold)
}

func TestMockSleep(t *testing.T) {
	clock := NewMock()
	checkRoughly(t, durationOf(func() { clock.Sleep(delay) }), delay, threshold)
}
