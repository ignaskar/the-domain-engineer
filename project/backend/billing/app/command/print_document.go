package command

import (
	"context"

	"eats/backend/billing/domain"
)

type documentPrinter interface {
	PrintDocument(ctx context.Context, doc *domain.Document) ([]byte, error)
}

type fileStorage interface {
	StoreFile(ctx context.Context, filename string, data []byte) (string, error)
}
