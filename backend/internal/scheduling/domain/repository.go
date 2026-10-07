package domain

import "context"

// Repository is the port for persisting appointments.
// TODO(step-2): implement in adapters/postgres.
type Repository interface {
	// Get returns ErrNotFound if the appointment does not exist.
	Get(ctx context.Context, id AppointmentID) (*Appointment, error)

	// Save inserts or updates. Returns ErrSlotTaken on overlap (exclusion
	// constraint) and ErrConcurrentModification on a stale version.
	Save(ctx context.Context, a *Appointment) error
}
