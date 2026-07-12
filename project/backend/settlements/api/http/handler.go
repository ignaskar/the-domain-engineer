package http

import (
	"context"
	"fmt"

	"eats/backend/settlements/app/command"
)

type Handler struct {
	commandHandler *command.Handlers
}

func NewHandler(commandHandler *command.Handlers) *Handler {
	if commandHandler == nil {
		panic("command handler is required")
	}

	return &Handler{
		commandHandler: commandHandler,
	}
}

// TODO: implement OnboardPartner.
//
// Steps:
// 1. Reject the request with 401 if request.Params.OperatorUUID.IsZero().
// 2. Parse Address, TaxID, IBAN value objects from request.Body. Each must construct successfully.
// 3. Build a *command.OnboardPartner from request.Body and the parsed VOs.
// 4. Dispatch with h.commandHandler.OnboardPartner.
// 5. Return OnboardPartner204Response{} on success.
func (h Handler) OnboardPartner(ctx context.Context, request OnboardPartnerRequestObject) (OnboardPartnerResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// TODO: implement CreatePlatformEntity.
//
// Same shape as OnboardPartner, but build *command.CreatePlatformEntity instead.
// Dispatch via h.commandHandler.CreatePlatformEntity, which returns the PlatformEntityUUID.
// On success, return CreatePlatformEntity201JSONResponse{PlatformEntityUuid: uuid}.
func (h Handler) CreatePlatformEntity(ctx context.Context, request CreatePlatformEntityRequestObject) (CreatePlatformEntityResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

func Register(e EchoRouter, commandHandlers *command.Handlers) error {
	handler := NewHandler(commandHandlers)

	RegisterHandlers(e, NewStrictHandler(handler, nil))

	return nil
}
