package domain

import "fmt"

// Status values are stored as-is in the database (step 2).
type Status string

const (
	StatusHeld      Status = "held"
	StatusConfirmed Status = "confirmed"
	StatusCancelled Status = "cancelled"
	StatusExpired   Status = "expired"
	StatusCompleted Status = "completed"
	StatusNoShow    Status = "no_show"
)

// ParseStatus is used by the repository when loading from storage.
func ParseStatus(value string) (Status, error) {
	switch s := Status(value); s {
	case StatusHeld, StatusConfirmed, StatusCancelled,
		StatusExpired, StatusCompleted, StatusNoShow:
		return s, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidStatus, value)
	}
}

// BlocksSlot reports whether this status occupies the time slot.
// Must match the exclusion constraint's WHERE clause (ADR 0003).
func (s Status) BlocksSlot() bool {
	return s == StatusHeld || s == StatusConfirmed
}

// IsFinal reports whether no further transition is possible.
func (s Status) IsFinal() bool {
	switch s {
	case StatusCancelled, StatusExpired, StatusCompleted, StatusNoShow:
		return true
	default:
		return false
	}
}
