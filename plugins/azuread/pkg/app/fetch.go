package app

import (
	"os"

	"github.com/aserto-dev/ds-load/plugins/azuread/pkg/fetch"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

//nolint:lll // user and group properties include complex defaults.
type FetchCmd struct {
	Tenant          string   `short:"a" help:"AzureAD tenant" env:"AZUREAD_TENANT" required:""`
	ClientID        string   `short:"i" help:"AzureAD Client ID" env:"AZUREAD_CLIENT_ID" required:""`
	ClientSecret    string   `short:"s" help:"AzureAD Client Secret" env:"AZUREAD_CLIENT_SECRET" required:""`
	RefreshToken    string   `short:"r" help:"AzureAD Refresh Token" env:"AZUREAD_REFRESH_TOKEN" optional:""`
	Groups          bool     `short:"g" help:"Include groups" env:"AZUREAD_INCLUDE_GROUPS" optional:""`
	UserProperties  []string `help:"User properties to query" env:"AZUREAD_USER_PROPERTIES" optional:"" default:"displayName,id,mail,createdDateTime,mobilePhone,userPrincipalName,accountEnabled"`
	GroupProperties []string `help:"Group properties to query" env:"AZUREAD_GROUP_PROPERTIES" optional:"" default:"displayName,id,mail,createdDateTime,mailNickname"`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	azureClient, err := createAzureAdClient(ctx.Context, cmd.Tenant, cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(ctx.Context, azureClient, cmd.UserProperties, cmd.GroupProperties)
	if err != nil {
		return err
	}

	return fetcher.WithGroups(cmd.Groups).Fetch(ctx.Context, os.Stdout, common.NewErrorWriter(os.Stderr))
}
