package fetch

import (
	"context"
	"encoding/json"
	"io"

	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/auth0/go-auth0/management"
	"github.com/pkg/errors"
)

type Fetcher struct {
	UserPID        string
	UserEmail      string
	ConnectionName string
	Roles          bool
	mgmt           *management.Management
}

func New(clientID, clientSecret, domain string) (*Fetcher, error) {
	options := []management.Option{
		management.WithClientCredentials(
			clientID,
			clientSecret,
		),
	}

	mgmt, err := management.New(
		domain,
		options...,
	)
	if err != nil {
		return nil, err
	}

	return &Fetcher{
		mgmt: mgmt,
	}, nil
}

func (f *Fetcher) WithEmail(email string) *Fetcher {
	f.UserEmail = email
	return f
}

func (f *Fetcher) WithUserPID(userPID string) *Fetcher {
	f.UserPID = userPID
	return f
}

func (f *Fetcher) WithRoles(roles bool) *Fetcher {
	f.Roles = roles
	return f
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	writer, err := js.NewJSONArrayWriter(outputWriter)
	if err != nil {
		return err
	}
	defer writer.Close()

	page := 0

	for {
		opts := []management.RequestOption{management.Page(page)}
		if f.ConnectionName != "" {
			opts = append(opts, management.Query(`identities.connection:"`+f.ConnectionName+`"`))
		}

		users, more, err := f.getUsers(opts)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
			return err
		}

		for _, user := range users {
			res, err := user.MarshalJSON()
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
				continue
			}
			var obj map[string]interface{}
			err = json.Unmarshal(res, &obj)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
				continue
			}
			if f.Roles {
				roles, err := f.getRoles(*user.ID)
				if err != nil {
					_, _ = errorWriter.Write([]byte(err.Error()))
					common.SetExitCode(1)
				} else {
					obj["roles"] = roles
				}
			}
			err = writer.Write(obj)
			if err != nil {
				return err
			}
		}
		if !more {
			break
		}
		page++
	}

	return nil
}

func (f *Fetcher) getUsers(opts []management.RequestOption) ([]*management.User, bool, error) {
	if f.UserPID != "" && f.UserEmail != "" {
		return nil, false, errors.New("only one of user-pid or user-email can be specified")
	}

	if f.UserPID != "" {
		// list only the user with the provided pid
		user, err := f.mgmt.User.Read(f.UserPID)
		if err != nil {
			return nil, false, err
		}
		if user == nil {
			return nil, false, errors.Wrapf(err, "failed to get user by pid %s", f.UserPID)
		}
		return []*management.User{user}, false, nil
	} else if f.UserEmail != "" {
		// List only users that have the provided email
		users, err := f.mgmt.User.ListByEmail(f.UserEmail)
		if err != nil {
			return nil, false, err
		}
		return users, false, nil
	} else {
		// List all users
		userList, err := f.mgmt.User.List(opts...)
		if err != nil {
			return nil, false, err
		}

		return userList.Users, userList.HasNext(), nil
	}
}

func (f *Fetcher) getRoles(uID string) ([]map[string]interface{}, error) {
	page := 0
	finished := false

	var results []map[string]interface{}

	for {
		if finished {
			break
		}

		reqOpts := management.Page(page)
		roles, err := f.mgmt.User.Roles(uID, reqOpts)
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
