package models

// TODO: design the Partner entity.
//
// You'll need to add:
//
// - PartnerType — an enum with values "restaurant" and "courier" (same shape
//   as LegalEntityType).
// - PartnerTypeRestaurant and PartnerTypeCourier — instances of PartnerType
//   for the two values.
// - Partner struct with two fields: a LegalEntity and a PlatformEntityUUID.
// - NewPartner constructor accepting a LegalEntity and a PlatformEntityUUID and
//   returning a Partner. No validation needed: LegalEntity already enforces
//   its own invariants.
