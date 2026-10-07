package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/antonpiat/clinic-platform/backend/internal/scheduling/domain"
)

func TestHold(t *testing.T) {
	policy := domain.DefaultPolicy()
	start := testNow.Add(24 * time.Hour)

	t.Run("creates a held appointment", func(t *testing.T) {
		p := holdParams(t, start)

		a, err := domain.Hold(p, testNow, policy)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if a.Status() != domain.StatusHeld {
			t.Errorf("status = %q, want held", a.Status())
		}
		if a.ID() != p.ID {
			t.Errorf("id = %s, want %s", a.ID(), p.ID)
		}
		if a.Slot() != p.Slot {
			t.Errorf("slot = %v, want %v", a.Slot(), p.Slot)
		}
		if want := testNow.Add(policy.HoldDuration); !a.HoldExpiresAt().Equal(want) {
			t.Errorf("holdExpiresAt = %s, want %s", a.HoldExpiresAt(), want)
		}
		if v := a.Snapshot().Version; v != 0 {
			t.Errorf("version = %d, want 0", v)
		}
	})

	t.Run("records AppointmentHeld", func(t *testing.T) {
		p := holdParams(t, start)
		a := must(domain.Hold(p, testNow, policy))

		events := a.PullEvents()
		if len(events) != 1 {
			t.Fatalf("got %d events, want 1", len(events))
		}
		e, ok := events[0].(domain.AppointmentHeld)
		if !ok {
			t.Fatalf("event is %T, want AppointmentHeld", events[0])
		}
		if e.EventType() != "scheduling.appointment.held.v1" {
			t.Errorf("EventType = %q", e.EventType())
		}
		if e.AppointmentID() != p.ID {
			t.Errorf("AppointmentID = %s", e.AppointmentID())
		}
		if !e.OccurredAt().Equal(testNow) {
			t.Errorf("OccurredAt = %s, want %s", e.OccurredAt(), testNow)
		}
		if e.PractitionerID != p.PractitionerID || e.PatientID != p.PatientID ||
			e.ServiceID != p.ServiceID || e.Slot != p.Slot || e.Price != p.Price {
			t.Errorf("event data does not match params: %+v", e)
		}
		if !e.HoldExpiresAt.Equal(a.HoldExpiresAt()) {
			t.Errorf("event HoldExpiresAt = %s, want %s", e.HoldExpiresAt, a.HoldExpiresAt())
		}
	})

	t.Run("normalizes now to UTC microseconds", func(t *testing.T) {
		newYork := must(time.LoadLocation("America/New_York"))
		localNow := testNow.In(newYork).Add(700 * time.Nanosecond)

		a := must(domain.Hold(holdParams(t, start), localNow, policy))

		want := testNow.Add(policy.HoldDuration)
		if got := a.HoldExpiresAt(); !got.Equal(want) || got.Location() != time.UTC {
			t.Errorf("holdExpiresAt = %s, want %s UTC", got, want)
		}
	})

	t.Run("slot one microsecond in the future is allowed", func(t *testing.T) {
		soon := testNow.Add(time.Microsecond)
		if _, err := domain.Hold(holdParams(t, soon), testNow, policy); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	failures := []struct {
		name    string
		mutate  func(p *domain.HoldParams)
		policy  domain.Policy
		wantErr error
	}{
		{
			name:    "slot starts exactly now",
			mutate:  func(p *domain.HoldParams) { *p = holdParams(t, testNow) },
			wantErr: domain.ErrSlotInPast,
		},
		{
			name:    "slot in the past",
			mutate:  func(p *domain.HoldParams) { *p = holdParams(t, testNow.Add(-time.Hour)) },
			wantErr: domain.ErrSlotInPast,
		},
		{
			name:    "duration shorter than service",
			mutate:  func(p *domain.HoldParams) { p.ServiceDuration = 90 * time.Minute },
			wantErr: domain.ErrDurationMismatch,
		},
		{
			name:    "duration longer than service",
			mutate:  func(p *domain.HoldParams) { p.ServiceDuration = 30 * time.Minute },
			wantErr: domain.ErrDurationMismatch,
		},
		{
			name:    "missing appointment id",
			mutate:  func(p *domain.HoldParams) { p.ID = domain.AppointmentID{} },
			wantErr: domain.ErrInvalidID,
		},
		{
			name:    "missing practitioner id",
			mutate:  func(p *domain.HoldParams) { p.PractitionerID = domain.PractitionerID{} },
			wantErr: domain.ErrInvalidID,
		},
		{
			name:    "missing patient id",
			mutate:  func(p *domain.HoldParams) { p.PatientID = domain.PatientID{} },
			wantErr: domain.ErrInvalidID,
		},
		{
			name:    "missing service id",
			mutate:  func(p *domain.HoldParams) { p.ServiceID = domain.ServiceID{} },
			wantErr: domain.ErrInvalidID,
		},
		{
			name:    "missing slot",
			mutate:  func(p *domain.HoldParams) { p.Slot = domain.TimeSlot{} },
			wantErr: domain.ErrInvalidTimeSlot,
		},
		{
			name:    "missing price",
			mutate:  func(p *domain.HoldParams) { p.Price = domain.Money{} },
			wantErr: domain.ErrInvalidMoney,
		},
		{
			name:    "zero hold duration",
			policy:  domain.Policy{HoldDuration: 0, PatientCancellationCutoff: time.Hour},
			wantErr: domain.ErrInvalidPolicy,
		},
	}
	for _, tt := range failures {
		t.Run(tt.name, func(t *testing.T) {
			p := holdParams(t, start)
			if tt.mutate != nil {
				tt.mutate(&p)
			}
			pol := policy
			if tt.policy != (domain.Policy{}) {
				pol = tt.policy
			}

			a, err := domain.Hold(p, testNow, pol)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if a != nil {
				t.Error("appointment must be nil on error")
			}
		})
	}
}

func TestNewHeldHelper(t *testing.T) {
	a := newHeld(t, testNow, testNow.Add(24*time.Hour))
	if a.Status() != domain.StatusHeld {
		t.Fatalf("status = %q", a.Status())
	}
	if n := len(a.PullEvents()); n != 0 {
		t.Errorf("newHeld should drain events, got %d", n)
	}
}

func TestConfirm(t *testing.T) {
	t.Skip("TODO(step-1): held -> confirmed; at exactly holdExpiresAt still ok; after expiry -> ErrHoldExpired; from every other status -> ErrNotHeld")
}

func TestCancel(t *testing.T) {
	t.Skip("TODO(step-1): patient cancels hold; patient outside cutoff; patient inside cutoff -> ErrCancellationWindowClosed; clinic inside cutoff ok; after start -> ErrAlreadyStarted; twice -> ErrAlreadyCancelled; final states -> ErrAlreadyFinished; invalid actor")
}

func TestExpire(t *testing.T) {
	t.Skip("TODO(step-1): held past expiry -> expired; before expiry -> ErrHoldNotExpired; confirmed -> ErrNotHeld")
}

func TestComplete(t *testing.T) {
	t.Skip("TODO(step-1): confirmed after end -> completed; before end -> ErrNotFinished; held -> ErrNotConfirmed")
}

func TestMarkNoShow(t *testing.T) {
	t.Skip("TODO(step-1): confirmed after start -> no_show; before start -> ErrNotStarted; held -> ErrNotConfirmed")
}

func TestPullEvents(t *testing.T) {
	t.Skip("TODO(step-1): returns events in order; second call returns empty")
}

func TestReconstitute(t *testing.T) {
	t.Skip("TODO(step-1): Snapshot -> Reconstitute -> Snapshot round-trips; no events recorded")
}
