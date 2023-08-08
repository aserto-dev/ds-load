package app

import (
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/verify"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	oktaClient, err := oktaclient.NewOktaClient(ctx.Context, v.Domain, v.APIToken, v.RequestTimeout)
	if err != nil {
		return err
	}

	verifier, err := verify.New(ctx.Context, oktaClient)
	if err != nil {
		return err
	}

	return verifier.Verify(ctx.Context)
}
