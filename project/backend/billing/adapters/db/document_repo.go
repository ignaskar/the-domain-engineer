package db

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"eats/backend/billing/adapters/db/dbmodels"
	"eats/backend/billing/domain"
	"eats/backend/common"
	"eats/backend/common/log"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) CreateDocument(
	ctx context.Context,
	series domain.DocumentSeries,
	createFunc func(documentNumber domain.DocumentNumber) (*domain.Document, error),
) (domain.DocumentUUID, error) {
	var docUUID domain.DocumentUUID
	var externalReference string

	// ReadCommitted is safe here: NextDocumentNumber uses "last_number = last_number + 1",
	// so PostgreSQL re-evaluates the expression on the current row value after acquiring
	// the row lock. There is no lost update risk because no value is read into Go memory
	// and written back. RepeatableRead would cause serialization errors under contention.
	err := common.UpdateInReadCommittedTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queries := dbmodels.New(tx)

		nextNumber, err := queries.NextDocumentNumber(ctx, series.String())
		if err != nil {
			return fmt.Errorf("error getting next document number: %w", err)
		}

		docNumber, err := domain.NewDocumentNumber(series, int(nextNumber))
		if err != nil {
			return fmt.Errorf("error creating document number: %w", err)
		}

		document, err := createFunc(docNumber)
		if err != nil {
			return fmt.Errorf("error creating document: %w", err)
		}
		docUUID = document.UUID()

		if document.ExternalReference() != nil {
			externalReference = *document.ExternalReference()
		}

		err = queries.SaveDocument(ctx, dbmodels.SaveDocumentParams{
			DocumentUuid:      document.UUID(),
			ExternalReference: document.ExternalReference(),
			DocumentNumber:    docNumber.String(),
			SeriesPrefix:      series.String(),
			DocumentType:      document.DocumentType(),
			IssueDate:         document.IssueDate(),
			Currency:          document.Currency(),
			TotalNetAmount:    document.Summary().NetAmount(),
			TotalTaxAmount:    document.Summary().TaxAmount(),
			TotalGrossAmount:  document.Summary().GrossAmount(),
		})
		if err != nil {
			return fmt.Errorf("error saving document: %w", err)
		}

		for _, li := range document.LineItems() {
			err = queries.SaveDocumentLineItem(ctx, dbmodels.SaveDocumentLineItemParams{
				LineItemUuid:    li.UUID(),
				DocumentUuid:    document.UUID(),
				Name:            li.Name(),
				Quantity:        int32(li.Quantity()),
				UnitNetAmount:   li.PriceBreakdown().UnitNetAmount(),
				UnitTaxAmount:   li.PriceBreakdown().UnitTaxAmount(),
				UnitGrossAmount: li.PriceBreakdown().UnitGrossAmount(),
				NetAmount:       li.PriceBreakdown().NetAmount(),
				TaxAmount:       li.PriceBreakdown().TaxAmount(),
				GrossAmount:     li.PriceBreakdown().GrossAmount(),
				TaxRate:         li.PriceBreakdown().TaxRate().Rate(),
				TaxType:         li.PriceBreakdown().TaxRate().TaxType(),
			})
			if err != nil {
				return fmt.Errorf("error saving document line item: %w", err)
			}
		}

		return nil
	})

	// We can't handle this with ON CONFLICT - we have to cancel the transaction not to generate a new document number
	if common.IsUniqueViolationError(err, "documents_external_reference_key") {
		logger := log.FromContext(ctx)
		logger.With("external_reference", externalReference).Info("Skipping document creation due to existing external reference")

		// This is outside of transaction, but it's okay - this is a read operation
		dbDoc, err := r.getDocumentByExternalReference(ctx, externalReference)
		if err != nil {
			return domain.DocumentUUID{}, fmt.Errorf("error retrieving existing document by external reference: %w", err)
		}

		return dbDoc.DocumentUuid, nil
	}
	if err != nil {
		return domain.DocumentUUID{}, err
	}

	return docUUID, nil
}

func (r *PostgresRepository) DocumentByUUID(ctx context.Context, uuid domain.DocumentUUID) (*domain.Document, error) {
	queries := dbmodels.New(r.db)

	dbDocument, err := queries.GetDocument(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("error getting document by uuid: %w", err)
	}

	dbLineItems, err := queries.GetDocumentLineItems(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("error getting document line items: %w", err)
	}

	var lineItems []domain.LineItem
	for _, li := range dbLineItems {
		breakdown := domain.UnmarshalPriceBreakdown(
			domain.UnmarshalTaxRate(li.TaxRate, li.TaxType),
			li.UnitNetAmount,
			li.UnitTaxAmount,
			li.UnitGrossAmount,
			li.NetAmount,
			li.TaxAmount,
			li.GrossAmount,
		)

		lineItems = append(lineItems, domain.UnmarshalLineItem(
			li.LineItemUuid,
			li.Name,
			breakdown,
			int(li.Quantity),
		))
	}

	documentSeries, err := domain.NewDocumentSeries(dbDocument.SeriesPrefix)
	if err != nil {
		return nil, fmt.Errorf("error creating document series: %w", err)
	}

	documentNumber, err := domain.UnmarshalDocumentNumber(
		documentSeries,
		dbDocument.DocumentNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating document number: %w", err)
	}

	breakdownSummary := domain.UnmarshalPriceBreakdownSummary(
		dbDocument.TotalNetAmount,
		dbDocument.TotalTaxAmount,
		dbDocument.TotalGrossAmount,
		nil,
	)

	return domain.UnmarshalDocument(
		dbDocument.DocumentUuid,
		dbDocument.ExternalReference,
		documentNumber,
		dbDocument.DocumentType,
		dbDocument.IssueDate,
		dbDocument.Currency,
		domain.LegalEntity{},
		domain.LegalEntity{},
		lineItems,
		breakdownSummary,
	), nil
}

func (r *PostgresRepository) getDocumentByExternalReference(ctx context.Context, externalRef string) (dbmodels.BillingDocument, error) {
	queries := dbmodels.New(r.db)

	dbDoc, err := queries.GetDocumentByExternalReference(ctx, &externalRef)
	if err != nil {
		return dbmodels.BillingDocument{}, fmt.Errorf("error getting document by external reference: %w", err)
	}

	return dbDoc, nil
}
