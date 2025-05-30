package app

import (
	"os"
	"strings"

	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/auth0client"
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/fetch"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type FetchCmd struct {
	Domain         string `name:"domain" short:"d" env:"AUTH0_DOMAIN" help:"auth0 domain" required:""`
	ClientID       string `name:"client-id" short:"i" env:"AUTH0_CLIENT_ID" help:"auth0 client id" required:""`
	ClientSecret   string `name:"client-secret" short:"s" env:"AUTH0_CLIENT_SECRET" help:"auth0 client secret" required:""`
	ConnectionName string `name:"connection-name" env:"AUTH0_CONNECTION_NAME" help:"auth0 connection name" optional:""`
	UserPID        string `name:"user-pid" env:"AUTH0_USER_PID" help:"auth0 user PID of the user you want to read" optional:""`
	UserEmail      string `name:"user-email" env:"AUTH0_USER_EMAIL" help:"auth0 user email of the user you want to read" optional:""`
	Roles          bool   `name:"roles" env:"AUTH0_ROLES" default:"false" negatable:"" help:"include roles"`
	Orgs           bool   `name:"orgs" env:"AUTH0_ORGS" default:"false" negatable:"" help:"include organizations"`
	RateLimit      bool   `name:"rate-limit" default:"true" help:"enable http client rate limiter" negatable:""`
	SAML           bool   `name:"saml" default:"false" help:"use SAML data loader"`
}

func (f *FetchCmd) Run(ctx *cc.CommonCtx) error {
	if f.UserPID != "" && !strings.Contains(f.UserPID, "|") {
		f.UserPID = auth0prefix + f.UserPID
	}

	client, err := auth0client.New(ctx.Context, f.ClientID, f.ClientSecret, f.Domain)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(ctx.Context, client)
	if err != nil {
		return err
	}

	fetcher = fetcher.
		WithUserPID(f.UserPID).
		WithEmail(f.UserEmail).
		WithRoles(f.Roles).
		WithOrgs(f.Orgs).
		WithSAML(f.SAML).
		WithConnectionName(f.ConnectionName)

	return fetcher.Fetch(ctx.Context, os.Stdout, common.NewErrorWriter(os.Stderr))
}
