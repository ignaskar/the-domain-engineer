package main

import "errors"

// Money is a monetary amount with a currency.
//
// Turn this into a proper value object:
//  1. Reject invalid inputs in NewMoney (negative amount, empty currency).
//  2. Enforce the same-currency invariant in Add.
//  3. Expose getters for `amount` and `currency` fields and make them unexported.
type Money struct {
	amount   int64
	currency string
}

func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}

	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}

	return Money{amount: amount, currency: currency}, nil
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("cannot add money with different currencies")
	}

	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}
