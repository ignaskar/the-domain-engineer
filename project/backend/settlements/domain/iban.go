package domain

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
