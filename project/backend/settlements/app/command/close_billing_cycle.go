package command

import (
	"context"
	"fmt"
	"time"

	"eats/backend/billing/api/module/client"
	"eats/backend/common/shared"
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
	_, _, err := h.billingCycleRepository.CloseBillingCycle(ctx, cmd.PartnerUUID)
	if err != nil {
		return err
	}

	unsettledCycles, err := h.billingCycleRepository.UnsettledClosedCycles(ctx, cmd.PartnerUUID)
	if err != nil {
		return err
	}

	for _, cycle := range unsettledCycles {
		orders, err := h.billingCycleRepository.BillingCycleOrders(ctx, cycle.UUID())
		if err != nil {
			return err
		}

		if len(orders) > 0 {
			switch cycle.PartnerType() {
			case domain.PartnerTypeRestaurant:
				err = h.issueRestaurantInvoices(ctx, cycle.UUID(), cycle.PartnerUUID())
				if err != nil {
					return err
				}
			case domain.PartnerTypeCourier:
				err = h.issueCourierInvoices(ctx, cycle.UUID())
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown partner type: %v", cycle.PartnerType())
			}
		}

		err = h.billingCycleRepository.SettleBillingCycle(ctx, cycle.UUID())
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handlers) issueRestaurantInvoices(ctx context.Context, billingCycleUUID domain.BillingCycleUUID, partnerUUID domain.LegalEntityUUID) error {
	partner, err := h.legalEntityRepository.PartnerByUUID(ctx, partnerUUID)
	if err != nil {
		return err
	}

	invoiceData, err := h.invoiceDataGenerator.CalculateCommissionInvoiceData(
		ctx,
		billingCycleUUID,
		partner.PlatformEntityUUID,
	)
	if err != nil {
		return err
	}

	_, err = h.issueInvoice(ctx, h.legalEntityRepository, invoiceData)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handlers) issueCourierInvoices(ctx context.Context, billingCycleUUID domain.BillingCycleUUID) error {
	invoiceData, err := h.invoiceDataGenerator.CalculateDeliveryInvoicesData(ctx, billingCycleUUID)
	if err != nil {
		return err
	}

	cachedFinder := app.NewCachedLegalEntityFinder(h.legalEntityRepository)

	for _, invoice := range invoiceData {
		_, err = h.issueInvoice(ctx, cachedFinder, invoice)
		if err != nil {
			return err
		}
	}

	return nil
}

type legalEntityFinder interface {
	LegalEntityByUUID(ctx context.Context, uuid domain.LegalEntityUUID) (models.LegalEntity, error)
}

func (h *Handlers) issueInvoice(ctx context.Context, legalEntityFinder legalEntityFinder, inv app.NewInvoiceData) (client.DocumentReadModel, error) {
	buyer, err := legalEntityFinder.LegalEntityByUUID(ctx, inv.BuyerUUID)
	if err != nil {
		return client.DocumentReadModel{}, err
	}

	seller, err := legalEntityFinder.LegalEntityByUUID(ctx, inv.SellerUUID)
	if err != nil {
		return client.DocumentReadModel{}, err
	}

	if buyer.Currency != seller.Currency {
		return client.DocumentReadModel{}, fmt.Errorf("buyer and seller have different currencies: %s vs %s", buyer.Currency, seller.Currency)
	}

	lineItems := make([]client.LineItem, len(inv.LineItems))
	for i, item := range inv.LineItems {
		lineItems[i] = client.LineItem{
			Name:       item.Name,
			Type:       item.Type,
			UnitAmount: shared.NewNetAmount(item.NetAmount),
			Quantity:   item.Quantity,
		}
	}

	return h.modules.IssueInvoice(ctx, client.IssueInvoiceRequest{
		Buyer:             newModuleLegalEntity(buyer),
		Seller:            newModuleLegalEntity(seller),
		ExternalReference: &inv.ExternalReference,
		IssueDate:         time.Now(),
		Currency:          buyer.Currency,
		LineItems:         lineItems,
	})
}

func newModuleLegalEntity(p models.LegalEntity) client.LegalEntity {
	return client.LegalEntity{
		Name:    p.BusinessName,
		Address: p.Address,
		TaxID:   &p.TaxID,
	}
}
