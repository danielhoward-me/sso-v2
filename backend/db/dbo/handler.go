package dbo

import (
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/danielhoward-me/sso-v2/backend/db"
)

type DBOHandler[
	PC string | uuid.UUID, // Primary Column
	TM interface{}, // Table Model
	FM interface{}, // Fetch Model
	DS DBOI[TM], // DBO Struct
] struct {
	dboStructMaker func(FM) DS
	table          postgres.WritableTable
	primaryColumn  postgres.ColumnString
	baseQuery      postgres.SelectStatement
}

func NewHandler[
	PC string | uuid.UUID,
	TM interface{},
	FM interface{},
	DS DBOI[TM],
](
	dboStructMaker func(FM) DS,
	table postgres.WritableTable,
	fetchTable postgres.ReadableTable,
	primaryColumn postgres.ColumnString,
	columns []postgres.Projection,
) (dboHandler DBOHandler[PC, TM, FM, DS]) {
	if len(columns) == 0 {
		panic(fmt.Errorf("at least 1 column must be selected to be fetched"))
	}

	dboHandler.dboStructMaker = dboStructMaker
	dboHandler.table = table
	dboHandler.primaryColumn = primaryColumn
	dboHandler.baseQuery = postgres.SELECT(columns[0], columns[1:]...).FROM(fetchTable)

	return
}

func (handler DBOHandler[PC, _, FM, DS]) New(id PC) (dbo DS, err error) {
	var rawData FM

	var idValue postgres.StringExpression
	var value string
	switch any(id).(type) {
	case string:
		stringId := any(id).(string)
		idValue = postgres.String(stringId)
		value = stringId
	case uuid.UUID:
		uuidValue := any(id).(uuid.UUID)
		idValue = postgres.UUID(uuidValue)
		value = uuidValue.String()
	}

	primaryKeyQuery := handler.primaryColumn.EQ(idValue)

	err = handler.baseQuery.WHERE(primaryKeyQuery).Query(db.DB, &rawData)

	if err != nil {
		err = fmt.Errorf("error occured when entry from %s table with column %s set to '%s': %s", handler.primaryColumn.TableName(), handler.primaryColumn.Name(), value, err)
		return
	}

	dbo = handler.dboStructMaker(rawData)

	dbo.setDBOOptions(
		handler.table,
		primaryKeyQuery,
	)

	return dbo, nil
}
