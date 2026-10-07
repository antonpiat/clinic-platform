package clock_test

import (
	"sync"
	"testing"
	"time"

	"github.com/antonpiat/clinic-platform/backend/internal/platform/clock"
)

func TestSystemReturnsUTC(t *testing.T) {
	before := time.Now()
	got := clock.System{}.Now()
	after := time.Now()

	if got.Location() != time.UTC {
		t.Errorf("location = %s, want UTC", got.Location())
	}
	if got.Before(before) || got.After(after) {
		t.Errorf("System.Now() = %s, not between %s and %s", got, before, after)
	}
}

func TestFixed(t *testing.T) {
	newYork, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	start := time.Date(2026, 11, 2, 5, 0, 0, 0, newYork) // 10:00 UTC
	want := time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC)

	c := clock.NewFixed(start)

	t.Run("returns the given time in UTC", func(t *testing.T) {
		if got := c.Now(); !got.Equal(want) || got.Location() != time.UTC {
			t.Errorf("Now() = %s, want %s", got, want)
		}
	})

	t.Run("does not move on its own", func(t *testing.T) {
		first := c.Now()
		time.Sleep(time.Millisecond)
		if got := c.Now(); !got.Equal(first) {
			t.Errorf("clock moved from %s to %s", first, got)
		}
	})

	t.Run("Advance", func(t *testing.T) {
		c.Advance(5 * time.Minute)
		if got := c.Now(); !got.Equal(want.Add(5 * time.Minute)) {
			t.Errorf("after Advance: %s", got)
		}
		c.Advance(-5 * time.Minute)
		if got := c.Now(); !got.Equal(want) {
			t.Errorf("after negative Advance: %s", got)
		}
	})

	t.Run("Set", func(t *testing.T) {
		target := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
		c.Set(target)
		if got := c.Now(); !got.Equal(target) {
			t.Errorf("after Set: %s", got)
		}
	})
}

// Run with -race: fails if Fixed is not safe for concurrent use.
func TestFixedConcurrentUse(t *testing.T) {
	c := clock.NewFixed(time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC))

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); _ = c.Now() }()
		go func() { defer wg.Done(); c.Advance(time.Second) }()
	}
	wg.Wait()

	want := time.Date(2026, 11, 2, 10, 0, 50, 0, time.UTC)
	if got := c.Now(); !got.Equal(want) {
		t.Errorf("after 50 concurrent Advance(1s): %s, want %s", got, want)
	}
}
