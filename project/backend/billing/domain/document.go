package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"eats/backend/common"
	"eats/backend/common/shared"
)

type DocumentFactory struct {
	taxRateProvider TaxRateProvider
}

func NewDocumentFactory(taxRateProvider TaxRateProvider) *DocumentFactory {
	if taxRateProvider == nil {
		panic("taxRateProvider must not be nil")
	}

	return &DocumentFactory{taxRateProvider}
}

func (f *DocumentFactory) NewReceiptBuilder(ctx context.Context, data NewDocumentData) (*DocumentBuilder, error) {
	if data.Buyer.TaxID() != nil {
		return nil, errors.New("receipts cannot be issued to buyers with a tax ID")
	}
	return f.newDocumentBuilder(ctx, DocumentTypeReceipt, data)
}

func (f *DocumentFactory) newDocumentBuilder(ctx context.Context, docType DocumentType, data NewDocumentData) (*DocumentBuilder, error) {
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
		lineItem, err := f.newLineItem(
			ctx,
			lid,
			data.Buyer.Address().CountryCode(),
			data.Buyer.TaxID(),
			data.Seller.Address().CountryCode(),
			data.Currency,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create line item for document: %w", err)
		}
		lineItems = append(lineItems, lineItem)
	}

	summary := summarizeLineItems(lineItems)

	return &DocumentBuilder{
		externalReference: data.ExternalReference,
		documentType:      DocumentTypeReceipt,
		issueDate:         data.IssueDate,
		currency:          data.Currency,
		seller:            data.Seller,
		buyer:             data.Buyer,
		lineItems:         lineItems,
		summary:           summary,
	}, nil
}

func (f *DocumentFactory) newLineItem(
	ctx context.Context,
	data NewLineItemData,
	buyerCountryCode shared.CountryCode,
	buyerTaxID *shared.TaxID,
	sellerCountryCode shared.CountryCode,
	currency shared.Currency,
) (LineItem, error) {
	if data.Name == "" {
		return LineItem{}, errors.New("name can't be empty")
	}

	if data.Quantity < 1 {
		return LineItem{}, errors.New("quantity must be positive")
	}

	if data.UnitAmount.Amount().IsNegative() {
		return LineItem{}, errors.New("unit amount can't be negative")
	}

	if data.LineItemType.IsZero() {
		return LineItem{}, errors.New("lineItemType can't be empty")
	}

	if buyerCountryCode.IsZero() {
		return LineItem{}, errors.New("buyerCountryCode can't be empty")
	}

	if sellerCountryCode.IsZero() {
		return LineItem{}, errors.New("sellerCountryCode can't be empty")
	}

	var priceBreakdown PriceBreakdown
	var err error

	taxRateRequest := TaxRateRequest{
		BuyerCountryCode:  buyerCountryCode,
		BuyerTaxID:        buyerTaxID,
		SellerCountryCode: sellerCountryCode,
		LineItemType:      data.LineItemType,
		TransactionDate:   time.Now(),
	}
	taxRate, err := f.taxRateProvider.GetTaxRate(ctx, taxRateRequest)
	if err != nil {
		return LineItem{}, fmt.Errorf("failed to get tax rate: %w", err)
	}

	if data.UnitAmount.IsGross() {
		priceBreakdown, err = NewPriceBreakdownFromGrossAmount(taxRate, data.UnitAmount.Amount(), currency, data.Quantity)
	} else {
		priceBreakdown, err = NewPriceBreakdownFromNetAmount(taxRate, data.UnitAmount.Amount(), currency, data.Quantity)
	}

	if err != nil {
		return LineItem{}, fmt.Errorf("failed to create price breakdown: %w", err)
	}

	return LineItem{
		uuid:         LineItemUUID{common.NewUUIDv7()},
		name:         data.Name,
		breakdown:    priceBreakdown,
		quantity:     data.Quantity,
		lineItemType: data.LineItemType,
	}, nil
}

type DocumentBuilder struct {
	externalReference *string
	documentType      DocumentType
	issueDate         time.Time
	currency          shared.Currency
	seller            LegalEntity
	buyer             LegalEntity
	lineItems         []LineItem
	summary           PriceBreakdownSummary
}

func (b *DocumentBuilder) Build(docNumber DocumentNumber) (*Document, error) {
	return &Document{
		uuid:              DocumentUUID{common.NewUUIDv7()},
		externalReference: b.externalReference,
		documentType:      b.documentType,
		issueDate:         b.issueDate,
		currency:          b.currency,
		documentNumber:    docNumber,
		seller:            b.seller,
		buyer:             b.buyer,
		lineItems:         b.lineItems,
		summary:           b.summary,
	}, nil
}

type DocumentType struct {
	common.Enum[DocumentTypeValues]
}

type DocumentTypeValues string

func (DocumentTypeValues) Values() []string {
	return []string{"receipt"}
}

var DocumentTypeReceipt = common.MustEnum[DocumentType]("receipt")

type DocumentRepository interface {
	DocumentByUUID(ctx context.Context, docUUID DocumentUUID) (*Document, error)
	CreateDocument(
		ctx context.Context,
		series DocumentSeries,
		createFunc func(documentNumber DocumentNumber) (*Document, error),
	) (DocumentUUID, error)
	UpdateFileUrl(ctx context.Context, docUUID DocumentUUID, fileUrl string) error
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
	Name         string
	Quantity     int
	UnitAmount   shared.LineAmount
	LineItemType shared.LineItemType
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
	summary           PriceBreakdownSummary
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

func (d *Document) Summary() PriceBreakdownSummary {
	return d.summary
}

type LineItemUUID struct {
	common.UUID
}

type LineItem struct {
	uuid         LineItemUUID
	name         string
	breakdown    PriceBreakdown
	quantity     int
	lineItemType shared.LineItemType
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

func (l LineItem) PriceBreakdown() PriceBreakdown {
	return l.breakdown
}

func (l LineItem) LineItemType() shared.LineItemType {
	return l.lineItemType
}
