package module

import (
	"context"
	"eats/backend/billing/api/module/client"
	"eats/backend/billing/app/command"
	"eats/backend/billing/domain"
	"eats/backend/common/shared"
)

type Billing struct {
	commandHandlers *command.Handlers
}

func New(
	commandHandlers *command.Handlers,
) *Billing {
	if commandHandlers == nil {
		panic("commandHandlers cannot be nil")
	}

	return &Billing{
		commandHandlers: commandHandlers,
	}
}

func (b *Billing) IssueReceipt(ctx context.Context, req client.IssueReceiptRequest) error {
	buyer, err := newDomainLegalEntityFromContract(req.Buyer)
	if err != nil {
		return err
	}

	seller, err := newDomainLegalEntityFromContract(req.Seller)
	if err != nil {
		return err
	}

	var lineItemData []domain.NewLineItemData
	for _, li := range req.LineItems {
		lid := domain.NewLineItemData{
			Name:       li.Name,
			Quantity:   li.Quantity,
			UnitAmount: li.UnitAmount,
		}
		lineItemData = append(lineItemData, lid)
	}

	dd := domain.NewDocumentData{
		ExternalReference: req.ExternalReference,
		IssueDate:         req.IssueDate,
		Currency:          req.Currency,
		Seller:            *seller,
		Buyer:             *buyer,
		LineItems:         lineItemData,
	}
	cmd := command.IssueReceipt{
		DocumentData: dd,
	}

	_, err = b.commandHandlers.IssueReceipt(ctx, cmd)
	return err
}

func newDomainLegalEntityFromContract(le client.LegalEntity) (*domain.LegalEntity, error) {
	address, err := shared.NewAddress(
		le.Address.Line1(),
		le.Address.Line2(),
		le.Address.PostalCode(),
		le.Address.City(),
		le.Address.CountryCode(),
	)
	if err != nil {
		return nil, err
	}

	domainLe, err := domain.NewLegalEntity(le.Name, address, le.TaxID)
	if err != nil {
		return nil, err
	}

	return &domainLe, nil
}
