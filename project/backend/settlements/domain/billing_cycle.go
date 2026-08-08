package domain

import (
	"time"

	"eats/backend/common"
)

type BillingCycleUUID struct {
	common.UUID
}

type BillingCycle struct {
	uuid        BillingCycleUUID
	partnerUUID LegalEntityUUID
	partnerType PartnerType
	number      int
	closed      bool
	settled     bool
	startDate   time.Time
	endDate     *time.Time
}

func NewInitialBillingCycle(partnerUUID LegalEntityUUID, partnerType PartnerType) (*BillingCycle, error) {
	// TODO: validate inputs (zero checks) and return a new BillingCycle with number=1.
	panic("not implemented")
}

func NewNextBillingCycle(previous *BillingCycle) (*BillingCycle, error) {
	// TODO: validate previous is non-nil and closed, then return a new BillingCycle
	// with number incremented and startDate set to now (UTC).
	panic("not implemented")
}

func (bc *BillingCycle) UUID() BillingCycleUUID {
	return bc.uuid
}

func (bc *BillingCycle) PartnerUUID() LegalEntityUUID {
	return bc.partnerUUID
}

func (bc *BillingCycle) PartnerType() PartnerType {
	return bc.partnerType
}

func (bc *BillingCycle) Number() int {
	return bc.number
}

func (bc *BillingCycle) Closed() bool {
	return bc.closed
}

func (bc *BillingCycle) StartDate() time.Time {
	return bc.startDate
}

func (bc *BillingCycle) EndDate() *time.Time {
	return bc.endDate
}

func (bc *BillingCycle) Settled() bool {
	return bc.settled
}

func (bc *BillingCycle) Close() error {
	// TODO: refuse if already closed; otherwise mark closed and set endDate to now (UTC).
	panic("not implemented")
}

func (bc *BillingCycle) Settle() error {
	// TODO: refuse if not closed or already settled; otherwise mark settled.
	panic("not implemented")
}

func UnmarshalBillingCycle(
	uuid BillingCycleUUID,
	partnerUUID LegalEntityUUID,
	partnerType PartnerType,
	number int,
	closed bool,
	settled bool,
	startDate time.Time,
	endDate *time.Time,
) *BillingCycle {
	return &BillingCycle{
		uuid:        uuid,
		partnerUUID: partnerUUID,
		partnerType: partnerType,
		number:      number,
		closed:      closed,
		settled:     settled,
		startDate:   startDate,
		endDate:     endDate,
	}
}
