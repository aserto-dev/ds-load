package verify

import (
	"context"

	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/openapi"

	"github.com/pkg/errors"
)

type Verifier struct {
	client *openapi.Client
}

func New(ctx context.Context, client *openapi.Client) (*Verifier, error) {
	return &Verifier{
		client: client,
	}, nil
}

func (v *Verifier) Verify(ctx context.Context) error {
	_, errReq := v.client.ListServices()

	if errReq != nil {
		return errors.Wrap(errReq, "failed to retrieve Services from OpenAPI source")
	}

	return nil
}
