package models

import (
	"context"
	"eats/backend/common"
	"eats/backend/common/shared"
	"eats/backend/settlements/domain"
	"errors"
)

type PlatformEntityUUID struct {
	domain.LegalEntityUUID
}

type LegalEntityType struct {
	common.Enum[LegalEntityTypeValues]
}

type LegalEntityTypeValues string

func (LegalEntityTypeValues) Values() []string {
	return []string{"platform", "partner"}
}

var LegalEntityPlatform = common.MustEnum[LegalEntityType]("platform")
var LegalEntityPartner = common.MustEnum[LegalEntityType]("partner")

type LegalEntityRepository interface {
	LegalEntityByUUID(ctx context.Context, uuid domain.LegalEntityUUID) (LegalEntity, error)
	PartnerByUUID(ctx context.Context, uuid domain.LegalEntityUUID) (Partner, error)
	SavePartner(ctx context.Context, partner Partner) error
	SavePlatformEntity(ctx context.Context, platform LegalEntity) error
}

type LegalEntity struct {
	UUID              domain.LegalEntityUUID
	Type              LegalEntityType
	BusinessName      string
	TaxID             shared.TaxID
	Address           shared.Address
	BankAccountNumber domain.IBAN
	Currency          shared.Currency
}

func NewLegalEntity(
	uuid domain.LegalEntityUUID,
	legalEntityType LegalEntityType,
	businessName string,
	taxID shared.TaxID,
	address shared.Address,
	bankAccountNumber domain.IBAN,
	currency shared.Currency,
) (LegalEntity, error) {
	if uuid.IsZero() {
		return LegalEntity{}, errors.New("legal entity UUID cannot be empty")
	}

	if legalEntityType.IsZero() {
		return LegalEntity{}, errors.New("legal entity Type cannot be empty")
	}

	if businessName == "" {
		return LegalEntity{}, errors.New("business name cannot be empty")
	}

	if taxID.IsZero() {
		return LegalEntity{}, errors.New("tax ID cannot be empty")
	}

	if address.IsZero() {
		return LegalEntity{}, errors.New("address cannot be empty")
	}

	if bankAccountNumber.IsZero() {
		return LegalEntity{}, errors.New("bank account number cannot be empty")
	}

	if currency.IsZero() {
		return LegalEntity{}, errors.New("currency cannot be empty")
	}

	return LegalEntity{
		UUID:              uuid,
		Type:              legalEntityType,
		BusinessName:      businessName,
		TaxID:             taxID,
		Address:           address,
		BankAccountNumber: bankAccountNumber,
		Currency:          currency,
	}, nil
}

// TODO: design the LegalEntity entity and the supporting legal-entity types.
//
// You'll need to add:
//
// - LegalEntityUUID — wraps common.UUID for type safety.
// - LegalEntityType — an enum with values "platform" and "partner". See
//   common.Enum and how IBAN-adjacent enums are structured elsewhere in the
//   codebase.
// - PlatformEntityUUID — wraps LegalEntityUUID, used when a value is
//   specifically a platform identifier.
// - LegalEntityPlatform and LegalEntityPartner — instances of LegalEntityType
//   for the two values.
// - LegalEntity struct with fields: UUID, Type, BusinessName, TaxID, Address,
//   BankAccountNumber (IBAN), Currency.
// - NewLegalEntity constructor that validates every field is non-zero (use
//   IsZero() on each VO, and check BusinessName is not the empty string).
//   Return a specific error per missing field.
//
// LegalEntity here is a settlements-side record, not a domain-layer entity
// (the billing LegalEntity is the encapsulated one). It's immutable and has
// no complex logic, so we pragmatically skip encapsulation: exported fields,
// no getters.
