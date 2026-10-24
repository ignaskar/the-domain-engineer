package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"eats/backend/billing/adapters/db/dbmodels"
	"eats/backend/billing/domain"
	"eats/backend/common"
)

type DocumentRecord struct {
	UUID domain.DocumentUUID
}

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
	createFunc func(documentNumber domain.DocumentNumber) (DocumentRecord, error),
) (domain.DocumentUUID, error) {
	var docUUID domain.DocumentUUID

	err := common.UpdateInReadCommittedTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queries := dbmodels.New(tx)

		n, err := queries.NextDocumentNumber(ctx, series.String())
		if err != nil {
			return fmt.Errorf("failed to get next document number: %w", err)
		}

		dn, err := domain.NewDocumentNumber(series, int(n))
		if err != nil {
			return fmt.Errorf("failed to create document number: %w", err)
		}

		dr, err := createFunc(dn)
		if err != nil {
			return fmt.Errorf("failed to create document record: %w", err)
		}

		err = queries.SaveDocument(ctx, dbmodels.SaveDocumentParams{
			DocumentUuid:   dr.UUID,
			DocumentNumber: dn.String(),
			SeriesPrefix:   series.String(),
		})
		if err != nil {
			return fmt.Errorf("failed to save document: %w", err)
		}

		docUUID = dr.UUID

		return nil
	})
	if err != nil {
		return domain.DocumentUUID{}, err
	}

	return docUUID, nil
}
