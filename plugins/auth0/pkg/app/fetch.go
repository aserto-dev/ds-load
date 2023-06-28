package app

import (
	"encoding/json"
	"net/http"

	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/httpclient"
	"github.com/aserto-dev/ds-load/sdk/plugin"

	"github.com/alecthomas/kong"
	"github.com/auth0/go-auth0/management"
)

type FetchCmd struct {
	Domain         string `name:"domain" short:"d" env:"AUTH0_DOMAIN" help:"auth0 domain" required:""`
	ClientID       string `name:"client-id" short:"i" env:"AUTH0_CLIENT_ID" help:"auth0 client id" required:""`
	ClientSecret   string `name:"client-secret" short:"s" env:"AUTH0_CLIENT_SECRET" help:"auth0 client secret" required:""`
	ConnectionName string `name:"connection-name" env:"AUTH0_CONNECTION_NAME" help:"auth0 connection name" optional:""`
	RateLimit      bool   `default:"true" help:"enable http client rate limiter" negatable:"" optional:""`
	Roles          bool   `env:"AUTH0_ROLES" default:"false" negatable:""`
}

func (cmd *FetchCmd) Run(context *kong.Context) error {
	options := []management.Option{
		management.WithClientCredentials(
			cmd.ClientID,
			cmd.ClientSecret,
		),
	}
	if cmd.RateLimit {
		client := http.DefaultClient
		client.Transport = httpclient.NewTransport(http.DefaultTransport)
		options = append(options, management.WithClient(client))
	}

	mgmt, err := management.New(
		cmd.Domain,
		options...,
	)
	if err != nil {
		return err
	}

	results := make(chan map[string]interface{}, 1)
	errors := make(chan error, 1)
	go func() {
		cmd.Fetch(mgmt, results, errors)
		close(results)
		close(errors)
	}()
	if err != nil {
		return err
	}

	return plugin.NewDSPlugin().WriteFetchOutput(results, errors, false)
}

func (cmd *FetchCmd) Fetch(mgmt *management.Management, results chan map[string]interface{}, errors chan error) {
	page := 0
	finished := false

	for {
		if finished {
			break
		}

		opts := []management.RequestOption{management.Page(page)}
		if cmd.ConnectionName != "" {
			opts = append(opts, management.Query(`identities.connection:"`+cmd.ConnectionName+`"`))
		}
		ul, err := mgmt.User.List(opts...)
		if err != nil {
			errors <- err
			return
		}

		for _, u := range ul.Users {
			res, err := u.MarshalJSON()
			if err != nil {
				errors <- err
			}
			var obj map[string]interface{}
			err = json.Unmarshal(res, &obj)
			if err != nil {
				errors <- err
			}
			if cmd.Roles {
				roles, err := getRoles(mgmt, *u.ID)
				if err != nil {
					errors <- err
				}
				obj["roles"] = roles
			}
			results <- obj
		}
		if !ul.HasNext() {
			finished = true
		}
		page++
	}
}

func getRoles(mgmt *management.Management, uID string) ([]map[string]interface{}, error) {
	page := 0
	finished := false

	var results []map[string]interface{}

	for {
		if finished {
			break
		}

		reqOpts := management.Page(page)
		roles, err := mgmt.User.Roles(uID, reqOpts)
		if err != nil {
			return nil, err
		}
		for _, role := range roles.Roles {
			res, err := json.Marshal(role)
			if err != nil {
				return nil, err
			}
			var obj map[string]interface{}
			err = json.Unmarshal(res, &obj)
			if err != nil {
				return nil, err
			}
			results = append(results, obj)
		}
		if !roles.HasNext() {
			finished = true
		}

		page++
	}
	return results, nil
}
