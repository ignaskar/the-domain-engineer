package domain

import (
	"eats/backend/common/shared"
	"errors"
)

type LegalEntity struct {
	name    string
	address shared.Address
	taxID   *shared.TaxID
}

func NewLegalEntity(name string, address shared.Address, taxID *shared.TaxID) (LegalEntity, error) {
	if name == "" {
		return LegalEntity{}, errors.New("name can't be empty")
	}

	if address.IsZero() {
		return LegalEntity{}, errors.New("address can't be empty")
	}

	return LegalEntity{
		name:    name,
		address: address,
		taxID:   taxID,
	}, nil
}

func (e LegalEntity) IsZero() bool {
	return e.name == "" && e.address.IsZero()
}

func (e LegalEntity) Name() string {
	return e.name
}

func (e LegalEntity) Address() shared.Address {
	return e.address
}

func (e LegalEntity) TaxID() *shared.TaxID {
	return e.taxID
}
