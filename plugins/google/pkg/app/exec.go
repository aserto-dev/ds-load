package app

import (
	"context"
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/google/pkg/fetch"
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
	fetcher, err := fetch.New(ctx, cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken, cmd.Customer)
	if err != nil {
		return err
	}
	fetcher = fetcher.WithGroups(cmd.Groups)

	templateContent, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}
	transformer := transform.NewGoTemplateTransform(templateContent)
	return cmd.exec(ctx, transformer, fetcher)
}

func (cmd *ExecCmd) exec(ctx context.Context, transformer plugin.Transformer, fetcher plugin.Fetcher) error {
	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()

	go func() {
		err := fetcher.Fetch(ctx, pipeWriter, os.Stderr)
		if err != nil {
			log.Printf("Could not fetch data %s", err.Error())
		}
		pipeWriter.Close()
	}()

	return transformer.Transform(ctx, pipeReader, os.Stdout, os.Stderr)
}
