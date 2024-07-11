package app

import (
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/openapi"
	"github.com/aserto-dev/ds-load/sdk/exec"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(ctx *cc.CommonCtx) error {
	openapiClient, err := openapi.New(cmd.Directory, cmd.URL, cmd.IDFormat, cmd.ServiceName)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(openapiClient)
	if err != nil {
		return err
	}
	fetcher = fetcher.WithDirectory(cmd.Directory).WithURL(cmd.URL).WithIDFormat(cmd.IDFormat).WithServiceName(cmd.ServiceName)

	templateContent, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}
	transformer := transform.NewGoTemplateTransform(templateContent)
	return exec.Execute(ctx.Context, ctx.Log, transformer, fetcher)
}
