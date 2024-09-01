package app

import "github.com/danielhoward-me/sso-v2/backend/internal/db"

func New(addr string) *App {
	return &App{
		addr: addr,
	}
}

type App struct {
	addr string
}

func (app *App) Start() {
	db.Connect()
	createAlertSender()

	r := createGinEngine()
	r.Run(app.addr)
}
