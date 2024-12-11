package app

import (
	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/openapi"
	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/verify"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	openapiClient, err := openapi.New(v.Directory, v.URL, v.IDFormat, v.ServiceName)
	if err != nil {
		return err
	}

	verifier, err := verify.New(ctx.Context, openapiClient)
	if err != nil {
		return err
	}

	return verifier.Verify(ctx.Context)
}
