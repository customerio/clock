package clock

import (
	"container/list"
	"sync"
	"time"
)

type waiter struct {
	mutex    sync.Mutex
	m        *mock
	wake     chan bool
	duration time.Duration
	frozen   bool
}

func (w *waiter) sleep() {
	frozen := w.frozen

	wakeAt := w.m.Now().Add(w.duration)

	for {
		d := wakeAt.Sub(w.m.Now())
		if d < 0 {
			break
		}
		if frozen {
			select {
			case frozen = <-w.wake:
			}
		} else {
			select {
			case frozen = <-w.wake:
			case <-time.After(d):
				return
			}
		}
	}
}

func (w *waiter) wakeup(freeze bool) {
	w.wake <- freeze
}

type mock struct {
	mutex  sync.Mutex
	base   time.Time
	last   time.Time
	frozen bool

	waiters *list.List
}

// NewMock returns a new manipulable Clock.
func NewMock() Mock {
	n := time.Now()
	return &mock{
		base:    n,
		last:    n,
		waiters: list.New(),
	}
}

func (c *mock) Now() time.Time {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.move()
	return c.base
}

func (c *mock) Set(t time.Time) Mock {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.base = t
	c.last = time.Now()
	c.wakeup()
	return c
}

func (c *mock) Add(d time.Duration) Mock {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.base = c.base.Add(d)
	c.wakeup()
	return c
}

func (c *mock) Freeze() Mock {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	// Freezing a frozen clock does nothing.
	if c.frozen {
		return c
	}

	c.move()
	c.frozen = true

	c.wakeup()
	return c
}

func (c *mock) IsFrozen() bool {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	return c.frozen
}

func (c *mock) Unfreeze() Mock {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.frozen = false
	c.last = time.Now()

	c.wakeup()
	return c
}

func (c *mock) clear(e *list.Element) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.waiters.Remove(e)
}

func (c *mock) Sleep(d time.Duration) {
	c.mutex.Lock()

	w := &waiter{
		m:        c,
		wake:     make(chan bool),
		frozen:   c.frozen,
		duration: d,
	}
	element := c.waiters.PushBack(w)
	c.mutex.Unlock()

	defer c.clear(element)

	w.sleep()
}

func (c *mock) wakeup() {
	for e := c.waiters.Front(); e != nil; e = e.Next() {
		w := e.Value.(*waiter)
		w.wakeup(c.frozen)
	}
}

func (c *mock) move() {
	if c.frozen {
		return
	}

	// Adjust the time by the amount of elapsed time since the last call.
	n := time.Now()
	diff := n.Sub(c.last)
	c.last = n
	c.base = c.base.Add(diff)
}

func (c *mock) Tick(d time.Duration) <-chan time.Time {
	c.mutex.Lock()

	w := &waiter{
		m:        c,
		wake:     make(chan bool),
		frozen:   c.frozen,
		duration: d,
	}
	element := c.waiters.PushBack(w)
	c.mutex.Unlock()

	ch := make(chan time.Time)

	go func() {
		// Not exactly correct since it doesn't account for slow receivers.
		for {
			w.sleep()
			ch <- c.Now()
		}
		c.clear(element)
	}()
	return ch
}

func (*mock) Ticker(d time.Duration) *time.Ticker {
	// TODO: make mockable
	return time.NewTicker(d)
}

func (c *mock) After(d time.Duration) <-chan time.Time {
	c.mutex.Lock()

	w := &waiter{
		m:        c,
		wake:     make(chan bool),
		frozen:   c.frozen,
		duration: d,
	}
	element := c.waiters.PushBack(w)
	c.mutex.Unlock()

	ch := make(chan time.Time)

	go func() {
		w.sleep()
		c.clear(element)
		ch <- c.Now()
	}()
	return ch
}
