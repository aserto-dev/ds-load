package app

import (
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/ldapclient"
	"github.com/aserto-dev/ds-load/sdk/exec"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(ctx *cc.CommonCtx) error {
	credentials := &ldapclient.Credentials{
		User:     cmd.User,
		Password: cmd.Password,
	}

	conOptions := &ldapclient.ConnectionOptions{
		Host:        cmd.Host,
		BaseDN:      cmd.BaseDn,
		UserFilter:  cmd.UserFilter,
		GroupFilter: cmd.GroupFilter,
		Insecure:    cmd.Insecure,
		UuidField:   cmd.UuidField,
	}

	client, err := ldapclient.NewLDAPClient(credentials, conOptions)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(client)
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
