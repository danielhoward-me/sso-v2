package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/danielhoward-me/sso/backend/oauth2"
)

func createGinEngine() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.CustomRecovery(recoveryHandler))
	engine.SetTrustedProxies(nil)

	if gin.IsDebugging() {
		engine.Use(gin.Logger())
	}

	// The nginx proxy only forwards requests if they start with "/api", otherwise
	// they are sent to the frontend
	apiRouter := engine.Group("/api")
	setupApiRouter(apiRouter)

	engine.NoRoute(noRouteHandler)

	return engine
}

func setupApiRouter(router *gin.RouterGroup) {
	router.GET("/admin/clients", func(c *gin.Context) {
		clients := oauth2.GetAllClients()

		out := []map[string]any{}
		for _, client := range clients {
			out = append(out, client.ToMap())
		}

		c.JSON(http.StatusOK, out)
	})

	router.GET("/admin/clients/:id", func(c *gin.Context) {
		id := c.Param("id")
		name := c.Query("name")

		client, _ := oauth2.NewClient(uuid.MustParse(id))
		client.UpdateName(name)
	})
}

func recoveryHandler(c *gin.Context, _ any) {
	handleError(http.StatusInternalServerError, c)
}

func noRouteHandler(c *gin.Context) {
	handleError(http.StatusNotFound, c)
}
