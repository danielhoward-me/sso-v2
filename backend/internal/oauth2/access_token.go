package oauth2

import (
	"fmt"
	"time"

	"github.com/danielhoward-me/sso-v2/backend/internal/db/dbo"
	"github.com/danielhoward-me/sso-v2/backend/internal/db/schema/model"
	"github.com/danielhoward-me/sso-v2/backend/internal/db/schema/table"
)

type AccessToken struct {
	dbo.DBO[model.AccessTokens]
	id           int32
	token        string
	refreshToken *RefreshToken
	created      time.Time
	expires      time.Time
	lastUsed     *time.Time
}

func makeAccessToken(rawAccessToken model.AccessTokens) (*AccessToken, int32) {
	refreshToken, err := NewRefreshToken(rawAccessToken.RefreshTokenID)
	if err != nil {
		panic(fmt.Errorf("failed to create client object when making auth code object: %s", err))
	}

	return &AccessToken{
		id:           rawAccessToken.ID,
		token:        rawAccessToken.Token,
		refreshToken: refreshToken,
		created:      rawAccessToken.Created,
		expires:      rawAccessToken.Expires,
		lastUsed:     rawAccessToken.LastUsed,
	}, rawAccessToken.ID
}

var AccessTokenDBOHandler = dbo.NewHandler(dbo.DBOHandlerOptions[model.AccessTokens, model.AccessTokens, *AccessToken]{
	DBOMaker: makeAccessToken,
	Table:    table.AccessTokens,
	IDColumn: table.AccessTokens.ID,
	Columns:  dbo.SelectColumnList{table.AccessTokens.AllColumns},
})

var NewAccessToken = AccessTokenDBOHandler.New
