package verify

import (
	"context"

	"github.com/aserto-dev/ds-load/plugins/jumpcloud/pkg/client"
	"github.com/pkg/errors"
)

type Verifier struct {
	client *client.JumpCloudClient
}

func New(ctx context.Context, client *client.JumpCloudClient) (*Verifier, error) {
	return &Verifier{
		client: client,
	}, nil
}

func (v *Verifier) Verify(ctx context.Context) error {
	_, errReq := v.client.ListUsers(ctx)

	if errReq != nil {
		return errors.Wrap(errReq, "failed to retrieve users from JumpCloud")
	}

	return nil
}
