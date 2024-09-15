package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(connectionString string) (err error) {
	if DB != nil {
		return
	}

	connection, err := sql.Open("postgres", connectionString)
	if err != nil {
		return
	}

	if err = connection.Ping(); err != nil {
		return
	}

	DB = connection
	return
}
