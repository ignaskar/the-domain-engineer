// This file contains tests that are executed to verify your solution.
// It's read-only, so all modifications will be ignored.
package app_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eats/backend/common"
	"eats/backend/common/shared"
	"eats/backend/settlements/app"
	"eats/backend/settlements/app/models"
	"eats/backend/settlements/domain"
)

const (
	platformName   = "Three Dots Eats"
	restaurantName = "Pizza Place"
	courierName    = "Speedy Couriers"

	burgerBarnName = "Burger Barn"
	sushiSpotName  = "Sushi Spot"

	restaurantCycleNumber = 7
	courierCycleNumber    = 3
)

// fakeLegalEntityRepository is an in-memory models.LegalEntityRepository.
// The statement service only reads (PartnerByUUID, LegalEntityByUUID), so the
// Save* methods are no-ops.
type fakeLegalEntityRepository struct {
	legalEntities map[domain.LegalEntityUUID]models.LegalEntity
	partners      map[domain.LegalEntityUUID]models.Partner
}

func (f *fakeLegalEntityRepository) LegalEntityByUUID(_ context.Context, uuid domain.LegalEntityUUID) (models.LegalEntity, error) {
	le, ok := f.legalEntities[uuid]
	if !ok {
		return models.LegalEntity{}, fmt.Errorf("legal entity %s not found", uuid)
	}
	return le, nil
}

func (f *fakeLegalEntityRepository) PartnerByUUID(_ context.Context, uuid domain.LegalEntityUUID) (models.Partner, error) {
	p, ok := f.partners[uuid]
	if !ok {
		return models.Partner{}, fmt.Errorf("partner %s not found", uuid)
	}
	return p, nil
}

func (f *fakeLegalEntityRepository) SavePartner(_ context.Context, _ models.Partner, _ *domain.BillingCycle) error {
	return nil
}

func (f *fakeLegalEntityRepository) SavePlatformEntity(_ context.Context, _ models.LegalEntity) error {
	return nil
}

// moneyCmpOpts lets cmp.Diff compare Money/Percentage values. decimal.Decimal
// and the currency enum both carry unexported fields, so we register comparers
// that short-circuit cmp's recursion into them.
func moneyCmpOpts() []cmp.Option {
	return []cmp.Option{
		cmp.Comparer(func(x, y decimal.Decimal) bool { return x.Equal(y) }),
		cmp.Comparer(func(x, y shared.Currency) bool { return x == y }),
	}
}

func newLegalEntity(
	t *testing.T,
	uuid domain.LegalEntityUUID,
	entityType models.LegalEntityType,
	businessName string,
	currency shared.Currency,
) models.LegalEntity {
	t.Helper()

	taxID, err := shared.NewTaxID("PL1234567890")
	require.NoError(t, err)

	address, err := shared.NewAddress("Main Street 1", "", "00-001", "Warsaw", shared.MustNewCountryCode("PL"))
	require.NoError(t, err)

	iban, err := domain.NewIBAN("DE89370400440532013000")
	require.NoError(t, err)

	le, err := models.NewLegalEntity(uuid, entityType, businessName, taxID, address, iban, currency)
	require.NoError(t, err)

	return le
}

func dec(t *testing.T, s string) decimal.Decimal {
	t.Helper()

	d, err := decimal.NewFromString(s)
	require.NoError(t, err)

	return d
}

func money(t *testing.T, amount string, currency shared.Currency) app.Money {
	t.Helper()

	return app.Money{Amount: dec(t, amount), Currency: currency}
}

func breakdown(t *testing.T, net, tax, gross string) models.AmountBreakdown {
	t.Helper()

	return models.AmountBreakdown{
		Net:   dec(t, net),
		Tax:   dec(t, tax),
		Gross: dec(t, gross),
	}
}

// restaurantScenario wires a closed restaurant billing cycle with two orders
// against a fully-populated fake repository.
type restaurantScenario struct {
	repo      *fakeLegalEntityRepository
	bc        *domain.BillingCycle
	orders    []models.Order
	usd       shared.Currency
	startDate time.Time
	endDate   time.Time
}

func newRestaurantScenario(t *testing.T) restaurantScenario {
	t.Helper()

	usd := shared.MustNewCurrency("USD")

	platformUUID := domain.LegalEntityUUID{UUID: common.NewUUIDv7()}
	restaurantUUID := domain.LegalEntityUUID{UUID: common.NewUUIDv7()}
	courierUUID := domain.LegalEntityUUID{UUID: common.NewUUIDv7()}

	platform := newLegalEntity(t, platformUUID, models.LegalEntityPlatform, platformName, usd)
	restaurant := newLegalEntity(t, restaurantUUID, models.LegalEntityPartner, restaurantName, usd)
	courier := newLegalEntity(t, courierUUID, models.LegalEntityPartner, courierName, usd)

	repo := &fakeLegalEntityRepository{
		legalEntities: map[domain.LegalEntityUUID]models.LegalEntity{
			platformUUID: platform,
			courierUUID:  courier,
		},
		partners: map[domain.LegalEntityUUID]models.Partner{
			restaurantUUID: models.NewPartner(restaurant, models.PlatformEntityUUID{LegalEntityUUID: platformUUID}),
		},
	}

	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	bc := domain.UnmarshalBillingCycle(
		domain.BillingCycleUUID{UUID: common.NewUUIDv7()},
		restaurantUUID,
		domain.PartnerTypeRestaurant,
		restaurantCycleNumber,
		true,  // closed
		false, // settled
		startDate,
		&endDate,
	)

	// Items gross feeds ordersGrossAmount; delivery breakdown feeds the
	// delivery summary. Numbers are picked so the sums are easy to verify:
	//   ordersGross   = 12.30 + 24.60 = 36.90
	//   deliveryGross =  4.92 +  6.15 = 11.07
	orders := []models.Order{
		models.UnmarshalOrder(
			models.OrderUUID{UUID: common.NewUUIDv7()},
			restaurantUUID,
			courierUUID,
			usd,
			breakdown(t, "10.00", "2.30", "12.30"),
			breakdown(t, "4.00", "0.92", "4.92"),
			breakdown(t, "14.00", "3.22", "17.22"),
			dec(t, "2.00"),
			startDate.Add(24*time.Hour),
		),
		models.UnmarshalOrder(
			models.OrderUUID{UUID: common.NewUUIDv7()},
			restaurantUUID,
			courierUUID,
			usd,
			breakdown(t, "20.00", "4.60", "24.60"),
			breakdown(t, "5.00", "1.15", "6.15"),
			breakdown(t, "25.00", "5.75", "30.75"),
			dec(t, "4.00"),
			startDate.Add(48*time.Hour),
		),
	}

	return restaurantScenario{
		repo:      repo,
		bc:        bc,
		orders:    orders,
		usd:       usd,
		startDate: startDate,
		endDate:   endDate,
	}
}

// App-layer tests like these usually aren't needed: domain and component tests
// normally cover this behavior. We add them here only to verify your solution,
// since the statement logic runs entirely in memory with nothing wired up to it yet.
func TestStatementService_GenerateRestaurantStatementSnapshot(t *testing.T) {
	s := newRestaurantScenario(t)

	commissionInvoice := app.CommissionInvoice{
		Number:      "COM/2026/01/0001",
		NetAmount:   dec(t, "6.00"),
		TaxAmount:   dec(t, "1.38"),
		GrossAmount: dec(t, "7.38"),
	}

	service := app.NewStatementService(nil, s.repo)

	snapshot, err := service.GenerateRestaurantStatementSnapshot(context.Background(), s.bc, s.orders, commissionInvoice)
	require.NoError(t, err)

	// Whole-summary comparison: accumulations, commission echo, and the payout
	// (payoutGross = ordersGross - commissionGross = 36.90 - 7.38 = 29.52).
	wantSummary := app.RestaurantStatementSummary{
		OrdersNetAmount:         money(t, "30.00", s.usd),
		OrdersTaxAmount:         money(t, "6.90", s.usd),
		OrdersGrossAmount:       money(t, "36.90", s.usd),
		CommissionInvoiceNumber: "COM/2026/01/0001",
		CommissionPercent:       app.Percentage{Value: models.CommissionRate},
		CommissionNetAmount:     money(t, "6.00", s.usd),
		CommissionTaxAmount:     money(t, "1.38", s.usd),
		CommissionGrossAmount:   money(t, "7.38", s.usd),
		DeliveryNetAmount:       money(t, "9.00", s.usd),
		DeliveryTaxAmount:       money(t, "2.07", s.usd),
		DeliveryGrossAmount:     money(t, "11.07", s.usd),
		PayoutGrossAmount:       money(t, "29.52", s.usd),
	}
	if diff := cmp.Diff(wantSummary, snapshot.Summary, moneyCmpOpts()...); diff != "" {
		t.Errorf("summary mismatch (-want +got):\n%s", diff)
	}

	// The headline number gets an explicit check for a clear failure message.
	assert.True(t,
		dec(t, "29.52").Equal(snapshot.Summary.PayoutGrossAmount.Amount),
		"payout gross %s", snapshot.Summary.PayoutGrossAmount.Amount,
	)

	// Currency and period reflect the inputs.
	assert.Equal(t, s.usd, snapshot.Currency)
	assert.Equal(t, app.NewDateRange(s.startDate, s.endDate), snapshot.Period)
	assert.Equal(t, s.bc.UUID(), snapshot.BillingCycleUUID)

	// DocumentName is non-empty and follows the "Restaurant Statement ..." format.
	assert.NotEmpty(t, snapshot.DocumentName)
	assert.True(t,
		strings.HasPrefix(snapshot.DocumentName, "Restaurant Statement "),
		"document name: %s", snapshot.DocumentName,
	)
	assert.Contains(t, snapshot.DocumentName, fmt.Sprintf("#%d", restaurantCycleNumber))

	// Legal entities are resolved through the repository.
	assert.Equal(t, restaurantName, snapshot.Restaurant.Name)
	assert.Equal(t, platformName, snapshot.Platform.Name)
	assert.False(t, snapshot.UUID.IsZero())

	// Every order resolves its courier name through the repository.
	require.Len(t, snapshot.Orders, 2)
	for i, o := range snapshot.Orders {
		assert.Equal(t, courierName, o.CourierName, "order %d courier name", i)
	}
}

func TestStatementService_GenerateRestaurantStatementSnapshot_Errors(t *testing.T) {
	t.Run("not_closed_cycle", func(t *testing.T) {
		service := app.NewStatementService(nil, &fakeLegalEntityRepository{})

		bc := domain.UnmarshalBillingCycle(
			domain.BillingCycleUUID{UUID: common.NewUUIDv7()},
			domain.LegalEntityUUID{UUID: common.NewUUIDv7()},
			domain.PartnerTypeRestaurant,
			1,
			false, // not closed
			false,
			time.Now().UTC(),
			nil,
		)

		_, err := service.GenerateRestaurantStatementSnapshot(context.Background(), bc, nil, app.CommissionInvoice{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not closed")
	})

	t.Run("not_a_restaurant", func(t *testing.T) {
		service := app.NewStatementService(nil, &fakeLegalEntityRepository{})

		endDate := time.Now().UTC()
		bc := domain.UnmarshalBillingCycle(
			domain.BillingCycleUUID{UUID: common.NewUUIDv7()},
			domain.LegalEntityUUID{UUID: common.NewUUIDv7()},
			domain.PartnerTypeCourier, // wrong partner type
			1,
			true,
			false,
			endDate.AddDate(0, 0, -7),
			&endDate,
		)

		_, err := service.GenerateRestaurantStatementSnapshot(context.Background(), bc, nil, app.CommissionInvoice{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a restaurant")
	})

	t.Run("commission_greater_or_equal_to_orders_gross", func(t *testing.T) {
		s := newRestaurantScenario(t)

		// ordersGross is 36.90; an equal commission gross must trip the sanity check.
		commissionInvoice := app.CommissionInvoice{
			Number:      "COM/2026/01/0001",
			NetAmount:   dec(t, "30.00"),
			TaxAmount:   dec(t, "6.90"),
			GrossAmount: dec(t, "36.90"),
		}

		service := app.NewStatementService(nil, s.repo)

		_, err := service.GenerateRestaurantStatementSnapshot(context.Background(), s.bc, s.orders, commissionInvoice)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "greater than or equal to orders gross amount")
	})
}

func TestStatementService_GenerateCourierStatementSnapshot(t *testing.T) {
	usd := shared.MustNewCurrency("USD")

	platformUUID := domain.LegalEntityUUID{UUID: common.NewUUIDv7()}
	courierPartnerUUID := domain.LegalEntityUUID{UUID: common.NewUUIDv7()}
	restaurant1UUID := domain.LegalEntityUUID{UUID: common.NewUUIDv7()}
	restaurant2UUID := domain.LegalEntityUUID{UUID: common.NewUUIDv7()}

	platform := newLegalEntity(t, platformUUID, models.LegalEntityPlatform, platformName, usd)
	courier := newLegalEntity(t, courierPartnerUUID, models.LegalEntityPartner, courierName, usd)
	restaurant1 := newLegalEntity(t, restaurant1UUID, models.LegalEntityPartner, burgerBarnName, usd)
	restaurant2 := newLegalEntity(t, restaurant2UUID, models.LegalEntityPartner, sushiSpotName, usd)

	repo := &fakeLegalEntityRepository{
		legalEntities: map[domain.LegalEntityUUID]models.LegalEntity{
			platformUUID:    platform,
			restaurant1UUID: restaurant1,
			restaurant2UUID: restaurant2,
		},
		partners: map[domain.LegalEntityUUID]models.Partner{
			courierPartnerUUID: models.NewPartner(courier, models.PlatformEntityUUID{LegalEntityUUID: platformUUID}),
		},
	}

	startDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)

	bc := domain.UnmarshalBillingCycle(
		domain.BillingCycleUUID{UUID: common.NewUUIDv7()},
		courierPartnerUUID,
		domain.PartnerTypeCourier,
		courierCycleNumber,
		true,
		false,
		startDate,
		&endDate,
	)

	// deliveryGross = 4.92 + 6.77 = 11.69 (the courier payout).
	orders := []models.Order{
		models.UnmarshalOrder(
			models.OrderUUID{UUID: common.NewUUIDv7()},
			restaurant1UUID,
			courierPartnerUUID,
			usd,
			models.NewAmountBreakdown(),
			breakdown(t, "4.00", "0.92", "4.92"),
			models.NewAmountBreakdown(),
			decimal.Zero,
			startDate.Add(24*time.Hour),
		),
		models.UnmarshalOrder(
			models.OrderUUID{UUID: common.NewUUIDv7()},
			restaurant2UUID,
			courierPartnerUUID,
			usd,
			models.NewAmountBreakdown(),
			breakdown(t, "5.50", "1.27", "6.77"),
			models.NewAmountBreakdown(),
			decimal.Zero,
			startDate.Add(48*time.Hour),
		),
	}

	service := app.NewStatementService(nil, repo)

	snapshot, err := service.GenerateCourierStatementSnapshot(context.Background(), bc, orders)
	require.NoError(t, err)

	wantSummary := app.CourierStatementSummary{
		TotalDeliveries:     2,
		DeliveryNetAmount:   money(t, "9.50", usd),
		DeliveryTaxAmount:   money(t, "2.19", usd),
		DeliveryGrossAmount: money(t, "11.69", usd),
		PayoutGrossAmount:   money(t, "11.69", usd),
	}
	if diff := cmp.Diff(wantSummary, snapshot.Summary, moneyCmpOpts()...); diff != "" {
		t.Errorf("summary mismatch (-want +got):\n%s", diff)
	}

	// Payout equals the sum of delivery gross, and TotalDeliveries equals the order count.
	assert.True(t,
		dec(t, "11.69").Equal(snapshot.Summary.PayoutGrossAmount.Amount),
		"payout gross %s", snapshot.Summary.PayoutGrossAmount.Amount,
	)
	assert.Equal(t, len(orders), snapshot.Summary.TotalDeliveries)

	// Currency and period reflect the inputs.
	assert.Equal(t, usd, snapshot.Currency)
	assert.Equal(t, app.NewDateRange(startDate, endDate), snapshot.Period)
	assert.Equal(t, bc.UUID(), snapshot.BillingCycleUUID)

	// DocumentName format.
	assert.NotEmpty(t, snapshot.DocumentName)
	assert.True(t,
		strings.HasPrefix(snapshot.DocumentName, "Courier Statement "),
		"document name: %s", snapshot.DocumentName,
	)
	assert.Contains(t, snapshot.DocumentName, fmt.Sprintf("#%d", courierCycleNumber))

	// Each delivery resolves its restaurant name through the repository.
	require.Len(t, snapshot.Deliveries, 2)
	assert.Equal(t, burgerBarnName, snapshot.Deliveries[0].RestaurantName)
	assert.Equal(t, sushiSpotName, snapshot.Deliveries[1].RestaurantName)
}

func TestStatementService_GenerateCourierStatementSnapshot_Errors(t *testing.T) {
	t.Run("not_closed_cycle", func(t *testing.T) {
		service := app.NewStatementService(nil, &fakeLegalEntityRepository{})

		bc := domain.UnmarshalBillingCycle(
			domain.BillingCycleUUID{UUID: common.NewUUIDv7()},
			domain.LegalEntityUUID{UUID: common.NewUUIDv7()},
			domain.PartnerTypeCourier,
			1,
			false, // not closed
			false,
			time.Now().UTC(),
			nil,
		)

		_, err := service.GenerateCourierStatementSnapshot(context.Background(), bc, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not closed")
	})

	t.Run("not_a_courier", func(t *testing.T) {
		service := app.NewStatementService(nil, &fakeLegalEntityRepository{})

		endDate := time.Now().UTC()
		bc := domain.UnmarshalBillingCycle(
			domain.BillingCycleUUID{UUID: common.NewUUIDv7()},
			domain.LegalEntityUUID{UUID: common.NewUUIDv7()},
			domain.PartnerTypeRestaurant, // wrong partner type
			1,
			true,
			false,
			endDate.AddDate(0, 0, -7),
			&endDate,
		)

		_, err := service.GenerateCourierStatementSnapshot(context.Background(), bc, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a courier")
	})
}
