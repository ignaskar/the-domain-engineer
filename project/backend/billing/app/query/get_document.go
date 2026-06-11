package query

import (
	"context"
	"eats/backend/billing/domain"
)

type GetDocumentByUUID struct {
	DocumentUUID domain.DocumentUUID
}

func (h *Handlers) GetDocumentByUUID(ctx context.Context, query GetDocumentByUUID) (*domain.Document, error) {
	return h.documentRepository.DocumentByUUID(ctx, query.DocumentUUID)
}
