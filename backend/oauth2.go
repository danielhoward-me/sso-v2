package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/danielhoward-me/sso-v2/backend/oauth2"
)

type TokenEndpointBodyBase struct {
	GrantType string `json:"grant_type"`
}
type TokenEndpointBodyAuthorisationCode struct {
	TokenEndpointBodyBase
	Code        string `json:"code"`
	RedirectUri string `json:"redirect_uri"`
}

type oauth2ErrorType string

const INVALID_REQUEST oauth2ErrorType = "invalid_request"
const INVALID_CLIENT oauth2ErrorType = "invalid_client"
const INVALID_GRANT oauth2ErrorType = "invalid_grant"
const UNAUTHORISED_CLIENT oauth2ErrorType = "unauthorized_client"
const UNSUPPORTED_GRANT_TYPE oauth2ErrorType = "unsupported_grant_type"

type oauth2Error struct {
	Error            oauth2ErrorType `json:"error"`
	ErrorDescription string          `json:"error_description"`
}

var INVALID_CLIENT_ERROR oauth2Error = oauth2Error{
	Error:            INVALID_CLIENT,
	ErrorDescription: "the specified client id or secret is not correct",
}
var INVALID_GRANT_ERROR oauth2Error = oauth2Error{
	Error:            INVALID_GRANT,
	ErrorDescription: "the specified grant type doesn't exist",
}

func handleTokenEndpoint(c *gin.Context) {
	clientId, clientSecret, exists := c.Request.BasicAuth()
	if !exists {
		c.JSON(401, INVALID_CLIENT_ERROR)
	}

	clientUuid, err := uuid.Parse(clientId)
	if err != nil {
		panic(fmt.Errorf("token endpoint was passed an invalid client id: %s", err))
	}

	client, exists := oauth2.NewClient(clientUuid)
	if !exists {
		c.JSON(401, INVALID_CLIENT_ERROR)
	}

	if !client.CheckSecret(clientSecret) {
		c.JSON(401, INVALID_CLIENT_ERROR)
	}

	body := mustParseBody[TokenEndpointBodyBase](c)

	switch body.GrantType {
	case "authorization_code":
		{
			fullBody := mustParseBody[TokenEndpointBodyAuthorisationCode](c)
		}
	default:
		{
			c.JSON(400, INVALID_GRANT_ERROR)
		}
	}
}

func mustParseBody[T any](c *gin.Context) T {
	var body T
	if err := c.BindJSON(&body); err != nil {
		panic(fmt.Errorf("token endpoint failed to pass request body: %s", err))
	}
	return body
}
