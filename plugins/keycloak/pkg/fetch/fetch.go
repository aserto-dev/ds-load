package fetch

import (
	"context"
	"encoding/json"
	"io"

	"github.com/aserto-dev/ds-load/plugins/keycloak/pkg/kc"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
)

type Fetcher struct {
	kcc    *kc.KeycloakClient
	Groups bool
	Roles  bool
}

func New(client *kc.KeycloakClient) (*Fetcher, error) {
	return &Fetcher{
		kcc: client,
	}, nil
}

func (f *Fetcher) WithGroups(groups bool) *Fetcher {
	f.Groups = groups
	return f
}

func (f *Fetcher) WithRoles(roles bool) *Fetcher {
	f.Roles = roles
	return f
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter io.Writer, errorWriter common.ErrorWriter) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	users, err := f.kcc.ListUsers(ctx)
	if err != nil {
		errorWriter.Error(err)
		return err
	}

	for _, user := range users {
		user.Type = "user"

		userBytes, err := json.Marshal(user)
		if err != nil {
			errorWriter.Error(err)
			continue
		}

		var obj map[string]any

		if err := json.Unmarshal(userBytes, &obj); err != nil {
			errorWriter.Error(err)
			continue
		}

		if err := writer.Write(obj); err != nil {
			errorWriter.Error(err)
		}
	}

	if f.Groups {
		if err := f.fetchGroups(ctx, writer, errorWriter); err != nil {
			return err
		}
	}

	if f.Roles {
		if err := f.fetchRoles(ctx, writer, errorWriter); err != nil {
			return err
		}
	}

	return nil
}

func (f *Fetcher) fetchGroups(ctx context.Context, //nolint:dupl
	writer *js.JSONArrayWriter,
	errorWriter common.ErrorWriter,
) error {
	groups, err := f.kcc.ListGroups(ctx)
	if err != nil {
		errorWriter.Error(err)
		return err
	}

	for _, group := range groups {
		group.Type = "group"
		groupBytes, err := json.Marshal(group)
		errorWriter.Error(err)

		var obj map[string]any
		if err := json.Unmarshal(groupBytes, &obj); err != nil {
			errorWriter.Error(err)
			continue
		}

		usersInGroup, err := f.kcc.GetUsersOfGroup(ctx, group.ID)
		errorWriter.Error(err)

		usersInGroupBytes, err := json.Marshal(usersInGroup)
		errorWriter.Error(err)

		var users []map[string]any
		if err := json.Unmarshal(usersInGroupBytes, &users); err != nil {
			errorWriter.Error(err)
		}

		obj["users"] = users

		err = writer.Write(obj)
		errorWriter.Error(err)
	}

	return nil
}

func (f *Fetcher) fetchRoles(ctx context.Context, //nolint:dupl
	writer *js.JSONArrayWriter,
	errorWriter common.ErrorWriter,
) error {
	roles, err := f.kcc.ListRoles(ctx)
	if err != nil {
		errorWriter.Error(err)
		return err
	}

	for _, role := range roles {
		role.Type = "role"
		roleBytes, err := json.Marshal(role)
		errorWriter.Error(err)

		var obj map[string]any
		if err := json.Unmarshal(roleBytes, &obj); err != nil {
			errorWriter.Error(err)
			continue
		}

		usersInRole, err := f.kcc.GetUsersOfRole(ctx, role.Name)
		errorWriter.Error(err)

		usersInRoleBytes, err := json.Marshal(usersInRole)
		errorWriter.Error(err)

		var users []map[string]any
		if err := json.Unmarshal(usersInRoleBytes, &users); err != nil {
			errorWriter.Error(err)
		}

		obj["users"] = users

		err = writer.Write(obj)
		errorWriter.Error(err)
	}

	return nil
}
