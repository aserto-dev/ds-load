package app

import (
	"github.com/aserto-dev/ds-load/plugins/azuread/pkg/verify"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	client, err := createAzureAdClient(ctx.Context, v.Tenant, v.ClientID, v.ClientSecret, v.RefreshToken)
	if err != nil {
		return err
	}

	verifier, err := verify.New(ctx.Context, client)
	if err != nil {
		return err
	}

	return verifier.WithGroups(v.Groups).Verify(ctx.Context)
}
