package fetch

import (
	"context"
	"encoding/json"
	"io"

	"github.com/aserto-dev/ds-load/plugins/jumpcloud/pkg/jc"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
)

type Fetcher struct {
	jcc    *jc.JumpCloudClient
	Groups bool
}

func New(client *jc.JumpCloudClient) (*Fetcher, error) {
	return &Fetcher{
		jcc: client,
	}, nil
}

func (f *Fetcher) WithGroups(groups bool) *Fetcher {
	f.Groups = groups
	return f
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter io.Writer, errorWriter common.ErrorWriter) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	users, err := f.jcc.ListUsers(ctx)
	if err != nil {
		errorWriter.Error(err)
		return err
	}

	idLookup := map[string]*jc.BaseUser{}

	for _, user := range users {
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

		idLookup[user.ID] = &user.BaseUser
	}

	if f.Groups {
		if err := f.fetchGroups(ctx, writer, errorWriter, idLookup); err != nil {
			return err
		}
	}

	return nil
}

func (f *Fetcher) fetchGroups(ctx context.Context,
	writer *js.JSONArrayWriter,
	errorWriter common.ErrorWriter,
	idLookup map[string]*jc.BaseUser,
) error {
	groups, err := f.jcc.ListGroups(ctx, jc.UserGroups)
	if err != nil {
		errorWriter.Error(err)
		return err
	}

	for _, group := range groups {
		groupBytes, err := json.Marshal(group)
		errorWriter.Error(err)

		var obj map[string]any
		if err := json.Unmarshal(groupBytes, &obj); err != nil {
			errorWriter.Error(err)
			continue
		}

		usersInGroup, err := f.jcc.ExpandUsersInGroup(ctx, group.ID, idLookup)
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
