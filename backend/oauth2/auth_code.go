package oauth2

import (
	"time"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/danielhoward-me/sso-v2/backend/db"
	"github.com/danielhoward-me/sso-v2/backend/db/schema/model"
	. "github.com/danielhoward-me/sso-v2/backend/db/schema/table"
	"github.com/danielhoward-me/sso-v2/backend/user"
)

type AuthCode struct {
	code        string
	client      Client
	user        user.User
	redirectUri string
	created     time.Time
	expires     time.Time
}

type rawCodeData = model.AuthCodes

func NewAuthCode(code string) (authCode AuthCode, exists bool, err error) {
	var rawCode rawCodeData
	err = SELECT(
		AuthCodes.AllColumns,
	).FROM(
		AuthCodes,
	).WHERE(
		AuthCodes.Code.EQ(String(code)),
	).Query(db.DB, &rawCode)
	return
}
