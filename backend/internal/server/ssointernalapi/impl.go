package ssointernalapi

import (
	"context"

	"github.com/danielhoward-me/sso-v2/backend/internal/oauth2"
)

type Api struct{}

var _ StrictServerInterface = (*Api)(nil)

func New() *Api {
	return &Api{}
}

func (Api) GetClientId(ctx context.Context, request GetClientIdRequestObject) (resp GetClientIdResponseObject, err error) {
	client, err := oauth2.NewClientFromUUID(request.Id)
	if err != nil {
		return
	}

	return GetClientId200JSONResponse(Client{
		Id:                     client.GetId(),
		Name:                   client.GetName(),
		ShowConfirmationPrompt: client.GetShowConfirmationPrompt(),
		Redirects:              client.GetRedirects(),
	}), nil
}
