package clock

import (
	"time"

	. "github.com/101loops/bdd"
)

var (
	fixedTime = time.Unix(1415926535, 0)
)

var _ = Describe("Mock Clock", func() {

	It("returns the time", func() {
		clock := NewMock()
		Check(clock.IsFrozen(), IsFalse)
		Check(timeDiff(clock), IsLessThan, 1*time.Millisecond)
	})

	It("sets time", func() {
		clock := NewMock().Set(fixedTime)
		Check(clock.Now().Sub(fixedTime), IsLessThan, 1*time.Millisecond)

		time.Sleep(delay)
		Check(clock.Now().Sub(fixedTime), IsRoughly, delay, threshold)
	})

	It("adds time", func() {
		clock := NewMock().Add(1 * time.Hour)
		Check(timeDiff(clock), IsRoughly, -1*time.Hour, threshold)
	})

	It("freezes", func() {
		clock := NewMock().Add(1 * time.Hour).Freeze()
		Check(clock.IsFrozen(), IsTrue)
		clockNow := clock.Now()

		time.Sleep(delay)

		Check(clock.Now(), IsSameTimeAs, clockNow)
	})

	It("freezes at passed-in time", func() {
		clock := NewMock().Freeze().Set(fixedTime)
		Check(clock.IsFrozen(), IsTrue)

		time.Sleep(delay)

		Check(clock.Now(), IsSameTimeAs, fixedTime)
	})

	It("unfreezes", func() {
		clock := NewMock().Freeze()
		Check(clock.IsFrozen(), IsTrue)
		old := clock.Now()

		time.Sleep(delay)

		clock.Unfreeze()
		Check(clock.IsFrozen(), IsFalse)
		Check(old.Sub(clock.Now()), IsRoughly, 0, threshold)
	})

	It("can sleep", func() {
		clock := NewMock()

		slept := durationOf(func() { clock.Sleep(delay) })
		Check(slept, IsRoughly, delay, threshold)
	})
})
