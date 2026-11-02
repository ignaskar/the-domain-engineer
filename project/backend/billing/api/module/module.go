package module

import (
	"context"
	"errors"

	"eats/backend/billing/api/module/client"
	"eats/backend/billing/app/command"
)

type Billing struct {
	commandHandlers *command.Handlers
}

func New(
	commandHandlers *command.Handlers,
) *Billing {
	if commandHandlers == nil {
		panic("commandHandlers cannot be nil")
	}

	return &Billing{
		commandHandlers: commandHandlers,
	}
}

func (b *Billing) IssueReceipt(ctx context.Context, req client.IssueReceiptRequest) error {
	return errors.New("not implemented")
}
