package app

import (
	"os"

	"github.com/aserto-dev/ds-load/plugins/azureadb2c/pkg/fetch"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

//nolint:lll // user and group properties include complex defaults.
type FetchCmd struct {
	Tenant          string   `short:"a" help:"AzureAD B2C tenant" env:"AZUREADB2C_TENANT" required:""`
	ClientID        string   `short:"i" help:"AzureAD B2C Client ID" env:"AZUREADB2C_CLIENT_ID" required:""`
	ClientSecret    string   `short:"s" help:"AzureAD B2C Client Secret" env:"AZUREADB2C_CLIENT_SECRET" required:""`
	RefreshToken    string   `short:"r" help:"AzureAD B2C Refresh Token" env:"AZUREADB2C_REFRESH_TOKEN" optional:""`
	Groups          bool     `short:"g" help:"Include groups" env:"AZUREADB2C_INCLUDE_GROUPS" optional:""`
	UserProperties  []string `help:"User properties to query" env:"AZUREADB2C_USER_PROPERTIES" optional:"" default:"displayName,id,mail,createdDateTime,mobilePhone,userPrincipalName,accountEnabled,identities,creationType"`
	GroupProperties []string `help:"Group properties to query" env:"AZUREADB2C_GROUP_PROPERTIES" optional:"" default:"displayName,id,mail,createdDateTime,mailNickname,members,transitiveMembers"`
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

	return fetcher.WithGroups(cmd.Groups).Fetch(ctx.Context, os.Stdout, common.ErrorWriter{Writer: os.Stderr})
}
