package app

import (
	"strings"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/fetch"
	"github.com/aserto-dev/ds-load/sdk/exec"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(ctx *cc.CommonCtx) error {
	if cmd.UserPID != "" && !strings.HasPrefix(cmd.UserPID, "auth0|") {
		cmd.UserPID = "auth0|" + cmd.UserPID
	}

	fetcher, err := fetch.New(cmd.ClientID, cmd.ClientSecret, cmd.Domain)
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
