package command

import (
	"context"
	"fmt"
	"time"

	"eats/backend/billing/api/module/client"
	"eats/backend/common/log"
	"eats/backend/common/shared"
	settlementsModule "eats/backend/settlements/api/module/client"
	"eats/backend/settlements/domain"
)

func (h *Handlers) StartSettlement(ctx context.Context, cmd settlementsModule.StartSettlementRequest) error {
	err := cmd.Validate()
	if err != nil {
		return err
	}

	restaurantUUID := domain.LegalEntityUUID{cmd.RestaurantUUID}

	restaurant, err := h.legalEntityRepository.LegalEntityByUUID(ctx, restaurantUUID)
	if err != nil {
		return fmt.Errorf("could not get restaurant entity: %w", err)
	}

	var lineItems []client.LineItem
	for _, l := range cmd.LineItems {
		lineItems = append(lineItems, client.LineItem{
			Name:       l.Name,
			Type:       l.Type,
			Quantity:   l.Quantity,
			UnitAmount: shared.NewGrossAmount(l.GrossAmount),
		})
	}

	externalReference := cmd.OrderUUID.String()

	err = h.modules.IssueReceipt(ctx, client.IssueReceiptRequest{
		ExternalReference: &externalReference,
		IssueDate:         time.Now(),
		Currency:          cmd.Currency,
		Seller: client.LegalEntity{
			Name:    restaurant.BusinessName,
			Address: restaurant.Address,
			TaxID:   &restaurant.TaxID,
		},
		Buyer: client.LegalEntity{
			Name:    cmd.CustomerName,
			Address: cmd.CustomerAddress,
		},
		LineItems: lineItems,
	})
	if err != nil {
		return fmt.Errorf("could not issue receipt: %w", err)
	}

	log.FromContext(ctx).Info(
		"Settlement started",
		"order", cmd.OrderUUID,
	)

	return nil
}
