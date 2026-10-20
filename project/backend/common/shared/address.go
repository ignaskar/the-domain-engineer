package shared

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type addressDbDTO struct {
	Line1       string      `json:"line_1,omitempty"`
	Line2       string      `json:"line_2,omitempty"`
	PostalCode  string      `json:"postal_code,omitempty"`
	City        string      `json:"city,omitempty"`
	CountryCode CountryCode `json:"country_code"`
}

type Address struct {
	line1       string
	line2       string
	postalCode  string
	city        string
	countryCode CountryCode
}

func NewAddress(line1, line2, postalCode, city string, countryCode CountryCode) (Address, error) {
	if line1 == "" {
		return Address{}, errors.New("address line 1 is required")
	}

	if postalCode == "" {
		return Address{}, errors.New("postal code is required")
	}

	if city == "" {
		return Address{}, errors.New("city is required")
	}

	if countryCode.IsZero() {
		return Address{}, errors.New("country code is required")
	}

	return Address{
		line1:       line1,
		line2:       line2,
		postalCode:  postalCode,
		city:        city,
		countryCode: countryCode,
	}, nil
}

func (a Address) IsZero() bool {
	return a == Address{}
}

func (a *Address) Scan(src any) error {
	var text []byte

	switch t := src.(type) {
	case string:
		text = []byte(t)
	case []byte:
		text = t
	default:
		return fmt.Errorf("invalid address type: %T", t)
	}

	err := json.Unmarshal(text, a)
	if err != nil {
		return fmt.Errorf("error unmarshalling %T from json: %w", a, err)
	}

	return nil
}

func (a Address) Value() (driver.Value, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("error marshalling %T to json: %w", a, err)
	}

	return string(data), nil
}

func (a Address) MarshalJSON() ([]byte, error) {
	dto := addressDbDTO{
		Line1:       a.line1,
		Line2:       a.line2,
		PostalCode:  a.postalCode,
		City:        a.city,
		CountryCode: a.countryCode,
	}

	return json.Marshal(dto)
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var dto addressDbDTO
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return fmt.Errorf("error unmarshalling %T from json: %w", a, err)
	}

	a.line1 = dto.Line1
	a.line2 = dto.Line2
	a.postalCode = dto.PostalCode
	a.city = dto.City
	a.countryCode = dto.CountryCode

	return nil
}

func (a Address) Line1() string {
	return a.line1
}

func (a Address) Line2() string {
	return a.line2
}

func (a Address) PostalCode() string {
	return a.postalCode
}

func (a Address) City() string {
	return a.city
}

func (a Address) CountryCode() CountryCode {
	return a.countryCode
}
