package app

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/azuread/pkg/azureclient"
	"github.com/aserto-dev/ds-load/sdk/common/version"
)

type CLI struct {
	Config          kong.ConfigFlag    `help:"Configuration file path" short:"c"`
	Version         VersionCmd         `cmd:"" help:"version information"`
	Fetch           FetchCmd           `cmd:"" help:"fetch Azure AD data"`
	Transform       TransformCmd       `cmd:"" help:"transform Azure AD data"`
	ExportTransform ExportTransformCmd `cmd:"" help:"export default transform template"`
	Exec            ExecCmd            `cmd:"" help:"fetch and transform Azure AD data" default:"withargs"`
	Verbosity       int                `short:"v" type:"counter" help:"Use to increase output verbosity."`
	Verify          VerifyCmd          `cmd:"verify" help:"verify fetcher configuration and credentials"`
}

type VersionCmd struct{}

func (cmd *VersionCmd) Run() error {
	fmt.Printf("%s - %s\n",
		AppName,
		version.GetInfo().String(),
	)

	return nil
}

func createAzureAdClient(
	ctx context.Context,
	tenant, clientID, clientSecret, refreshToken string,
) (*azureclient.AzureADClient, error) {
	if refreshToken != "" {
		return azureclient.NewAzureADClientWithRefreshToken(
			ctx,
			tenant,
			clientID,
			clientSecret,
			refreshToken)
	}

	return azureclient.NewAzureADClient(
		ctx,
		tenant,
		clientID,
		clientSecret)
}
