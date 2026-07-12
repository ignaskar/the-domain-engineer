package command

import (
	"context"
	"fmt"

	"eats/backend/common/shared"
	"eats/backend/settlements/app/models"
	"eats/backend/settlements/domain"
)

type CreatePlatformEntity struct {
	PlatformEntityUUID domain.LegalEntityUUID
	BusinessName       string
	TaxID              shared.TaxID
	Address            shared.Address
	BankAccountNumber  domain.IBAN
	Currency           shared.Currency
}

func (h *Handlers) CreatePlatformEntity(ctx context.Context, cmd CreatePlatformEntity) (models.PlatformEntityUUID, error) {
	// TODO: implement
	// 1. Create a new LegalEntity with type Platform using models.NewLegalEntity.
	// 2. Save it using legalEntityRepository.SavePlatformEntity.
	// 3. Return PlatformEntityUUID wrapping the legal entity's UUID.
	return models.PlatformEntityUUID{}, fmt.Errorf("not implemented")
}
