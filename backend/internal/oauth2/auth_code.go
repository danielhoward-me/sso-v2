package oauth2

import (
	"time"

	"github.com/danielhoward-me/sso-v2/backend/internal/db/dbo"
	"github.com/danielhoward-me/sso-v2/backend/internal/db/schema/model"
	"github.com/danielhoward-me/sso-v2/backend/internal/db/schema/table"
	"github.com/danielhoward-me/sso-v2/backend/internal/user"
)

type AuthCode struct {
	dbo.DBO[model.AuthCodes]
	id          int32
	code        string
	client      *Client
	user        *user.User
	redirectUri string
	created     time.Time
	expires     time.Time
}

func makeAuthCode(rawAuthCode model.AuthCodes) (authCode *AuthCode, id int32, err error) {
	client, err := NewClient(rawAuthCode.ClientID)
	if err != nil {
		return
	}

	user, err := user.New(rawAuthCode.UserID)
	if err != nil {
		return
	}

	return &AuthCode{
		id:          rawAuthCode.ID,
		code:        rawAuthCode.Code,
		client:      client,
		user:        user,
		redirectUri: rawAuthCode.RedirectURI,
		created:     rawAuthCode.Created,
		expires:     rawAuthCode.Expires,
	}, rawAuthCode.ID, nil
}

var authCodeDBOHandler = dbo.NewHandler(dbo.DBOHandlerOptions[model.AuthCodes, model.AuthCodes, *AuthCode]{
	DBOMaker: makeAuthCode,
	Table:    table.AuthCodes,
	IDColumn: table.AuthCodes.ID,
	Columns:  dbo.SelectColumnList{table.AuthCodes.AllColumns},
})

var NewAuthCode = authCodeDBOHandler.New
