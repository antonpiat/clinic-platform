package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/antonpiat/clinic-platform/backend/internal/scheduling/domain"
)

func TestIDs(t *testing.T) {
	const valid = "3f2504e0-4f89-11d3-9a0c-0305e82c3301"

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid lowercase", input: valid, want: valid},
		{name: "uppercase is normalized", input: "3F2504E0-4F89-11D3-9A0C-0305E82C3301", want: valid},
		{name: "empty", input: "", wantErr: true},
		{name: "too short", input: valid[:35], wantErr: true},
		{name: "too long", input: valid + "0", wantErr: true},
		{name: "non-hex character", input: "3f2504e0-4f89-11d3-9a0c-0305e82c330g", wantErr: true},
		{name: "missing hyphens", input: "3f2504e04f8911d39a0c0305e82c3301xxxx", wantErr: true},
		{name: "hyphen in wrong place", input: "3f2504e-04f89-11d3-9a0c-0305e82c3301", wantErr: true},
		{name: "braces", input: "{3f2504e0-4f89-11d3-9a0c-0305e82c33}", wantErr: true},
	}

	// Every ID type shares the same rules.
	constructors := map[string]func(string) (string, error){
		"AppointmentID": func(s string) (string, error) {
			id, err := domain.NewAppointmentID(s)
			return id.String(), err
		},
		"PractitionerID": func(s string) (string, error) {
			id, err := domain.NewPractitionerID(s)
			return id.String(), err
		},
		"PatientID": func(s string) (string, error) {
			id, err := domain.NewPatientID(s)
			return id.String(), err
		},
		"ServiceID": func(s string) (string, error) {
			id, err := domain.NewServiceID(s)
			return id.String(), err
		},
	}

	for typeName, newID := range constructors {
		for _, tt := range tests {
			t.Run(typeName+"/"+tt.name, func(t *testing.T) {
				got, err := newID(tt.input)
				if tt.wantErr {
					if !errors.Is(err, domain.ErrInvalidID) {
						t.Fatalf("want ErrInvalidID, got %v", err)
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got != tt.want {
					t.Errorf("got %q, want %q", got, tt.want)
				}
			})
		}
	}
}

func TestIDIsZero(t *testing.T) {
	var zero domain.AppointmentID
	if !zero.IsZero() {
		t.Error("zero value should report IsZero")
	}
	id, err := domain.NewAppointmentID("3f2504e0-4f89-11d3-9a0c-0305e82c3301")
	if err != nil {
		t.Fatal(err)
	}
	if id.IsZero() {
		t.Error("valid ID should not report IsZero")
	}
}

func TestTimeSlot(t *testing.T) {
	base := time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC)
	newyork, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("valid", func(t *testing.T) {
		s, err := domain.NewTimeSlot(base, base.Add(time.Hour))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !s.Start().Equal(base) || !s.End().Equal(base.Add(time.Hour)) {
			t.Errorf("got [%s, %s)", s.Start(), s.End())
		}
		if s.Duration() != time.Hour {
			t.Errorf("Duration = %s, want 1h", s.Duration())
		}
	})

	t.Run("converts to UTC", func(t *testing.T) {
		start := time.Date(2026, 11, 2, 5, 0, 0, 0, newyork) // 10:00 UTC
		s, err := domain.NewTimeSlot(start, start.Add(time.Hour))
		if err != nil {
			t.Fatal(err)
		}
		if s.Start().Location() != time.UTC {
			t.Errorf("location = %s, want UTC", s.Start().Location())
		}
		if !s.Start().Equal(base) {
			t.Errorf("start = %s, want %s", s.Start(), base)
		}
	})

	t.Run("truncates to microseconds", func(t *testing.T) {
		s, err := domain.NewTimeSlot(base.Add(1500*time.Nanosecond), base.Add(time.Hour))
		if err != nil {
			t.Fatal(err)
		}
		if want := base.Add(time.Microsecond); !s.Start().Equal(want) {
			t.Errorf("start = %s, want %s", s.Start(), want)
		}
	})

	invalid := []struct {
		name       string
		start, end time.Time
	}{
		{"zero start", time.Time{}, base},
		{"zero end", base, time.Time{}},
		{"start equals end", base, base},
		{"start after end", base.Add(time.Hour), base},
		{"equal after truncation", base, base.Add(500 * time.Nanosecond)},
	}
	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewTimeSlot(tt.start, tt.end)
			if !errors.Is(err, domain.ErrInvalidTimeSlot) {
				t.Fatalf("want ErrInvalidTimeSlot, got %v", err)
			}
		})
	}
}

func TestTimeSlotOverlaps(t *testing.T) {
	base := time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC)
	slot := func(fromMin, toMin int) domain.TimeSlot {
		t.Helper()
		s, err := domain.NewTimeSlot(
			base.Add(time.Duration(fromMin)*time.Minute),
			base.Add(time.Duration(toMin)*time.Minute),
		)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}

	a := slot(0, 60) // 10:00-11:00
	tests := []struct {
		name  string
		other domain.TimeSlot
		want  bool
	}{
		{"identical", slot(0, 60), true},
		{"overlaps start", slot(-30, 30), true},
		{"overlaps end", slot(30, 90), true},
		{"contained", slot(15, 45), true},
		{"contains", slot(-30, 90), true},
		{"back-to-back after", slot(60, 120), false},
		{"back-to-back before", slot(-60, 0), false},
		{"separate", slot(120, 180), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := a.Overlaps(tt.other); got != tt.want {
				t.Errorf("a.Overlaps = %v, want %v", got, tt.want)
			}
			if got := tt.other.Overlaps(a); got != tt.want {
				t.Errorf("Overlaps is not symmetric: other.Overlaps(a) = %v", got)
			}
		})
	}
}

func TestMoney(t *testing.T) {
	tests := []struct {
		name     string
		cents    int64
		currency string
		wantErr  bool
	}{
		{name: "valid", cents: 4500, currency: "EUR"},
		{name: "zero is allowed", cents: 0, currency: "USD"},
		{name: "negative", cents: -1, currency: "EUR", wantErr: true},
		{name: "lowercase currency", cents: 100, currency: "eur", wantErr: true},
		{name: "two letters", cents: 100, currency: "EU", wantErr: true},
		{name: "four letters", cents: 100, currency: "EURO", wantErr: true},
		{name: "empty currency", cents: 100, currency: "", wantErr: true},
		{name: "digits", cents: 100, currency: "E1R", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := domain.NewMoney(tt.cents, tt.currency)
			if tt.wantErr {
				if !errors.Is(err, domain.ErrInvalidMoney) {
					t.Fatalf("want ErrInvalidMoney, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if m.Cents() != tt.cents || m.Currency() != tt.currency {
				t.Errorf("got %d %s", m.Cents(), m.Currency())
			}
		})
	}
}

func TestStatus(t *testing.T) {
	tests := []struct {
		status     domain.Status
		blocksSlot bool
		isFinal    bool
	}{
		{domain.StatusHeld, true, false},
		{domain.StatusConfirmed, true, false},
		{domain.StatusCancelled, false, true},
		{domain.StatusExpired, false, true},
		{domain.StatusCompleted, false, true},
		{domain.StatusNoShow, false, true},
	}
	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			parsed, err := domain.ParseStatus(string(tt.status))
			if err != nil {
				t.Fatalf("ParseStatus: %v", err)
			}
			if parsed != tt.status {
				t.Errorf("ParseStatus = %q, want %q", parsed, tt.status)
			}
			if got := tt.status.BlocksSlot(); got != tt.blocksSlot {
				t.Errorf("BlocksSlot = %v, want %v", got, tt.blocksSlot)
			}
			if got := tt.status.IsFinal(); got != tt.isFinal {
				t.Errorf("IsFinal = %v, want %v", got, tt.isFinal)
			}
		})
	}

	for _, bad := range []string{"", "HELD", "pending", "noshow"} {
		t.Run("invalid "+bad, func(t *testing.T) {
			if _, err := domain.ParseStatus(bad); !errors.Is(err, domain.ErrInvalidStatus) {
				t.Fatalf("want ErrInvalidStatus, got %v", err)
			}
		})
	}
}
