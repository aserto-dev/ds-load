package app

import (
	"os"

	"github.com/aserto-dev/ds-load/plugins/okta/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type FetchCmd struct {
	Domain         string `env:"DS_OKTA_DOMAIN"`
	APIToken       string `env:"DS_OKTA_TOKEN"`
	Groups         bool   `env:"DS_OKTA_GROUPS" default:"true" negatable:""`
	Roles          bool   `env:"DS_OKTA_ROLES" default:"true" negatable:""`
	RequestTimeout int64  `default:"0" optional:""`
}

func (f *FetchCmd) Run(ctx *cc.CommonCtx) error {
	oktaClient, err := oktaclient.NewOktaClient(f.Domain, f.APIToken, f.RequestTimeout)
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(ctx.Context, oktaClient)
	if err != nil {
		return err
	}

	return fetcher.WithGroups(f.Groups).WithRoles(f.Roles).Fetch(ctx.Context, os.Stdout, os.Stderr)
}
