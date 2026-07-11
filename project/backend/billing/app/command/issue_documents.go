package command

import (
	"context"

	"eats/backend/billing/domain"
)

type IssueReceipt struct {
	DocumentData domain.NewDocumentData
}

func (h *Handlers) IssueReceipt(ctx context.Context, cmd IssueReceipt) (domain.DocumentUUID, error) {
	builder, err := h.documentFactory.NewReceiptBuilder(ctx, cmd.DocumentData)
	if err != nil {
		return domain.DocumentUUID{}, err
	}

	return h.documentRepository.CreateDocument(
		ctx,
		domain.DocumentSeriesReceipt,
		func(documentNumber domain.DocumentNumber) (*domain.Document, error) {
			return builder.Build(documentNumber)
		},
	)
}
