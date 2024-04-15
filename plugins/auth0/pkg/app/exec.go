package app

import (
	"strings"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/auth0client"
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/fetch"
	"github.com/aserto-dev/ds-load/sdk/exec"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(ctx *cc.CommonCtx) error {
	if cmd.UserPID != "" && !strings.Contains(cmd.UserPID, "|") {
		cmd.UserPID = auth0prefix + cmd.UserPID
	}

	client, err := auth0client.New(ctx.Context, cmd.ClientID, cmd.ClientSecret, cmd.Domain)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(ctx.Context, client)
	if err != nil {
		return err
	}
	fetcher = fetcher.WithUserPID(cmd.UserPID).WithEmail(cmd.UserEmail).WithRoles(cmd.Roles)

	templateContent, err := cmd.getTemplateContent()
	if err != nil {
		return err
	}
	transformer := transform.NewGoTemplateTransform(templateContent)
	return exec.Execute(ctx.Context, ctx.Log, transformer, fetcher)
}
