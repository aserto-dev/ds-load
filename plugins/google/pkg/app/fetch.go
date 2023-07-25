package app

import (
	"context"
	"encoding/json"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/google/pkg/googleclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/plugin"
)

type FetchCmd struct {
	ClientID     string `short:"i" help:"Google Client ID" env:"GOOGLE_CLIENT_ID" required:""`
	ClientSecret string `short:"s" help:"Google Client Secret" env:"GOOGLE_CLIENT_SECRET" required:""`
	RefreshToken string `short:"r" help:"Google Refresh Token" env:"GOOGLE_REFRESH_TOKEN" required:""`
	Groups       bool   `short:"g" help:"Retrieve Google groups" env:"GOOGLE_GROUPS" default:"false"`
	Customer     string `help:"Google Workspace Customer field" env:"GOOGLE_CUSTOMER" default:"my_customer"`
}

func (cmd *FetchCmd) Run(ctx *kong.Context) error {
	googleClient, err := createGoogleClient(cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken, cmd.Customer)
	if err != nil {
		return err
	}

	results := make(chan map[string]interface{}, 1)
	errors := make(chan error, 1)
	go func() {
		Fetch(googleClient, cmd.Groups, results, errors)
		close(results)
		close(errors)
	}()

	return plugin.NewDSPlugin().WriteFetchOutput(results, errors)
}

func Fetch(googleClient *googleclient.GoogleClient, fetchGroups bool, results chan map[string]interface{}, errors chan error) {
	users, err := googleClient.ListUsers()
	if err != nil {
		errors <- err
		common.SetExitCode(1)
		return
	}

	for _, user := range users {
		userBytes, err := json.Marshal(user)
		if err != nil {
			errors <- err
			common.SetExitCode(1)
			continue
		}
		var obj map[string]interface{}
		err = json.Unmarshal(userBytes, &obj)
		if err != nil {
			errors <- err
			common.SetExitCode(1)
			continue
		}

		results <- obj
	}

	if fetchGroups {
		groups, err := googleClient.ListGroups()
		if err != nil {
			errors <- err
			common.SetExitCode(1)
			return
		}

		for _, group := range groups {
			groupBytes, err := json.Marshal(group)
			if err != nil {
				errors <- err
				common.SetExitCode(1)
				continue
			}
			var obj map[string]interface{}
			err = json.Unmarshal(groupBytes, &obj)
			if err != nil {
				errors <- err
				common.SetExitCode(1)
				continue
			}

			usersInGroup, err := googleClient.GetUsersInGroup(group.Id)
			if err != nil {
				errors <- err
				common.SetExitCode(1)
			} else {
				usersInGroupBytes, err := json.Marshal(usersInGroup)
				if err != nil {
					errors <- err
					common.SetExitCode(1)
				} else {
					var users []map[string]interface{}
					err = json.Unmarshal(usersInGroupBytes, &users)
					if err != nil {
						errors <- err
						common.SetExitCode(1)
					}
					obj["users"] = users
				}
			}

			results <- obj
		}
	}
}

func createGoogleClient(clientID, clientSecret, refrestToken, customer string) (googleClient *googleclient.GoogleClient, err error) {
	return googleclient.NewGoogleClient(
		context.Background(),
		clientID,
		clientSecret,
		refrestToken,
		customer)
}