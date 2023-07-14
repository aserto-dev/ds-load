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
	googleClient, err := createGoogleClient(cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken, cmd.Customer)
	if err != nil {
		return err
	}

	results := make(chan map[string]interface{}, 1)
	errCh := make(chan error, 1)
	go func() {
		Fetch(googleClient, cmd.Groups, results, errCh)
		close(results)
		close(errCh)
	}()
	content, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}

	return plugin.NewDSPlugin(plugin.WithTemplate(content), plugin.WithMaxChunkSize(cmd.MaxChunkSize)).WriteFetchOutput(results, errCh, true)
}
