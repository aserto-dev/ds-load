package app

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds/plugins/okta/pkg/oktaclient"
	"github.com/aserto-dev/ds/sdk/common/js"
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
	if err != nil {
		return err
	}

	go func() {
		for err := range errors {
			os.Stderr.WriteString(err.Error())
			os.Stderr.WriteString("\n")
		}
	}()

	writer, err := js.NewJSONArrayWriter(os.Stdout)
	if err != nil {
		return err
	}
	defer writer.Close()
	for o := range results {
		err := writer.Write(o)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cmd *FetchCmd) Fetch(oktaClient oktaclient.OktaClient, results chan map[string]interface{}, errors chan error) {
	users, resp, err := oktaClient.ListUsers(context.Background(), nil)
	if err != nil {
		errors <- err
		return
	}
	err = handleResponse(resp)
	if err != nil {
		errors <- err
		return
	}

	for _, user := range users {
		userBytes, err := json.Marshal(user)
		if err != nil {
			errors <- err
			return
		}
		var obj map[string]interface{}
		err = json.Unmarshal(userBytes, &obj)
		if err != nil {
			errors <- err
		}
		results <- obj

		if cmd.Groups {
			// Write all user groups
			groups, resp, err := oktaClient.ListUserGroups(context.Background(), user.Id)
			if err != nil {
				errors <- err
				return
			}
			err = handleResponse(resp)
			if err != nil {
				errors <- err
				return
			}
			for _, group := range groups {
				groupBytes, err := json.Marshal(group)
				if err != nil {
					errors <- err
					return
				}
				var obj map[string]interface{}
				err = json.Unmarshal(groupBytes, &obj)
				if err != nil {
					errors <- err
				}
				results <- obj
			}
		}

		if cmd.Roles {
			// Write all user roles
			roles, resp, err := oktaClient.ListAssignedRolesForUser(context.Background(), user.Id, nil)
			if err != nil {
				errors <- err
				return
			}
			err = handleResponse(resp)
			if err != nil {
				errors <- err
				return
			}
			for _, role := range roles {
				roleBytes, err := json.Marshal(role)
				if err != nil {
					os.Stderr.WriteString(err.Error())
				}
				var obj map[string]interface{}
				err = json.Unmarshal(roleBytes, &obj)
				if err != nil {
					errors <- err
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
