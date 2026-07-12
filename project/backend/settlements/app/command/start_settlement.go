package command

import (
	"context"
	"fmt"

	settlementsModule "eats/backend/settlements/api/module/client"
)

func (h *Handlers) StartSettlement(ctx context.Context, cmd settlementsModule.StartSettlementRequest) error {
	// TODO: implement
	// 1. Validate the request via cmd.Validate(). Fail fast on invalid input.
	// 2. Fetch the restaurant legal entity by UUID via h.legalEntityRepository.LegalEntityByUUID.
	// 3. Build a billing IssueReceiptRequest:
	//    - ExternalReference: derive from cmd.OrderUUID (e.g., "receipt-<orderUUID>").
	//    - IssueDate: time.Now().
	//    - Currency: cmd.Currency.
	//    - Seller: built from the restaurant's BusinessName, Address, TaxID.
	//    - Buyer: cmd.CustomerName, cmd.CustomerAddress.
	//    - LineItems: convert each cmd.LineItems entry into client.LineItem.
	// 4. Call h.modules.IssueReceipt.
	// 5. Log on success.
	return fmt.Errorf("not implemented")
}
