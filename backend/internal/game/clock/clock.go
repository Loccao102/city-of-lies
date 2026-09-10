package clock

import (
	"sync"
	"time"
)

// Clock provides a time source for the simulation engine.
type Clock interface {
	Now() time.Time
	GameSecond() int64
	AdvanceSeconds(seconds int64)
}

type RealClock struct {
	mu         sync.RWMutex
	gameSecond int64
}

func NewRealClock(initialGameSecond int64) *RealClock {
	return &RealClock{gameSecond: initialGameSecond}
}

func (c *RealClock) Now() time.Time {
	return time.Now()
}

func (c *RealClock) GameSecond() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.gameSecond
}

func (c *RealClock) AdvanceSeconds(seconds int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.gameSecond += seconds
}
