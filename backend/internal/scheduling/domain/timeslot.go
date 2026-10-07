package domain

import (
	"fmt"
	"time"
)

// TimeSlot is a half-open interval [start, end), stored in UTC at
// microsecond precision (what Postgres timestamptz keeps), so a slot
// survives a database round-trip unchanged.
type TimeSlot struct {
	start time.Time
	end   time.Time
}

func NewTimeSlot(start, end time.Time) (TimeSlot, error) {
	if start.IsZero() || end.IsZero() {
		return TimeSlot{}, fmt.Errorf("%w: start and end are required", ErrInvalidTimeSlot)
	}
	start = start.UTC().Truncate(time.Microsecond)
	end = end.UTC().Truncate(time.Microsecond)
	if !start.Before(end) {
		return TimeSlot{}, fmt.Errorf("%w: start %s must be before end %s",
			ErrInvalidTimeSlot, start.Format(time.RFC3339), end.Format(time.RFC3339))
	}
	return TimeSlot{start: start, end: end}, nil
}

func (s TimeSlot) Start() time.Time { return s.start }
func (s TimeSlot) End() time.Time   { return s.end }

func (s TimeSlot) Duration() time.Duration {
	return s.end.Sub(s.start)
}

// Overlaps reports whether two half-open slots share any instant.
// Back-to-back slots (a.end == b.start) do not overlap.
func (s TimeSlot) Overlaps(other TimeSlot) bool {
	return s.start.Before(other.end) && other.start.Before(s.end)
}
