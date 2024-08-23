package main

import (
	"fmt"

	"github.com/danielhoward-me/sso-v2/backend/db"
	"github.com/danielhoward-me/sso-v2/backend/utils"
)

var PORT = utils.GetEnv("PORT", "3001")

func main() {
	db.Connect()

	r := createGinEngine()
	r.Run(fmt.Sprintf(":%s", PORT))
}
