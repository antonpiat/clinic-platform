package domain

// Money is an amount in minor units (cents) with an ISO 4217 currency code.
type Money struct {
	cents    int64
	currency string
}

func NewMoney(cents int64, currency string) (Money, error) {
	// TODO(step-1): reject negative amounts and currency not matching ^[A-Z]{3}$ (ErrInvalidMoney)
	panic("not implemented")
}

func (m Money) Cents() int64     { return m.cents }
func (m Money) Currency() string { return m.currency }
