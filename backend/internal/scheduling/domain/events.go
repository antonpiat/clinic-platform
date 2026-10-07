package domain

import "time"

// Event is a fact recorded by the aggregate. Mapped to an explicit
// payload by the outbox adapter (step 2); never serialized directly.
type Event interface {
	EventType() string
	AppointmentID() AppointmentID
	OccurredAt() time.Time
}

type meta struct {
	appointmentID AppointmentID
	occurredAt    time.Time
}

func (m meta) AppointmentID() AppointmentID { return m.appointmentID }
func (m meta) OccurredAt() time.Time        { return m.occurredAt }

type AppointmentHeld struct {
	meta
	PractitionerID PractitionerID
	PatientID      PatientID
	ServiceID      ServiceID
	Slot           TimeSlot
	Price          Money
	HoldExpiresAt  time.Time
}

type AppointmentConfirmed struct {
	meta
	PractitionerID PractitionerID
	PatientID      PatientID
	Slot           TimeSlot
}

type AppointmentCancelled struct {
	meta
	PractitionerID PractitionerID
	PatientID      PatientID
	Slot           TimeSlot
	By             Actor
}

type AppointmentExpired struct {
	meta
	PatientID PatientID
	Slot      TimeSlot
}

type AppointmentCompleted struct {
	meta
	PatientID PatientID
}

type AppointmentNoShow struct {
	meta
	PatientID PatientID
	Slot      TimeSlot
}

func (AppointmentHeld) EventType() string      { return "scheduling.appointment.held.v1" }
func (AppointmentConfirmed) EventType() string { return "scheduling.appointment.confirmed.v1" }
func (AppointmentCancelled) EventType() string { return "scheduling.appointment.cancelled.v1" }
func (AppointmentExpired) EventType() string   { return "scheduling.appointment.expired.v1" }
func (AppointmentCompleted) EventType() string { return "scheduling.appointment.completed.v1" }
func (AppointmentNoShow) EventType() string    { return "scheduling.appointment.no_show.v1" }
