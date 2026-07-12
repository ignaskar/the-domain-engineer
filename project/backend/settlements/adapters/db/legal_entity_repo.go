package db

import (
	"context"
	"eats/backend/settlements/adapters/db/dbmodels"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"eats/backend/settlements/app/models"
	"eats/backend/settlements/domain"
)

type LegalEntityRepository struct {
	db *pgxpool.Pool
}

func NewLegalEntityRepository(db *pgxpool.Pool) *LegalEntityRepository {
	if db == nil {
		panic("db is nil")
	}

	return &LegalEntityRepository{
		db: db,
	}
}

func (r *LegalEntityRepository) LegalEntityByUUID(ctx context.Context, uuid domain.LegalEntityUUID) (models.LegalEntity, error) {
	queries := dbmodels.New(r.db)

	le, err := queries.LegalEntityByUUID(ctx, uuid)
	if err != nil {
		return models.LegalEntity{}, fmt.Errorf("error getting legal entity by uuid: %w", err)
	}

	return legalEntityFromDB(le), nil
}

func (r *LegalEntityRepository) SavePlatformEntity(ctx context.Context, platform models.LegalEntity) error {
	queries := dbmodels.New(r.db)

	err := queries.SaveLegalEntity(ctx, dbmodels.SaveLegalEntityParams{
		LegalEntityUuid:   platform.UUID,
		LegalEntityType:   models.LegalEntityPlatform,
		BusinessName:      platform.BusinessName,
		Address:           platform.Address,
		TaxID:             platform.TaxID,
		BankAccountNumber: platform.BankAccountNumber.String(),
		Currency:          platform.Currency,
	})
	if err != nil {
		return fmt.Errorf("error saving legal entity: %w", err)
	}

	return nil
}

func legalEntityFromDB(l dbmodels.SettlementsLegalEntity) models.LegalEntity {
	return models.LegalEntity{
		UUID:              l.LegalEntityUuid,
		Type:              l.LegalEntityType,
		BusinessName:      l.BusinessName,
		TaxID:             l.TaxID,
		Address:           l.Address,
		BankAccountNumber: domain.UnmarshalIBAN(l.BankAccountNumber),
		Currency:          l.Currency,
	}
}
