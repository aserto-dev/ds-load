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
	cognitoClient, err := createCognitoClient(cmd.AccessKey, cmd.SecretKey, cmd.UserPoolID, cmd.Region)
	if err != nil {
		return err
	}

	results := make(chan map[string]interface{}, 1)
	errCh := make(chan error, 1)
	go func() {
		Fetch(cognitoClient, cmd.Groups, results, errCh)
		close(results)
		close(errCh)
	}()
	content, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}

	return plugin.NewDSPlugin(plugin.WithTemplate(content)).WriteFetchOutput(results, errCh, true)
}
