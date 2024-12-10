package app

import (
	"fmt"

	"github.com/aserto-dev/ds-load/plugins/google/pkg/googleclient"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type GetTokenCmd struct {
	ClientID     string `short:"i" help:"Google Client ID" env:"GOOGLE_CLIENT_ID" required:""`
	ClientSecret string `short:"s" help:"Google Client Secret" env:"GOOGLE_CLIENT_SECRET" required:""`
	Port         int    `short:"p" help:"Port number to run callback server on" env:"GOOGLE_PORT" default:"8761"`
}

func (cmd *GetTokenCmd) Run(ctx *cc.CommonCtx) error {
	refreshToken, err := googleclient.GetRefreshToken(
		ctx.Context,
		cmd.ClientID,
		cmd.ClientSecret,
		cmd.Port)
	if err != nil {
		return err
	}

	fmt.Println("Refresh token: ", refreshToken)
	return nil
}
