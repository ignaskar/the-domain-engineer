package command

import (
	"context"

	"eats/backend/billing/api/module/client"
	"eats/backend/settlements/app"
	"eats/backend/settlements/app/models"
	"eats/backend/settlements/domain"
)

type invoiceDataGenerator interface {
	CalculateDeliveryInvoicesData(ctx context.Context, billingCycleUUID domain.BillingCycleUUID) ([]app.NewInvoiceData, error)
	CalculateCommissionInvoiceData(ctx context.Context, billingCycleUUID domain.BillingCycleUUID, platformUUID models.PlatformEntityUUID) (app.NewInvoiceData, error)
}

type CloseBillingCycle struct {
	PartnerUUID domain.LegalEntityUUID
}

func (h *Handlers) CloseBillingCycle(ctx context.Context, cmd CloseBillingCycle) error {
	// TODO: implement the orchestration.
	// Close the cycle, then settle all unsettled closed cycles.
	// For each cycle with orders, dispatch invoices by partner type before settling.
	panic("not implemented")
}

func (h *Handlers) issueRestaurantInvoices(ctx context.Context, billingCycleUUID domain.BillingCycleUUID, partnerUUID domain.LegalEntityUUID) error {
	// TODO: resolve the partner, calculate commission invoice data, and issue it.
	panic("not implemented")
}

func (h *Handlers) issueCourierInvoices(ctx context.Context, billingCycleUUID domain.BillingCycleUUID) error {
	// TODO: calculate delivery invoices data. Use a cached legal entity finder since the
	// same courier seller appears for every restaurant on the cycle.
	panic("not implemented")
}

type legalEntityFinder interface {
	LegalEntityByUUID(ctx context.Context, uuid domain.LegalEntityUUID) (models.LegalEntity, error)
}

func (h *Handlers) issueInvoice(ctx context.Context, legalEntityFinder legalEntityFinder, inv app.NewInvoiceData) (client.DocumentReadModel, error) {
	// TODO: resolve buyer and seller, map line items, and call IssueInvoice on the billing module.
	// Reject mismatched currencies (cross-border invoices are out of scope).
	panic("not implemented")
}

func newModuleLegalEntity(p models.LegalEntity) client.LegalEntity {
	return client.LegalEntity{
		Name:    p.BusinessName,
		Address: p.Address,
		TaxID:   &p.TaxID,
	}
}
