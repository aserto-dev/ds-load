package app

import (
	"context"
	"io"
	"os"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/google/pkg/fetch"
	"github.com/aserto-dev/ds-load/sdk/plugin"
	"github.com/aserto-dev/ds-load/sdk/transform"
	"github.com/rs/zerolog"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(ctx *cc.CommonCtx) error {
	fetcher, err := fetch.New(ctx.Context, cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken, cmd.Customer)
	if err != nil {
		return err
	}
	fetcher = fetcher.WithGroups(cmd.Groups)

	templateContent, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}
	transformer := transform.NewGoTemplateTransform(templateContent)
	return cmd.exec(ctx.Context, ctx.Log, transformer, fetcher)
}

func (cmd *ExecCmd) exec(ctx context.Context, log *zerolog.Logger, transformer plugin.Transformer, fetcher plugin.Fetcher) error {
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
