package app

import (
	"os"

	"github.com/aserto-dev/ds-load/plugins/zitadel/pkg/client"
	"github.com/aserto-dev/ds-load/plugins/zitadel/pkg/fetch"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type FetchCmd struct {
	client.ZitadelConfig
	Groups bool `short:"g" help:"Retrieve zitadel groups" env:"ZITADEL_GROUPS" default:"false"`
	Roles  bool `short:"r"  help:"Retrieve zitadel roles" env:"ZITADEL_ROLES" default:"false"`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	kcClient, err := client.NewZitadelClient(ctx.Context, &cmd.ZitadelConfig)
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
