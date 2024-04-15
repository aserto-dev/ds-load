package verify

import (
	"context"

	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/auth0client"
	"github.com/auth0/go-auth0/management"
	"github.com/pkg/errors"
)

type Verifier struct {
	client *auth0client.Auth0Client
}

func New(ctx context.Context, client *auth0client.Auth0Client) (*Verifier, error) {
	return &Verifier{
		client: client,
	}, nil

}

func (v *Verifier) Verify(ctx context.Context) error {
	_, err := v.client.Mgmt.User.List(
		ctx,
		management.Page(0),
		management.PerPage(1),
	)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Auth0 and list users")
	}

	return nil
}
