package domain

import "fmt"

// Money is an amount in minor units (cents) with an ISO 4217 currency code.
type Money struct {
	cents    int64
	currency string
}

func NewMoney(cents int64, currency string) (Money, error) {
	if cents < 0 {
		return Money{}, fmt.Errorf("%w: negative amount %d", ErrInvalidMoney, cents)
	}
	if !isCurrencyCode(currency) {
		return Money{}, fmt.Errorf("%w: currency %q must be 3 uppercase letters", ErrInvalidMoney, currency)
	}
	return Money{cents: cents, currency: currency}, nil
}

func (m Money) Cents() int64     { return m.cents }
func (m Money) Currency() string { return m.currency }

func isCurrencyCode(s string) bool {
	if len(s) != 3 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < 'A' || s[i] > 'Z' {
			return false
		}
	}
	return true
}
