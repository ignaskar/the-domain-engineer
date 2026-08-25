package app

import (
	"time"

	"eats/backend/common/shared"
	"eats/backend/settlements/app/models"
	"eats/backend/settlements/domain"
)

type LegalEntity struct {
	UUID    domain.LegalEntityUUID `json:"uuid"`
	Name    string                 `json:"name"`
	Address Address                `json:"address"`
	TaxID   string                 `json:"tax_id"`
}

func NewLegalEntityFromModel(legalEntity models.LegalEntity) LegalEntity {
	return LegalEntity{
		UUID: legalEntity.UUID,
		Name: legalEntity.BusinessName,
		Address: Address{
			Line1:      legalEntity.Address.Line1(),
			Line2:      legalEntity.Address.Line2(),
			City:       legalEntity.Address.City(),
			PostalCode: legalEntity.Address.PostalCode(),
			Country:    legalEntity.Address.CountryCode().String(),
		},
		TaxID: legalEntity.TaxID.String(),
	}
}

type Address struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type DateRange struct {
	Name string    `json:"name"`
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

func NewDateRange(from, to time.Time) DateRange {
	from = from.Truncate(24 * time.Hour)
	to = to.Truncate(24 * time.Hour)

	name := from.Format("January 2, 2006") + " - " + to.Format("January 2, 2006")

	return DateRange{
		Name: name,
		From: from,
		To:   to,
	}
}

type RestaurantStatement struct {
	Snapshot RestaurantStatementSnapshot
	Payment  StatementPayment
}

type RestaurantStatementSnapshot struct {
	UUID             models.StatementUUID       `json:"uuid"`
	BillingCycleUUID domain.BillingCycleUUID    `json:"billing_cycle_uuid"`
	DocumentName     string                     `json:"document_name"`
	Currency         shared.Currency            `json:"currency"`
	Period           DateRange                  `json:"period"`
	Platform         LegalEntity                `json:"platform"`
	Restaurant       LegalEntity                `json:"restaurant"`
	Orders           []RestaurantStatementOrder `json:"orders"`
	Summary          RestaurantStatementSummary `json:"summary"`
}

type RestaurantStatementSummary struct {
	OrdersNetAmount   Money `json:"orders_net_amount"`
	OrdersTaxAmount   Money `json:"orders_tax_amount"`
	OrdersGrossAmount Money `json:"orders_gross_amount"`

	CommissionInvoiceNumber string     `json:"commission_invoice_number"`
	CommissionPercent       Percentage `json:"commission_percent"`
	CommissionNetAmount     Money      `json:"commission_net_amount"`
	CommissionTaxAmount     Money      `json:"commission_tax_amount"`
	CommissionGrossAmount   Money      `json:"commission_gross_amount"`

	DeliveryNetAmount   Money `json:"delivery_net_amount"`
	DeliveryTaxAmount   Money `json:"delivery_tax_amount"`
	DeliveryGrossAmount Money `json:"delivery_gross_amount"`

	PayoutGrossAmount Money `json:"payout_gross_amount"`
}

type RestaurantStatementOrder struct {
	Date        time.Time `json:"date"`
	ShortID     string    `json:"short_id"`
	CourierName string    `json:"courier_name"`

	ItemsNetAmount      Money `json:"items_net_amount"`
	ItemsTaxAmount      Money `json:"items_tax_amount"`
	DeliveryNetAmount   Money `json:"delivery_net_amount"`
	DeliveryTaxAmount   Money `json:"delivery_tax_amount"`
	TotalGrossAmount    Money `json:"total_gross_amount"`
	CommissionNetAmount Money `json:"commission_net_amount"`
}

type StatementPayment struct {
	PaidAt               time.Time `json:"paid_at"`
	PaymentAccountNumber string    `json:"payment_account_number"`
	PaymentReference     string    `json:"payment_reference"`
}
