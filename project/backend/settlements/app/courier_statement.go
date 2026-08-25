package app

import (
	"time"

	"github.com/shopspring/decimal"

	"eats/backend/common/shared"
	"eats/backend/settlements/app/models"
	"eats/backend/settlements/domain"
)

type CourierStatement struct {
	Snapshot CourierStatementSnapshot
	Payment  StatementPayment
}

type CourierStatementSnapshot struct {
	UUID             models.StatementUUID    `json:"uuid"`
	BillingCycleUUID domain.BillingCycleUUID `json:"billing_cycle_uuid"`
	DocumentName     string                  `json:"document_name"`
	Currency         shared.Currency         `json:"currency"`
	Period           DateRange               `json:"period"`

	Platform LegalEntity `json:"platform"`
	Courier  LegalEntity `json:"courier"`

	Deliveries []CourierStatementDelivery `json:"deliveries"`
	Summary    CourierStatementSummary    `json:"summary"`
}

type CourierStatementDelivery struct {
	Date           time.Time `json:"date"`
	ShortID        string    `json:"short_id"`
	RestaurantName string    `json:"restaurant_name"`

	DeliveryNetAmount   Money `json:"delivery_net_amount"`
	DeliveryTaxAmount   Money `json:"delivery_tax_amount"`
	DeliveryGrossAmount Money `json:"delivery_gross_amount"`
}

type CourierStatementSummary struct {
	TotalDeliveries int `json:"total_deliveries"`

	DeliveryNetAmount   Money `json:"delivery_net_amount"`
	DeliveryTaxAmount   Money `json:"delivery_tax_amount"`
	DeliveryGrossAmount Money `json:"delivery_gross_amount"`

	PayoutGrossAmount Money `json:"payout_gross_amount"`
}

type Money struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency shared.Currency `json:"currency"`
}

func (m Money) String() string {
	return m.Amount.StringFixed(int32(m.Currency.DecimalPlaces()))
}

type Percentage struct {
	Value decimal.Decimal `json:"value"`
}

func (p Percentage) String() string {
	return p.Value.Mul(decimal.NewFromInt(100)).String() + "%"
}
