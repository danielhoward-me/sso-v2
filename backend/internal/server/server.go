package server

import (
	"fmt"
	"net/http"

	"github.com/danielhoward-me/sso-v2/backend/internal/alert"
	"github.com/danielhoward-me/sso-v2/backend/internal/server/accountapi"
	"github.com/danielhoward-me/sso-v2/backend/internal/server/ssoapi"
	"github.com/danielhoward-me/sso-v2/backend/internal/server/ssointernalapi"
	"github.com/danielhoward-me/sso-v2/backend/internal/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/hostrouter"
	oapiMiddleware "github.com/oapi-codegen/nethttp-middleware"
)

type Options struct {
	Alerter             *alert.Alerter
	AccountHostname     string
	SsoHostname         string
	SsoInternalHostname string
}

type Api struct {
	router chi.Router
	opts   Options
}

func New(opts Options) *Api {
	api := &Api{opts: opts}
	api.initialiseRouter()
	return api
}

func (api Api) Start(addr string) {
	api.alertIfExists("sso backend server has started")
	http.ListenAndServe(addr, api.router)
}

func (api Api) alertIfExists(message string) {
	if api.opts.Alerter != nil {
		api.opts.Alerter.Alert(message)
	}
}

func (api *Api) initialiseRouter() {
	r := chi.NewRouter()

	if utils.IsDevelopment() {
		r.Use(middleware.Logger)
	}

	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Mount("/api", api.getHostRouter())

	api.router = r
}

func (api Api) getHostRouter() hostrouter.Routes {
	hr := hostrouter.New()

	hr.Map(api.opts.AccountHostname, getAccountHandler())
	hr.Map(api.opts.SsoHostname, getSsoHandler())
	hr.Map(api.opts.SsoInternalHostname, getSsoInternalHandler())

	return hr
}

func getAccountHandler() chi.Router {
	swagger := getSwagger("account", accountapi.GetSwagger)
	api := accountapi.New()
	apiHandler := accountapi.NewStrictHandler(api, nil)

	r := chi.NewRouter()
	r.Use(oapiMiddleware.OapiRequestValidatorWithOptions(swagger, &oapiMiddleware.Options{
		SilenceServersWarning: true,
	}))
	accountapi.HandlerFromMux(apiHandler, r)

	return r
}

func getSsoHandler() chi.Router {
	swagger := getSwagger("sso", ssoapi.GetSwagger)
	api := ssoapi.New()
	apiHandler := ssoapi.NewStrictHandler(api, nil)

	r := chi.NewRouter()
	r.Use(oapiMiddleware.OapiRequestValidatorWithOptions(swagger, &oapiMiddleware.Options{
		SilenceServersWarning: true,
	}))
	ssoapi.HandlerFromMux(apiHandler, r)

	return r
}

func getSsoInternalHandler() chi.Router {
	swagger := getSwagger("ssointernal", ssointernalapi.GetSwagger)
	api := ssointernalapi.New()
	apiHandler := ssointernalapi.NewStrictHandler(api, nil)

	r := chi.NewRouter()
	r.Use(oapiMiddleware.OapiRequestValidatorWithOptions(swagger, &oapiMiddleware.Options{
		SilenceServersWarning: true,
	}))
	ssointernalapi.HandlerFromMux(apiHandler, r)

	return r
}

func getSwagger(name string, swaggerFunc func() (*openapi3.T, error)) *openapi3.T {
	swagger, err := swaggerFunc()
	if err != nil {
		panic(fmt.Errorf("failed to load %s swagger: %s", name, err))
	}

	// Add /api as a base path to all servers
	server := openapi3.Server{URL: "/api"}
	swagger.Servers = openapi3.Servers([]*openapi3.Server{&server})

	return swagger
}
