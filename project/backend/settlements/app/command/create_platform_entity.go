package command

import (
	"context"
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
	le, err := models.NewLegalEntity(
		cmd.PlatformEntityUUID,
		models.LegalEntityPlatform,
		cmd.BusinessName,
		cmd.TaxID,
		cmd.Address,
		cmd.BankAccountNumber,
		cmd.Currency,
	)
	if err != nil {
		return models.PlatformEntityUUID{}, err
	}

	err = h.legalEntityRepository.SavePlatformEntity(ctx, le)
	if err != nil {
		return models.PlatformEntityUUID{}, err
	}

	return models.PlatformEntityUUID{le.UUID}, err
}
