package openapi

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
)

const (
	Base64    = "base64"
	Canonical = "canonical"
)

type Client struct {
	docs     []*openapi3.T
	idFormat string
}

type API struct {
	Type        string `json:"type"`
	Service     string `json:"service"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	DisplayName string `json:"displayName"` //nolint:tagliatelle // keep open api format
	ID          string `json:"id"`
	ServiceID   string `json:"serviceID"` //nolint:tagliatelle // keep open api format
}

type Service struct {
	Type        string `json:"type"`
	DisplayName string `json:"displayName"` //nolint:tagliatelle // keep open api format
	ID          string `json:"id"`
}

func New(directory, specURL, idFormat, serviceName string) (*Client, error) {
	c := &Client{}
	c.idFormat = idFormat
	c.docs = make([]*openapi3.T, 0)

	if specURL == "" {
		return nil, errors.Errorf("spec URL must be set")
	}

	parsedURL, err := url.Parse(specURL)
	if err != nil {
		return nil, errors.Wrapf(err, "url not parsed: %s", specURL)
	}

	doc, err := openapi3.NewLoader().LoadFromURI(parsedURL)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot load OpenAPI spec from URL : %s", specURL)
	}

	if serviceName != "" {
		if doc.Info.Extensions == nil {
			doc.Info.Extensions = make(map[string]any, 0)
		}

		doc.Info.Extensions["ServiceName"] = canonicalizeServiceName(serviceName, Canonical)
	}

	c.docs = append(c.docs, doc)

	if directory == "" {
		return nil, errors.Errorf("directory must be set")
	}

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

	return c, nil
}

func (c *Client) ListServices() ([]Service, error) {
	services := make([]Service, 0)

	for _, service := range c.docs {
		id := ""

		if service.Info.Extensions["ServiceName"] != "" {
			serviceName, ok := service.Info.Extensions["ServiceName"].(string)

			if !ok {
				return services, errors.Errorf("service name cannot be cast to string")
			}

			id = serviceName
		}

		svc := newService(service.Info.Title, id, c.idFormat)
		services = append(services, *svc)
	}

	return services, nil
}

func (c *Client) ListAPIs() ([]API, error) {
	apis := make([]API, 0)

	for _, service := range c.docs {
		apiList := c.ListAPIsInService(service, c.idFormat)
		apis = append(apis, apiList...)
	}

	return apis, nil
}

func (c *Client) ListAPIsInService(service *openapi3.T, idFormat string) []API {
	apis := make([]API, 0)

	serviceID, ok := service.Info.Extensions["ServiceName"].(string)
	if serviceID == "" || !ok {
		serviceID = service.Info.Title
	}

	for pathKey, pathItem := range service.Paths.Map() {
		if pathItem.Get != nil {
			api := newAPI(serviceID, service.Info.Title, "GET", pathKey, idFormat)
			apis = append(apis, *api)
		}

		if pathItem.Post != nil {
			api := newAPI(serviceID, service.Info.Title, "POST", pathKey, idFormat)
			apis = append(apis, *api)
		}

		if pathItem.Put != nil {
			api := newAPI(serviceID, service.Info.Title, "PUT", pathKey, idFormat)
			apis = append(apis, *api)
		}

		if pathItem.Patch != nil {
			api := newAPI(serviceID, service.Info.Title, "PATCH", pathKey, idFormat)
			apis = append(apis, *api)
		}

		if pathItem.Delete != nil {
			api := newAPI(serviceID, service.Info.Title, "DELETE", pathKey, idFormat)
			apis = append(apis, *api)
		}

		if pathItem.Options != nil {
			api := newAPI(serviceID, service.Info.Title, "OPTIONS", pathKey, idFormat)
			apis = append(apis, *api)
		}
	}

	return apis
}

func newService(name, id, idFormat string) *Service {
	service := &Service{}
	service.DisplayName = name
	service.Type = "service"
	service.ID = id

	if id == "" {
		service.ID = canonicalizeServiceName(name, idFormat)
	}

	return service
}

func newAPI(serviceID, serviceName, method, path, idFormat string) *API {
	api := &API{}
	api.Type = "endpoint"
	api.ServiceID = serviceID
	api.Service = serviceName
	api.Method = method
	api.Path = path
	api.DisplayName = fmt.Sprintf("%s %s", method, path)
	api.ID = canonicalizeEndpoint(api.ServiceID, method, path, idFormat)

	return api
}

func canonicalizePath(uri string) string {
	parts := strings.Split(uri, "/")
	return strings.Join(parts[1:], ".")
}

func canonicalizeEndpoint(service, method, path, idFormat string) string {
	parts := []string{service, method}

	switch idFormat {
	case Base64:
		parts = append(parts, path)
		return base64.StdEncoding.EncodeToString([]byte(strings.Join(parts, "")))
	case Canonical:
		parts = append(parts, canonicalizePath(path))
		return strings.Join(parts, ":")
	default:
		parts = append(parts, path)
		return strings.Join(parts, ":")
	}
}

func canonicalizeServiceName(serviceName, idFormat string) string {
	switch idFormat {
	case Base64:
		return base64.StdEncoding.EncodeToString([]byte(serviceName))
	case Canonical:
		return strings.ReplaceAll(serviceName, " ", "_")
	default:
		return strings.ReplaceAll(serviceName, " ", "_")
	}
}
