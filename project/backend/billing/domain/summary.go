package domain

import (
	"github.com/shopspring/decimal"
)

type PriceBreakdownSummary struct {
	netAmount   decimal.Decimal
	taxAmount   decimal.Decimal
	grossAmount decimal.Decimal
	taxes       []TaxSummary
}

func summarizeLineItems(lineItems []LineItem) PriceBreakdownSummary {
	var netAmount, taxAmount, grossAmount decimal.Decimal
	for _, lineItem := range lineItems {
		pb := lineItem.PriceBreakdown()
		netAmount = netAmount.Add(pb.NetAmount())
		taxAmount = taxAmount.Add(pb.TaxAmount())
		grossAmount = grossAmount.Add(pb.GrossAmount())
	}

	summaries := make(map[taxRateKey]*TaxSummary)
	for _, lineItem := range lineItems {
		pb := lineItem.PriceBreakdown()
		key := pb.rate.key()

		summary, ok := summaries[key]
		if !ok {
			summaries[key] = newTaxSummary(pb)
			continue
		}

		summary.add(pb)
	}

	taxes := make([]TaxSummary, 0, len(summaries))
	for _, summary := range summaries {
		taxes = append(taxes, *summary)
	}

	return PriceBreakdownSummary{
		netAmount:   netAmount,
		taxAmount:   taxAmount,
		grossAmount: grossAmount,
		taxes:       taxes,
	}
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

type taxRateKey struct {
	rate    string
	taxType string
}

func (t TaxRate) key() taxRateKey {
	return taxRateKey{
		rate:    t.rate.String(),
		taxType: t.taxType.String(),
	}
}

func newTaxSummary(price PriceBreakdown) *TaxSummary {
	return &TaxSummary{
		taxRate:   price.TaxRate(),
		netAmount: price.NetAmount(),
		taxAmount: price.TaxAmount(),
	}
}

func (t *TaxSummary) add(price PriceBreakdown) {
	t.netAmount = t.netAmount.Add(price.NetAmount())
	t.taxAmount = t.taxAmount.Add(price.TaxAmount())
}
