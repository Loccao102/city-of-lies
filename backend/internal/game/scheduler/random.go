package scheduler

import (
	"math/rand"
	"sync"
)

// SeededRNG provides a deterministic, thread-safe pseudo-random source.
type SeededRNG struct {
	mu  sync.Mutex
	src *rand.Rand
}

func NewSeededRNG(seed int64) *SeededRNG {
	return &SeededRNG{
		src: rand.New(rand.NewSource(seed)),
	}
}

func (r *SeededRNG) Float64() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.src.Float64()
}

func (r *SeededRNG) Intn(n int) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.src.Intn(n)
}

func (r *SeededRNG) Perm(n int) []int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.src.Perm(n)
}
