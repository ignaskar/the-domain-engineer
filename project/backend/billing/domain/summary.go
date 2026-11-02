package domain

import "github.com/shopspring/decimal"

type PriceBreakdownSummary struct {
	netAmount   decimal.Decimal
	taxAmount   decimal.Decimal
	grossAmount decimal.Decimal
	taxes       []TaxSummary
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

func summarizeLineItems(lineItems []LineItem) PriceBreakdownSummary {
	var netAmount decimal.Decimal
	var taxAmount decimal.Decimal
	var grossAmount decimal.Decimal

	taxSummaryGroups := make(map[taxRateKey]*TaxSummary)

	for _, li := range lineItems {
		netAmount = netAmount.Add(li.PriceBreakdown().NetAmount())
		taxAmount = taxAmount.Add(li.PriceBreakdown().TaxAmount())
		grossAmount = grossAmount.Add(li.PriceBreakdown().GrossAmount())

		taxRateKey := li.PriceBreakdown().TaxRate().key()

		taxSummary, ok := taxSummaryGroups[taxRateKey]
		if !ok {
			taxSummaryGroups[taxRateKey] = &TaxSummary{
				taxRate:   li.PriceBreakdown().TaxRate(),
				netAmount: li.PriceBreakdown().NetAmount(),
				taxAmount: li.PriceBreakdown().TaxAmount(),
			}
			continue
		}

		taxSummary.add(li.PriceBreakdown())
	}

	var taxSummaries []TaxSummary
	for _, g := range taxSummaryGroups {
		taxSummaries = append(taxSummaries, TaxSummary{
			taxRate:   g.taxRate,
			netAmount: g.netAmount,
			taxAmount: g.taxAmount,
		})
	}

	return PriceBreakdownSummary{
		netAmount:   netAmount,
		taxAmount:   taxAmount,
		grossAmount: grossAmount,
		taxes:       taxSummaries,
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

func (t *TaxSummary) add(breakdown PriceBreakdown) {
	t.netAmount = t.netAmount.Add(breakdown.NetAmount())
	t.taxAmount = t.taxAmount.Add(breakdown.TaxAmount())
}
