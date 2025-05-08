package kc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	defaultConnectionTimeout = 30 * time.Second
)

type KeycloakClientConfig struct {
	TokenURL     string `name:"token-url" short:"u" env:"KEYCLOAK_TOKEN_URL" help:"keycloak token URL" required:""`
	ClientID     string `name:"client-id" short:"i" env:"KEYCLOAK_CLIENT_ID" help:"keycloak client id" required:""`
	ClientSecret string `name:"client-secret" short:"s" env:"KEYCLOAK_CLIENT_SECRET" help:"keycloak client secret" required:""`
}

type KeycloakClient struct {
	config   *KeycloakClientConfig
	token    *oauth2.Token
	adminURL string
	headers  map[string]string
	timeout  time.Duration
}

func NewKeycloakClient(ctx context.Context, cfg *KeycloakClientConfig) (*KeycloakClient, error) {
	token, err := accessToken(ctx, cfg)
	if err != nil {
		return nil, err
	}

	issuer, err := extractIssuerFromToken(token)
	if err != nil {
		return nil, err
	}

	base, err := url.Parse(issuer)
	if err != nil {
		return nil, err
	}

	adminURL := url.URL{
		Scheme: base.Scheme,
		Host:   base.Host,
		Path:   "/admin" + base.Path,
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"Authorization": "Bearer " + token.AccessToken,
	}

	return &KeycloakClient{
		config:   cfg,
		token:    token,
		adminURL: adminURL.String(),
		headers:  headers,
		timeout:  defaultConnectionTimeout,
	}, nil
}

func accessToken(ctx context.Context, cfg *KeycloakClientConfig) (*oauth2.Token, error) {
	cc := clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     cfg.TokenURL,
	}

	return cc.Token(ctx)
}

func extractIssuerFromToken(token *oauth2.Token) (string, error) {
	accessToken := token.AccessToken

	// Parse without verifying the signature
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	claims := jwt.MapClaims{}

	_, _, err := parser.ParseUnverified(accessToken, claims)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse token")
	}

	iss, ok := claims["iss"].(string)
	if !ok {
		return "", errors.Errorf("issuer not found or not a string")
	}

	return iss, nil
}

func (c *KeycloakClient) ListUsers(ctx context.Context) ([]*User, error) {
	url := c.adminURL + "/users"

	users := []*User{}

	if err := makeHTTPRequest(ctx, url, http.MethodGet, c.headers, nil, nil, &users); err != nil {
		return []*User{}, err
	}

	return users, nil
}

func (c *KeycloakClient) ListGroups(ctx context.Context) ([]*Group, error) {
	url := c.adminURL + "/groups"

	groups := []*Group{}

	if err := makeHTTPRequest(ctx, url, http.MethodGet, c.headers, nil, nil, &groups); err != nil {
		return []*Group{}, err
	}

	return groups, nil
}

// realm roles.
func (c *KeycloakClient) ListRoles(ctx context.Context) ([]*Role, error) {
	url := c.adminURL + "/roles"

	roles := []*Role{}

	if err := makeHTTPRequest(ctx, url, http.MethodGet, c.headers, nil, nil, &roles); err != nil {
		return []*Role{}, err
	}

	return roles, nil
}

// get roles of user.
func (c *KeycloakClient) GetRolesOfUser(ctx context.Context, id string) (*RealmMappings, error) {
	url := c.adminURL + "/users/" + id + "/role-mappings"

	realmMappings := &RealmMappings{}

	if err := makeHTTPRequest(ctx, url, http.MethodGet, c.headers, nil, nil, &realmMappings); err != nil {
		return &RealmMappings{}, err
	}

	return realmMappings, nil
}

// get roles of group.
func (c *KeycloakClient) GetRolesOfGroup(ctx context.Context, id string) (*RealmMappings, error) {
	url := c.adminURL + "/groups/" + id + "/role-mappings"

	realmMappings := &RealmMappings{}

	if err := makeHTTPRequest(ctx, url, http.MethodGet, c.headers, nil, nil, &realmMappings); err != nil {
		return &RealmMappings{}, err
	}

	return realmMappings, nil
}

// get users of role.
func (c *KeycloakClient) GetUsersOfRole(ctx context.Context, role string) ([]*User, error) {
	url := c.adminURL + "/roles/" + role + "/users"

	users := []*User{}

	if err := makeHTTPRequest(ctx, url, http.MethodGet, c.headers, nil, nil, &users); err != nil {
		return []*User{}, err
	}

	return users, nil
}

// get users of group.
func (c *KeycloakClient) GetUsersOfGroup(ctx context.Context, id string) ([]*User, error) {
	url := c.adminURL + "/groups/" + id + "/members"

	users := []*User{}

	if err := makeHTTPRequest(ctx, url, http.MethodGet, c.headers, nil, nil, &users); err != nil {
		return []*User{}, err
	}

	return users, nil
}

// get groups of user.
func (c *KeycloakClient) GetGroupsOfUser(ctx context.Context, id string) ([]*Group, error) {
	url := c.adminURL + "/users/" + id + "/groups"

	groups := []*Group{}

	if err := makeHTTPRequest(ctx, url, http.MethodGet, c.headers, nil, nil, &groups); err != nil {
		return []*Group{}, err
	}

	return groups, nil
}

var (
	ErrEmptyResponse = errors.New("empty response")
	ErrStatusNotOK   = errors.New("status not OK")
)

func makeHTTPRequest[T any](
	ctx context.Context,
	reqURL, method string,
	headers map[string]string,
	queryParams url.Values,
	body io.Reader,
	resp T,
) error {
	client := http.Client{}

	u, err := url.Parse(reqURL)
	if err != nil {
		return err
	}

	if method == http.MethodGet && queryParams != nil {
		q := u.Query()

		for k, v := range queryParams {
			q.Set(k, strings.Join(v, ","))
		}

		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res == nil {
		return errors.Wrapf(ErrEmptyResponse, "req %s", u.String())
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.Wrapf(ErrStatusNotOK, "req: %s status: %s response: %s", u.String(), res.Status, buf)
	}

	if err := json.Unmarshal(buf, &resp); err != nil {
		return err
	}

	return nil
}
