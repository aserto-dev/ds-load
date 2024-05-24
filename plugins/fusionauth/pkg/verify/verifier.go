package verify

import (
	"context"

	"github.com/aserto-dev/ds-load/plugins/fusionauth/pkg/fusionauthclient"
	"github.com/pkg/errors"
)

type Verifier struct {
	client *fusionauthclient.FusionAuthClient
}

func New(ctx context.Context, client *fusionauthclient.FusionAuthClient) (*Verifier, error) {
	return &Verifier{
		client: client,
	}, nil

}

func (v *Verifier) Verify(ctx context.Context) error {
	_, errReq := v.client.ListUsers(ctx)

	if errReq != nil {
		return errors.Wrap(errReq, "failed to retrieve users from FusionAuth")
	}

	return nil
}
