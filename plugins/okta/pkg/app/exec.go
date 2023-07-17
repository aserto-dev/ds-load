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

func (cmd *ExecCmd) Run(kongCtx *kong.Context) error {
	ctx := context.Background()
	oktaClient, err := oktaclient.NewOktaClient(ctx, cmd.Domain, cmd.APIToken, cmd.RequestTimeout)
	if err != nil {
		return err
	}

	cmd.oktaClient = oktaClient
	results := make(chan map[string]interface{}, 1)
	errCh := make(chan error, 1)

	go func() {
		cmd.Fetch(ctx, results, errCh)
		close(results)
		close(errCh)
	}()

	content, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}

	return plugin.NewDSPlugin(plugin.WithTemplate(content), plugin.WithMaxChunkSize(cmd.MaxChunkSize)).WriteFetchOutput(results, errCh, true)
}
