package clock

import (
	"sync"
	"time"
)

// FakeClock provides controlled deterministic time for tests.
type FakeClock struct {
	mu         sync.RWMutex
	nowTime    time.Time
	gameSecond int64
}

func NewFakeClock(start time.Time, initialSecond int64) *FakeClock {
	return &FakeClock{
		nowTime:    start,
		gameSecond: initialSecond,
	}
}

func (f *FakeClock) Now() time.Time {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.nowTime
}

func (f *FakeClock) GameSecond() int64 {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.gameSecond
}

func (f *FakeClock) AdvanceSeconds(seconds int64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.gameSecond += seconds
	f.nowTime = f.nowTime.Add(time.Duration(seconds) * time.Second)
}
