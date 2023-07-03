package app

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
	"github.com/aserto-dev/ds-load/sdk/plugin"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

type FetchCmd struct {
	Domain   string `env:"DS_OKTA_DOMAIN"`
	APIToken string `env:"DS_OKTA_TOKEN"`
	Groups   bool   `env:"DS_OKTA_GROUPS" default:"true" negatable:""`
	Roles    bool   `env:"DS_OKTA_ROLES" default:"true" negatable:""`
}

type OktaPager func(context.Context, *okta.Response, *[]*okta.User) (*okta.Response, error)

func (cmd *FetchCmd) Run(ctx *kong.Context) error {
	oktaClient, err := oktaclient.NewOktaClient(context.Background(), cmd.Domain, cmd.APIToken)
	if err != nil {
		return err
	}
	results := make(chan map[string]interface{}, 1)
	errors := make(chan error, 1)
	go func() {
		cmd.Fetch(oktaClient, results, errors)
		close(results)
		close(errors)
	}()
	return plugin.NewDSPlugin().WriteFetchOutput(results, errors, false)
}

func (cmd *FetchCmd) Fetch(oktaClient oktaclient.OktaClient, results chan map[string]interface{}, errors chan error) {
	users, resp, err := oktaClient.ListUsers(context.Background(), nil)
	if err != nil {
		errors <- err
		SetExitCode(1)
		return
	}
	err = handleResponse(resp)
	if err != nil {
		errors <- err
		SetExitCode(1)
		return
	}

	for _, user := range users {
		userBytes, err := json.Marshal(user)
		if err != nil {
			errors <- err
			SetExitCode(1)
			return
		}
		var obj map[string]interface{}
		err = json.Unmarshal(userBytes, &obj)
		if err != nil {
			errors <- err
			SetExitCode(1)
		}
		results <- obj

		if cmd.Groups {
			// Write all user groups
			groups, resp, err := oktaClient.ListUserGroups(context.Background(), user.Id)
			if err != nil {
				errors <- err
				SetExitCode(1)
				return
			}
			err = handleResponse(resp)
			if err != nil {
				errors <- err
				SetExitCode(1)
				return
			}
			for _, group := range groups {
				groupBytes, err := json.Marshal(group)
				if err != nil {
					errors <- err
					SetExitCode(1)
					return
				}
				var obj map[string]interface{}
				err = json.Unmarshal(groupBytes, &obj)
				if err != nil {
					errors <- err
					SetExitCode(1)
				}
				results <- obj
			}
		}

		if cmd.Roles {
			// Write all user roles
			roles, resp, err := oktaClient.ListAssignedRolesForUser(context.Background(), user.Id, nil)
			if err != nil {
				errors <- err
				SetExitCode(1)
				return
			}
			err = handleResponse(resp)
			if err != nil {
				errors <- err
				SetExitCode(1)
				return
			}
			for _, role := range roles {
				roleBytes, err := json.Marshal(role)
				if err != nil {
					errors <- err
					SetExitCode(1)
				}
				var obj map[string]interface{}
				err = json.Unmarshal(roleBytes, &obj)
				if err != nil {
					errors <- err
					SetExitCode(1)
				}
				results <- obj
			}
		}
	}
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
