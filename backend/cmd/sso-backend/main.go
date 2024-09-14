package main

import (
	"fmt"

	"github.com/danielhoward-me/sso-v2/backend/app"
	"github.com/danielhoward-me/sso-v2/backend/internal/utils"
)

var PORT = utils.GetEnv("PORT", "3001")

func main() {
	server := app.New(fmt.Sprintf(":%s", PORT))
	server.Start()
}
