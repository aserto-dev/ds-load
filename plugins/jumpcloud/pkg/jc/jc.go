package jc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

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
	}

	return c, nil
}

func (c *JumpCloudClient) ListDirectories() ([]any, error) {
	url := baseURL + "/v2/directories"

	var directories []any
	resp, err := makeHTTPRequest(url, http.MethodGet, c.headers, nil, nil, directories)
	if err != nil {
		return []any{}, err
	}

	return resp, nil
}

func (c *JumpCloudClient) ListUsers() ([]*User, error) {
	url := baseURL + "/search/systemusers"

	users := struct {
		Results    []*User `json:"results"`
		TotalCount int     `json:"totalCount"`
	}{}

	resp, err := makeHTTPRequest(url, http.MethodPost, c.headers, nil, nil, users)
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}

type GroupType int

const (
	AllGroups GroupType = iota + 1
	SystemGroups
	UserGroups
)

func (c *JumpCloudClient) ListGroups(groupType GroupType) ([]*Group, error) {
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
	resp, err := makeHTTPRequest(url, http.MethodGet, c.headers, nil, nil, groups)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JumpCloudClient) GetUsersInGroup(groupID string) ([]*User, error) {
	url := baseURL + "/v2/usergroups/" + groupID + "/members"

	members := []struct {
		To struct {
			Id         string `json:"id"`
			Type       string `json:"type"`
			Attributes any    `json:"attributes"`
		}
		Attributes any `json:"attributes"`
	}{}

	resp, err := makeHTTPRequest(url, http.MethodGet, c.headers, nil, nil, members)
	if err != nil {
		return nil, err
	}

	users := []*User{}
	for _, member := range resp {
		user, err := c.GetUserByID(member.To.Id)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (c *JumpCloudClient) GetUserByID(id string) (*User, error) {
	url := baseURL + "/Systemusers/" + id + "?limit=1000&skip=0"

	var user User
	resp, err := makeHTTPRequest(url, http.MethodGet, c.headers, nil, nil, &user)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *JumpCloudClient) GetUserByEmail(email string) (*User, error) {
	return nil, status.Error(codes.Unimplemented, "GetUserByEmail not implemented")
}

func makeHTTPRequest[T any](fullUrl string, method string, headers map[string]string, queryParameters url.Values, body io.Reader, responseType T) (T, error) {
	client := http.Client{}
	u, err := url.Parse(fullUrl)
	if err != nil {
		return responseType, err
	}

	if method == http.MethodGet && queryParameters != nil {
		q := u.Query()

		for k, v := range queryParameters {
			q.Set(k, strings.Join(v, ","))
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return responseType, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		return responseType, err
	}

	if res == nil {
		return responseType, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return responseType, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return responseType, fmt.Errorf("error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData)
	}

	var responseObject T
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return responseType, err
	}

	return responseObject, nil
}
