package app

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/sdk/plugin"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(context *kong.Context) error {
	azureClient, err := createAzureAdClient(cmd.Tenant, cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken)
	if err != nil {
		return err
	}

	results := make(chan map[string]interface{}, 1)
	errCh := make(chan error, 1)
	go func() {
		Fetch(azureClient, results, errCh)
		close(results)
		close(errCh)
	}()
	content, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}

	return plugin.NewDSPlugin(plugin.WithTemplate(content)).WriteFetchOutput(results, errCh, true)
}
