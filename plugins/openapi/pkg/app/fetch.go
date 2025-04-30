package app

import (
	"os"

	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/openapi"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type FetchCmd struct {
	Directory   string `short:"d" help:"OpenAPI Spec Directory" env:"OPENAPI_DIRECTORY"`
	URL         string `short:"u" help:"OpenAPI Spec URL" env:"OPENAPI_URL"`
	IDFormat    string `short:"f" help:"ID Format (base64, canonical, default)" env:"OPENAPI_IDFORMAT"`
	ServiceName string `short:"n" help:"Service name when importing from a URL" env:"OPENAPI_SERVICE_NAME"`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	openapiClient, err := openapi.New(cmd.Directory, cmd.URL, cmd.IDFormat, cmd.ServiceName)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(openapiClient)
	if err != nil {
		return err
	}

	fetcher = fetcher.WithDirectory(cmd.Directory).WithURL(cmd.URL).WithIDFormat(cmd.IDFormat).WithServiceName(cmd.ServiceName)

	return fetcher.Fetch(ctx.Context, os.Stdout, common.NewErrorWriter(os.Stderr))
}
