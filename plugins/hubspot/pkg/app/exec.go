package app

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/hubspot/pkg/hubspotclient"
	"github.com/aserto-dev/ds-load/sdk/plugin"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(context *kong.Context) error {
	var hubspotClient *hubspotclient.HubspotClient
	var err error
	if cmd.PrivateAccessToken != "" {
		hubspotClient, _ = createHubspotClient(cmd.PrivateAccessToken)
	} else {
		hubspotClient, err = createHubspotOAuth2Client(cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken)
		if err != nil {
			return err
		}
	}

	results := make(chan map[string]interface{}, 1)
	errCh := make(chan error, 1)
	go func() {
		Fetch(hubspotClient, cmd.Contacts, cmd.Companies, results, errCh)
		close(results)
		close(errCh)
	}()
	content, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}

	return plugin.NewDSPlugin(plugin.WithTemplate(content), plugin.WithMaxChunkSize(cmd.MaxChunkSize)).WriteFetchOutput(results, errCh, true)
}
