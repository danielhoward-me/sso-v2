package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/danielhoward-me/sso/backend/utils"
)

var DB *sql.DB

var PGUSER = utils.GetEnv("PGUSER")
var PGPASSWORD = utils.GetEnv("PGPASSWORD")
var PGHOST = utils.GetEnv("PGHOST")
var PGPORT = utils.GetEnv("PGPORT", "5432")
var PGDATABASE = utils.GetEnv("PGDATABASE")
var PSSSLMODE = utils.GetEnv("PSSSLMODE", "disable")

var CONNECTION_STRING = fmt.Sprintf(
	"postgres://%s:%s@%s:%s/%s?sslmode=%s",
	PGUSER,
	PGPASSWORD,
	PGHOST,
	PGPORT,
	PGDATABASE,
	PSSSLMODE,
)

func Connect() {
	if DB != nil {
		return
	}

	fmt.Printf("Connecting to database %s\n", PGDATABASE)

	connection, err := sql.Open("postgres", CONNECTION_STRING)
	if err != nil {
		panic(fmt.Errorf("failed to create database connection: %s", err))
	}

	if err = connection.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping database: %s", err))
	}

	fmt.Printf("Connected to database %s\n", PGDATABASE)

	DB = connection
}
