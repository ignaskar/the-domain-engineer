package domain

import "github.com/shopspring/decimal"

type PriceBreakdownSummary struct {
	netAmount   decimal.Decimal
	taxAmount   decimal.Decimal
	grossAmount decimal.Decimal
	taxes       []TaxSummary
}

func summarizeLineItems(lineItems []LineItem) PriceBreakdownSummary {
	// TODO implement me
	return PriceBreakdownSummary{}
}

func (p PriceBreakdownSummary) NetAmount() decimal.Decimal {
	return p.netAmount
}

func (p PriceBreakdownSummary) TaxAmount() decimal.Decimal {
	return p.taxAmount
}

func (p PriceBreakdownSummary) GrossAmount() decimal.Decimal {
	return p.grossAmount
}

func (p PriceBreakdownSummary) Taxes() []TaxSummary {
	return p.taxes
}

type TaxSummary struct {
	taxRate   TaxRate
	netAmount decimal.Decimal
	taxAmount decimal.Decimal
}

func (t TaxSummary) TaxRate() TaxRate {
	return t.taxRate
}

func (t TaxSummary) NetAmount() decimal.Decimal {
	return t.netAmount
}

func (t TaxSummary) TaxAmount() decimal.Decimal {
	return t.taxAmount
}
