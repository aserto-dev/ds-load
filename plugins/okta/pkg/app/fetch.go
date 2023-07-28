package app

import (
	"os"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/fetch"
)

type FetchCmd struct {
	Domain         string `env:"DS_OKTA_DOMAIN"`
	APIToken       string `env:"DS_OKTA_TOKEN"`
	Groups         bool   `env:"DS_OKTA_GROUPS" default:"true" negatable:""`
	Roles          bool   `env:"DS_OKTA_ROLES" default:"true" negatable:""`
	RequestTimeout int64  `default:"0" optional:""`
}

func (f *FetchCmd) Run(ctx *cc.CommonCtx) error {
	fetcher, err := fetch.New(ctx.Context, f.Domain, f.APIToken, f.RequestTimeout, f.Groups, f.Roles)
	if err != nil {
		return err
	}

	return fetcher.Fetch(ctx.Context, os.Stdout, os.Stderr)
}
