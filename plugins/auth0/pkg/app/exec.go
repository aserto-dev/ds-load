package app

import (
	"context"
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/app/fetch"
	"github.com/aserto-dev/ds-load/sdk/plugin"
	"github.com/aserto-dev/ds-load/sdk/transform"
	"io"
	"log"
	"os"
	"strings"
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

	fetcher, err := fetch.New(cmd.UserPID, cmd.ClientID, cmd.ClientSecret, cmd.Domain)
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
