package app

import (
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/auth0client"
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/verify"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	client, err := auth0client.New(ctx.Context, v.ClientID, v.ClientSecret, v.Domain)
	if err != nil {
		return err
	}

	verifier, err := verify.New(ctx.Context, client)
	if err != nil {
		return err
	}

	return verifier.Verify(ctx.Context)
}
