package verify

import (
	"context"

	"github.com/aserto-dev/ds-load/plugins/azureadb2c/pkg/azureclient"
	"github.com/pkg/errors"
)

type Verifier struct {
	client *azureclient.AzureADClient
	Groups bool
	B2C    bool
}

func New(ctx context.Context, client *azureclient.AzureADClient) (*Verifier, error) {
	return &Verifier{
		client: client,
	}, nil

}

func (v *Verifier) WithGroups(groups bool) *Verifier {
	v.Groups = groups
	return v
}

func (f *Verifier) WithB2C(b2c bool) *Verifier {
	f.B2C = b2c
	return f
}

func (v *Verifier) Verify(ctx context.Context) error {
	_, errReq := v.client.ListUsers(ctx, v.Groups)

	if errReq != nil {
		return errors.Wrap(errReq, "failed to retrieve users from AzureAD")
	}

	return nil
}
