package verify

import (
	"context"

	"github.com/aserto-dev/ds-load/plugins/azuread/pkg/azureclient"
	"github.com/pkg/errors"
)

type Verifier struct {
	client *azureclient.AzureADClient
}

func New(ctx context.Context, client *azureclient.AzureADClient) (*Verifier, error) {
	return &Verifier{
		client: client,
	}, nil

}

func (v *Verifier) Verify(ctx context.Context) error {
	_, errReq := v.client.ListUsers(ctx)

	if errReq != nil {
		return errors.Wrap(errReq, "failed to retrieve users from AzureAD")
	}

	return nil
}
