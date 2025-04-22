package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/auth0client"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/auth0/go-auth0/management"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Fetcher struct {
	UserPID        string
	UserEmail      string
	ConnectionName string
	Roles          bool
	Orgs           bool
	SAML           bool
	client         *auth0client.Auth0Client
}

const defaultMaxUsers = 100

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

func (f *Fetcher) WithSAML(saml bool) *Fetcher {
	f.SAML = saml
	return f
}

func (f *Fetcher) WithConnectionName(connectionName string) *Fetcher {
	f.ConnectionName = connectionName
	return f
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	if f.Roles {
		err := f.fetchGroups(ctx, writer, errorWriter)
		if err != nil {
			return err
		}
	}

	return f.fetchUsers(ctx, writer, errorWriter)
}

func (f *Fetcher) fetchUsers(ctx context.Context, outputWriter *js.JSONArrayWriter, errorWriter io.Writer) error {
	page := 0

	for {
		opts := []management.RequestOption{management.Page(page)}
		if f.ConnectionName != "" {
			opts = append(opts, management.Query(f.getConnectionQuery()))
		}

		users, more, err := f.fetchUsersList(ctx, opts)
		if err != nil {
			common.WriteErrorWithExitCode(errorWriter, err, 1)
			return err
		}

		for _, user := range users {
			obj, err := f.buildOutputObjects(ctx, user)
			if err != nil {
				common.WriteErrorWithExitCode(errorWriter, err, 1)
			}

			if obj == nil {
				continue
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

// return a object map to output or a boolean to skip current user.
func (f *Fetcher) buildOutputObjects(ctx context.Context, user *management.User) (map[string]any, error) {
	var obj map[string]any

	res, err := user.MarshalJSON()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(res, &obj); err != nil {
		return nil, err
	}

	obj["email_verified"] = user.GetEmailVerified()
	obj["object_type"] = "user"

	if f.Roles {
		roles, err := f.fetchUserRoles(ctx, *user.ID)
		if err != nil {
			return nil, err
		} else {
			obj["roles"] = roles
		}
	}

	if f.Orgs {
		orgs, err := f.fetchOrgs(ctx, *user.ID)
		if err != nil {
			return nil, err
		} else {
			obj["orgs"] = orgs
		}
	}

	return obj, nil
}

func (f *Fetcher) fetchGroups(ctx context.Context, outputWriter *js.JSONArrayWriter, errorWriter io.Writer) error {
	page := 0

	for f.Roles {
		opts := []management.RequestOption{management.Page(page)}

		if f.ConnectionName != "" {
			opts = append(opts, management.Query(f.getConnectionQuery()))
		}

		roles, more, err := f.fetchRoles(ctx, opts)
		if err != nil {
			common.WriteErrorWithExitCode(errorWriter, err, 1)
			return err
		}

		for _, role := range roles {
			res := role.String()

			var obj map[string]interface{}

			if err := json.Unmarshal([]byte(res), &obj); err != nil {
				common.WriteErrorWithExitCode(errorWriter, err, 1)
				continue
			}

			obj["object_type"] = "role"
			if err := outputWriter.Write(obj); err != nil {
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

func (f *Fetcher) fetchUsersList(ctx context.Context, opts []management.RequestOption) ([]*management.User, bool, error) {
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
	}

	if f.UserEmail != "" {
		// List only users that have the provided email
		users, err := f.client.Mgmt.User.ListByEmail(ctx, f.UserEmail)
		if err != nil {
			return nil, false, err
		}

		return users, false, nil
	}

	// List all users
	if !f.SAML {
		userList, err := f.client.Mgmt.User.List(ctx, opts...)
		if err != nil {
			return nil, false, err
		}

		return userList.Users, userList.HasNext(), nil
	}

	// Use special SAML user list, to avoid known unmarshal errors, see notes below.
	ul := &UserList{}
	if err := FetchUsers(ctx, f.client.Mgmt, &ul, opts...); err != nil {
		return nil, false, err
	}

	return ul.UserList(), ul.HasNext(), nil
}

func (f *Fetcher) fetchRoles(ctx context.Context, opts []management.RequestOption) ([]*management.Role, bool, error) {
	roles, err := f.client.Mgmt.Role.List(ctx, opts...)
	if err != nil {
		return nil, false, err
	}

	if roles == nil {
		return nil, false, errors.Wrap(err, "failed to get roles")
	}

	return roles.Roles, roles.HasNext(), nil
}

func (f *Fetcher) fetchUserRoles(ctx context.Context, uID string) ([]map[string]any, error) {
	page := 0
	finished := false

	var results []map[string]any

	for !finished {
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

			var obj map[string]any
			if err := json.Unmarshal(res, &obj); err != nil {
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

func (f *Fetcher) fetchOrgs(ctx context.Context, uID string) ([]map[string]any, error) {
	page := 0
	finished := false

	var results []map[string]any

	for !finished {
		reqOpts := management.Page(page)

		orgs, err := f.client.Mgmt.User.Organizations(ctx, uID, reqOpts)
		if err != nil {
			return nil, err
		}

		for _, org := range orgs.Organizations {
			obj, err := formatObject(org)
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

func formatObject(org any) (map[string]any, error) {
	res, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}

	var obj map[string]any
	if err := json.Unmarshal(res, &obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (f *Fetcher) getConnectionQuery() string {
	if f.ConnectionName == "" {
		return ""
	}

	return `identities.connection:"` + f.ConnectionName + `"`
}

// Specialized SAML user list function
//
// The Auth0 golang SDK does not properly handle the unmarshal of the returned payload into a management.UserList.
//
// The returned payload contains:
// "email":"user@domain.com",
// "emailVerified":"true",
// "email_verified":"user@domain.com"
//
// Which results in an unmarshal error when calling
// `func (m *UserManager) List(ctx context.Context, opts ...RequestOption) (ul *UserList, err error)`
// resulting in an error `strconv.ParseBool: parsing "user@domain.com": invalid syntax`
//
// The implementation below works around the issues by using custom JSON marshaling to map the values into the management.User instances.
type User struct {
	management.User
}

type UserList struct {
	management.List
	Users []*User `json:"users"`
}

func (ul UserList) UserList() []*management.User {
	return lo.Map(ul.Users, func(v *User, i int) *management.User { return &v.User })
}

func (u *User) UnmarshalJSON(b []byte) error {
	var raw map[string]any

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	verified := false
	if emailVerified, ok := raw["emailVerified"].(string); ok {
		verified, _ = strconv.ParseBool(emailVerified)
	}

	delete(raw, "emailVerified")
	delete(raw, "email_verified")

	buf, err := json.Marshal(raw)
	if err != nil {
		return err
	}

	type tTmpUser User

	var tmpUser tTmpUser

	if err := json.Unmarshal(buf, &tmpUser); err != nil {
		return err
	}

	tmpUser.VerifyEmail = &verified

	*u = User(tmpUser)

	return nil
}

func FetchUsers(ctx context.Context, m *management.Management, payload any, options ...management.RequestOption) error {
	options = append(options,
		management.PerPage(defaultMaxUsers),
		management.IncludeTotals(true),
	)

	request, err := m.NewRequest(ctx, "GET", m.URI("users"), payload, options...)
	if err != nil {
		return fmt.Errorf("failed to create a new request: %w", err)
	}

	response, err := m.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send the request: %w", err)
	}
	defer response.Body.Close()

	// If the response contains a client or a server error then return the error.
	if response.StatusCode >= http.StatusBadRequest {
		return err // newError(response) //TODO create correct error based on response
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read the response body: %w", err)
	}

	if len(responseBody) > 0 && string(responseBody) != "{}" {
		if err := json.Unmarshal(responseBody, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal response payload: %w", err)
		}
	}

	return nil
}
