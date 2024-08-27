package app

import (
	"os"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/azureadb2c/pkg/fetch"
)

type FetchCmd struct {
	Tenant       string `short:"a" help:"AzureAD B2C tenant" env:"AZUREADB2C_TENANT" required:""`
	ClientID     string `short:"i" help:"AzureAD B2C Client ID" env:"AZUREADB2C_CLIENT_ID" required:""`
	ClientSecret string `short:"s" help:"AzureAD B2C Client Secret" env:"AZUREADB2C_CLIENT_SECRET" required:""`
	RefreshToken string `short:"r" help:"AzureAD B2C Refresh Token" env:"AZUREADB2C_REFRESH_TOKEN" optional:""`
	Groups       bool   `short:"g" help:"Include groups" env:"AZUREADB2C_INCLUDE_GROUPS" optional:""`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	azureClient, err := createAzureAdClient(ctx.Context, cmd.Tenant, cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(ctx.Context, azureClient)
	if err != nil {
		return err
	}

	return fetcher.WithGroups(cmd.Groups).Fetch(ctx.Context, os.Stdout, os.Stderr)
}
