package app

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aserto-dev/ds-load/common/js"
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/httpclient"

	"github.com/alecthomas/kong"
	"gopkg.in/auth0.v5/management"
)

type FetchCmd struct {
	Domain       string `name:"domain" short:"d" env:"AUTH0_DOMAIN" help:"auth0 domain" required:""`
	ClientID     string `name:"client-id" short:"i" env:"AUTH0_CLIENT_ID" help:"auth0 client id" required:""`
	ClientSecret string `name:"client-secret" short:"s" env:"AUTH0_CLIENT_SECRET" help:"auth0 client secret" required:""`
	RateLimit    bool   `default:"true" help:"enable http client rate limiter" negatable:"" optional:""`
}

func (cmd *FetchCmd) Run(context *kong.Context) error {
	options := []management.ManagementOption{
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
		Fetch(mgmt, results, errors)
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

	writer := js.NewJSONArrayWriter(os.Stdout)
	defer writer.Close()
	for o := range results {
		err := writer.Write(o)
		if err != nil {
			return err
		}
	}
	return nil
}

func Fetch(mgmt *management.Management, results chan map[string]interface{}, errors chan error) {
	page := 0
	finished := false

	for {
		if finished {
			break
		}

		opts := management.Page(page)
		ul, err := mgmt.User.List(opts)
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
			roles, err := getRoles(mgmt, *u.ID)
			if err != nil {
				errors <- err
			}
			obj["roles"] = roles
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
