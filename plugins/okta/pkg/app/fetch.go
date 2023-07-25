package app

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/plugin"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

type FetchCmd struct {
	Domain         string `env:"DS_OKTA_DOMAIN"`
	APIToken       string `env:"DS_OKTA_TOKEN"`
	Groups         bool   `env:"DS_OKTA_GROUPS" default:"true" negatable:""`
	Roles          bool   `env:"DS_OKTA_ROLES" default:"true" negatable:""`
	RequestTimeout int64  `default:"0" optional:""`

	oktaClient oktaclient.OktaClient `kong:"-"`
}

func (fetcher *FetchCmd) Run(kongCtx *kong.Context) error {
	ctx := context.Background()
	oktaClient, err := oktaclient.NewOktaClient(ctx, fetcher.Domain, fetcher.APIToken, fetcher.RequestTimeout)
	if err != nil {
		return err
	}
	fetcher.oktaClient = oktaClient
	results := make(chan map[string]interface{}, 1)
	errors := make(chan error, 1)
	go func() {
		fetcher.Fetch(ctx, results, errors)
		close(results)
		close(errors)
	}()
	return plugin.NewDSPlugin().WriteFetchOutput(results, errors)
}

func (fetcher *FetchCmd) Fetch(ctx context.Context, results chan map[string]interface{}, errors chan error) {
	var response *okta.Response
	var users []*okta.User
	var err error

	users, response, err = fetcher.oktaClient.ListUsers(ctx, nil)
	if err != nil {
		errors <- err
		common.SetExitCode(1)
		return
	}

	for {
		err = handleResponse(response)
		if err != nil {
			errors <- err
			common.SetExitCode(1)
			return
		}

		for _, user := range users {
			userResult, err := fetcher.processUser(ctx, user)
			if err != nil {
				errors <- err
				common.SetExitCode(1)
			}
			results <- userResult
		}

		if response != nil && response.HasNextPage() {
			response, err = response.Next(ctx, &users)
			if err != nil {
				errors <- err
				common.SetExitCode(1)
			}
		} else {
			break
		}
	}
}

func (fetcher *FetchCmd) processUser(ctx context.Context, user *okta.User) (map[string]interface{}, error) {
	userBytes, err := json.Marshal(user)
	if err != nil {
		common.SetExitCode(1)
		return nil, err
	}
	var userResult map[string]interface{}
	err = json.Unmarshal(userBytes, &userResult)
	if err != nil {
		common.SetExitCode(1)
		return nil, err
	}

	if fetcher.Groups {
		// Write all user groups
		groups, err := fetcher.getGroups(ctx, user.Id)
		if err != nil {
			common.SetExitCode(1)
			return nil, err
		}
		userResult["groups"] = groups
	}

	if fetcher.Roles {
		// Write all user roles
		roles, err := fetcher.getRoles(ctx, user.Id)
		if err != nil {
			common.SetExitCode(1)
			return nil, err
		}
		userResult["roles"] = roles
	}
	return userResult, nil
}

func (fetcher *FetchCmd) getGroups(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	var response *okta.Response
	var result []map[string]interface{}
	var groups []*okta.Group
	var err error

	groups, response, err = fetcher.oktaClient.ListUserGroups(ctx, userID)
	if err != nil {
		return nil, err
	}

	for {
		err = handleResponse(response)
		if err != nil {
			return nil, err
		}
		for _, group := range groups {
			groupBytes, err := json.Marshal(group)
			if err != nil {
				return nil, err
			}
			var obj map[string]interface{}
			err = json.Unmarshal(groupBytes, &obj)
			if err != nil {
				return nil, err
			}
			result = append(result, obj)
		}

		if response != nil && response.HasNextPage() {
			response, err = response.Next(ctx, &groups)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return result, nil
}

func (fetcher *FetchCmd) getRoles(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	var response *okta.Response
	var result []map[string]interface{}
	var roles []*okta.Role
	var err error

	roles, response, err = fetcher.oktaClient.ListAssignedRolesForUser(ctx, userID, nil)
	if err != nil {
		return nil, err
	}

	for {
		err = handleResponse(response)
		if err != nil {
			return nil, err
		}

		for _, role := range roles {
			roleBytes, err := json.Marshal(role)
			if err != nil {
				return nil, err
			}
			var obj map[string]interface{}
			err = json.Unmarshal(roleBytes, &obj)
			if err != nil {
				return nil, err
			}
			result = append(result, obj)
		}

		if response != nil && response.HasNextPage() {
			response, err = response.Next(ctx, &roles)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return result, nil
}

func handleResponse(resp *okta.Response) error {
	if resp.Response != nil && resp.StatusCode == http.StatusTooManyRequests {
		response, err := json.Marshal(resp) //nolint:staticcheck
		if err != nil {
			return err
		}
		os.Stdout.Write(response)
	}
	return nil
}
