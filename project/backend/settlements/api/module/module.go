package module

import (
	"context"
	"eats/backend/settlements/api/module/client"
	"eats/backend/settlements/app/command"
	"eats/backend/settlements/app/models"
	"eats/backend/settlements/domain"
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
	return s.commandHandler.StartSettlement(ctx, cmd)
}

func (s Settlements) GetPlatformEntity(ctx context.Context, req client.GetPlatformEntityRequest) (client.GetPlatformEntityResponse, error) {
	p, err := s.legalEntityRepository.PartnerByUUID(ctx, domain.LegalEntityUUID{req.PartnerUUID})
	if err != nil {
		return client.GetPlatformEntityResponse{}, err
	}

	return client.GetPlatformEntityResponse{
		PlatformUUID: p.PlatformEntityUUID.UUID,
	}, nil
}
