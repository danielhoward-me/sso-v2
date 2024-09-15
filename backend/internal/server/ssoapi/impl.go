package ssoapi

import (
	"context"

	"github.com/danielhoward-me/sso-v2/backend/internal/oauth2"
)

type Api struct{}

var _ StrictServerInterface = (*Api)(nil)

func New() *Api {
	return &Api{}
}

func (Api) GetAdminClients(ctx context.Context, request GetAdminClientsRequestObject) (resp GetAdminClientsResponseObject, err error) {
	clients, err := oauth2.GetAllClients()
	if err != nil {
		return
	}

	clientsData := []Client{}
	for _, client := range clients {
		clientsData = append(clientsData, Client{
			Id:                     client.GetId(),
			Name:                   client.GetName(),
			ShowConfirmationPrompt: client.GetShowConfirmationPrompt(),
			Redirects:              client.GetRedirects(),
		})
	}

	return GetAdminClients200JSONResponse(clientsData), nil
}
