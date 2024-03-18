package fetch

import (
	"context"
	"encoding/json"
	"io"

	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/auth0client"
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
	client         *auth0client.Auth0Client
}

func New(ctx context.Context, client *auth0client.Auth0Client) (*Fetcher, error) {
	return &Fetcher{
		client: client,
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

	if f.Roles {
		err := f.fetchGroups(writer, errorWriter)
		if err != nil {
			return err
		}
	}

	return f.fetchUsers(writer, errorWriter)
}

func (f *Fetcher) fetchUsers(outputWriter *js.JSONArrayWriter, errorWriter io.Writer) error {
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
			obj["email_verified"] = user.GetEmailVerified()
			obj["object_type"] = "user"
			if f.Roles {
				roles, err := f.getUserRoles(*user.ID)
				if err != nil {
					_, _ = errorWriter.Write([]byte(err.Error()))
					common.SetExitCode(1)
				} else {
					obj["roles"] = roles
				}
			}
			err = outputWriter.Write(obj)
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

func (f *Fetcher) fetchGroups(outputWriter *js.JSONArrayWriter, errorWriter io.Writer) error {
	page := 0

	for f.Roles {
		opts := []management.RequestOption{management.Page(page)}
		if f.ConnectionName != "" {
			opts = append(opts, management.Query(`identities.connection:"`+f.ConnectionName+`"`))
		}

		roles, more, err := f.getRoles(opts)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
			return err
		}

		for _, role := range roles {
			res := role.String()
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(res), &obj)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
				continue
			}
			obj["object_type"] = "role"
			err = outputWriter.Write(obj)
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
		user, err := f.client.Mgmt.User.Read(f.UserPID)
		if err != nil {
			return nil, false, err
		}
		if user == nil {
			return nil, false, errors.Wrapf(err, "failed to get user by pid %s", f.UserPID)
		}
		return []*management.User{user}, false, nil
	} else if f.UserEmail != "" {
		// List only users that have the provided email
		users, err := f.client.Mgmt.User.ListByEmail(f.UserEmail)
		if err != nil {
			return nil, false, err
		}
		return users, false, nil
	} else {
		// List all users
		userList, err := f.client.Mgmt.User.List(opts...)
		if err != nil {
			return nil, false, err
		}

		return userList.Users, userList.HasNext(), nil
	}
}

func (f *Fetcher) getRoles(opts []management.RequestOption) ([]*management.Role, bool, error) {
	roles, err := f.client.Mgmt.Role.List(opts...)
	if err != nil {
		return nil, false, err
	}
	if roles == nil {
		return nil, false, errors.Wrap(err, "failed to get roles")
	}

	return roles.Roles, roles.HasNext(), nil
}

func (f *Fetcher) getUserRoles(uID string) ([]map[string]interface{}, error) {
	page := 0
	finished := false

	var results []map[string]interface{}

	for {
		if finished {
			break
		}

		reqOpts := management.Page(page)
		roles, err := f.client.Mgmt.User.Roles(uID, reqOpts)
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
