package main

import (
	"fmt"
	"os"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/gin-gonic/gin"
)

var DEFAULT_PARAMS = &types.Params{
	"title": getTitle(),
}

var alerterActive = false
var alerter *router.ServiceRouter

func createAlertSender() {
	alertUrl := os.Getenv("ALERT_URL")
	if alertUrl == "" {
		fmt.Println("ALERT_URL env variable isn't set, not alerting on server errors")
		return
	}

	sender, err := shoutrrr.CreateSender(alertUrl)
	if err != nil {
		panic(fmt.Errorf("error when creating shoutrrr sender: %s", err))
	}

	alerter = sender
	alerterActive = true
}

func getTitle() string {
	baseTitle := "SSO Backend"

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when fetching system hostname: %s\n", err)
		return baseTitle
	}

	return fmt.Sprintf("%s on %s", baseTitle, hostname)
}

func alert(err any) {
	if !alerterActive {
		return
	}

	message := fmt.Sprintf("%sInternal server error has occured:\n%s", getDevMessage(), err)
	alerter.Send(message, DEFAULT_PARAMS)
}

func getDevMessage() string {
	if gin.IsDebugging() {
		return "DEVELOPMENT MODE\n\n"
	} else {
		return ""
	}
}
