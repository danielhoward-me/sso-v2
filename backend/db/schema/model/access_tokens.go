//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type AccessTokens struct {
	Token    string `sql:"primary_key"`
	ClientID uuid.UUID
	UserID   uuid.UUID
	Created  time.Time
	LastUsed *time.Time
}
