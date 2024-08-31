package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/danielhoward-me/sso-v2/backend/internal/db/dbo"
	"github.com/danielhoward-me/sso-v2/backend/internal/db/schema/model"
	"github.com/danielhoward-me/sso-v2/backend/internal/db/schema/table"
)

type ProfilePicture = model.ProfilePicture
type User struct {
	dbo.DBO[model.Users]
	id             int32
	uuid           uuid.UUID
	username       string
	password       string
	email          string
	profilePicture ProfilePicture
	created        time.Time
	lastUpdated    time.Time
}

func makeUser(rawUser model.Users) (*User, int32) {
	return &User{
		id:             rawUser.ID,
		uuid:           rawUser.UUID,
		username:       rawUser.Username,
		password:       rawUser.Password,
		email:          rawUser.Email,
		profilePicture: rawUser.ProfilePicture,
		created:        rawUser.Created,
		lastUpdated:    rawUser.LastUpdated,
	}, rawUser.ID
}

var dboHandler = dbo.NewHandler(dbo.DBOHandlerOptions[model.Users, model.Users, *User]{
	DBOMaker:   makeUser,
	Table:      table.Users,
	IDColumn:   table.Users.ID,
	UUIDColumn: table.Users.UUID,
	Columns:    dbo.SelectColumnList{table.Users.AllColumns},
})

var New = dboHandler.New
var NewFromUUID = dboHandler.NewFromUUID
