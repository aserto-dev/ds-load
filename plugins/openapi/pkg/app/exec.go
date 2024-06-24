package app

import (
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/openapiclient"
	"github.com/aserto-dev/ds-load/sdk/exec"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(ctx *cc.CommonCtx) error {

	openapiClient, err := openapiclient.NewOpenAPIClient(cmd.Directory, cmd.URL)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(openapiClient)
	if err != nil {
		return err
	}
	fetcher = fetcher.WithDirectory(cmd.Directory).WithURL(cmd.URL)

	templateContent, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}
	transformer := transform.NewGoTemplateTransform(templateContent)
	return exec.Execute(ctx.Context, ctx.Log, transformer, fetcher)
}
