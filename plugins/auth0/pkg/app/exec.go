package app

import (
	"context"
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/app/fetcher"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/sdk/plugin"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(kongCtx *kong.Context) error {
	ctx := context.Background()

	fetcher, err := fetcher.New(cmd.UserPID, cmd.ClientID, cmd.ClientSecret, cmd.Domain)
	if err != nil {
		return err
	}

	templateContent, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}
	transformer := transform.NewGoTemplateTransform(templateContent)
	return cmd.exec(ctx, transformer, fetcher)
}

func (cmd *ExecCmd) exec(ctx context.Context, transformer plugin.Transformer, fetcher plugin.Fetcher) error {
	if cmd.UserPID != "" && !strings.HasPrefix(cmd.UserPID, "auth0|") {
		cmd.UserPID = "auth0|" + cmd.UserPID
	}

	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()

	go func() {
		fetcher.Fetch(ctx, pipeWriter, os.Stderr)
		pipeWriter.Close()
	}()

	return transformer.Transform(ctx, pipeReader, os.Stdout, os.Stderr)
}
