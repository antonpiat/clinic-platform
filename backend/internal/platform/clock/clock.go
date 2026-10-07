// Package clock provides the time source for the application layer.
// The domain never reads time itself: it receives `now` as a parameter.
package clock

import (
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
}

// System returns the real current time in UTC.
type System struct{}

func (System) Now() time.Time {
	return time.Now().UTC()
}

// Fixed is a controllable clock for tests. Safe for concurrent use
// (step 2's concurrency test calls it from many goroutines).
type Fixed struct {
	mu  sync.Mutex
	now time.Time
}

func NewFixed(t time.Time) *Fixed {
	return &Fixed{now: t.UTC()}
}

func (f *Fixed) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.now
}

// Advance moves the clock forward (or backward, with a negative d).
func (f *Fixed) Advance(d time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.now = f.now.Add(d)
}

// Set jumps the clock to t.
func (f *Fixed) Set(t time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.now = t.UTC()
}

// Compile-time checks that both types satisfy Clock.
var (
	_ Clock = System{}
	_ Clock = (*Fixed)(nil)
)
