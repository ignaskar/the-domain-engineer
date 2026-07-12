package command

import (
	"context"
	"errors"

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
	if cmd.PlatformEntityUUID.IsZero() {
		return errors.New("empty platform entity uuid")
	}

	le, err := models.NewLegalEntity(
		cmd.PartnerUUID,
		models.LegalEntityPartner,
		cmd.BusinessName,
		cmd.TaxID,
		cmd.Address,
		cmd.BankAccountNumber,
		cmd.Currency,
	)
	if err != nil {
		return err
	}

	partner := models.NewPartner(le, cmd.PlatformEntityUUID)

	return h.legalEntityRepository.SavePartner(ctx, partner)
}
