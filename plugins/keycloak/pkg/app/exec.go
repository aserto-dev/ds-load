package app

import (
	"github.com/aserto-dev/ds-load/plugins/keycloak/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/keycloak/pkg/kc"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
	"github.com/aserto-dev/ds-load/sdk/exec"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(ctx *cc.CommonCtx) error {
	gClient, err := kc.NewKeycloakClient(ctx.Context, &cmd.KeycloakClientConfig)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(gClient)
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
