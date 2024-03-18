package verify

import (
	"context"

	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/pkg/errors"
)

type Verifier struct {
	client *oktaclient.OktaClient
}

func New(ctx context.Context, client *oktaclient.OktaClient) (*Verifier, error) {
	return &Verifier{
		client: client,
	}, nil

}

func (v *Verifier) Verify(ctx context.Context) error {
	filter := query.NewQueryParams(query.WithLimit(1))
	_, _, err := v.client.User.ListUsers(ctx, filter)

	if err != nil {
		return errors.Wrap(err, "failed to retrieve user from Okta")
	}

	return nil
}
