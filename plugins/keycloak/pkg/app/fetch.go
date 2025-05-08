package app

import (
	"os"

	"github.com/aserto-dev/ds-load/plugins/keycloak/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/keycloak/pkg/kc"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type FetchCmd struct {
	kc.KeycloakClientConfig
	Groups bool `short:"g" help:"Retrieve keycloak groups" env:"KEYCLOAK_GROUPS" default:"false"`
	Roles  bool `short:"r"  help:"Retrieve keycloak roles" env:"KEYCLOAK_ROLES" default:"false"`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	kcClient, err := kc.NewKeycloakClient(ctx.Context, &cmd.KeycloakClientConfig)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(kcClient)
	if err != nil {
		return err
	}

	fetcher = fetcher.WithGroups(cmd.Groups).WithRoles(cmd.Roles)

	return fetcher.Fetch(ctx.Context, os.Stdout, common.NewErrorWriter(os.Stderr))
}
