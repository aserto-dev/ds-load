package app

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/httpclient"
	"github.com/aserto-dev/ds-load/sdk/plugin"
	"github.com/pkg/errors"

	"github.com/alecthomas/kong"
	"github.com/auth0/go-auth0/management"
)

type FetchCmd struct {
	Domain         string `name:"domain" short:"d" env:"AUTH0_DOMAIN" help:"auth0 domain" required:""`
	ClientID       string `name:"client-id" short:"i" env:"AUTH0_CLIENT_ID" help:"auth0 client id" required:""`
	ClientSecret   string `name:"client-secret" short:"s" env:"AUTH0_CLIENT_SECRET" help:"auth0 client secret" required:""`
	ConnectionName string `name:"connection-name" env:"AUTH0_CONNECTION_NAME" help:"auth0 connection name" optional:""`
	UserPID        string `name:"user-pid" env:"AUTH0_USER_PID" help:"auth0 user PID of the user you want to read" optional:""`
	UserEmail      string `name:"user-email" env:"AUTH0_USER_EMAIL" help:"auth0 user email of the user you want to read" optional:""`
	Roles          bool   `name:"roles" env:"AUTH0_ROLES" default:"false" negatable:"" help:"include roles"`
	RateLimit      bool   `name:"rate-limit" default:"true" help:"enable http client rate limiter" negatable:""`

	mgmt *management.Management `kong:"-"`
}

func (fetcher *FetchCmd) Run(context *kong.Context) error {
	if fetcher.UserPID != "" && !strings.HasPrefix(fetcher.UserPID, "auth0|") {
		fetcher.UserPID = "auth0|" + fetcher.UserPID
	}

	options := []management.Option{
		management.WithClientCredentials(
			fetcher.ClientID,
			fetcher.ClientSecret,
		),
	}
	if fetcher.RateLimit {
		client := http.DefaultClient
		client.Transport = httpclient.NewTransport(http.DefaultTransport)
		options = append(options, management.WithClient(client))
	}

	mgmt, err := management.New(
		fetcher.Domain,
		options...,
	)
	if err != nil {
		return err
	}

	fetcher.mgmt = mgmt

	results := make(chan map[string]interface{}, 1)
	errCh := make(chan error, 1)
	go func() {
		fetcher.Fetch(results, errCh)
		close(results)
		close(errCh)
	}()
	if err != nil {
		return err
	}

	return plugin.NewDSPlugin().WriteFetchOutput(results, errCh, false)
}

func (fetcher *FetchCmd) Fetch(results chan map[string]interface{}, errCh chan error) {
	page := 0
	finished := false

	if fetcher.UserPID != "" {
		user, err := fetcher.readByPID()
		if err != nil {
			errCh <- err
			return
		}
		results <- user
		return
	}

	if fetcher.UserEmail != "" {
		users, err := fetcher.readByEmail()
		if err != nil {
			errCh <- err
			return
		}
		for _, user := range users {
			results <- user
		}
		return
	}

	for {
		if finished {
			break
		}

		opts := []management.RequestOption{management.Page(page)}
		if fetcher.ConnectionName != "" {
			opts = append(opts, management.Query(`identities.connection:"`+fetcher.ConnectionName+`"`))
		}
		ul, err := fetcher.mgmt.User.List(opts...)
		if err != nil {
			errCh <- err
			return
		}

		for _, u := range ul.Users {
			res, err := u.MarshalJSON()
			if err != nil {
				errCh <- err
				continue
			}
			var obj map[string]interface{}
			err = json.Unmarshal(res, &obj)
			if err != nil {
				errCh <- err
				continue
			}
			if fetcher.Roles {
				roles, err := fetcher.getRoles(*u.ID)
				if err != nil {
					errCh <- err
				} else {
					obj["roles"] = roles
				}
			}
			results <- obj
		}
		if !ul.HasNext() {
			finished = true
		}
		page++
	}
}

func (fetcher *FetchCmd) getRoles(uID string) ([]map[string]interface{}, error) {
	page := 0
	finished := false

	var results []map[string]interface{}

	for {
		if finished {
			break
		}

		reqOpts := management.Page(page)
		roles, err := fetcher.mgmt.User.Roles(uID, reqOpts)
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

func (fetcher *FetchCmd) readByPID() (map[string]interface{}, error) {

	user, err := fetcher.mgmt.User.Read(fetcher.UserPID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.Wrapf(err, "failed to get user by pid %s", fetcher.UserPID)
	}
	res, err := user.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	err = json.Unmarshal(res, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (fetcher *FetchCmd) readByEmail() ([]map[string]interface{}, error) {
	var users []map[string]interface{}

	auth0Users, err := fetcher.mgmt.User.ListByEmail(fetcher.UserEmail)
	if err != nil {
		return nil, err
	}
	if len(auth0Users) < 1 {
		return nil, errors.Wrapf(err, "failed to get user by emal %s", fetcher.UserEmail)
	}

	for _, user := range auth0Users {
		res, err := user.MarshalJSON()
		if err != nil {
			return nil, err
		}
		var obj map[string]interface{}
		err = json.Unmarshal(res, &obj)
		if err != nil {
			return nil, err
		}
		users = append(users, obj)
	}

	return users, nil
}
