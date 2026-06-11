package command

import (
	"context"
	"fmt"
	"path"

	"eats/backend/billing/domain"
)

type documentPrinter interface {
	PrintDocument(ctx context.Context, doc *domain.Document) ([]byte, error)
}

type fileStorage interface {
	StoreFile(ctx context.Context, filename string, data []byte) (string, error)
}

type PrintDocument struct {
	DocumentUUID domain.DocumentUUID
}

func (h *Handlers) PrintDocument(ctx context.Context, cmd PrintDocument) error {
	document, err := h.documentRepository.DocumentByUUID(ctx, cmd.DocumentUUID)
	if err != nil {
		return fmt.Errorf("failed to get document by uuid: %w", err)
	}

	bytes, err := h.documentPrinter.PrintDocument(ctx, document)
	if err != nil {
		return fmt.Errorf("failed to print document: %w", err)
	}

	var subdir string
	switch document.DocumentType() {
	case domain.DocumentTypeReceipt:
		subdir = "receipts"
	default:
		return fmt.Errorf("document type %s not recognized", document.DocumentType())
	}

	p := path.Join("documents", subdir, document.DocumentNumber().String()+".html")

	url, err := h.fileStorage.StoreFile(ctx, p, bytes)
	if err != nil {
		return fmt.Errorf("failed to store file: %w", err)
	}

	err = h.documentRepository.UpdateFileUrl(ctx, cmd.DocumentUUID, url)
	if err != nil {
		return fmt.Errorf("failed to update document url: %w", err)
	}

	return nil
}
