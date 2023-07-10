package app

import (
	"net/http"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/httpclient"
	"github.com/aserto-dev/ds-load/sdk/plugin"
	"github.com/auth0/go-auth0/management"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(context *kong.Context) error {
	if cmd.UserPID != "" && !strings.HasPrefix(cmd.UserPID, "auth0|") {
		cmd.UserPID = "auth0|" + cmd.UserPID
	}

	options := []management.Option{
		management.WithClientCredentials(
			cmd.ClientID,
			cmd.ClientSecret,
		),
	}
	if cmd.RateLimit {
		client := http.DefaultClient
		client.Transport = httpclient.NewTransport(http.DefaultTransport)
		options = append(options, management.WithClient(client))
	}

	mgmt, err := management.New(
		cmd.Domain,
		options...,
	)
	if err != nil {
		return err
	}

	cmd.mgmt = mgmt

	results := make(chan map[string]interface{}, 1)
	errCh := make(chan error, 1)
	go func() {
		cmd.Fetch(results, errCh)
		close(results)
		close(errCh)
	}()

	content, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}

	return plugin.NewDSPlugin(plugin.WithTemplate(content)).WriteFetchOutput(results, errCh, true)
}
