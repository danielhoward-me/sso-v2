package ssoapi

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

func (Api) PutAdminClientsId(ctx context.Context, request PutAdminClientsIdRequestObject) (resp PutAdminClientsIdResponseObject, err error) {
	client, err := oauth2.NewClientFromUUID(request.Id)
	if err != nil {
		if utils.ErrIsNoRows(err) {
			return PutAdminClientsId404JSONResponse{
				Message: fmt.Sprintf("The client with ID %s does not exist", request.Id.String()),
				Status:  404,
			}, nil
		}
		return
	}

	if request.Body.Name != nil {
		if err = client.UpdateName(*request.Body.Name); err != nil {
			return
		}
	}

	if request.Body.ShowConfirmationPrompt != nil {
		if err = client.UpdateShowConfirmationPrompt(*request.Body.ShowConfirmationPrompt); err != nil {
			return
		}
	}

	if request.Body.Redirects != nil {
		if err = client.UpdateRedirects(*request.Body.Redirects); err != nil {
			return
		}
	}

	return PutAdminClientsId200JSONResponse{
		Message: "Client updated successfully",
	}, nil
}
