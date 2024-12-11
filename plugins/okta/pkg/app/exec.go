package app

import (
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
	"github.com/aserto-dev/ds-load/sdk/exec"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(ctx *cc.CommonCtx) error {
	oktaClient, err := oktaclient.NewOktaClient(cmd.Domain, cmd.APIToken, cmd.RequestTimeout)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(ctx.Context, oktaClient)
	if err != nil {
		return err
	}

	fetcher = fetcher.WithGroups(cmd.Groups).WithRoles(cmd.Roles)

	templateContent, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}
	transformer := transform.NewGoTemplateTransform(templateContent)
	return exec.Execute(ctx.Context, ctx.Log, transformer, fetcher)
}
