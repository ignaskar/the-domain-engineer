package app

import (
	"context"

	"github.com/shopspring/decimal"

	"eats/backend/settlements/app/models"
	"eats/backend/settlements/domain"
)

type StatementService struct {
	repository            models.BillingCycleRepository
	legalEntityRepository models.LegalEntityRepository
}

func NewStatementService(
	repository models.BillingCycleRepository,
	legalEntityRepository models.LegalEntityRepository,
) *StatementService {
	return &StatementService{
		repository:            repository,
		legalEntityRepository: legalEntityRepository,
	}
}

type CommissionInvoice struct {
	Number      string
	NetAmount   decimal.Decimal
	TaxAmount   decimal.Decimal
	GrossAmount decimal.Decimal
}

func (s *StatementService) GenerateRestaurantStatementSnapshot(ctx context.Context, bc *domain.BillingCycle, orders []models.Order, commissionInvoice CommissionInvoice) (RestaurantStatementSnapshot, error) {
	// TODO: implement
	return RestaurantStatementSnapshot{}, nil
}

func (s *StatementService) GenerateCourierStatementSnapshot(ctx context.Context, bc *domain.BillingCycle, orders []models.Order) (CourierStatementSnapshot, error) {
	// TODO: implement
	return CourierStatementSnapshot{}, nil
}
