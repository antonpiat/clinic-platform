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

// newConfirmed returns a confirmed appointment (held and confirmed at now),
// starting at start, with its events already drained.
func newConfirmed(t *testing.T, now, start time.Time) *domain.Appointment {
	t.Helper()
	a := newHeld(t, now, start)
	if err := a.Confirm(now); err != nil {
		t.Fatalf("setup: confirm: %v", err)
	}
	a.PullEvents()
	return a
}

// withStatus rebuilds a copy of a in the given status, so tests can reach
// any state without depending on other transitions being implemented.
func withStatus(a *domain.Appointment, status domain.Status) *domain.Appointment {
	s := a.Snapshot()
	s.Status = status
	if status != domain.StatusHeld {
		s.HoldExpiresAt = time.Time{}
	}
	return domain.Reconstitute(s)
}

// allStatusesExcept lists every status but the given ones.
func allStatusesExcept(except ...domain.Status) []domain.Status {
	all := []domain.Status{
		domain.StatusHeld, domain.StatusConfirmed, domain.StatusCancelled,
		domain.StatusExpired, domain.StatusCompleted, domain.StatusNoShow,
	}
	var out []domain.Status
	for _, s := range all {
		skip := false
		for _, e := range except {
			if s == e {
				skip = true
			}
		}
		if !skip {
			out = append(out, s)
		}
	}
	return out
}

// assertUnchanged fails if a failed transition changed state or recorded events.
func assertUnchanged(t *testing.T, a *domain.Appointment, before domain.Snapshot) {
	t.Helper()
	if after := a.Snapshot(); after != before {
		t.Errorf("state changed on failure:\nbefore %+v\nafter  %+v", before, after)
	}
	if n := len(a.PullEvents()); n != 0 {
		t.Errorf("recorded %d events on failure, want 0", n)
	}
}
