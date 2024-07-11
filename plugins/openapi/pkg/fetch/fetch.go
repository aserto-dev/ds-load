package fetch

import (
	"context"
	"io"

	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/openapi"
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

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	services, err := f.client.ListServices()
	if err != nil {
		_, _ = errorWriter.Write([]byte(err.Error()))
	}

	for _, service := range services {
		err = writer.Write(service)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
		}
	}

	apis, err := f.client.ListAPIs()
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
