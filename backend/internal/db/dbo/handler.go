package dbo

import (
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/danielhoward-me/sso-v2/backend/internal/db"
)

type SelectColumnList = []postgres.Projection

type DBOHandler[
	TM interface{}, // Table Model
	FM interface{}, // Fetch Model
	DS DBOI[TM], // DBO Struct
] struct {
	dboMaker     func(FM) (DS, int32)
	table        postgres.Table
	idColumn     postgres.ColumnInteger
	uuidColumn   postgres.ColumnString
	stringColumn postgres.ColumnString
	baseQuery    postgres.SelectStatement
}

type DBOHandlerOptions[
	TM interface{}, // Table Model
	FM interface{}, // Fetch Model
	DS DBOI[TM], // DBO Struct
] struct {
	// The function called when an object instance is created, it should return a DBO object with
	// the correct fields populated and the id of the entry which can be used in further queries
	DBOMaker func(FM) (DS, int32)
	// The table that any updated data should be written to
	Table postgres.Table
	// The table that the data is fetched from, defaults to the value in Table if not set
	FetchTable postgres.ReadableTable
	// The primary column in the database
	IDColumn postgres.ColumnInteger
	// The unique column in the table that stores the uuid for of that row
	UUIDColumn postgres.ColumnString
	// The unique column in the table that stores an identifier values of that row
	StringColumn postgres.ColumnString
	// The columns that need to be fetched from the database
	Columns SelectColumnList
}

func NewHandler[
	TM interface{}, // Table Model
	FM interface{}, // Fetch Model
	DS DBOI[TM], // DBO Struct
](options DBOHandlerOptions[TM, FM, DS]) (dboHandler DBOHandler[TM, FM, DS]) {
	if len(options.Columns) == 0 {
		panic(fmt.Errorf("at least 1 column must be selected to be fetched"))
	}

	if options.FetchTable == nil {
		options.FetchTable = options.Table
	}
	baseQuery := postgres.SELECT(options.Columns[0], options.Columns[1:]...).FROM(options.FetchTable)

	dboHandler.dboMaker = options.DBOMaker
	dboHandler.table = options.Table
	dboHandler.idColumn = options.IDColumn
	dboHandler.uuidColumn = options.UUIDColumn
	dboHandler.stringColumn = options.StringColumn
	dboHandler.baseQuery = baseQuery

	return
}

func (handler DBOHandler[_, FM, DS]) New(id int32) (dbo DS, err error) {
	selectorKeyQuery := handler.idColumn.EQ(postgres.Int32(id))
	return handler.makeDboStruct(selectorKeyQuery)
}

func (handler DBOHandler[_, FM, DS]) NewFromUUID(uuid uuid.UUID) (dbo DS, err error) {
	if handler.uuidColumn == nil {
		panic(fmt.Errorf("uuid column must be set in order to use NewFromUUID"))
	}

	selectorKeyQuery := handler.uuidColumn.EQ(postgres.UUID(uuid))
	return handler.makeDboStruct(selectorKeyQuery)
}

func (handler DBOHandler[_, FM, DS]) NewFromString(value string) (dbo DS, err error) {
	if handler.stringColumn == nil {
		panic(fmt.Errorf("string column must be set in order to use NewFromString"))
	}

	selectorKeyQuery := handler.uuidColumn.EQ(postgres.String(value))
	return handler.makeDboStruct(selectorKeyQuery)
}

func (handler DBOHandler[_, FM, DS]) makeDboStruct(selectorKeyQuery postgres.BoolExpression) (dbo DS, err error) {
	var rawData FM

	err = handler.baseQuery.WHERE(selectorKeyQuery).Query(db.DB, &rawData)

	if err != nil {
		return
	}

	dbo, id := handler.dboMaker(rawData)

	dbo.setDBOOptions(
		handler.table,
		handler.idColumn.EQ(postgres.Int32(id)),
	)

	return dbo, nil
}
