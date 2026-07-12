package db

import (
	"context"
	"fmt"

	"eats/backend/common"

	pgx "github.com/jackc/pgx/v5"
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
	queries := dbmodels.New(r.db)

	le, err := queries.LegalEntityByUUID(ctx, uuid)
	if err != nil {
		return models.Partner{}, fmt.Errorf("error getting legal entity by uuid: %w", err)
	}

	pid, err := queries.PlatformByPartnerUUID(ctx, uuid)
	if err != nil {
		return models.Partner{}, fmt.Errorf("error getting partner by uuid: %w", err)
	}

	return models.Partner{
		LegalEntity:        legalEntityFromDB(le),
		PlatformEntityUUID: pid,
	}, nil
}

func (r *LegalEntityRepository) SavePartner(ctx context.Context, partner models.Partner) error {
	// TODO: implement
	// In a single transaction (use common.UpdateInTx):
	// 1. Save the partner's legal entity.
	// 2. Verify the platform entity exists and is of type "platform".
	// 3. Save the partner-platform mapping.
	return common.UpdateInTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queries := dbmodels.New(tx)

		err := queries.SaveLegalEntity(ctx, dbmodels.SaveLegalEntityParams{
			LegalEntityUuid:   partner.LegalEntity.UUID,
			LegalEntityType:   models.LegalEntityPartner,
			BusinessName:      partner.LegalEntity.BusinessName,
			TaxID:             partner.LegalEntity.TaxID,
			Address:           partner.LegalEntity.Address,
			BankAccountNumber: partner.LegalEntity.BankAccountNumber.String(),
			Currency:          partner.LegalEntity.Currency,
		})
		if err != nil {
			return fmt.Errorf("error saving legal entity: %w", err)
		}

		le, err := queries.LegalEntityByUUID(ctx, partner.PlatformEntityUUID.LegalEntityUUID)
		if err != nil {
			return fmt.Errorf("error getting legal entity by uuid: %w", err)
		}

		if le.LegalEntityType != models.LegalEntityPlatform {
			return fmt.Errorf("legal entity %s is not a platform", partner.PlatformEntityUUID)
		}

		err = queries.SavePartnerPlatformMapping(ctx, dbmodels.SavePartnerPlatformMappingParams{
			PartnerUuid:        partner.LegalEntity.UUID,
			PlatformEntityUuid: partner.PlatformEntityUUID,
		})
		if err != nil {
			return fmt.Errorf("error saving partner platform mapping: %w", err)
		}

		return nil
	})
}
