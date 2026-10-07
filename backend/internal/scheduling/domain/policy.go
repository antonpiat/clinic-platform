package domain

import "time"

// Policy holds the clinic's booking rules. Passed into the aggregate so
// rules are configurable per clinic later without touching domain code.
type Policy struct {
	// How long a slot stays held before it must be confirmed.
	HoldDuration time.Duration
	// Patients cannot cancel a confirmed appointment closer than this to its start.
	PatientCancellationCutoff time.Duration
}

func DefaultPolicy() Policy {
	return Policy{
		HoldDuration:              5 * time.Minute,
		PatientCancellationCutoff: 24 * time.Hour,
	}
}

// Actor identifies who performs an action with different permissions.
type Actor string

const (
	ActorPatient Actor = "patient"
	ActorClinic  Actor = "clinic"
)

func (a Actor) valid() bool {
	return a == ActorPatient || a == ActorClinic
}
