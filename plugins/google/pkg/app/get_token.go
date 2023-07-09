package app

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/google/pkg/googleclient"
)

type GetTokenCmd struct {
	ClientID     string `short:"i" help:"Google Client ID" env:"GOOGLE_CLIENT_ID" required:""`
	ClientSecret string `short:"s" help:"Google Client Secret" env:"GOOGLE_CLIENT_SECRET" required:""`
	Port         int    `short:"p" help:"Port number to run callback server on" env:"GOOGLE_PORT" default:"8761"`
}

func (cmd *GetTokenCmd) Run(ctx *kong.Context) error {
	refreshToken, err := googleclient.GetRefreshToken(
		context.Background(),
		cmd.ClientID,
		cmd.ClientSecret,
		cmd.Port)
	if err != nil {
		return err
	}

	fmt.Println("Refresh token: ", refreshToken)
	return nil
}
