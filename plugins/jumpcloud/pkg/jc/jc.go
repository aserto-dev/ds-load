package jc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	baseURL      string = "https://console.jumpcloud.com/api"
	apiKeyHeader string = "x-api-key"
)

type JumpCloudClient struct {
	baseURL *url.URL
	apiKey  string
	headers map[string]string
	timeout time.Duration
}

func NewJumpCloudClient(ctx context.Context, apiKey string) (*JumpCloudClient, error) {
	base, _ := url.Parse(baseURL)

	c := &JumpCloudClient{
		apiKey:  apiKey,
		baseURL: base,
		headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			apiKeyHeader:   apiKey,
		},
		timeout: 30 * time.Second,
	}

	return c, nil
}

func (c *JumpCloudClient) ListDirectories(ctx context.Context) ([]any, error) {
	url := baseURL + "/v2/directories"

	var directories []any

	if err := makeHTTPRequest(ctx, url, http.MethodGet, c.headers, nil, nil, &directories); err != nil {
		return []any{}, err
	}

	return directories, nil
}

func (c *JumpCloudClient) ListUsers(ctx context.Context) ([]*User, error) {
	url := baseURL + "/search/systemusers"

	users := struct {
		Results    []*User `json:"results"`
		TotalCount int     `json:"totalCount"`
	}{}

	if err := makeHTTPRequest(ctx, url, http.MethodPost, c.headers, nil, nil, &users); err != nil {
		return []*User{}, err
	}

	lo.ForEach(users.Results, func(item *User, index int) { item.Type = TypeUser })

	return users.Results, nil
}

type GroupType int

const (
	AllGroups GroupType = iota + 1
	SystemGroups
	UserGroups
)

func (c *JumpCloudClient) ListGroups(ctx context.Context, groupType GroupType) ([]*Group, error) {
	var url string
	switch groupType {
	case AllGroups:
		url = baseURL + "/v2/groups"
	case SystemGroups:
		url = baseURL + "/v2/groups?filter=type:eq:system_group"
	case UserGroups:
		url = baseURL + "/v2/groups?filter=type:eq:user_group"
	}

	groups := []*Group{}

	if err := makeHTTPRequest(ctx, url, http.MethodGet, c.headers, nil, nil, &groups); err != nil {
		return nil, err
	}

	lo.ForEach(groups, func(item *Group, index int) { item.Name = strings.ReplaceAll(item.Name, " ", "_") })

	return groups, nil
}

const batchSize int = 10

type Members []struct {
	To struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"to"`
}

func (c *JumpCloudClient) GetUsersInGroup(ctx context.Context, groupID string) ([]*BaseUser, error) {
	u, err := url.Parse(baseURL + "/v2/usergroups/" + groupID + "/members")
	if err != nil {
		return nil, err
	}

	qv := u.Query()
	qv.Add("limit", strconv.FormatInt(int64(batchSize), 10))
	qv.Add("skip", strconv.FormatInt(0, 10))

	u.RawQuery = qv.Encode()

	members := []struct {
		To struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Attributes any    `json:"attributes"`
		}
		Attributes any `json:"attributes"`
	}{}

	idList := []string{}

	for {
		if err := makeHTTPRequest(ctx, u.String(), http.MethodGet, c.headers, nil, nil, &members); err != nil {
			return nil, err
		}

		for _, v := range members {
			idList = append(idList, v.To.ID)
		}

		if len(members) != batchSize {
			break
		}

		qv := u.Query()
		qv.Set("skip", strconv.FormatInt(int64(len(idList)), 10))
		u.RawQuery = qv.Encode()
	}

	users := []*BaseUser{}

	for _, id := range idList {
		user, err := c.GetBaseUserByID(ctx, id)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (c *JumpCloudClient) GetBaseUserByID(ctx context.Context, id string) (*BaseUser, error) {
	url := baseURL + "/Systemusers/" + id + "?limit=1&skip=0"

	user := &BaseUser{}

	if err := makeHTTPRequest(ctx, url, http.MethodGet, c.headers, nil, nil, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (c *JumpCloudClient) GetUserByID(ctx context.Context, id string) (*User, error) {
	url := baseURL + "/Systemusers/" + id + "?limit=1&skip=0"

	user := &User{}

	if err := makeHTTPRequest(ctx, url, http.MethodGet, c.headers, nil, nil, &user); err != nil {
		return nil, err
	}

	return user, nil
}

func (c *JumpCloudClient) GetUserByEmail(_ context.Context, _ string) (*User, error) {
	return nil, status.Error(codes.Unimplemented, "GetUserByEmail not implemented")
}

var (
	ErrEmptyResponse = errors.New("empty response")
	ErrStatusNotOK   = errors.New("status not OK")
)

func makeHTTPRequest[T any](ctx context.Context, reqURL, method string, headers map[string]string, queryParams url.Values, body io.Reader, resp T) error {
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

	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return err
	}

	return nil
}
