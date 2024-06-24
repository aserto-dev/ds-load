package openapiclient

import (
	"context"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPIClient struct {
	openapiDocument *openapi3.T
}

type API struct {
	Method string
	Path   string
}

func NewOpenAPIClient(directory, specurl string) (*OpenAPIClient, error) {
	c := &OpenAPIClient{}

	specURL, err := url.Parse(specurl)
	if err != nil {
		return nil, err
	}

	doc, err := openapi3.NewLoader().LoadFromURI(specURL)
	if err != nil {
		return nil, err
	}

	c.openapiDocument = doc
	return c, nil
}

func (c *OpenAPIClient) ListAPIs(ctx context.Context) ([]API, error) {
	apis := make([]API, 0)

	for uri, path := range c.openapiDocument.Paths.Map() {

		if path.Get != nil {
			api := &API{}
			api.Method = "GET"
			api.Path = uri
			apis = append(apis, *api)
		}
		if path.Post != nil {
			api := &API{}
			api.Method = "POST"
			api.Path = uri
			apis = append(apis, *api)
		}
		if path.Put != nil {
			api := &API{}
			api.Method = "PUT"
			api.Path = uri
			apis = append(apis, *api)
		}
		if path.Delete != nil {
			api := &API{}
			api.Method = "DELETE"
			api.Path = uri
			apis = append(apis, *api)
		}
		if path.Options != nil {
			api := &API{}
			api.Method = "OPTIONS"
			api.Path = uri
			apis = append(apis, *api)
		}
	}

	return apis, nil
}
