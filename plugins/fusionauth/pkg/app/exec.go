package app

import (
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/fusionauth/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/fusionauth/pkg/fusionauthclient"
	"github.com/aserto-dev/ds-load/sdk/exec"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(ctx *cc.CommonCtx) error {

	fusionauthClient, err := fusionauthclient.NewFusionAuthClient(cmd.HostURL, cmd.ApiKey)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(fusionauthClient)
	if err != nil {
		return err
	}
	fetcher = fetcher.WithGroups(cmd.Groups)

	templateContent, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}
	transformer := transform.NewGoTemplateTransform(templateContent)
	return exec.Execute(ctx.Context, ctx.Log, transformer, fetcher)
}
