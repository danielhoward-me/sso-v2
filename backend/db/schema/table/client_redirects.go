//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var ClientRedirects = newClientRedirectsTable("public", "client_redirects", "")

type clientRedirectsTable struct {
	postgres.Table

	// Columns
	ClientID postgres.ColumnString
	Redirect postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ClientRedirectsTable struct {
	clientRedirectsTable

	EXCLUDED clientRedirectsTable
}

// AS creates new ClientRedirectsTable with assigned alias
func (a ClientRedirectsTable) AS(alias string) *ClientRedirectsTable {
	return newClientRedirectsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ClientRedirectsTable with assigned schema name
func (a ClientRedirectsTable) FromSchema(schemaName string) *ClientRedirectsTable {
	return newClientRedirectsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ClientRedirectsTable with assigned table prefix
func (a ClientRedirectsTable) WithPrefix(prefix string) *ClientRedirectsTable {
	return newClientRedirectsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ClientRedirectsTable with assigned table suffix
func (a ClientRedirectsTable) WithSuffix(suffix string) *ClientRedirectsTable {
	return newClientRedirectsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newClientRedirectsTable(schemaName, tableName, alias string) *ClientRedirectsTable {
	return &ClientRedirectsTable{
		clientRedirectsTable: newClientRedirectsTableImpl(schemaName, tableName, alias),
		EXCLUDED:             newClientRedirectsTableImpl("", "excluded", ""),
	}
}

func newClientRedirectsTableImpl(schemaName, tableName, alias string) clientRedirectsTable {
	var (
		ClientIDColumn = postgres.StringColumn("client_id")
		RedirectColumn = postgres.StringColumn("redirect")
		allColumns     = postgres.ColumnList{ClientIDColumn, RedirectColumn}
		mutableColumns = postgres.ColumnList{}
	)

	return clientRedirectsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ClientID: ClientIDColumn,
		Redirect: RedirectColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
