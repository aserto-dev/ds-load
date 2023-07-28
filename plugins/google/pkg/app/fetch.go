package app

import (
	"os"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/google/pkg/fetch"
)

type FetchCmd struct {
	ClientID     string `short:"i" help:"Google Client ID" env:"GOOGLE_CLIENT_ID" required:""`
	ClientSecret string `short:"s" help:"Google Client Secret" env:"GOOGLE_CLIENT_SECRET" required:""`
	RefreshToken string `short:"r" help:"Google Refresh Token" env:"GOOGLE_REFRESH_TOKEN" required:""`
	Groups       bool   `short:"g" help:"Retrieve Google groups" env:"GOOGLE_GROUPS" default:"false"`
	Customer     string `help:"Google Workspace Customer field" env:"GOOGLE_CUSTOMER" default:"my_customer"`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	fetcher, err := fetch.New(ctx.Context, cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken, cmd.Customer)
	if err != nil {
		return err
	}
	fetcher = fetcher.WithGroups(cmd.Groups)

	return fetcher.Fetch(ctx.Context, os.Stdout, os.Stderr)
}
