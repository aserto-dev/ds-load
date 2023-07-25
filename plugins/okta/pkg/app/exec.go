package app

import (
	"context"
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
	"github.com/aserto-dev/ds-load/sdk/plugin"
	"github.com/aserto-dev/ds-load/sdk/transform"
	"github.com/rs/zerolog/log"
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

	templateContent, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}

	pipeReader, pipeWriter := io.Pipe()
	transformer := transform.NewGoTemplateTransform(templateContent)

	go func() {
		err = plugin.NewDSPlugin(plugin.WithOutputWriter(pipeWriter)).WriteFetchOutput(results, errCh)
		if err != nil {
			log.Printf("Could not write fetcher output %s", err.Error())
		}

		pipeWriter.Close()
	}()

	defer pipeReader.Close()

	return cmd.exec(ctx, transformer, pipeReader)
}

func (cmd *ExecCmd) exec(ctx context.Context, transformer plugin.Transformer, pipeReader io.Reader) error {
	return transformer.Transform(ctx, pipeReader, os.Stdout, os.Stderr)
}
