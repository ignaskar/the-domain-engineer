package db

import (
	"context"
	"fmt"

	"eats/backend/billing/adapters/db/dbmodels"
	"eats/backend/common"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"eats/backend/billing/domain"
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
			return fmt.Errorf("next document number failed: %w", err)
		}

		docNumber, err := domain.NewDocumentNumber(series, int(n))
		if err != nil {
			return err
		}

		docRecord, err := createFunc(docNumber)
		if err != nil {
			return err
		}

		err = queries.SaveDocument(ctx, dbmodels.SaveDocumentParams{
			DocumentUuid:   docRecord.UUID,
			DocumentNumber: docNumber.String(),
			SeriesPrefix:   series.String(),
		})

		docUUID = docRecord.UUID
		return nil
	})
	if err != nil {
		return domain.DocumentUUID{}, err
	}

	return docUUID, nil
}
