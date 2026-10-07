package domain

import "time"

// Appointment is the aggregate root of the scheduling context.
// It guards lifecycle rules. Overlap between appointments is NOT checked
// here: the database exclusion constraint owns that rule (ADR 0003).
type Appointment struct {
	id             AppointmentID
	practitionerID PractitionerID
	patientID      PatientID
	serviceID      ServiceID
	slot           TimeSlot
	price          Money
	status         Status
	holdExpiresAt  time.Time // zero unless status == held
	version        int       // optimistic locking, incremented by the repository
	events         []Event
}

// HoldParams are the inputs for a new hold. ServiceDuration and Price come
// from the catalog (step 3); until then the caller supplies them.
type HoldParams struct {
	ID              AppointmentID
	PractitionerID  PractitionerID
	PatientID       PatientID
	ServiceID       ServiceID
	Slot            TimeSlot
	ServiceDuration time.Duration
	Price           Money
}

// Hold creates a temporary reservation that must be confirmed within policy.HoldDuration.
func Hold(p HoldParams, now time.Time, policy Policy) (*Appointment, error) {
	// TODO(step-1): slot must start after now, else ErrSlotInPast
	// TODO(step-1): slot duration must equal p.ServiceDuration, else ErrDurationMismatch
	// TODO(step-1): status = held, holdExpiresAt = now + policy.HoldDuration
	// TODO(step-1): record AppointmentHeld
	panic("not implemented")
}

// Confirm turns a hold into a confirmed appointment.
func (a *Appointment) Confirm(now time.Time) error {
	// TODO(step-1): only from held, else ErrNotHeld
	// TODO(step-1): now must not be after holdExpiresAt, else ErrHoldExpired
	// TODO(step-1): status = confirmed, clear holdExpiresAt, record AppointmentConfirmed
	panic("not implemented")
}

// Cancel cancels a held or confirmed appointment. Patients cannot cancel a
// confirmed appointment inside policy.PatientCancellationCutoff; the clinic can.
func (a *Appointment) Cancel(now time.Time, by Actor, policy Policy) error {
	// TODO(step-1): validate actor, else ErrInvalidActor
	// TODO(step-1): cancelled -> ErrAlreadyCancelled; expired/completed/no_show -> ErrAlreadyFinished
	// TODO(step-1): now >= slot start -> ErrAlreadyStarted
	// TODO(step-1): patient + confirmed + (start - now) < cutoff -> ErrCancellationWindowClosed
	// TODO(step-1): status = cancelled, clear holdExpiresAt, record AppointmentCancelled
	panic("not implemented")
}

// Expire releases a hold whose window has passed. Called by the worker job (step 7).
func (a *Appointment) Expire(now time.Time) error {
	// TODO(step-1): only from held, else ErrNotHeld
	// TODO(step-1): now must be after holdExpiresAt, else ErrHoldNotExpired
	// TODO(step-1): status = expired, record AppointmentExpired
	panic("not implemented")
}

// Complete marks a confirmed appointment as attended.
func (a *Appointment) Complete(now time.Time) error {
	// TODO(step-1): only from confirmed, else ErrNotConfirmed
	// TODO(step-1): now must be >= slot end, else ErrNotFinished
	// TODO(step-1): status = completed, record AppointmentCompleted
	panic("not implemented")
}

// MarkNoShow records that the patient did not attend.
func (a *Appointment) MarkNoShow(now time.Time) error {
	// TODO(step-1): only from confirmed, else ErrNotConfirmed
	// TODO(step-1): now must be >= slot start, else ErrNotStarted
	// TODO(step-1): status = no_show, record AppointmentNoShow
	panic("not implemented")
}

// PullEvents returns recorded events and clears them.
func (a *Appointment) PullEvents() []Event {
	events := a.events
	a.events = nil
	return events
}

// === Persistence support ===

// Snapshot is the full state, used by the repository to save (step 2).
type Snapshot struct {
	ID             AppointmentID
	PractitionerID PractitionerID
	PatientID      PatientID
	ServiceID      ServiceID
	Slot           TimeSlot
	Price          Money
	Status         Status
	HoldExpiresAt  time.Time
	Version        int
}

func (a *Appointment) Snapshot() Snapshot {
	return Snapshot{
		ID:             a.id,
		PractitionerID: a.practitionerID,
		PatientID:      a.patientID,
		ServiceID:      a.serviceID,
		Slot:           a.slot,
		Price:          a.price,
		Status:         a.status,
		HoldExpiresAt:  a.holdExpiresAt,
		Version:        a.version,
	}
}

// Reconstitute rebuilds an aggregate from storage. No rules, no events.
func Reconstitute(s Snapshot) *Appointment {
	return &Appointment{
		id:             s.ID,
		practitionerID: s.PractitionerID,
		patientID:      s.PatientID,
		serviceID:      s.ServiceID,
		slot:           s.Slot,
		price:          s.Price,
		status:         s.Status,
		holdExpiresAt:  s.HoldExpiresAt,
		version:        s.Version,
	}
}

// === Read accessors ===

func (a *Appointment) ID() AppointmentID        { return a.id }
func (a *Appointment) Status() Status           { return a.status }
func (a *Appointment) Slot() TimeSlot           { return a.slot }
func (a *Appointment) HoldExpiresAt() time.Time { return a.holdExpiresAt }
