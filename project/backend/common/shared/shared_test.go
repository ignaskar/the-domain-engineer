// This file contains tests that are executed to verify your solution.
// It's read-only, so all modifications will be ignored.
package shared

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAddress_ValidAddress(t *testing.T) {
	countryCode := MustNewCountryCode("US")

	addr, err := NewAddress("123 Main St", "Apt 4", "12345", "New York", countryCode)

	require.NoError(t, err)
	require.Equal(t, "123 Main St", addr.Line1())
	require.Equal(t, "Apt 4", addr.Line2())
	require.Equal(t, "12345", addr.PostalCode())
	require.Equal(t, "New York", addr.City())
	require.Equal(t, "US", addr.CountryCode().Code())
}

func TestNewAddress_ValidAddressWithoutLine2(t *testing.T) {
	countryCode := MustNewCountryCode("US")

	addr, err := NewAddress("123 Main St", "", "12345", "New York", countryCode)

	require.NoError(t, err)
	require.Equal(t, "123 Main St", addr.Line1())
	require.Equal(t, "", addr.Line2())
}

func TestNewAddress_MissingLine1(t *testing.T) {
	countryCode := MustNewCountryCode("US")

	_, err := NewAddress("", "Apt 4", "12345", "New York", countryCode)

	require.Error(t, err)
}

func TestNewAddress_MissingPostalCode(t *testing.T) {
	countryCode := MustNewCountryCode("US")

	_, err := NewAddress("123 Main St", "Apt 4", "", "New York", countryCode)

	require.Error(t, err)
}

func TestNewAddress_MissingCity(t *testing.T) {
	countryCode := MustNewCountryCode("US")

	_, err := NewAddress("123 Main St", "Apt 4", "12345", "", countryCode)

	require.Error(t, err)
}

func TestNewAddress_MissingCountryCode(t *testing.T) {
	_, err := NewAddress("123 Main St", "Apt 4", "12345", "New York", CountryCode{})

	require.Error(t, err)
}

func TestAddress_MarshalJSON(t *testing.T) {
	countryCode := MustNewCountryCode("US")

	addr, err := NewAddress("123 Main St", "Apt 4", "12345", "New York", countryCode)
	require.NoError(t, err)

	data, err := json.Marshal(addr)
	require.NoError(t, err)

	expected, err := json.Marshal(map[string]any{
		"line_1":       "123 Main St",
		"line_2":       "Apt 4",
		"postal_code":  "12345",
		"city":         "New York",
		"country_code": "US",
	})
	require.NoError(t, err)

	require.JSONEq(t, string(expected), string(data))
}

func TestAddress_Scan(t *testing.T) {
	input, err := json.Marshal(map[string]any{
		"line_1":       "123 Main St",
		"line_2":       "Apt 4",
		"postal_code":  "12345",
		"city":         "New York",
		"country_code": "US",
	})
	require.NoError(t, err)

	var addr Address
	err = addr.Scan(string(input))

	require.NoError(t, err)
	require.Equal(t, "123 Main St", addr.Line1())
	require.Equal(t, "Apt 4", addr.Line2())
	require.Equal(t, "12345", addr.PostalCode())
	require.Equal(t, "New York", addr.City())
	require.Equal(t, "US", addr.CountryCode().Code())
}
