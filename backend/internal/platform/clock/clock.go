// Package clock provides the time source for the application layer.
// The domain never reads time itself: it receives `now` as a parameter.
package clock

import "time"

type Clock interface {
	Now() time.Time
}

// System returns the real current time in UTC.
type System struct{}

func (System) Now() time.Time {
	// TODO(step-1): time.Now().UTC()
	panic("not implemented")
}

// Fixed is a controllable clock for tests. Safe for concurrent use
// (step 2's concurrency test calls it from many goroutines).
type Fixed struct {
	// TODO(step-1): sync.Mutex + current time
}

func NewFixed(t time.Time) *Fixed {
	// TODO(step-1)
	panic("not implemented")
}

func (f *Fixed) Now() time.Time {
	// TODO(step-1)
	panic("not implemented")
}

func (f *Fixed) Advance(d time.Duration) {
	// TODO(step-1)
	panic("not implemented")
}
