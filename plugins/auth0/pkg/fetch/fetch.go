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
	Orgs           bool
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

func (f *Fetcher) WithOrgs(orgs bool) *Fetcher {
	f.Orgs = orgs
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

		users, more, err := f.getUsers(ctx, opts)
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
			if f.Roles {
				roles, err := f.getRoles(ctx, *user.ID)
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

func (f *Fetcher) getUsers(ctx context.Context, opts []management.RequestOption) ([]*management.User, bool, error) {
	if f.UserPID != "" && f.UserEmail != "" {
		return nil, false, errors.New("only one of user-pid or user-email can be specified")
	}

	if f.UserPID != "" {
		// list only the user with the provided pid
		user, err := f.client.Mgmt.User.Read(ctx, f.UserPID)
		if err != nil {
			return nil, false, err
		}
		if user == nil {
			return nil, false, errors.Wrapf(err, "failed to get user by pid %s", f.UserPID)
		}
		return []*management.User{user}, false, nil
	} else if f.UserEmail != "" {
		// List only users that have the provided email
		users, err := f.client.Mgmt.User.ListByEmail(ctx, f.UserEmail)
		if err != nil {
			return nil, false, err
		}
		return users, false, nil
	} else {
		// List all users
		userList, err := f.client.Mgmt.User.List(ctx, opts...)
		if err != nil {
			return nil, false, err
		}

		return userList.Users, userList.HasNext(), nil
	}
}

func (f *Fetcher) getRoles(ctx context.Context, uID string) ([]map[string]interface{}, error) {
	page := 0
	finished := false

	var results []map[string]interface{}

	for {
		if finished {
			break
		}

		reqOpts := management.Page(page)
		roles, err := f.client.Mgmt.User.Roles(ctx, uID, reqOpts)
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

func (f *Fetcher) getOrgs(ctx context.Context, uID string) ([]map[string]interface{}, error) {
	page := 0
	finished := false

	var results []map[string]interface{}

	for {
		if finished {
			break
		}

		reqOpts := management.Page(page)
		orgs, err := f.client.Mgmt.User.Organizations(ctx, uID, reqOpts)
		if err != nil {
			return nil, err
		}
		for _, org := range orgs.Organizations {
			res, err := json.Marshal(org)
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
		if !orgs.HasNext() {
			finished = true
		}

		page++
	}

	return results, nil
}
