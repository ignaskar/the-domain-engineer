package models

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"eats/backend/billing/api/module/client"
	"eats/backend/common"
	"eats/backend/common/shared"
	"eats/backend/settlements/domain"
)

// CommissionRate is the platform's cut of items.Net (food + beverage).
var CommissionRate = decimal.New(20, -2) // 20%

type OrderUUID struct {
	common.UUID
}

type OrderRepository interface {
	SaveOrder(ctx context.Context, order Order) error
}

// Order is the settlements-side aggregate that captures everything we learn
// when StartSettlement runs: the parties involved, the receipt totals, and
// the platform commission.
type Order struct {
	// TODO: add fields for orderUUID, restaurantUUID, courierUUID, currency,
	// itemsBreakdown, deliveryBreakdown, totalBreakdown, commissionNetAmount,
	// and orderedAt.
}

// AmountBreakdown holds a triple of (net, tax, gross) amounts.
type AmountBreakdown struct {
	Net   decimal.Decimal
	Tax   decimal.Decimal
	Gross decimal.Decimal
}

func NewAmountBreakdown() AmountBreakdown {
	// TODO: return a zero-valued breakdown.
	panic("not implemented")
}

func (o AmountBreakdown) Add(net, tax, gross decimal.Decimal) AmountBreakdown {
	// TODO: return a new AmountBreakdown with the amounts added. Don't mutate o.
	panic("not implemented")
}

func NewOrder(
	orderUUID OrderUUID,
	restaurantUUID domain.LegalEntityUUID,
	courierUUID domain.LegalEntityUUID,
	currency shared.Currency,
	orderedAt time.Time,
	receipt client.DocumentReadModel,
) (Order, error) {
	// TODO: implement.
	// The implementation must:
	//   - Reject zero values for orderUUID, restaurantUUID, courierUUID, currency,
	//     and orderedAt.
	//   - Aggregate Food and Beverage line items into itemsBreakdown.
	//   - Aggregate Delivery line items into deliveryBreakdown.
	//   - Set totalBreakdown from receipt.NetTotal/TaxTotal/GrossTotal.
	//   - Compute commissionNetAmount as itemsBreakdown.Net * CommissionRate,
	//     rounded to 2 decimal places.
	panic("not implemented")
}

// TODO: add getter for each field (UUID, RestaurantUUID, CourierUUID,
// Currency, ItemsBreakdown, DeliveryBreakdown, TotalBreakdown,
// CommissionNetAmount, OrderedAt) plus a ShortID helper that returns the
// last 8 chars of the orderUUID.

// UnmarshalOrder rebuilds an Order from already-validated state.
func UnmarshalOrder(
	orderUUID OrderUUID,
	restaurantUUID domain.LegalEntityUUID,
	courierUUID domain.LegalEntityUUID,
	currency shared.Currency,
	itemsBreakdown AmountBreakdown,
	deliveryBreakdown AmountBreakdown,
	totalBreakdown AmountBreakdown,
	commissionNetAmount decimal.Decimal,
	orderedAt time.Time,
) Order {
	return Order{
		orderUUID:           orderUUID,
		restaurantUUID:      restaurantUUID,
		courierUUID:         courierUUID,
		currency:            currency,
		commissionNetAmount: commissionNetAmount,
		itemsBreakdown:      itemsBreakdown,
		deliveryBreakdown:   deliveryBreakdown,
		totalBreakdown:      totalBreakdown,
		orderedAt:           orderedAt,
	}
}
