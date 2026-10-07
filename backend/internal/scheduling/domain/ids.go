package domain

import "fmt"

// IDs are UUID strings. They are generated outside the domain (app layer)
// and validated here, so the domain stays stdlib-only and deterministic.

type AppointmentID struct{ value string }

type PractitionerID struct{ value string }

type PatientID struct{ value string }

type ServiceID struct{ value string }

func NewAppointmentID(value string) (AppointmentID, error) {
	v, err := parseUUID(value)
	if err != nil {
		return AppointmentID{}, err
	}
	return AppointmentID{value: v}, nil
}

func NewPractitionerID(value string) (PractitionerID, error) {
	v, err := parseUUID(value)
	if err != nil {
		return PractitionerID{}, err
	}
	return PractitionerID{value: v}, nil
}

func NewPatientID(value string) (PatientID, error) {
	v, err := parseUUID(value)
	if err != nil {
		return PatientID{}, err
	}
	return PatientID{value: v}, nil
}

func NewServiceID(value string) (ServiceID, error) {
	v, err := parseUUID(value)
	if err != nil {
		return ServiceID{}, err
	}
	return ServiceID{value: v}, nil
}

func (id AppointmentID) String() string  { return id.value }
func (id PractitionerID) String() string { return id.value }
func (id PatientID) String() string      { return id.value }
func (id ServiceID) String() string      { return id.value }

// IsZero reports whether the ID was never set (zero value of the struct).
func (id AppointmentID) IsZero() bool  { return id.value == "" }
func (id PractitionerID) IsZero() bool { return id.value == "" }
func (id PatientID) IsZero() bool      { return id.value == "" }
func (id ServiceID) IsZero() bool      { return id.value == "" }

// parseUUID validates the canonical 8-4-4-4-12 hex form and returns it
// lowercased, so IDs compare equal regardless of input case
// (Postgres always returns lowercase).
func parseUUID(value string) (string, error) {
	if len(value) != 36 {
		return "", fmt.Errorf("%w: %q is not a UUID", ErrInvalidID, value)
	}
	out := []byte(value)
	for i, c := range out {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return "", fmt.Errorf("%w: %q is not a UUID", ErrInvalidID, value)
			}
		default:
			switch {
			case c >= '0' && c <= '9', c >= 'a' && c <= 'f':
			case c >= 'A' && c <= 'F':
				out[i] = c + ('a' - 'A')
			default:
				return "", fmt.Errorf("%w: %q is not a UUID", ErrInvalidID, value)
			}
		}
	}
	return string(out), nil
}
