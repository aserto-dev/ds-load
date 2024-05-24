package app

import (
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/fusionauth/pkg/fusionauthclient"
	"github.com/aserto-dev/ds-load/plugins/fusionauth/pkg/verify"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	fusionauthClient, err := fusionauthclient.NewFusionAuthClient(v.HostURL, v.APIKey)
	if err != nil {
		return err
	}

	verifier, err := verify.New(ctx.Context, fusionauthClient)
	if err != nil {
		return err
	}

	return verifier.Verify(ctx.Context)
}
