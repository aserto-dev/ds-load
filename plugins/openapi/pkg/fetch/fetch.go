package fetch

import (
	"context"
	"io"

	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/openapi"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
)

type Fetcher struct {
	client      *openapi.Client
	directory   string
	idFormat    string
	specURL     string
	serviceName string
}

func New(client *openapi.Client) (*Fetcher, error) {
	return &Fetcher{
		client: client,
	}, nil
}

func (f *Fetcher) WithDirectory(directory string) *Fetcher {
	f.directory = directory
	return f
}

func (f *Fetcher) WithURL(url string) *Fetcher {
	f.specURL = url
	return f
}

func (f *Fetcher) WithIDFormat(idFormat string) *Fetcher {
	f.idFormat = idFormat
	return f
}

func (f *Fetcher) WithServiceName(serviceName string) *Fetcher {
	f.serviceName = serviceName
	return f
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter io.Writer, errorWriter common.ErrorWriter) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	services, err := f.client.ListServices()
	errorWriter.Error(err)

	for _, service := range services {
		err := writer.Write(service)
		errorWriter.Error(err)
	}

	apis, err := f.client.ListAPIs()
	errorWriter.Error(err)

	for _, api := range apis {
		err := writer.Write(api)
		errorWriter.Error(err)
	}

	return nil
}
