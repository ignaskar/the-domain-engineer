package domain

import (
	"errors"
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
	if partnerUUID.IsZero() {
		return nil, errors.New("partner UUID cannot be empty")
	}
	if partnerType.IsZero() {
		return nil, errors.New("partner type cannot be empty")
	}

	return &BillingCycle{
		uuid:        BillingCycleUUID{common.NewUUIDv7()},
		partnerUUID: partnerUUID,
		partnerType: partnerType,
		number:      1,
		closed:      false,
		settled:     false,
		startDate:   time.Now().UTC(),
		endDate:     nil,
	}, nil
}

func NewNextBillingCycle(previous *BillingCycle) (*BillingCycle, error) {
	if previous == nil {
		return nil, errors.New("previous billing cycle cannot be empty")
	}
	if !previous.Closed() {
		return nil, errors.New("cannot start a new billing cycle if previous one is not closed")
	}

	return &BillingCycle{
		uuid:        BillingCycleUUID{common.NewUUIDv7()},
		partnerUUID: previous.partnerUUID,
		partnerType: previous.partnerType,
		number:      previous.number + 1,
		closed:      false,
		settled:     false,
		startDate:   time.Now().UTC(),
		endDate:     nil,
	}, nil
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
	if bc.closed {
		return errors.New("billing cycle already closed")
	}

	bc.closed = true
	bc.endDate = common.ToPtr(time.Now().UTC())

	return nil
}

func (bc *BillingCycle) Settle() error {
	if !bc.closed {
		return errors.New("cannot settle an unclosed billing cycle")
	}
	if bc.settled {
		return errors.New("billing cycle already settled")
	}

	bc.settled = true

	return nil
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
