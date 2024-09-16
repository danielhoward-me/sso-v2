package server

import (
	"fmt"
	"net/http"

	"github.com/danielhoward-me/sso-v2/backend/internal/alert"
	"github.com/danielhoward-me/sso-v2/backend/internal/server/accountapi"
	"github.com/danielhoward-me/sso-v2/backend/internal/server/internalapi"
	"github.com/danielhoward-me/sso-v2/backend/internal/server/ssoapi"
	"github.com/danielhoward-me/sso-v2/backend/internal/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/hostrouter"
	"github.com/google/uuid"
	oapiMiddleware "github.com/oapi-codegen/nethttp-middleware"
)

type Options struct {
	Alerter          *alert.Alerter
	AccountHostname  string
	InternalHostname string
	SsoHostname      string
}

type Api struct {
	router chi.Router
	opts   Options
}

func New(opts Options) *Api {
	initialiseOapiFormats()

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
	// r.Use(middleware.Recoverer)

	r.Mount("/api", api.getHostRouter())

	api.router = r
}

func (api Api) getHostRouter() hostrouter.Routes {
	hr := hostrouter.New()

	hr.Map(api.opts.AccountHostname, getAccountHandler())
	hr.Map(api.opts.InternalHostname, getInternalHandler())
	hr.Map(api.opts.SsoHostname, getSsoHandler())

	return hr
}

func getAccountHandler() chi.Router {
	swagger := getSwagger("account", accountapi.GetSwagger)
	api := accountapi.New()
	apiHandler := accountapi.NewStrictHandler(api, nil)

	r := chi.NewRouter()
	r.Use(oapiMiddleware.OapiRequestValidatorWithOptions(swagger, getMiddlewareOptions()))
	accountapi.HandlerFromMux(apiHandler, r)

	return r
}

func getInternalHandler() chi.Router {
	swagger := getSwagger("internal", internalapi.GetSwagger)
	api := internalapi.New()
	apiHandler := internalapi.NewStrictHandler(api, nil)

	r := chi.NewRouter()
	r.Use(oapiMiddleware.OapiRequestValidatorWithOptions(swagger, getMiddlewareOptions()))
	internalapi.HandlerFromMux(apiHandler, r)

	return r
}

func getSsoHandler() chi.Router {
	swagger := getSwagger("sso", ssoapi.GetSwagger)
	api := ssoapi.New()
	apiHandler := ssoapi.NewStrictHandler(api, nil)

	r := chi.NewRouter()
	r.Use(oapiMiddleware.OapiRequestValidatorWithOptions(swagger, getMiddlewareOptions()))
	ssoapi.HandlerFromMux(apiHandler, r)

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

func getMiddlewareOptions() *oapiMiddleware.Options {
	return &oapiMiddleware.Options{
		SilenceServersWarning: true,
		MultiErrorHandler:     multiErrorHandler,
		ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
			h := w.Header()
			h.Del("Content-Length")
			h.Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(statusCode)

			switch statusCode {
			case http.StatusBadRequest:
				// message already contains the json
			case http.StatusNotFound:
				message = utils.MustMarshal(Error{
					Status:  http.StatusNotFound,
					Message: "This api endpoint does not exist",
				})
			default:
				message = utils.MustMarshal(Error{
					Status:  statusCode,
					Message: "There was an unkown error",
				})
			}
			fmt.Fprint(w, message)
		},
		Options: openapi3filter.Options{
			MultiError: true,
		},
	}
}

type inputError struct {
	errorDetails []ErrorDetail
}

func (inputError inputError) Error() string {
	outputErr := ErrorWithDetails{
		Status:  http.StatusBadRequest,
		Message: "There was atleast one validation error in your request",
		Details: inputError.errorDetails,
	}

	return utils.MustMarshal(outputErr)
}

func multiErrorHandler(errs openapi3.MultiError) (int, error) {
	fmt.Println(len(errs))
	outputErrors := []ErrorDetail{}
	for _, err := range errs {
		var errorDetail ErrorDetail

		switch e := err.(type) {
		case *openapi3filter.RequestError:
			errorDetail = ErrorDetail{
				Message: getRequestErrorMessage(e),
			}

			if e.RequestBody != nil {
				input := "request body"
				errorDetail.Input = &input
			}
			if e.Parameter != nil && e.Parameter.Name != "" {
				errorDetail.Input = &e.Parameter.Name
			}
		default:
			errorDetail = ErrorDetail{
				Message: e.Error(),
			}
		}

		outputErrors = append(outputErrors, errorDetail)
	}
	return http.StatusBadRequest, inputError{
		errorDetails: outputErrors,
	}
}

func getRequestErrorMessage(err *openapi3filter.RequestError) string {
	if err.Reason != "" {
		return err.Reason
	}
	if err.Err.Error() != "" {
		return err.Err.Error()
	}
	return err.Error()
}

func initialiseOapiFormats() {
	openapi3.DefineStringFormatCallback("uuid", func(value string) (err error) {
		_, err = uuid.Parse(value)
		if err != nil {
			fmt.Println(err)
			err = fmt.Errorf("'%s' is not a valid UUID", value)
		}
		return
	})
}
