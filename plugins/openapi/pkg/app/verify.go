package app

import (
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/openapi"
	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/verify"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	openapiClient, err := openapi.New(v.Directory, v.URL, v.IDFormat)
	if err != nil {
		return err
	}

	verifier, err := verify.New(ctx.Context, openapiClient)
	if err != nil {
		return err
	}

	return verifier.Verify(ctx.Context)
}
