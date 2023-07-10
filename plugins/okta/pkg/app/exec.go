package app

import (
	"context"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
	"github.com/aserto-dev/ds-load/sdk/plugin"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(ctx *kong.Context) error {
	oktaClient, err := oktaclient.NewOktaClient(context.Background(), cmd.Domain, cmd.APIToken)
	if err != nil {
		return err
	}

	results := make(chan map[string]interface{}, 1)
	errCh := make(chan error, 1)
	go func() {
		cmd.Fetch(oktaClient, results, errCh)
		close(results)
		close(errCh)
	}()

	content, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}

	return plugin.NewDSPlugin(plugin.WithTemplate(content)).WriteFetchOutput(results, errCh, true)
}
