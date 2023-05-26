package app

import (
	"encoding/json"
	"os"

	"github.com/alecthomas/kong"
	"gopkg.in/auth0.v5/management"
)

type FetchCmd struct {
	Domain       string `name:"domain" short:"d" env:"AUTH0_DOMAIN" help:"auth0 domain" required:""`
	ClientID     string `name:"client-id" short:"i" env:"AUTH0_CLIENT_ID" help:"auth0 client id" required:""`
	ClientSecret string `name:"client-secret" short:"s" env:"AUTH0_CLIENT_SECRET" help:"auth0 client secret" required:""`
}

func (cmd *FetchCmd) Run(context *kong.Context) error {
	mgmt, err := management.New(
		cmd.Domain,
		management.WithClientCredentials(
			cmd.ClientID,
			cmd.ClientSecret,
		))
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

	for o := range results {
		res, err := json.Marshal(o)
		if err != nil {
			return err
		}
		os.Stdout.Write(res)
		os.Stdout.WriteString("\n")
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
			results <- obj
		}
		if !ul.HasNext() {
			finished = true
		}
		page++
	}
}
