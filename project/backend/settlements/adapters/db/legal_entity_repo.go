package db

import (
	"context"
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
	// TODO: implement
	// 1. Use dbmodels.New(r.db) to get a Queries handle.
	// 2. Call queries.LegalEntityByUUID and translate the row into models.LegalEntity.
	return models.LegalEntity{}, fmt.Errorf("not implemented")
}

func (r *LegalEntityRepository) SavePlatformEntity(ctx context.Context, platform models.LegalEntity) error {
	// TODO: implement
	// Save the platform legal entity using queries.SaveLegalEntity.
	// Set legal_entity_type to models.LegalEntityPlatform.
	return fmt.Errorf("not implemented")
}
