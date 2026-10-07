package domain

import (
	"fmt"
	"time"
)

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
	if err := p.validate(); err != nil {
		return nil, err
	}
	if policy.HoldDuration <= 0 {
		return nil, fmt.Errorf("%w: hold duration must be positive", ErrInvalidPolicy)
	}

	now = normalize(now)
	if !p.Slot.Start().After(now) {
		return nil, fmt.Errorf("%w: slot starts at %s, now is %s",
			ErrSlotInPast, p.Slot.Start().Format(time.RFC3339), now.Format(time.RFC3339))
	}
	if p.Slot.Duration() != p.ServiceDuration {
		return nil, fmt.Errorf("%w: slot is %s, service is %s",
			ErrDurationMismatch, p.Slot.Duration(), p.ServiceDuration)
	}

	a := &Appointment{
		id:             p.ID,
		practitionerID: p.PractitionerID,
		patientID:      p.PatientID,
		serviceID:      p.ServiceID,
		slot:           p.Slot,
		price:          p.Price,
		status:         StatusHeld,
		holdExpiresAt:  now.Add(policy.HoldDuration),
	}
	a.record(AppointmentHeld{
		meta:           meta{appointmentID: a.id, occurredAt: now},
		PractitionerID: a.practitionerID,
		PatientID:      a.patientID,
		ServiceID:      a.serviceID,
		Slot:           a.slot,
		Price:          a.price,
		HoldExpiresAt:  a.holdExpiresAt,
	})
	return a, nil
}

// validate rejects zero values: every field must come from a constructor.
func (p HoldParams) validate() error {
	switch {
	case p.ID.IsZero():
		return fmt.Errorf("%w: appointment id is required", ErrInvalidID)
	case p.PractitionerID.IsZero():
		return fmt.Errorf("%w: practitioner id is required", ErrInvalidID)
	case p.PatientID.IsZero():
		return fmt.Errorf("%w: patient id is required", ErrInvalidID)
	case p.ServiceID.IsZero():
		return fmt.Errorf("%w: service id is required", ErrInvalidID)
	case p.Slot == (TimeSlot{}):
		return fmt.Errorf("%w: slot is required", ErrInvalidTimeSlot)
	case p.Price == (Money{}):
		return fmt.Errorf("%w: price is required", ErrInvalidMoney)
	}
	return nil
}

func (a *Appointment) record(e Event) {
	a.events = append(a.events, e)
}

// normalize matches TimeSlot's precision so timestamps survive a DB round-trip.
func normalize(t time.Time) time.Time {
	return t.UTC().Truncate(time.Microsecond)
}

// Confirm turns a hold into a confirmed appointment.
// Confirming at exactly holdExpiresAt is still allowed.
func (a *Appointment) Confirm(now time.Time) error {
	if a.status != StatusHeld {
		return fmt.Errorf("%w: status is %s", ErrNotHeld, a.status)
	}
	now = normalize(now)
	if now.After(a.holdExpiresAt) {
		return fmt.Errorf("%w: expired at %s", ErrHoldExpired, a.holdExpiresAt.Format(time.RFC3339))
	}

	a.status = StatusConfirmed
	a.holdExpiresAt = time.Time{}
	a.record(AppointmentConfirmed{
		meta:           meta{appointmentID: a.id, occurredAt: now},
		PractitionerID: a.practitionerID,
		PatientID:      a.patientID,
		Slot:           a.slot,
	})
	return nil
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
// At exactly holdExpiresAt the hold is still valid (Confirm would succeed),
// so it cannot expire yet.
func (a *Appointment) Expire(now time.Time) error {
	if a.status != StatusHeld {
		return fmt.Errorf("%w: status is %s", ErrNotHeld, a.status)
	}
	now = normalize(now)
	if !now.After(a.holdExpiresAt) {
		return fmt.Errorf("%w: expires at %s", ErrHoldNotExpired, a.holdExpiresAt.Format(time.RFC3339))
	}

	a.status = StatusExpired
	a.holdExpiresAt = time.Time{}
	a.record(AppointmentExpired{
		meta:      meta{appointmentID: a.id, occurredAt: now},
		PatientID: a.patientID,
		Slot:      a.slot,
	})
	return nil
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
