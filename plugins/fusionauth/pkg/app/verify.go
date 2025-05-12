package app

import (
	"github.com/aserto-dev/ds-load/plugins/fusionauth/pkg/client"
	"github.com/aserto-dev/ds-load/plugins/fusionauth/pkg/verify"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	fusionauthClient, err := client.NewFusionAuthClient(v.HostURL, v.APIKey)
	if err != nil {
		return err
	}

	verifier, err := verify.New(ctx.Context, fusionauthClient)
	if err != nil {
		return err
	}

	return verifier.Verify(ctx.Context)
}
