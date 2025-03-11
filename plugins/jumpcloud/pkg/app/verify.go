package app

import (
	"github.com/aserto-dev/ds-load/plugins/jumpcloud/pkg/jc"
	"github.com/aserto-dev/ds-load/plugins/jumpcloud/pkg/verify"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	gClient, err := jc.NewJumpCloudClient(ctx.Context, v.APIKey)
	if err != nil {
		return err
	}

	verifier, err := verify.New(ctx.Context, gClient)
	if err != nil {
		return err
	}

	return verifier.Verify(ctx.Context)
}
