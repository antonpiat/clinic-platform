package domain_test

import "testing"

func TestHold(t *testing.T) {
	t.Skip("TODO(step-1): valid hold sets status held + holdExpiresAt; slot in past; slot starting exactly now; duration mismatch; records AppointmentHeld")
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
