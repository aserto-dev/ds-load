package app

import (
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/azuread/pkg/fetch"
	"github.com/aserto-dev/ds-load/sdk/exec"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(ctx *cc.CommonCtx) error {
	fetcher, err := fetch.New(ctx.Context, cmd.Tenant, cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken)
	if err != nil {
		return err
	}

	templateContent, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}
	transformer := transform.NewGoTemplateTransform(templateContent)
	return exec.Execute(ctx.Context, ctx.Log, transformer, fetcher)
}
