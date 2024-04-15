package verify

import (
	"context"

	"github.com/aserto-dev/ds-load/plugins/cognito/pkg/cognitoclient"
	"github.com/pkg/errors"
)

type Verifier struct {
	client *cognitoclient.CognitoClient
}

func New(ctx context.Context, client *cognitoclient.CognitoClient) (*Verifier, error) {
	return &Verifier{
		client: client,
	}, nil

}

func (v *Verifier) Verify(ctx context.Context) error {
	_, errReq := v.client.ListUsers(ctx)

	if errReq != nil {
		return errors.Wrap(errReq, "failed to retrieve users from Cognito")
	}

	return nil
}
