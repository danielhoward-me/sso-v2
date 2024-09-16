package oauth2

import (
	"github.com/danielhoward-me/sso-v2/backend/internal/db"
	"github.com/danielhoward-me/sso-v2/backend/internal/db/dbo"
	"github.com/danielhoward-me/sso-v2/backend/internal/db/schema/model"
	"github.com/danielhoward-me/sso-v2/backend/internal/db/schema/table"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
)

type Client struct {
	dbo.DBO[model.Clients]
	id                     int32
	uuid                   uuid.UUID
	name                   string
	secret                 string
	showConfirmationPrompt bool
	redirects              []string
}

type rawClientData struct {
	model.Clients
	ClientRedirects []model.ClientRedirects
}

func makeClient(rawClient rawClientData) (*Client, int32, error) {
	redirects := []string{}
	for _, redirect := range rawClient.ClientRedirects {
		redirects = append(redirects, redirect.Redirect)
	}

	return &Client{
		id:                     rawClient.ID,
		uuid:                   rawClient.UUID,
		name:                   rawClient.Name,
		secret:                 rawClient.Secret,
		showConfirmationPrompt: rawClient.ShowConfirmationPrompt,
		redirects:              redirects,
	}, rawClient.ID, nil
}

var clientDBOHandler = dbo.NewHandler(dbo.DBOHandlerOptions[model.Clients, rawClientData, *Client]{
	DBOMaker: makeClient,
	Table:    table.Clients,
	FetchTable: table.Clients.LEFT_JOIN(
		table.ClientRedirects,
		table.Clients.ID.EQ(table.ClientRedirects.ClientID),
	),
	IDColumn:   table.Clients.ID,
	UUIDColumn: table.Clients.UUID,
	Columns:    dbo.SelectColumnList{table.Clients.AllColumns, table.ClientRedirects.Redirect},
})

var NewClient = clientDBOHandler.New
var NewClientFromUUID = clientDBOHandler.NewFromUUID

func (client Client) GetId() uuid.UUID {
	return client.uuid
}

func (client Client) GetName() string {
	return client.name
}

func (client Client) GetShowConfirmationPrompt() bool {
	return client.showConfirmationPrompt
}

func (client Client) GetRedirects() []string {
	return client.redirects
}

func (client Client) CheckSecret(secret string) bool {
	matches, err := argon2.VerifyEncoded([]byte(secret), []byte(client.secret))
	if err != nil {
		// Errors occur when the secrect hasn't been encoded properly so
		// authentication should just fail
		return false
	}
	return matches
}

func (client *Client) UpdateName(name string) (err error) {
	if err = client.Update(
		dbo.UpdateColumnList{table.Clients.Name},
		model.Clients{Name: name},
	); err != nil {
		return
	}

	client.name = name
	return
}

func (client *Client) UpdateShowConfirmationPrompt(showConfirmationPrompt bool) (err error) {
	if err = client.Update(
		dbo.UpdateColumnList{table.Clients.ShowConfirmationPrompt},
		model.Clients{ShowConfirmationPrompt: showConfirmationPrompt},
	); err != nil {
		return
	}

	client.showConfirmationPrompt = showConfirmationPrompt
	return
}

func (client *Client) UpdateRedirects(redirects []string) (err error) {
	_, err = table.ClientRedirects.DELETE().WHERE(table.ClientRedirects.ClientID.EQ(postgres.Int32(client.id))).Exec(db.DB)
	if err != nil {
		return
	}

	client.redirects = []string{}

	if len(redirects) == 0 {
		return
	}

	insertStatement := table.ClientRedirects.INSERT(table.ClientRedirects.AllColumns)
	for _, redirect := range redirects {
		insertStatement = insertStatement.MODEL(model.ClientRedirects{
			ClientID: client.id,
			Redirect: redirect,
		})
	}
	if _, err = insertStatement.Exec(db.DB); err != nil {
		return
	}

	client.redirects = redirects
	return
}

func GetAllClients() (clients []*Client, err error) {
	var clientIds []struct {
		ID int32 `json:"id"`
	}
	if err = postgres.SELECT(table.Clients.ID.AS("id")).FROM(table.Clients).Query(db.DB, &clientIds); err != nil {
		return
	}

	for _, client := range clientIds {
		client, _ := NewClient(client.ID)
		clients = append(clients, client)
	}

	return
}
