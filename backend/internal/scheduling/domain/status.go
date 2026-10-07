package domain

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
	// TODO(step-1): accept only the constants above, else ErrInvalidStatus
	panic("not implemented")
}

// BlocksSlot reports whether this status occupies the time slot.
// Must match the exclusion constraint's WHERE clause (ADR 0003).
func (s Status) BlocksSlot() bool {
	// TODO(step-1): true for held and confirmed
	panic("not implemented")
}

// IsFinal reports whether no further transition is possible.
func (s Status) IsFinal() bool {
	// TODO(step-1): true for cancelled, expired, completed, no_show
	panic("not implemented")
}
