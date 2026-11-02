// This file contains tests that are executed to verify your solution.
// It's read-only, so all modifications will be ignored.
package query_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eats/backend/billing/app/query"
	"eats/backend/billing/domain"
	"eats/backend/common"
	"eats/backend/common/shared"
)

func TestGetDocumentByUUID_returns_document(t *testing.T) {
	doc := newTestDocument(t)

	repo := &mockDocumentRepository{
		documents: map[domain.DocumentUUID]*domain.Document{
			doc.UUID(): doc,
		},
	}
	handlers := query.NewHandlers(repo)

	q := query.GetDocumentByUUID{
		DocumentUUID: doc.UUID(),
	}

	got, err := handlers.GetDocumentByUUID(context.Background(), q)
	require.NoError(t, err)
	assert.Equal(t, doc.UUID(), got.UUID())
}

func TestGetDocumentByUUID_returns_error_when_not_found(t *testing.T) {
	repo := &mockDocumentRepository{}
	handlers := query.NewHandlers(repo)

	q := query.GetDocumentByUUID{
		DocumentUUID: domain.DocumentUUID{UUID: common.NewUUIDv7()},
	}

	_, err := handlers.GetDocumentByUUID(context.Background(), q)
	require.Error(t, err)
}

type mockDocumentRepository struct {
	documents map[domain.DocumentUUID]*domain.Document
}

func (m *mockDocumentRepository) CreateDocument(
	_ context.Context,
	series domain.DocumentSeries,
	createFunc func(documentNumber domain.DocumentNumber) (*domain.Document, error),
) (domain.DocumentUUID, error) {
	docNumber, err := domain.NewDocumentNumber(series, 1)
	if err != nil {
		return domain.DocumentUUID{}, err
	}

	doc, err := createFunc(docNumber)
	if err != nil {
		return domain.DocumentUUID{}, err
	}

	if m.documents == nil {
		m.documents = make(map[domain.DocumentUUID]*domain.Document)
	}
	m.documents[doc.UUID()] = doc

	return doc.UUID(), nil
}

func (m *mockDocumentRepository) DocumentByUUID(_ context.Context, docUUID domain.DocumentUUID) (*domain.Document, error) {
	if m.documents == nil {
		return nil, errors.New("document not found")
	}
	doc, ok := m.documents[docUUID]
	if !ok {
		return nil, errors.New("document not found")
	}
	return doc, nil
}

func newTestDocument(t *testing.T) *domain.Document {
	t.Helper()

	sellerTaxID, err := shared.NewTaxID("1234567890")
	require.NoError(t, err)

	addr, err := shared.NewAddress("123 Main St", "", "12345", "New York", shared.MustNewCountryCode("US"))
	require.NoError(t, err)

	seller, err := domain.NewLegalEntity("Food Delivery Inc.", addr, &sellerTaxID)
	require.NoError(t, err)

	buyer, err := domain.NewLegalEntity("John Doe", addr, nil)
	require.NoError(t, err)

	series, err := domain.NewDocumentSeries("R")
	require.NoError(t, err)

	docNumber, err := domain.NewDocumentNumber(series, 1)
	require.NoError(t, err)

	doc, err := domain.NewReceipt(domain.NewDocumentData{
		ExternalReference: common.ToPtr("EXT-REF-001"),
		IssueDate:         time.Date(2025, 3, 14, 0, 0, 0, 0, time.UTC),
		Currency:          shared.MustNewCurrency("USD"),
		Seller:            seller,
		Buyer:             buyer,
		LineItems: []domain.NewLineItemData{
			{
				Name:       "Cheeseburger",
				Quantity:   1,
				UnitAmount: shared.NewNetAmount(decimal.NewFromFloat(10.00)),
			},
		},
	}, docNumber)
	require.NoError(t, err)

	return doc
}
