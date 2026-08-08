package models

import (
	"context"
	"errors"
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
	orderUUID           OrderUUID
	restaurantUUID      domain.LegalEntityUUID
	courierUUID         domain.LegalEntityUUID
	currency            shared.Currency
	itemsBreakdown      AmountBreakdown
	deliveryBreakdown   AmountBreakdown
	totalBreakdown      AmountBreakdown
	commissionNetAmount decimal.Decimal
	orderedAt           time.Time
}

// AmountBreakdown holds a triple of (net, tax, gross) amounts.
type AmountBreakdown struct {
	Net   decimal.Decimal
	Tax   decimal.Decimal
	Gross decimal.Decimal
}

func NewAmountBreakdown() AmountBreakdown {
	return AmountBreakdown{
		Net:   decimal.Zero,
		Tax:   decimal.Zero,
		Gross: decimal.Zero,
	}
}

func (o AmountBreakdown) Add(net, tax, gross decimal.Decimal) AmountBreakdown {
	return AmountBreakdown{
		Net:   o.Net.Add(net),
		Tax:   o.Tax.Add(tax),
		Gross: o.Gross.Add(gross),
	}
}

func NewOrder(
	orderUUID OrderUUID,
	restaurantUUID domain.LegalEntityUUID,
	courierUUID domain.LegalEntityUUID,
	currency shared.Currency,
	orderedAt time.Time,
	receipt client.DocumentReadModel,
) (Order, error) {
	if orderUUID.IsZero() {
		return Order{}, errors.New("orderUUID cannot be empty")
	}
	if restaurantUUID.IsZero() {
		return Order{}, errors.New("restaurantUUID cannot be empty")
	}
	if courierUUID.IsZero() {
		return Order{}, errors.New("courierUUID cannot be empty")
	}
	if currency.IsZero() {
		return Order{}, errors.New("currency cannot be empty")
	}
	if orderedAt.IsZero() {
		return Order{}, errors.New("orderedAt cannot be empty")
	}

	itemsBreakdown := NewAmountBreakdown()
	deliveryBreakdown := NewAmountBreakdown()
	for _, item := range receipt.LineItems {
		switch item.Type {
		case shared.LineItemTypeFood, shared.LineItemTypeBeverage:
			itemsBreakdown = itemsBreakdown.Add(item.NetAmount, item.TaxAmount, item.GrossAmount)
		case shared.LineItemTypeDelivery:
			deliveryBreakdown = deliveryBreakdown.Add(item.NetAmount, item.TaxAmount, item.GrossAmount)
		}
	}

	totalBreakdown := NewAmountBreakdown()
	totalBreakdown = totalBreakdown.Add(receipt.NetTotal, receipt.TaxTotal, receipt.GrossTotal)
	commissionNetAmount := itemsBreakdown.Net.Mul(CommissionRate).Round(2)

	return Order{
		orderUUID:           orderUUID,
		restaurantUUID:      restaurantUUID,
		courierUUID:         courierUUID,
		currency:            currency,
		itemsBreakdown:      itemsBreakdown,
		deliveryBreakdown:   deliveryBreakdown,
		totalBreakdown:      totalBreakdown,
		commissionNetAmount: commissionNetAmount,
		orderedAt:           orderedAt,
	}, nil
}

// TODO: add getter for each field (UUID, RestaurantUUID, CourierUUID,
// Currency, ItemsBreakdown, DeliveryBreakdown, TotalBreakdown,
// CommissionNetAmount, OrderedAt) plus a ShortID helper that returns the
// last 8 chars of the orderUUID.
func (o Order) UUID() OrderUUID {
	return o.orderUUID
}
func (o Order) RestaurantUUID() domain.LegalEntityUUID {
	return o.restaurantUUID
}
func (o Order) CourierUUID() domain.LegalEntityUUID {
	return o.courierUUID
}

func (o Order) Currency() shared.Currency {
	return o.currency
}
func (o Order) ItemsBreakdown() AmountBreakdown {
	return o.itemsBreakdown
}
func (o Order) DeliveryBreakdown() AmountBreakdown {
	return o.deliveryBreakdown
}
func (o Order) TotalBreakdown() AmountBreakdown {
	return o.totalBreakdown
}
func (o Order) CommissionNetAmount() decimal.Decimal {
	return o.commissionNetAmount
}
func (o Order) OrderedAt() time.Time {
	return o.orderedAt
}
func (o Order) ShortID() string {
	return o.orderUUID.String()[len(o.orderUUID.String())-8:]
}

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
