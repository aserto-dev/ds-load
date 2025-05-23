package app

import (
	"github.com/aserto-dev/ds-load/plugins/keycloak/pkg/kc"
	"github.com/aserto-dev/ds-load/plugins/keycloak/pkg/verify"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	gClient, err := kc.NewKeycloakClient(ctx.Context, &v.KeycloakClientConfig)
	if err != nil {
		return err
	}

	verifier, err := verify.New(ctx.Context, gClient)
	if err != nil {
		return err
	}

	return verifier.Verify(ctx.Context)
}
