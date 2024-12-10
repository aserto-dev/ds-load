package app

import (
	"github.com/aserto-dev/ds-load/plugins/google/pkg/googleclient"
	"github.com/aserto-dev/ds-load/plugins/google/pkg/verify"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	gClient, err := googleclient.NewGoogleClient(ctx.Context, v.ClientID, v.ClientSecret, v.RefreshToken, v.Customer)
	if err != nil {
		return err
	}

	verifier, err := verify.New(ctx.Context, gClient)
	if err != nil {
		return err
	}

	return verifier.Verify(ctx.Context)
}
