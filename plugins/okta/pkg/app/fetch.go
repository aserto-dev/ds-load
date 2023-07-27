package app

import (
	"context"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/app/fetch"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
)

type FetchCmd struct {
	Domain         string `env:"DS_OKTA_DOMAIN"`
	APIToken       string `env:"DS_OKTA_TOKEN"`
	Groups         bool   `env:"DS_OKTA_GROUPS" default:"true" negatable:""`
	Roles          bool   `env:"DS_OKTA_ROLES" default:"true" negatable:""`
	RequestTimeout int64  `default:"0" optional:""`

	oktaClient oktaclient.OktaClient `kong:"-"`
}

func (f *FetchCmd) Run(kongCtx *kong.Context) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	fetcher, err := fetch.New(timeoutCtx, f.Domain, f.APIToken, f.RequestTimeout, f.Groups, f.Roles)
	if err != nil {
		return err
	}

	return fetcher.Fetch(timeoutCtx, os.Stdout, os.Stderr)
}
