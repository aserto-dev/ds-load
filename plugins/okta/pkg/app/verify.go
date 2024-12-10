package app

import (
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/verify"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	oktaClient, err := oktaclient.NewOktaClient(v.Domain, v.APIToken, v.RequestTimeout)
	if err != nil {
		return err
	}

	verifier, err := verify.New(ctx.Context, oktaClient)
	if err != nil {
		return err
	}

	return verifier.Verify(ctx.Context)
}
