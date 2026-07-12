package command

import (
	"context"

	"eats/backend/common/shared"
	"eats/backend/settlements/app/models"
	"eats/backend/settlements/domain"
)

type OnboardPartner struct {
	PartnerUUID        domain.LegalEntityUUID
	PlatformEntityUUID models.PlatformEntityUUID
	PartnerType        domain.PartnerType
	BusinessName       string
	TaxID              shared.TaxID
	Address            shared.Address
	BankAccountNumber  domain.IBAN
	Currency           shared.Currency
}

func (h *Handlers) OnboardPartner(ctx context.Context, cmd OnboardPartner) error {
	// TODO: implement
	// 1. Validate PlatformEntityUUID is not zero.
	// 2. Create a new LegalEntity with type Partner using models.NewLegalEntity.
	// 3. Create a Partner using models.NewPartner.
	// 4. Save everything atomically using legalEntityRepository.SavePartner.
	return nil
}
