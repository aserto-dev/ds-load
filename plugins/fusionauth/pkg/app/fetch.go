package app

import (
	"os"

	"github.com/aserto-dev/ds-load/plugins/fusionauth/pkg/client"
	"github.com/aserto-dev/ds-load/plugins/fusionauth/pkg/fetch"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type FetchCmd struct {
	HostURL string `short:"u" help:"FusionAuth Host URL" env:"FUSIONAUTH_HOST_URL" required:""`
	APIKey  string `short:"k" help:"FusionAuth API Key" env:"FUSIONAUTH_API_KEY" required:""`
	Groups  bool   `short:"g" help:"Retrieve FusionAuth groups" env:"FUSIONAUTH_GROUPS" default:"false" negatable:""`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	fusionauthClient, err := client.NewFusionAuthClient(cmd.HostURL, cmd.APIKey)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(fusionauthClient)
	if err != nil {
		return err
	}

	fetcher = fetcher.WithGroups(cmd.Groups).WithHost(cmd.HostURL)

	return fetcher.Fetch(ctx.Context, os.Stdout, common.NewErrorWriter(os.Stderr))
}
