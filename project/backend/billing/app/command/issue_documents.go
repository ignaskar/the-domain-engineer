package command

import (
	"context"
	"eats/backend/billing/domain"
)

type IssueReceipt struct {
	DocumentData domain.NewDocumentData
}

func (h *Handlers) IssueReceipt(ctx context.Context, cmd IssueReceipt) (domain.DocumentUUID, error) {
	docUUID, err := h.documentRepository.CreateDocument(
		ctx,
		domain.DocumentSeriesReceipt,
		func(documentNumber domain.DocumentNumber) (*domain.Document, error) {
			return domain.NewReceipt(cmd.DocumentData, documentNumber)
		},
	)
	if err != nil {
		return domain.DocumentUUID{}, err
	}

	return docUUID, nil
}
