package domain

import (
	"errors"

	"github.com/shopspring/decimal"

	"eats/backend/common"
	"eats/backend/common/shared"
)

type TaxType struct {
	common.Enum[TaxTypeValues]
}

func (t TaxType) DisplayName() string {
	switch t.String() {
	case "vat":
		return "VAT"
	case "gst":
		return "GST"
	case "sales-tax":
		return "Sales Tax"
	default:
		return "Unknown"
	}
}

type TaxTypeValues string

func (TaxTypeValues) Values() []string {
	return []string{"vat", "gst", "sales-tax"}
}

var (
	TaxTypeVAT      = common.MustEnum[TaxType]("vat")
	TaxTypeGST      = common.MustEnum[TaxType]("gst")
	TaxTypeSalesTax = common.MustEnum[TaxType]("sales-tax")
)

type PriceBreakdown struct {
	rate TaxRate

	unitNetAmount   decimal.Decimal
	unitTaxAmount   decimal.Decimal
	unitGrossAmount decimal.Decimal

	netAmount   decimal.Decimal
	taxAmount   decimal.Decimal
	grossAmount decimal.Decimal
}

func NewPriceBreakdownFromNetAmount(
	rate TaxRate,
	unitNetAmount decimal.Decimal,
	currency shared.Currency,
	quantity int,
) (PriceBreakdown, error) {
	// TODO implement me
	return PriceBreakdown{}, errors.New("not implemented")
}

func NewPriceBreakdownFromGrossAmount(
	rate TaxRate,
	unitGrossAmount decimal.Decimal,
	currency shared.Currency,
	quantity int,
) (PriceBreakdown, error) {
	// TODO implement me
	return PriceBreakdown{}, errors.New("not implemented")
}

func (t PriceBreakdown) TaxRate() TaxRate {
	return t.rate
}

func (t PriceBreakdown) UnitNetAmount() decimal.Decimal {
	return t.unitNetAmount
}

func (t PriceBreakdown) UnitTaxAmount() decimal.Decimal {
	return t.unitTaxAmount
}

func (t PriceBreakdown) UnitGrossAmount() decimal.Decimal {
	return t.unitGrossAmount
}

func (t PriceBreakdown) NetAmount() decimal.Decimal {
	return t.netAmount
}

func (t PriceBreakdown) TaxAmount() decimal.Decimal {
	return t.taxAmount
}

func (t PriceBreakdown) GrossAmount() decimal.Decimal {
	return t.grossAmount
}

type TaxRate struct {
	rate    decimal.Decimal
	taxType TaxType
}

func NewTaxRate(rate decimal.Decimal, taxType TaxType) (TaxRate, error) {
	if rate.IsNegative() {
		return TaxRate{}, common.NewInvalidInputError("tax-rate-negative", "tax rate cannot be negative")
	}
	if taxType.IsZero() {
		return TaxRate{}, common.NewInvalidInputError("tax-type-zero", "tax type cannot be zero value")
	}

	return TaxRate{
		rate:    rate,
		taxType: taxType,
	}, nil
}

func (t TaxRate) IsZero() bool {
	return t.taxType.IsZero()
}

func (t TaxRate) Equal(other TaxRate) bool {
	return t.rate.Equal(other.rate) && t.taxType.String() == other.taxType.String()
}

func (t TaxRate) Rate() decimal.Decimal {
	return t.rate
}

func (t TaxRate) TaxType() TaxType {
	return t.taxType
}

// roundInCurrency rounds the given amount according to the decimal places of the currency.
// Different currencies have different rules for decimal places (e.g., JPY has 0, USD has 2).
func roundInCurrency(amount decimal.Decimal, currency shared.Currency) decimal.Decimal {
	return amount.Round(int32(currency.DecimalPlaces()))
}
