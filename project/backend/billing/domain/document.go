package domain

import (
	"errors"
	"fmt"
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
	if data.Buyer.TaxID() != nil {
		return nil, errors.New("receipts cannot be issued to buyers with a tax ID")
	}
	if data.Buyer.IsZero() {
		return nil, errors.New("buyer can't be empty")
	}
	if data.Seller.IsZero() {
		return nil, errors.New("seller can't be empty")
	}
	if data.Currency.IsZero() {
		return nil, errors.New("currency can't be empty")
	}
	if data.IssueDate.IsZero() {
		return nil, errors.New("issue date can't be empty")
	}
	tomorrow := time.Now().Truncate(24*time.Hour).AddDate(0, 0, 1)
	if !data.IssueDate.Before(tomorrow) {
		return nil, errors.New("issue date can't be in the future")
	}
	if data.Seller.TaxID() == nil {
		return nil, errors.New("seller must have a tax ID to issue billing documents")
	}
	if len(data.LineItems) == 0 {
		return nil, errors.New("document must have at least one line item")
	}

	lineItems := make([]LineItem, 0, len(data.LineItems))
	for _, lid := range data.LineItems {
		lineItem, err := newLineItem(lid)
		if err != nil {
			return nil, fmt.Errorf("failed to create line item for document: %w", err)
		}
		lineItems = append(lineItems, lineItem)
	}

	return &Document{
		uuid:              DocumentUUID{common.NewUUIDv7()},
		externalReference: data.ExternalReference,
		documentType:      DocumentTypeReceipt,
		issueDate:         data.IssueDate,
		documentNumber:    docNumber,
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

type Document struct {
	uuid              DocumentUUID
	externalReference *string
	documentType      DocumentType
	issueDate         time.Time
	currency          shared.Currency
	documentNumber    DocumentNumber
	seller            LegalEntity
	buyer             LegalEntity
	lineItems         []LineItem
}

func (d *Document) UUID() DocumentUUID {
	return d.uuid
}

func (d *Document) ExternalReference() *string {
	return d.externalReference
}

func (d *Document) DocumentType() DocumentType {
	return d.documentType
}

func (d *Document) DocumentNumber() DocumentNumber {
	return d.documentNumber
}

func (d *Document) IssueDate() time.Time {
	return d.issueDate
}

func (d *Document) Currency() shared.Currency {
	return d.currency
}

func (d *Document) Seller() LegalEntity {
	return d.seller
}

func (d *Document) Buyer() LegalEntity {
	return d.buyer
}

func (d *Document) LineItems() []LineItem {
	return d.lineItems
}

type LineItemUUID struct {
	common.UUID
}

type LineItem struct {
	uuid     LineItemUUID
	name     string
	quantity int
}

func newLineItem(data NewLineItemData) (LineItem, error) {
	if data.Name == "" {
		return LineItem{}, errors.New("name can't be empty")
	}

	if data.Quantity < 1 {
		return LineItem{}, errors.New("quantity must be positive")
	}

	return LineItem{
		uuid:     LineItemUUID{common.NewUUIDv7()},
		name:     data.Name,
		quantity: data.Quantity,
	}, nil
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
