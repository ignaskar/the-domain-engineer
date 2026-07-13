package command

import (
	"context"
	billingClient "eats/backend/billing/api/module/client"
	"eats/backend/common"
	"eats/backend/common/log"
	"eats/backend/common/shared"
	settlementsClient "eats/backend/settlements/api/module/client"
	"eats/backend/settlements/app/models"
	"eats/backend/settlements/domain"
	"fmt"
	"time"

	settlementsModule "eats/backend/settlements/api/module/client"
)

func (h *Handlers) StartSettlement(ctx context.Context, cmd settlementsModule.StartSettlementRequest) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	le, err := h.legalEntityRepository.LegalEntityByUUID(ctx, domain.LegalEntityUUID{cmd.RestaurantUUID})
	if err != nil {
		return err
	}

	externalRef := common.ToPtr(fmt.Sprintf("receipt-%s", cmd.OrderUUID))
	req := billingClient.IssueReceiptRequest{
		ExternalReference: externalRef,
		IssueDate:         time.Now(),
		Currency:          cmd.Currency,
		Seller:            toBillingLegalEntity(le),
		Buyer: billingClient.LegalEntity{
			Name:    cmd.CustomerName,
			Address: cmd.CustomerAddress,
			TaxID:   nil,
		},
		LineItems: toBillingLineItems(cmd.LineItems),
	}

	err = h.modules.IssueReceipt(ctx, req)
	if err != nil {
		return err
	}

	log.FromContext(ctx).With("order_id", cmd.OrderUUID).Info("Settlement started successfully")
	return nil
}

func toBillingLineItems(lis []settlementsClient.LineItem) []billingClient.LineItem {
	out := make([]billingClient.LineItem, len(lis))
	for _, li := range lis {
		out = append(out, toBillingLineItem(li))
	}
	return out
}

func toBillingLineItem(li settlementsClient.LineItem) billingClient.LineItem {
	return billingClient.LineItem{
		Name:       li.Name,
		Type:       li.Type,
		UnitAmount: shared.NewGrossAmount(li.GrossAmount),
		Quantity:   li.Quantity,
	}
}

func toBillingLegalEntity(le models.LegalEntity) billingClient.LegalEntity {
	return billingClient.LegalEntity{
		Name:    le.BusinessName,
		Address: le.Address,
		TaxID:   &le.TaxID,
	}
}
