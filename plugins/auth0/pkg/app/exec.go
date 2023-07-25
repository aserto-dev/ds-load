package app

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/httpclient"
	"github.com/aserto-dev/ds-load/sdk/plugin"
	"github.com/aserto-dev/ds-load/sdk/transform"
	"github.com/auth0/go-auth0/management"
	"github.com/rs/zerolog/log"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(kongCtx *kong.Context) error {
	ctx := context.Background()
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
