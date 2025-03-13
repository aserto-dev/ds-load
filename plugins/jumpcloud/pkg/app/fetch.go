package app

import (
	"os"

	"github.com/aserto-dev/ds-load/plugins/jumpcloud/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/jumpcloud/pkg/jc"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type FetchCmd struct {
	APIKey string `short:"k" help:"JumpCloud API Key" env:"JC_API_KEY" required:""`
	Groups bool   `short:"g" help:"Retrieve JumpCloud groups" env:"JC_GROUPS" default:"false"`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	jcClient, err := jc.NewJumpCloudClient(ctx.Context, cmd.APIKey)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(jcClient)
	if err != nil {
		return err
	}

	fetcher = fetcher.WithGroups(cmd.Groups)

	return fetcher.Fetch(ctx.Context, os.Stdout, os.Stderr)
}
