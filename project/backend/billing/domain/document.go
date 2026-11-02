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

type Document struct {
	uuid              DocumentUUID
	externalReference *string
	documentType      DocumentType
	number            DocumentNumber
	issueDate         time.Time
	currency          shared.Currency
	seller            LegalEntity
	buyer             LegalEntity
	lineItems         []LineItem
}

func (d Document) UUID() DocumentUUID {
	return d.uuid
}

func (d Document) ExternalReference() *string {
	return d.externalReference
}

func (d Document) DocumentType() DocumentType {
	return d.documentType
}

func (d Document) DocumentNumber() DocumentNumber {
	return d.number
}

func (d Document) IssueDate() time.Time {
	return d.issueDate
}

func (d Document) Currency() shared.Currency {
	return d.currency
}

func (d Document) Buyer() LegalEntity {
	return d.buyer
}

func (d Document) Seller() LegalEntity {
	return d.seller
}

func (d Document) LineItems() []LineItem {
	return d.lineItems
}

type LineItem struct {
	uuid     LineItemUUID
	name     string
	quantity int
}

func (l LineItem) UUID() LineItemUUID {
	return l.uuid
}

func (l LineItem) Name() string {
	return l.name
}

func (l LineItem) Quantity() int {
	return l.quantity
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
		return nil, errors.New("at least one line item")
	}
	for _, li := range data.LineItems {
		if li.Name == "" {
			return nil, errors.New("name can't be empty")
		}
		if li.Quantity <= 0 {
			return nil, errors.New("quantity must be positive")
		}
	}

	lineItems := make([]LineItem, len(data.LineItems))
	for idx, li := range data.LineItems {
		lineItems[idx] = LineItem{
			uuid:     LineItemUUID{common.NewUUIDv7()},
			name:     li.Name,
			quantity: li.Quantity,
		}
	}

	return &Document{
		uuid:              DocumentUUID{common.NewUUIDv7()},
		externalReference: data.ExternalReference,
		documentType:      DocumentTypeReceipt,
		number:            docNumber,
		issueDate:         data.IssueDate,
		currency:          data.Currency,
		seller:            data.Seller,
		buyer:             data.Buyer,
		lineItems:         lineItems,
	}, nil
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
