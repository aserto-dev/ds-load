package app

import (
	"context"
	"encoding/json"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/google/pkg/googleclient"
	"github.com/aserto-dev/ds-load/sdk/plugin"
)

type FetchCmd struct {
	ClientID     string `short:"i" help:"Google Client ID" env:"GOOGLE_CLIENT_ID" required:""`
	ClientSecret string `short:"s" help:"Google Client Secret" env:"GOOGLE_CLIENT_SECRET" required:""`
	RefreshToken string `short:"r" help:"Google Refresh Token" env:"GOOGLE_REFRESH_TOKEN" required:""`
	Groups       bool   `short:"g" help:"Retrieve Google groups" env:"GOOGLE_GROUPS" default:"false" negatable:""`
}

func (cmd *FetchCmd) Run(ctx *kong.Context) error {
	googleClient, err := createGoogleClient(cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken)
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
	return plugin.NewDSPlugin().WriteFetchOutput(results, errors, false)
}

func Fetch(googleClient *googleclient.GoogleClient, fetchGroups bool, results chan map[string]interface{}, errors chan error) {
	users, err := googleClient.ListUsers()
	if err != nil {
		errors <- err
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
	}

	if fetchGroups {
		groups, err := googleClient.ListGroups()
		if err != nil {
			errors <- err
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

			usersInGroup, err := googleClient.GetUsersInGroups(group.Id)
			if err != nil {
				errors <- err
			} else {
				usersInGroupBytes, err := json.Marshal(usersInGroup)
				if err != nil {
					errors <- err
				} else {
					var users []map[string]interface{}
					err = json.Unmarshal(usersInGroupBytes, &users)
					if err != nil {
						errors <- err
					}
					obj["users"] = users
				}
			}

			results <- obj
		}
	}
}

func createGoogleClient(clientID, clientSecret, refrestToken string) (googleClient *googleclient.GoogleClient, err error) {
	return googleclient.NewGoogleClient(
		context.Background(),
		clientID,
		clientSecret,
		refrestToken)
}
