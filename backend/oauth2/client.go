package oauth2

import (
	"errors"
	"fmt"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"

	"github.com/danielhoward-me/sso-v2/backend/db"
	"github.com/danielhoward-me/sso-v2/backend/db/schema/model"
	. "github.com/danielhoward-me/sso-v2/backend/db/schema/table"
)

type Client struct {
	id                     uuid.UUID
	name                   string
	secret                 string
	showConfirmationPrompt bool
	redirects              []string
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
	client.id = rawClient.ID
	client.name = rawClient.Name
	client.secret = rawClient.Secret
	client.showConfirmationPrompt = rawClient.ShowConfirmationPrompt

	client.redirects = []string{}
	for _, redirect := range rawClient.ClientRedirects {
		client.redirects = append(client.redirects, redirect.Redirect)
	}

	return
}

func (client *Client) ToMap() map[string]any {
	return map[string]any{
		"id":                     client.id,
		"name":                   client.name,
		"showConfirmationPrompt": client.showConfirmationPrompt,
		"redirects":              client.redirects,
	}
}

func (client *Client) CheckSecret(secret string) bool {
	matches, err := argon2.VerifyEncoded([]byte(secret), []byte(client.secret))
	if err != nil {
		// Errors occur when the secrect hasn't been encoded properly so
		// authentication should just fail
		return false
	}
	return matches
}

func (client *Client) UpdateName(name string) {
	_, err := Clients.UPDATE(Clients.Name).
		MODEL(model.Clients{Name: name}).
		WHERE(Clients.ID.EQ(UUID(client.id))).
		Exec(db.DB)

	if err != nil {
		panic(fmt.Errorf("error occured when updating the name of client with ID '%s': %s", client.id.String(), err))
	}

	client.name = name
}
