package domain

import (
	"errors"
	"unicode"
)

// TODO: design the IBAN value object.
//
// You'll need to add:
//
// - IBAN struct that wraps the raw string. The raw string should not be
//   exported. Callers must go through the constructor.
// - NewIBAN(string) (IBAN, error) constructor that validates the input.
//   Validation rules:
//     - empty string returns an error
//     - length below 15 returns an error (Norway has the shortest IBANs)
//     - length above 34 returns an error (max IBAN length per ISO 13616)
//     - first 2 characters must be letters (the country code)
//     - the rest must be alphanumeric
// - IsZero() bool method that reports whether the IBAN is the zero value.
// - String() string method that returns the raw IBAN.
// - UnmarshalIBAN(string) IBAN constructor used by the persistence layer
//   when loading an already-validated IBAN from the database (no validation,
//   no error). Look at how other value objects in the project expose this.

type IBAN struct {
	iban string
}

func NewIBAN(iban string) (IBAN, error) {
	if iban == "" {
		return IBAN{}, errors.New("IBAN empty")
	}

	if len(iban) < 15 {
		return IBAN{}, errors.New("IBAN too short")
	}

	if len(iban) > 34 {
		return IBAN{}, errors.New("IBAN too long")
	}

	for i, s := range iban {
		if i == 0 || i == 1 {
			if !unicode.IsLetter(s) {
				return IBAN{}, errors.New("first two IBAN characters must be letters")
			}
			continue
		}

		if !unicode.IsDigit(s) && !unicode.IsLetter(s) {
			return IBAN{}, errors.New("remaining IBAN characters must be alphanumeric")
		}
	}

	return IBAN{iban: iban}, nil
}

func (i IBAN) IsZero() bool {
	return i.iban == ""
}

func (i IBAN) String() string {
	return i.iban
}

func UnmarshalIBAN(iban string) IBAN {
	return IBAN{iban: iban}
}
