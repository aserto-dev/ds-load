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

func New(userPid, clientID, clientSecret, domain string) (*Fetcher, error) {
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
		mgmt:    mgmt,
		UserPID: userPid,
	}, nil
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	writer, err := js.NewJSONArrayWriter(outputWriter)
	if err != nil {
		return err
	}
	defer writer.Close()

	page := 0
	finished := false

	for {
		if finished {
			break
		}

		opts := []management.RequestOption{management.Page(page)}
		if f.ConnectionName != "" {
			opts = append(opts, management.Query(`identities.connection:"`+f.ConnectionName+`"`))
		}
		ul, err := f.mgmt.User.List(opts...)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
			return err
		}

		for _, u := range ul.Users {
			res, err := u.MarshalJSON()
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
				roles, err := f.getRoles(*u.ID)
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
		if !ul.HasNext() {
			finished = true
		}
		page++
	}

	return nil
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

func (f *Fetcher) FetchUserByID(ctx context.Context, id string, outputWriter, errorWriter io.Writer) error {
	writer, err := js.NewJSONArrayWriter(outputWriter)
	if err != nil {
		return err
	}
	defer writer.Close()

	user, err := f.readByPID()
	if err != nil {
		_, _ = errorWriter.Write([]byte(err.Error()))
		common.SetExitCode(1)
		return err
	}
	return writer.Write(user)
}

func (f *Fetcher) readByPID() (map[string]interface{}, error) {
	user, err := f.mgmt.User.Read(f.UserPID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.Wrapf(err, "failed to get user by pid %s", f.UserPID)
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

func (f *Fetcher) FetchUserByEmail(ctx context.Context, email string, outputWriter, errorWriter io.Writer) error {
	writer, err := js.NewJSONArrayWriter(outputWriter)
	if err != nil {
		return err
	}
	defer writer.Close()

	users, err := f.readByEmail()
	if err != nil {
		_, _ = errorWriter.Write([]byte(err.Error()))
		common.SetExitCode(1)
		return err
	}
	for _, user := range users {
		err = writer.Write(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Fetcher) readByEmail() ([]map[string]interface{}, error) {
	var users []map[string]interface{}

	auth0Users, err := f.mgmt.User.ListByEmail(f.UserEmail)
	if err != nil {
		return nil, err
	}
	if len(auth0Users) < 1 {
		return nil, errors.Wrapf(err, "failed to get user by emal %s", f.UserEmail)
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
