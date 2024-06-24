package fetch

import (
	"context"
	"io"

	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/openapiclient"
	"github.com/aserto-dev/ds-load/sdk/common/js"
)

type Fetcher struct {
	openapiClient *openapiclient.OpenAPIClient
	directory     string
	specurl       string
}

func New(client *openapiclient.OpenAPIClient) (*Fetcher, error) {
	return &Fetcher{
		openapiClient: client,
	}, nil
}

func (f *Fetcher) WithDirectory(directory string) *Fetcher {
	f.directory = directory
	return f
}

func (f *Fetcher) WithURL(url string) *Fetcher {
	f.specurl = url
	return f
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	apis, err := f.openapiClient.ListAPIs(ctx)
	if err != nil {
		_, _ = errorWriter.Write([]byte(err.Error()))
	}

	for _, api := range apis {
		err = writer.Write(api)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
		}
	}

	return nil
}
