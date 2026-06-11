// This file contains tests that are executed to verify your solution.
// It's read-only, so all modifications will be ignored.
package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eats/backend/billing/app/command"
	"eats/backend/billing/domain"
	"eats/backend/common"
	"eats/backend/common/shared"
)

func TestIssueReceipt_creates_document(t *testing.T) {
	repo := &mockDocumentRepository{}
	handlers := command.NewHandlers(repo)

	cmd := command.IssueReceipt{
		DocumentData: validDocumentData(t),
	}

	docUUID, err := handlers.IssueReceipt(context.Background(), cmd)
	require.NoError(t, err)
	assert.False(t, docUUID.IsZero(), "returned UUID should not be zero")
	assert.Len(t, repo.documents, 1, "one document should be stored")
}

type mockDocumentRepository struct {
	documents map[domain.DocumentUUID]*domain.Document
}

func (m *mockDocumentRepository) CreateDocument(
	ctx context.Context,
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
	doc, ok := m.documents[docUUID]
	if !ok {
		return nil, errors.New("document not found")
	}
	return doc, nil
}

func validDocumentData(t *testing.T) domain.NewDocumentData {
	t.Helper()

	sellerTaxID, err := shared.NewTaxID("1234567890")
	require.NoError(t, err)

	addr, err := shared.NewAddress("123 Main St", "", "12345", "New York", shared.MustNewCountryCode("US"))
	require.NoError(t, err)

	seller, err := domain.NewLegalEntity("Food Delivery Inc.", addr, &sellerTaxID)
	require.NoError(t, err)

	buyer, err := domain.NewLegalEntity("John Doe", addr, nil)
	require.NoError(t, err)

	return domain.NewDocumentData{
		ExternalReference: common.ToPtr("EXT-REF-001"),
		IssueDate:         time.Date(2025, 3, 14, 0, 0, 0, 0, time.UTC),
		Currency:          shared.MustNewCurrency("USD"),
		Seller:            seller,
		Buyer:             buyer,
		LineItems: []domain.NewLineItemData{
			{
				Name:       "Cheeseburger",
				Quantity:   2,
				UnitAmount: shared.NewNetAmount(decimal.NewFromFloat(10.00)),
			},
		},
	}
}
