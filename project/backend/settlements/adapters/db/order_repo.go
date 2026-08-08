package db

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"eats/backend/common"
	"eats/backend/settlements/adapters/db/dbmodels"
	"eats/backend/settlements/app/models"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (o *OrderRepository) SaveOrder(ctx context.Context, order models.Order) error {
	return common.UpdateInTx(ctx, o.db, func(ctx context.Context, tx pgx.Tx) error {
		queries := dbmodels.New(tx)

		err := queries.SaveOrder(ctx, dbmodels.SaveOrderParams{
			OrderUuid:           order.UUID(),
			RestaurantUuid:      order.RestaurantUUID(),
			CourierUuid:         order.CourierUUID(),
			Currency:            order.Currency(),
			CommissionNetAmount: order.CommissionNetAmount(),
			OrderedAt:           order.OrderedAt(),
		})
		if err != nil {
			return fmt.Errorf("error saving order: %w", err)
		}

		err = queries.SaveOrderBreakdown(ctx, dbmodels.SaveOrderBreakdownParams{
			OrderUuid:     order.UUID(),
			BreakdownType: dbmodels.SettlementsBreakdownTypeDelivery,
			NetAmount:     order.DeliveryBreakdown().Net,
			TaxAmount:     order.DeliveryBreakdown().Tax,
			GrossAmount:   order.DeliveryBreakdown().Gross,
		})
		if err != nil {
			return fmt.Errorf("error saving order delivery breakdown: %w", err)
		}

		err = queries.SaveOrderBreakdown(ctx, dbmodels.SaveOrderBreakdownParams{
			OrderUuid:     order.UUID(),
			BreakdownType: dbmodels.SettlementsBreakdownTypeItems,
			NetAmount:     order.ItemsBreakdown().Net,
			TaxAmount:     order.ItemsBreakdown().Tax,
			GrossAmount:   order.ItemsBreakdown().Gross,
		})
		if err != nil {
			return fmt.Errorf("error saving order items breakdown: %w", err)
		}

		err = queries.SaveOrderBreakdown(ctx, dbmodels.SaveOrderBreakdownParams{
			OrderUuid:     order.UUID(),
			BreakdownType: dbmodels.SettlementsBreakdownTypeTotal,
			NetAmount:     order.TotalBreakdown().Net,
			TaxAmount:     order.TotalBreakdown().Tax,
			GrossAmount:   order.TotalBreakdown().Gross,
		})
		if err != nil {
			return fmt.Errorf("error saving order total breakdown: %w", err)
		}

		return nil
	})
}
