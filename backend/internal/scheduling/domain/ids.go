package domain

// IDs are UUID strings. They are generated outside the domain (app layer)
// and validated here, so the domain stays stdlib-only and deterministic.

type AppointmentID struct{ value string }

type PractitionerID struct{ value string }

type PatientID struct{ value string }

type ServiceID struct{ value string }

func NewAppointmentID(value string) (AppointmentID, error) {
	// TODO(step-1): validate canonical UUID format (36 chars, hex, hyphens at 8-13-18-23), else ErrInvalidID
	panic("not implemented")
}

func NewPractitionerID(value string) (PractitionerID, error) {
	// TODO(step-1): same validation as NewAppointmentID
	panic("not implemented")
}

func NewPatientID(value string) (PatientID, error) {
	// TODO(step-1): same validation as NewAppointmentID
	panic("not implemented")
}

func NewServiceID(value string) (ServiceID, error) {
	// TODO(step-1): same validation as NewAppointmentID
	panic("not implemented")
}

func (id AppointmentID) String() string  { return id.value }
func (id PractitionerID) String() string { return id.value }
func (id PatientID) String() string      { return id.value }
func (id ServiceID) String() string      { return id.value }
