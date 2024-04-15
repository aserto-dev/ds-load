package app

import (
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/cognito/pkg/cognitoclient"
	"github.com/aserto-dev/ds-load/plugins/cognito/pkg/verify"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	cognitoClient, err := cognitoclient.NewCognitoClient(v.AccessKey, v.SecretKey, v.UserPoolID, v.Region)
	if err != nil {
		return err
	}

	verifier, err := verify.New(ctx.Context, cognitoClient)
	if err != nil {
		return err
	}

	return verifier.Verify(ctx.Context)
}
