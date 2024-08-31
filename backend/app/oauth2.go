package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/danielhoward-me/sso-v2/backend/internal/oauth2"
	"github.com/danielhoward-me/sso-v2/backend/internal/utils"
)

type tokenEndpointBodyBase struct {
	GrantType string `json:"grant_type"`
}
type tokenEndpointBodyAuthorisationCode struct {
	tokenEndpointBodyBase
	Code        string `json:"code"`
	RedirectUri string `json:"redirect_uri"`
}

type oauth2ErrorType string

const invalidRequest oauth2ErrorType = "invalid_request"
const invalidClient oauth2ErrorType = "invalid_client"
const invalidGrant oauth2ErrorType = "invalid_grant"
const unauthorizedClient oauth2ErrorType = "unauthorized_client"
const unsupportedGrantType oauth2ErrorType = "unsupported_grant_type"

type oauth2Error struct {
	Error            oauth2ErrorType `json:"error"`
	ErrorDescription string          `json:"error_description"`
}

var invalidClientError oauth2Error = oauth2Error{
	Error:            invalidClient,
	ErrorDescription: "the specified client id or secret is not correct",
}
var invalidGrantError oauth2Error = oauth2Error{
	Error:            invalidGrant,
	ErrorDescription: "the specified grant type doesn't exist",
}

func handleTokenEndpoint(c *gin.Context) {
	clientId, clientSecret, exists := c.Request.BasicAuth()
	if !exists {
		c.JSON(401, invalidClientError)
	}

	clientUuid, err := uuid.Parse(clientId)
	if err != nil {
		panic(fmt.Errorf("token endpoint was passed an invalid client id: %s", err))
	}

	client, err := oauth2.NewClientFromUUID(clientUuid)
	if err != nil {
		if utils.ErrIsNoRows(err) {
			c.JSON(401, invalidClientError)
			return
		}

		panic(err)
	}

	if !client.CheckSecret(clientSecret) {
		c.JSON(401, invalidClientError)
	}

	body := mustParseBody[tokenEndpointBodyBase](c)

	switch body.GrantType {
	case "authorization_code":
		{
			// fullBody := mustParseBody[TokenEndpointBodyAuthorisationCode](c)
		}
	default:
		{
			c.JSON(400, invalidGrantError)
		}
	}
}

func mustParseBody[T interface{}](c *gin.Context) T {
	var body T
	if err := c.BindJSON(&body); err != nil {
		panic(fmt.Errorf("token endpoint failed to pass request body: %s", err))
	}
	return body
}
