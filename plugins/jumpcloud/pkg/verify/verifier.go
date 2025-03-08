package verify

import (
	"context"

	"github.com/aserto-dev/ds-load/plugins/jumpcloud/pkg/jcclient"
	"github.com/pkg/errors"
)

type Verifier struct {
	client *jcclient.JumpCloudClient
}

func New(ctx context.Context, client *jcclient.JumpCloudClient) (*Verifier, error) {
	return &Verifier{
		client: client,
	}, nil
}

func (v *Verifier) Verify(ctx context.Context) error {
	_, errReq := v.client.ListUsers()

	if errReq != nil {
		return errors.Wrap(errReq, "failed to retrieve users from Google")
	}

	return nil
}
