package command

import (
	"context"

	"eats/backend/billing/api/module/client"
	"eats/backend/settlements/app/models"
)

type ModulesContract interface {
	IssueReceipt(ctx context.Context, req client.IssueReceiptRequest) error
}

type Handlers struct {
	legalEntityRepository models.LegalEntityRepository

	modules ModulesContract
}

func NewHandlers(
	legalEntityRepository models.LegalEntityRepository,
	modules ModulesContract,
) *Handlers {
	if legalEntityRepository == nil {
		panic("legalEntityRepository is required")
	}
	if modules == nil {
		panic("modules is required")
	}

	return &Handlers{
		legalEntityRepository: legalEntityRepository,
		modules:               modules,
	}
}
