package verify

import (
	"context"

	"github.com/aserto-dev/ds-load/plugins/google/pkg/googleclient"
	"github.com/pkg/errors"
)

type Verifier struct {
	client *googleclient.GoogleClient
}

func New(ctx context.Context, client *googleclient.GoogleClient) (*Verifier, error) {
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
