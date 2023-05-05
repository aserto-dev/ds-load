package app

import (
	"encoding/json"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/config"
	"gopkg.in/auth0.v5/management"
)

type FetchCmd struct {
	Config       string `short:"c" type:"path" help:"Path to the config file. Any argument provided to the CLI will take precedence."`
	Domain       string `cmd:""`
	ClientID     string `cmd:""`
	ClientSecret string `cmd:""`
}

func (cmd *FetchCmd) Run(context *kong.Context) error {
	cfg, err := config.NewConfig(cmd.Config)
	if err != nil {
		return err
	}
	if cmd.Domain != "" {
		cfg.Domain = cmd.Domain
	}
	if cmd.ClientID != "" {
		cfg.ClientID = cmd.ClientID
	}
	if cmd.ClientSecret != "" {
		cfg.ClientSecret = cmd.ClientSecret
	}

	mgmt, err := management.New(
		cfg.Domain,
		management.WithClientCredentials(
			cfg.ClientID,
			cfg.ClientSecret,
		))
	if err != nil {
		return err
	}

	page := 0
	finished := false

	for {
		if finished {
			break
		}

		opts := management.Page(page)
		ul, err := mgmt.User.List(opts)
		if err != nil {
			return err
		}

		for _, u := range ul.Users {
			res, err := json.Marshal(u)
			if err != nil {
				return err
			}
			os.Stdout.Write(res)
			os.Stdout.WriteString("\n")
		}
		if !ul.HasNext() {
			finished = true
		}
		page++
	}
	return nil
}
