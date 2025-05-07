package app

import (
	"os"

	"github.com/aserto-dev/ds-load/plugins/keycloak/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/keycloak/pkg/kc"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type FetchCmd struct {
	APIKey string `short:"k" help:"keycloak API Key" env:"KC_API_KEY" required:""`
	Groups bool   `short:"g" help:"Retrieve keycloak groups" env:"KC_GROUPS" default:"false"`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	kcClient, err := kc.NewKeyCloudClient(ctx.Context, cmd.APIKey)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(kcClient)
	if err != nil {
		return err
	}

	fetcher = fetcher.WithGroups(cmd.Groups)

	return fetcher.Fetch(ctx.Context, os.Stdout, common.NewErrorWriter(os.Stderr))
}
