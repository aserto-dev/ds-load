package verify

import (
	"context"

	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/openapiclient"

	"github.com/pkg/errors"
)

type Verifier struct {
	client *openapiclient.OpenAPIClient
}

func New(ctx context.Context, client *openapiclient.OpenAPIClient) (*Verifier, error) {
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
