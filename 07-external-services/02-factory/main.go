package main

type PaymentDetails struct {
	Amount   int
	Currency string
}

type PaymentsClient interface {
	GetPaymentDetails(nonce string) (PaymentDetails, error)
}

type Order struct {
	nonce    string
	amount   int
	currency string
}

func (o *Order) Nonce() string    { return o.nonce }
func (o *Order) Amount() int      { return o.amount }
func (o *Order) Currency() string { return o.currency }

type OrderFactory struct {
	paymentsClient PaymentsClient
}

func NewOrderFactory(paymentsClient PaymentsClient) OrderFactory {
	if paymentsClient == nil {
		panic("payments client is required")
	}
	return OrderFactory{
		paymentsClient: paymentsClient,
	}
}
