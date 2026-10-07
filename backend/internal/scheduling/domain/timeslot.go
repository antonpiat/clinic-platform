package domain

import "time"

// TimeSlot is a half-open interval [start, end), always stored in UTC.
type TimeSlot struct {
	start time.Time
	end   time.Time
}

func NewTimeSlot(start, end time.Time) (TimeSlot, error) {
	// TODO(step-1): reject zero times and start >= end (ErrInvalidTimeSlot)
	// TODO(step-1): normalize both to UTC
	panic("not implemented")
}

func (s TimeSlot) Start() time.Time { return s.start }
func (s TimeSlot) End() time.Time   { return s.end }

func (s TimeSlot) Duration() time.Duration {
	// TODO(step-1)
	panic("not implemented")
}

// Overlaps reports whether two half-open slots share any instant.
// Back-to-back slots (a.end == b.start) do not overlap.
func (s TimeSlot) Overlaps(other TimeSlot) bool {
	// TODO(step-1)
	panic("not implemented")
}
