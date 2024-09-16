package main

import (
	"fmt"
	"net"

	"github.com/danielhoward-me/sso-v2/backend/internal/alert"
	"github.com/danielhoward-me/sso-v2/backend/internal/db"
	"github.com/danielhoward-me/sso-v2/backend/internal/server"
	"github.com/danielhoward-me/sso-v2/backend/internal/utils"
)

func main() {
	pgUser := utils.RequireEnv("PGUSER")
	pgPassword := utils.RequireEnv("PGPASSWORD")
	pgHost := utils.RequireEnv("PGHOST")
	pgPort := utils.GetEnv("PGPORT", "5432")
	pgDatabase := utils.RequireEnv("PGDATABASE")
	pgSslMode := utils.GetEnv("PGSSLMODE", "disable")

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		pgUser,
		pgPassword,
		pgHost,
		pgPort,
		pgDatabase,
		pgSslMode,
	)

	fmt.Printf("Connecting to database %s\n", pgDatabase)
	if err := db.Connect(connectionString); err != nil {
		panic(fmt.Errorf("failed to connect to database: %s", err))
	}
	fmt.Println("Connected to database")

	apiOptions := server.Options{
		AccountHostname:  utils.RequireEnv("ACCOUNT_HOSTNAME"),
		InternalHostname: utils.RequireEnv("INTERNAL_HOSTNAME"),
		SsoHostname:      utils.RequireEnv("SSO_HOSTNAME"),
	}

	alerter, exists := getAlerter()
	if exists {
		apiOptions.Alerter = alerter
	}

	fmt.Println("Creating API handler")
	api := server.New(apiOptions)
	fmt.Println("Successfully created API handler")

	port := utils.GetEnv("PORT", "3001")
	fmt.Printf("Server avaliable on port %s\n", port)
	api.Start(net.JoinHostPort("0.0.0.0", port))
}

func getAlerter() (*alert.Alerter, bool) {
	alertUrl := utils.GetEnv("ALERT_URL")
	if alertUrl == "" {
		fmt.Println("ALERT_URL env variable isn't set, not alerting on server errors")
		return nil, false
	}

	fmt.Println("Creating alerter")
	alerter, err := alert.New(alertUrl)
	if err != nil {
		panic(fmt.Errorf("failed to create alerter: %s", err))
	}
	fmt.Println("Successfully created alerter")

	return alerter, true
}
