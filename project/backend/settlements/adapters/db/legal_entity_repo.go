package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"eats/backend/settlements/adapters/db/dbmodels"
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

	legalEntity, err := queries.LegalEntityByUUID(ctx, uuid)
	if err != nil {
		return models.LegalEntity{}, fmt.Errorf("error getting legal entity by uuid: %w", err)
	}

	return legalEntityFromDB(legalEntity), nil
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

func (r *LegalEntityRepository) SavePlatformEntity(ctx context.Context, platform models.LegalEntity) error {
	queries := dbmodels.New(r.db)

	err := queries.SaveLegalEntity(ctx, dbmodels.SaveLegalEntityParams{
		LegalEntityUuid:   platform.UUID,
		LegalEntityType:   models.LegalEntityPlatform,
		BusinessName:      platform.BusinessName,
		TaxID:             platform.TaxID,
		Address:           platform.Address,
		BankAccountNumber: platform.BankAccountNumber.String(),
		Currency:          platform.Currency,
	})
	if err != nil {
		return fmt.Errorf("error saving legal entity: %w", err)
	}

	return nil
}

func (r *LegalEntityRepository) PartnerByUUID(ctx context.Context, uuid domain.LegalEntityUUID) (models.Partner, error) {
	// TODO: implement
	// 1. Query legal entity by UUID.
	// 2. Query platform UUID via the partner_platform_mappings table.
	// 3. Return a Partner combining both.
	return models.Partner{}, fmt.Errorf("not implemented")
}

func (r *LegalEntityRepository) SavePartner(ctx context.Context, partner models.Partner) error {
	// TODO: implement
	// In a single transaction (use common.UpdateInTx):
	// 1. Save the partner's legal entity.
	// 2. Verify the platform entity exists and is of type "platform".
	// 3. Save the partner-platform mapping.
	return fmt.Errorf("not implemented")
}
