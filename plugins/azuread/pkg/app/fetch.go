package app

import (
	"context"
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/azuread/pkg/app/fetch"
	"os"
	"time"
)

type FetchCmd struct {
	Tenant       string `short:"a" help:"AzureAD tenant" env:"AZUREAD_TENANT" required:""`
	ClientID     string `short:"i" help:"AzureAD Client ID" env:"AZUREAD_CLIENT_ID" required:""`
	ClientSecret string `short:"s" help:"AzureAD Client Secret" env:"AZUREAD_CLIENT_SECRET" required:""`
	RefreshToken string `short:"r" help:"AzureAD Refresh Token" env:"AZUREAD_REFRESH_TOKEN" optional:""`
}

func (cmd *FetchCmd) Run(kongCtx *kong.Context) error {
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Hour)
	defer cancel()

	fetcher, err := fetch.New(ctx, cmd.Tenant, cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken)
	if err != nil {
		return err
	}

	return fetcher.Fetch(timeoutCtx, os.Stdout, os.Stderr)
}
