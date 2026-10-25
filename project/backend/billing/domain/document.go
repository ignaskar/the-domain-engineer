package domain

import (
	"errors"
	"time"

	"eats/backend/common"
	"eats/backend/common/shared"
)

type DocumentType struct {
	common.Enum[DocumentTypeValues]
}

type DocumentTypeValues string

func (DocumentTypeValues) Values() []string {
	return []string{"receipt"}
}

var DocumentTypeReceipt = common.MustEnum[DocumentType]("receipt")

func NewReceipt(data NewDocumentData, docNumber DocumentNumber) (*Document, error) {
	// TODO implement me
	return nil, errors.New("not implemented")
}

type NewDocumentData struct {
	ExternalReference *string
	IssueDate         time.Time
	Currency          shared.Currency
	Seller            LegalEntity
	Buyer             LegalEntity
	LineItems         []NewLineItemData
}

type NewLineItemData struct {
	Name     string
	Quantity int
}

type DocumentUUID struct {
	common.UUID
}

type LineItemUUID struct {
	common.UUID
}
