package oauth2

import (
	"fmt"
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

func makeAuthCode(rawAuthCode model.AuthCodes) (*AuthCode, int32) {
	client, err := NewClient(rawAuthCode.ClientID)
	if err != nil {
		panic(fmt.Errorf("failed to create client object when making auth code object: %s", err))
	}

	user, err := user.New(rawAuthCode.UserID)
	if err != nil {
		panic(fmt.Errorf("failed to create user object when making auth code object: %s", err))
	}

	return &AuthCode{
		id:          rawAuthCode.ID,
		code:        rawAuthCode.Code,
		client:      client,
		user:        user,
		redirectUri: rawAuthCode.RedirectURI,
		created:     rawAuthCode.Created,
		expires:     rawAuthCode.Expires,
	}, rawAuthCode.ID
}

var authCodeDBOHandler = dbo.NewHandler(dbo.DBOHandlerOptions[model.AuthCodes, model.AuthCodes, *AuthCode]{
	DBOMaker: makeAuthCode,
	Table:    table.AuthCodes,
	IDColumn: table.AuthCodes.ID,
	Columns:  dbo.SelectColumnList{table.AuthCodes.AllColumns},
})

var NewAuthCode = authCodeDBOHandler.New
