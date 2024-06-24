package app

import (
	"os"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/openapiclient"
)

type FetchCmd struct {
	Directory string `short:"d" help:"OpenAPI Spec Directory" env:"OPENAPI_DIRECTORY" required:""`
	URL       string `short:"u" help:"OpenAPI Spec URL" env:"OPENAPI_URL" required:""`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	openapiClient, err := openapiclient.NewOpenAPIClient(cmd.Directory, cmd.URL)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(openapiClient)
	if err != nil {
		return err
	}
	fetcher = fetcher.WithDirectory(cmd.Directory).WithURL(cmd.URL)

	return fetcher.Fetch(ctx.Context, os.Stdout, os.Stderr)
}
