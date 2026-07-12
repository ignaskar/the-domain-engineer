package models

import (
	"context"
)

type LegalEntityRepository interface {
	LegalEntityByUUID(ctx context.Context, uuid LegalEntityUUID) (LegalEntity, error)
	PartnerByUUID(ctx context.Context, uuid LegalEntityUUID) (Partner, error)
	SavePartner(ctx context.Context, partner Partner) error
	SavePlatformEntity(ctx context.Context, platform LegalEntity) error
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
