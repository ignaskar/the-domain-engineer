package tax

import (
	"context"
	"eats/backend/billing/domain"
	"eats/backend/common"
	"eats/backend/common/shared"
	"fmt"

	"github.com/ThreeDotsLabs/the-domain-engineer/clients"
	"github.com/ThreeDotsLabs/the-domain-engineer/clients/tax"
)

type Client struct {
	clients *clients.Clients
}

func NewClient(clients *clients.Clients) *Client {
	if clients == nil {
		panic("nil clients")
	}
	return &Client{clients: clients}
}

func (c *Client) GetTaxRate(ctx context.Context, input domain.TaxRateRequest) (domain.TaxRate, error) {
	taxClass, err := toTaxClass(input.LineItemType)
	if err != nil {
		return domain.TaxRate{}, err
	}

	var taxID *string
	if input.BuyerTaxID != nil {
		taxID = common.ToPtr(input.BuyerTaxID.String())
	}

	req := tax.GetTaxRateJSONRequestBody{
		BuyerCountryCode:  input.BuyerCountryCode.Code(),
		BuyerTaxId:        taxID,
		SellerCountryCode: input.SellerCountryCode.Code(),
		TaxClass:          &taxClass,
		TransactionDate:   input.TransactionDate,
	}
	resp, err := c.clients.Tax.GetTaxRateWithResponse(ctx, req)
	if err != nil {
		return domain.TaxRate{}, fmt.Errorf("failed to get tax rate: %w", err)
	}

	taxType, err := fromTaxClass(resp.JSON200.TaxType)
	if err != nil {
		return domain.TaxRate{}, fmt.Errorf("failed to parse tax type: %w", err)
	}

	taxRate, err := domain.NewTaxRate(resp.JSON200.Rate, taxType)
	if err != nil {
		return domain.TaxRate{}, fmt.Errorf("failed to build tax rate: %w", err)
	}

	return taxRate, nil
}

func toTaxClass(lit shared.LineItemType) (tax.TaxRateRequestTaxClass, error) {
	switch lit {
	case shared.LineItemTypeBeverage:
		return tax.TaxRateRequestTaxClassBEVERAGE, nil
	case shared.LineItemTypeDelivery:
		return tax.TaxRateRequestTaxClassDELIVERY, nil
	case shared.LineItemTypeFood:
		return tax.TaxRateRequestTaxClassFOOD, nil
	case shared.LineItemTypeService:
		return tax.TaxRateRequestTaxClassSERVICE, nil
	default:
		return "", fmt.Errorf("unknown line item type: %v", lit)
	}
}

func fromTaxClass(tc tax.TaxRateResponseTaxType) (domain.TaxType, error) {
	switch tc {
	case tax.GST:
		return domain.TaxTypeGST, nil
	case tax.VAT:
		return domain.TaxTypeVAT, nil
	case tax.SALES:
		return domain.TaxTypeSalesTax, nil
	default:
		return domain.TaxType{}, fmt.Errorf("unknown tax rate response tax type: %v", tc)
	}
}
