package verify

import (
	"context"

	"github.com/aserto-dev/ds-load/plugins/keycloak/pkg/kc"
	"github.com/pkg/errors"
)

type Verifier struct {
	client *kc.KeycloakClient
}

func New(ctx context.Context, client *kc.KeycloakClient) (*Verifier, error) {
	return &Verifier{
		client: client,
	}, nil
}

func (v *Verifier) Verify(ctx context.Context) error {
	_, errReq := v.client.ListUsers(ctx)

	if errReq != nil {
		return errors.Wrap(errReq, "failed to retrieve users from Google")
	}

	return nil
}
