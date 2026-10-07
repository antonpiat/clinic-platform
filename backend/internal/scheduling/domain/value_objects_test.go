package domain_test

import "testing"

func TestIDs(t *testing.T) {
	t.Skip("TODO(step-1): valid UUID; empty; wrong length; non-hex; misplaced hyphens")
}

func TestTimeSlot(t *testing.T) {
	t.Skip("TODO(step-1): valid; start == end; start > end; zero times; converts to UTC; Duration; Overlaps incl. back-to-back = false")
}

func TestMoney(t *testing.T) {
	t.Skip("TODO(step-1): valid; zero ok; negative; lowercase/short/long currency")
}

func TestStatus(t *testing.T) {
	t.Skip("TODO(step-1): ParseStatus valid/invalid; BlocksSlot only held+confirmed; IsFinal")
}
