package app

import (
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/openapiclient"
	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/verify"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	openapiClient, err := openapiclient.NewOpenAPIClient(v.Directory, v.URL)
	if err != nil {
		return err
	}

	verifier, err := verify.New(ctx.Context, openapiClient)
	if err != nil {
		return err
	}

	return verifier.Verify(ctx.Context)
}
