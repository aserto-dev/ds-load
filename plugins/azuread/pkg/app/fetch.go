package app

import (
	"os"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/azuread/pkg/fetch"
)

type FetchCmd struct {
	Tenant       string `short:"a" help:"AzureAD tenant" env:"AZUREAD_TENANT" required:""`
	ClientID     string `short:"i" help:"AzureAD Client ID" env:"AZUREAD_CLIENT_ID" required:""`
	ClientSecret string `short:"s" help:"AzureAD Client Secret" env:"AZUREAD_CLIENT_SECRET" required:""`
	RefreshToken string `short:"r" help:"AzureAD Refresh Token" env:"AZUREAD_REFRESH_TOKEN" optional:""`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	fetcher, err := fetch.New(ctx.Context, cmd.Tenant, cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken)
	if err != nil {
		return err
	}

	return fetcher.Fetch(ctx.Context, os.Stdout, os.Stderr)
}
