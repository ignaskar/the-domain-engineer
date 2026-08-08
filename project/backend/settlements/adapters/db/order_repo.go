package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

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
	// TODO: open a transaction (use common.UpdateInTx) and write to
	// settlements.orders + settlements.order_breakdowns via the sqlc-generated
	// queries.
	panic("not implemented")
}
