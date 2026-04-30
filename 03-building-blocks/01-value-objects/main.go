package main

import "errors"

// Money is a monetary amount with a currency.
//
// Turn this into a proper value object:
//  1. Reject invalid inputs in NewMoney (negative amount, empty currency).
//  2. Enforce the same-currency invariant in Add.
//
// Encapsulation (unexported fields) and immutability (Add returns a new Money)
// are already in place. The tests in main_test.go describe the target behavior.
type Money struct {
	amount   int64
	currency string
}

func NewMoney(amount int64, currency string) (Money, error) {
	// TODO: validate amount and currency.
	if amount < 0 {
		return Money{}, errors.New("amount must be greater than zero")
	}

	if currency == "" {
		return Money{}, errors.New("currency must be provided")
	}

	return Money{amount: amount, currency: currency}, nil
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("currency must match")
	}

	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}
