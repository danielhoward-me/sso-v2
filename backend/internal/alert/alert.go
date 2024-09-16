package alert

import (
	"fmt"
	"os"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/danielhoward-me/sso-v2/backend/internal/utils"
)

var defaultParams = &types.Params{
	"title": getTitle(),
}

type Alerter struct {
	sender *router.ServiceRouter
}

func New(alertUrl string) (alerter *Alerter, err error) {
	sender, err := shoutrrr.CreateSender(alertUrl)
	if err != nil {
		return
	}

	return &Alerter{
		sender: sender,
	}, nil
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

func (alerter Alerter) Alert(message string) {
	fullMessage := fmt.Sprintf("%s%s", getDevMessage(), message)
	alerter.sender.Send(fullMessage, defaultParams)
}

func getDevMessage() string {
	if utils.IsDevelopment() {
		return "DEVELOPMENT MODE\n\n"
	} else {
		return ""
	}
}
