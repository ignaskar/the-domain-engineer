package module

import (
	"context"
	"errors"

	"eats/backend/settlements/api/module/client"
	"eats/backend/settlements/app/command"
	"eats/backend/settlements/app/models"
)

type Settlements struct {
	commandHandler        *command.Handlers
	legalEntityRepository models.LegalEntityRepository
}

func New(commandHandler *command.Handlers, legalEntityRepository models.LegalEntityRepository) *Settlements {
	if commandHandler == nil {
		panic("commandHandler is nil")
	}
	if legalEntityRepository == nil {
		panic("legalEntityRepository is nil")
	}

	return &Settlements{
		commandHandler:        commandHandler,
		legalEntityRepository: legalEntityRepository,
	}
}

func (s Settlements) StartSettlement(ctx context.Context, cmd client.StartSettlementRequest) error {
	return errors.New("not implemented")
}

func (s Settlements) GetPlatformEntity(ctx context.Context, req client.GetPlatformEntityRequest) (client.GetPlatformEntityResponse, error) {
	return client.GetPlatformEntityResponse{}, errors.New("not implemented")
}
