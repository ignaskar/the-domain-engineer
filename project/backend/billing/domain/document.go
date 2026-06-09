package domain

import (
	"errors"
	"time"

	"eats/backend/common"
	"eats/backend/common/shared"
)

type Document struct {
	documentUUID      DocumentUUID
	documentType      DocumentType
	documentNumber    DocumentNumber
	externalReference *string
	issueDate         time.Time
	currency          shared.Currency
	seller            LegalEntity
	buyer             LegalEntity
	lineItems         []LineItem
}

type DocumentType struct {
	common.Enum[DocumentTypeValues]
}

type DocumentTypeValues string

func (DocumentTypeValues) Values() []string {
	return []string{"receipt"}
}

var DocumentTypeReceipt = common.MustEnum[DocumentType]("receipt")

func (d Document) UUID() DocumentUUID {
	return d.documentUUID
}

func (d Document) DocumentType() DocumentType {
	return d.documentType
}

func (d Document) DocumentNumber() DocumentNumber {
	return d.documentNumber
}

func (d Document) ExternalReference() *string {
	return d.externalReference
}

func (d Document) IssueDate() time.Time {
	return d.issueDate
}

func (d Document) Currency() shared.Currency {
	return d.currency
}

func (d Document) Seller() LegalEntity {
	return d.seller
}

func (d Document) Buyer() LegalEntity {
	return d.buyer
}

func (d Document) LineItems() []LineItem {
	return d.lineItems
}

func NewReceipt(data NewDocumentData, docNumber DocumentNumber) (*Document, error) {
	if data.Buyer.IsZero() {
		return nil, errors.New("buyer can't be empty")
	}

	if data.Seller.IsZero() {
		return nil, errors.New("seller can't be empty")
	}

	if data.Seller.TaxID() == nil {
		return nil, errors.New("seller must have a tax ID")
	}

	if data.Buyer.TaxID() != nil {
		return nil, errors.New("receipts cannot be issued to buyers with a tax ID")
	}

	if data.Currency.IsZero() {
		return nil, errors.New("currency can't be empty")
	}

	if data.IssueDate.IsZero() {
		return nil, errors.New("issue date can't be empty")
	}

	if data.IssueDate.After(time.Now()) {
		return nil, errors.New("issue date can't be in the future")
	}

	if len(data.LineItems) == 0 {
		return nil, errors.New("there must be at least one line item")
	}

	lineItems := make([]LineItem, 0, len(data.LineItems))
	for _, lineItem := range data.LineItems {
		if lineItem.Name == "" {
			return nil, errors.New("name can't be empty")
		}

		if lineItem.Quantity <= 0 {
			return nil, errors.New("quantity must be positive")
		}

		lineItems = append(lineItems, LineItem{
			lineItemUUID: LineItemUUID{common.NewUUIDv7()},
			name:         lineItem.Name,
			quantity:     lineItem.Quantity,
		})
	}

	document := &Document{
		documentUUID:      DocumentUUID{common.NewUUIDv7()},
		documentType:      DocumentTypeReceipt,
		documentNumber:    docNumber,
		externalReference: data.ExternalReference,
		issueDate:         data.IssueDate,
		currency:          data.Currency,
		seller:            data.Seller,
		buyer:             data.Buyer,
		lineItems:         lineItems,
	}

	return document, nil
}

type NewDocumentData struct {
	ExternalReference *string
	IssueDate         time.Time
	Currency          shared.Currency
	Seller            LegalEntity
	Buyer             LegalEntity
	LineItems         []NewLineItemData
}

type LineItem struct {
	lineItemUUID LineItemUUID
	name         string
	quantity     int
}

func (l LineItem) UUID() LineItemUUID {
	return l.lineItemUUID
}

func (l LineItem) Name() string {
	return l.name
}

func (l LineItem) Quantity() int {
	return l.quantity
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
