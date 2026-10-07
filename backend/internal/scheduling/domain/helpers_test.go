package domain_test

import (
	"testing"
	"time"

	"github.com/antonpiat/clinic-platform/backend/internal/scheduling/domain"
)

// Fixed IDs so failures are easy to read.
const (
	appointmentUUID  = "a0000000-0000-4000-8000-000000000001"
	practitionerUUID = "b0000000-0000-4000-8000-000000000001"
	patientUUID      = "c0000000-0000-4000-8000-000000000001"
	serviceUUID      = "d0000000-0000-4000-8000-000000000001"

	serviceDuration = 60 * time.Minute
)

// testNow is "the present" in every test: Monday 2026-11-02 09:00 UTC.
var testNow = time.Date(2026, 11, 2, 9, 0, 0, 0, time.UTC)

// must unwraps a (value, error) pair in test setup. It panics on error,
// which fails the test with a stack trace pointing at the bad setup line.
// (Go only allows must(f()) when f's results are must's only arguments,
// so it cannot also take *testing.T.)
func must[T any](v T, err error) T {
	if err != nil {
		panic("test setup: " + err.Error())
	}
	return v
}

// holdParams returns valid params for a 60-minute slot starting at start.
// Tests change single fields to provoke one failure at a time.
func holdParams(t *testing.T, start time.Time) domain.HoldParams {
	t.Helper()
	return domain.HoldParams{
		ID:              must(domain.NewAppointmentID(appointmentUUID)),
		PractitionerID:  must(domain.NewPractitionerID(practitionerUUID)),
		PatientID:       must(domain.NewPatientID(patientUUID)),
		ServiceID:       must(domain.NewServiceID(serviceUUID)),
		Slot:            must(domain.NewTimeSlot(start, start.Add(serviceDuration))),
		ServiceDuration: serviceDuration,
		Price:           must(domain.NewMoney(4500, "EUR")),
	}
}

// newHeld returns a valid held appointment created at now, starting at start,
// with its creation event already drained.
func newHeld(t *testing.T, now, start time.Time) *domain.Appointment {
	t.Helper()
	a := must(domain.Hold(holdParams(t, start), now, domain.DefaultPolicy()))
	a.PullEvents()
	return a
}
