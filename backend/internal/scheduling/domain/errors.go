package domain

import "errors"

var (
	// Value objects
	ErrInvalidID       = errors.New("scheduling: invalid id")
	ErrInvalidTimeSlot = errors.New("scheduling: invalid time slot")
	ErrInvalidMoney    = errors.New("scheduling: invalid money")
	ErrInvalidStatus   = errors.New("scheduling: invalid status")
	ErrInvalidActor    = errors.New("scheduling: invalid actor")
	ErrInvalidPolicy   = errors.New("scheduling: invalid policy")

	// Hold
	ErrSlotInPast       = errors.New("scheduling: slot starts in the past")
	ErrDurationMismatch = errors.New("scheduling: slot duration does not match the service")

	// Lifecycle
	ErrNotHeld                  = errors.New("scheduling: appointment is not held")
	ErrHoldExpired              = errors.New("scheduling: hold has expired")
	ErrHoldNotExpired           = errors.New("scheduling: hold has not expired yet")
	ErrNotConfirmed             = errors.New("scheduling: appointment is not confirmed")
	ErrAlreadyCancelled         = errors.New("scheduling: appointment is already cancelled")
	ErrAlreadyFinished          = errors.New("scheduling: appointment is already finished")
	ErrAlreadyStarted           = errors.New("scheduling: appointment has already started")
	ErrNotStarted               = errors.New("scheduling: appointment has not started yet")
	ErrNotFinished              = errors.New("scheduling: appointment has not ended yet")
	ErrCancellationWindowClosed = errors.New("scheduling: too late for the patient to cancel")

	// Persistence (returned by Repository implementations, step 2)
	ErrNotFound               = errors.New("scheduling: appointment not found")
	ErrSlotTaken              = errors.New("scheduling: slot is already taken")
	ErrConcurrentModification = errors.New("scheduling: appointment was modified concurrently")
)
