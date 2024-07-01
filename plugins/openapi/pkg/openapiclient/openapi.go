package openapiclient

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
)

type OpenAPIClient struct {
	docs     []*openapi3.T
	idFormat string
}

type API struct {
	Type        string `json:"type"`
	Service     string `json:"service"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	DisplayName string `json:"displayName"`
	ID          string `json:"id"`
	ServiceID   string `json:"serviceID"`
}

type Service struct {
	Type        string `json:"type"`
	DisplayName string `json:"displayName"`
	ID          string `json:"id"`
}

func NewOpenAPIClient(directory, specurl string, idFormat string) (*OpenAPIClient, error) {
	c := &OpenAPIClient{}
	c.idFormat = idFormat
	c.docs = make([]*openapi3.T, 0)

	if specurl != "" {
		specURL, err := url.Parse(specurl)
		if err != nil {
			return nil, errors.Wrapf(err, "url not parsed: %s", specurl)
		}
		doc, err := openapi3.NewLoader().LoadFromURI(specURL)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot load OpenAPI spec from URL : %s", specurl)
		}
		c.docs = append(c.docs, doc)
	}

	if directory != "" {
		if _, err := os.Stat(directory); errors.Is(err, os.ErrNotExist) {
			return nil, errors.Wrapf(err, "directory not found: %s", directory)
		}
		files, err := os.ReadDir(directory)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot read directory: %s", directory)
		}
		for _, file := range files {
			if !file.IsDir() {
				filename := fmt.Sprintf("%s/%s", directory, file.Name())
				doc, err := openapi3.NewLoader().LoadFromFile(filename)
				if err != nil {
					return nil, errors.Wrapf(err, "cannot open file: %s", file.Name())
				}
				c.docs = append(c.docs, doc)
			}
		}
	}

	return c, nil
}

func (c *OpenAPIClient) ListServices() ([]Service, error) {
	services := make([]Service, 0)
	for _, service := range c.docs {
		svc := newService(service.Info.Title, c.idFormat)
		services = append(services, *svc)
	}
	return services, nil
}

func (c *OpenAPIClient) ListAPIs() ([]API, error) {
	apis := make([]API, 0)
	for _, service := range c.docs {
		apiList := c.ListAPIsInService(service, c.idFormat)
		apis = append(apis, apiList...)
	}
	return apis, nil
}

func (c *OpenAPIClient) ListAPIsInService(service *openapi3.T, idFormat string) []API {

	apis := make([]API, 0)
	for pathKey, pathItem := range service.Paths.Map() {

		if pathItem.Get != nil {
			api := newAPI(service.Info.Title, "GET", pathKey, idFormat)
			apis = append(apis, *api)
		}
		if pathItem.Post != nil {
			api := newAPI(service.Info.Title, "POST", pathKey, idFormat)
			apis = append(apis, *api)
		}
		if pathItem.Put != nil {
			api := newAPI(service.Info.Title, "PUT", pathKey, idFormat)
			apis = append(apis, *api)
		}
		if pathItem.Patch != nil {
			api := newAPI(service.Info.Title, "PATCH", pathKey, idFormat)
			apis = append(apis, *api)
		}
		if pathItem.Delete != nil {
			api := newAPI(service.Info.Title, "DELETE", pathKey, idFormat)
			apis = append(apis, *api)
		}
		if pathItem.Options != nil {
			api := newAPI(service.Info.Title, "OPTIONS", pathKey, idFormat)
			apis = append(apis, *api)
		}
	}

	return apis
}

func newService(name string, idFormat string) *Service {
	service := &Service{}
	service.DisplayName = name
	service.Type = "service"
	service.ID = canonicalizeServiceName(name, idFormat)
	return service
}

func newAPI(service, method, path string, idFormat string) *API {
	api := &API{}
	api.Type = "endpoint"
	api.Service = service
	api.Method = method
	api.Path = path
	api.DisplayName = fmt.Sprintf("%s %s", method, path)
	api.ServiceID = canonicalizeServiceName(api.Service, idFormat)
	api.ID = canonicalizeEndpoint(api.ServiceID, method, path, idFormat)
	return api
}

func canonicalizePath(uri string) string {
	parts := strings.Split(uri, "/")
	return strings.Join(parts[1:], ".")
}

func canonicalizeEndpoint(service, method, path string, idFormat string) string {
	parts := []string{service, method}
	switch idFormat {
	case "canonical":
		parts = append(parts, canonicalizePath(path))
		return strings.Join(parts, ":")
	case "base64":
		parts = append(parts, path)
		return base64.StdEncoding.EncodeToString([]byte(strings.Join(parts, "")))
	default:
		parts = append(parts, path)
		return strings.Join(parts, ":")
	}
}

func canonicalizeServiceName(serviceName string, idFormat string) string {
	switch idFormat {
	case "canonical":
		return strings.Replace(serviceName, " ", "_", -1)
	case "base64":
		return base64.StdEncoding.EncodeToString([]byte(serviceName))
	default:
		return strings.Replace(serviceName, " ", "_", -1)
	}
}
