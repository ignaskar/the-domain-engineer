package shared

import (
	"database/sql/driver"
	"fmt"
)

type TaxID struct {
	taxID string
}

func (t TaxID) IsZero() bool {
	return t.taxID == ""
}

func (t TaxID) String() string {
	return t.taxID
}

// Scan accepts any value from the database without validation.
// Upside: future changes to validation rules won't break existing data in the database.
// Downside: if invalid data somehow ends up in the database, it won't be caught here.
// Alternative: use NewTaxID in Scan to enforce validation on read.
func (t *TaxID) Scan(src any) error {
	if src == nil {
		*t = TaxID{}
		return nil
	}

	text, ok := src.(string)
	if !ok {
		return fmt.Errorf("invalid type for TaxID, expected string, got %T", src)
	}

	*t = TaxID{taxID: text}
	return nil
}

func (t TaxID) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.taxID, nil
}
