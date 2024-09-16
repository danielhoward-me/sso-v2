package internalapi

import (
	"context"
	"fmt"

	"github.com/danielhoward-me/sso-v2/backend/internal/oauth2"
	"github.com/danielhoward-me/sso-v2/backend/internal/utils"
)

type Api struct{}

var _ StrictServerInterface = (*Api)(nil)

func New() *Api {
	return &Api{}
}

func (Api) GetClientsId(ctx context.Context, request GetClientsIdRequestObject) (resp GetClientsIdResponseObject, err error) {
	client, err := oauth2.NewClientFromUUID(request.Id)
	if err != nil {
		if utils.ErrIsNoRows(err) {
			return GetClientsId404JSONResponse{
				Message: fmt.Sprintf("The client with ID %s does not exist", request.Id.String()),
				Status:  404,
			}, nil
		}
		return
	}

	return GetClientsId200JSONResponse(Client{
		Id:                     client.GetId(),
		Name:                   client.GetName(),
		ShowConfirmationPrompt: client.GetShowConfirmationPrompt(),
		Redirects:              client.GetRedirects(),
	}), nil
}
