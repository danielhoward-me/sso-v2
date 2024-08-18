package oauth2

import (
	"errors"
	"fmt"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"

	"github.com/danielhoward-me/sso/backend/db"
	"github.com/danielhoward-me/sso/backend/db/schema/model"
	. "github.com/danielhoward-me/sso/backend/db/schema/table"
)

type Client struct {
	ID                     uuid.UUID
	Name                   string
	Secret                 string
	ShowConfirmationPrompt bool
	Redirects              []string
}

type rawClientData struct {
	model.Clients
	ClientRedirects []model.ClientRedirects
}

func GetAllClients() (clients []Client) {
	var clientIds []struct{ ID uuid.UUID }
	err := SELECT(Clients.ID.AS("id")).FROM(Clients).Query(db.DB, &clientIds)
	if err != nil {
		panic(fmt.Errorf("error occured when fetching all client IDs: %s", err))
	}

	for _, client := range clientIds {
		client, _ := NewClient(client.ID)
		clients = append(clients, client)
	}

	return
}

func NewClient(id uuid.UUID) (client Client, exists bool) {
	var rawClient rawClientData
	err := SELECT(
		Clients.AllColumns,
		ClientRedirects.Redirect,
	).FROM(
		Clients.LEFT_JOIN(
			ClientRedirects,
			Clients.ID.EQ(ClientRedirects.ClientID),
		),
	).WHERE(
		Clients.ID.EQ(UUID(id)),
	).Query(db.DB, &rawClient)

	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			exists = false
			return
		}

		panic(fmt.Errorf("error occured when fetching client with ID '%s': %s", id.String(), err))
	}

	client = processRawClient(rawClient)
	exists = true
	return
}

func processRawClient(rawClient rawClientData) (client Client) {
	client.ID = rawClient.ID
	client.Name = rawClient.Name
	client.Secret = rawClient.Secret
	client.ShowConfirmationPrompt = rawClient.ShowConfirmationPrompt

	client.Redirects = []string{}
	for _, redirect := range rawClient.ClientRedirects {
		client.Redirects = append(client.Redirects, redirect.Redirect)
	}

	return
}

func (client *Client) ToMap() map[string]any {
	return map[string]any{
		"id":                     client.ID,
		"name":                   client.Name,
		"showConfirmationPrompt": client.ShowConfirmationPrompt,
		"redirects":              client.Redirects,
	}
}

func (client *Client) UpdateName(name string) {
	_, err := Clients.UPDATE(Clients.Name).
		MODEL(model.Clients{Name: name}).
		WHERE(Clients.ID.EQ(UUID(client.ID))).
		Exec(db.DB)

	if err != nil {
		panic(fmt.Errorf("error occured when updating the name of client with ID '%s': %s", client.ID.String(), err))
	}

	client.Name = name
}
