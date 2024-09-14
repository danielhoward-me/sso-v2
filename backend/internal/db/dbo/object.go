package dbo

import (
	"fmt"

	"github.com/danielhoward-me/sso-v2/backend/internal/db"
	"github.com/go-jet/jet/v2/postgres"
)

type UpdateColumnList = postgres.ColumnList

type DBOI[TM interface{}] interface {
	setDBOOptions(postgres.WritableTable, postgres.BoolExpression)

	Update(columns UpdateColumnList, data TM) error
}

type DBO[TM interface{}] struct {
	table           postgres.WritableTable
	primaryKeyQuery postgres.BoolExpression
}

// lint:ignore U1000 This function is used in handler.go but in a way that the static check can't pick up
func (dbo *DBO[_]) setDBOOptions(table postgres.WritableTable, primaryKeyQuery postgres.BoolExpression) {
	dbo.table = table
	dbo.primaryKeyQuery = primaryKeyQuery
}

func (dbo DBO[TM]) Update(columns UpdateColumnList, data TM) error {
	_, err := dbo.table.UPDATE(columns).
		MODEL(data).
		WHERE(dbo.primaryKeyQuery).
		Exec(db.DB)

	if err != nil {
		return fmt.Errorf("error occured when updating dbo: %s", err)
	}

	return nil
}
