package oauth2

import (
	"time"

	"github.com/danielhoward-me/sso-v2/backend/internal/db/dbo"
	"github.com/danielhoward-me/sso-v2/backend/internal/db/schema/model"
	"github.com/danielhoward-me/sso-v2/backend/internal/db/schema/table"
	"github.com/danielhoward-me/sso-v2/backend/internal/user"
)

type RefreshToken struct {
	dbo.DBO[model.RefreshTokens]
	id      int32
	token   string
	client  *Client
	user    *user.User
	created time.Time
	expires time.Time
}

func makeRefreshToken(rawRefreshToken model.RefreshTokens) (refreshToken *RefreshToken, id int32, err error) {
	client, err := NewClient(rawRefreshToken.ClientID)
	if err != nil {
		return
	}

	user, err := user.New(rawRefreshToken.UserID)
	if err != nil {
		return
	}

	return &RefreshToken{
		id:      rawRefreshToken.ID,
		token:   rawRefreshToken.Token,
		client:  client,
		user:    user,
		created: rawRefreshToken.Created,
		expires: rawRefreshToken.Expires,
	}, rawRefreshToken.ID, nil
}

var RefreshTokenDBOHandler = dbo.NewHandler(dbo.DBOHandlerOptions[model.RefreshTokens, model.RefreshTokens, *RefreshToken]{
	DBOMaker: makeRefreshToken,
	Table:    table.RefreshTokens,
	IDColumn: table.RefreshTokens.ID,
	Columns:  dbo.SelectColumnList{table.RefreshTokens.AllColumns},
})

var NewRefreshToken = RefreshTokenDBOHandler.New
