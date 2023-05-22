package app

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

// TODO: add flags to allow picking each our of users,groups and roles
type FetchCmd struct {
	Domain    string `cmd:"" env:"DS_OKTA_DOMAIN"`
	APIToken  string `cmd:"" env:"DS_OKTA_TOKEN"`
	UserPID   string `cmd:"" env:"DS_OKTA_USER_PID"`
	UserEmail string `cmd:"" env:"DS_OKTA_USER_EMAIL"`
	Users     bool   `cmd:"" env:"DS_OKTA_USERS"`
	Groups    bool   `cmd:"" env:"DS_OKTA_GROUPS"`
	Roles     bool   `cmd:"" env:"DS_OKTA_ROLES"`
}

type OktaPager func(context.Context, *okta.Response, *[]*okta.User) (*okta.Response, error)

func (cmd *FetchCmd) Run(ctx *kong.Context) error {

	oktaClient, err := oktaclient.NewOktaClient(context.Background(), cmd.Domain, cmd.APIToken)
	if err != nil {
		return err
	}

	users, resp, err := oktaClient.ListUsers(context.Background(), nil)
	if err != nil {
		return err
	}
	err = handleResponse(resp)
	if err != nil {
		return err
	}

	for _, user := range users {
		userBytes, err := json.Marshal(user)
		if err != nil {
			return err
		}
		os.Stdout.Write(userBytes)
		os.Stdout.WriteString("\n")

		// Write all user groups
		groups, resp, err := oktaClient.ListUserGroups(context.Background(), user.Id)
		if err != nil {
			return err
		}
		err = handleResponse(resp)
		if err != nil {
			return err
		}
		for _, group := range groups {
			groupBytes, err := json.Marshal(group)
			if err != nil {
				os.Stderr.Write([]byte(err.Error()))
			}
			os.Stdout.Write(groupBytes)
			os.Stdout.WriteString("\n")
		}

		// Write all user roles
		roles, resp, err := oktaClient.ListAssignedRolesForUser(context.Background(), user.Id, nil)
		if err != nil {
			return err
		}
		err = handleResponse(resp)
		if err != nil {
			return err
		}
		for _, role := range roles {
			roleBytes, err := json.Marshal(role)
			if err != nil {
				os.Stderr.WriteString(err.Error())
			}
			os.Stdout.Write(roleBytes)
			os.Stdout.WriteString("\n")
		}
	}

	return nil
}

func handleResponse(resp *okta.Response) error {
	if resp.Response != nil && resp.StatusCode == http.StatusTooManyRequests {
		response, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		os.Stdout.Write(response)
	}
	return nil
}
