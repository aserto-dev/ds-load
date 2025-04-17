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

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	users, err := f.jcc.ListUsers(ctx)
	if err != nil {
		common.WriteErrorWithExitCode(errorWriter, err, 1)
		return err
	}

	idLookup := map[string]*jc.BaseUser{}

	for _, user := range users {
		userBytes, err := json.Marshal(user)
		if err != nil {
			common.WriteErrorWithExitCode(errorWriter, err, 1)
			continue
		}

		var obj map[string]interface{}

		if err := json.Unmarshal(userBytes, &obj); err != nil {
			common.WriteErrorWithExitCode(errorWriter, err, 1)
			continue
		}

		if err := writer.Write(obj); err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
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
	errorWriter io.Writer,
	idLookup map[string]*jc.BaseUser,
) error {
	groups, err := f.jcc.ListGroups(ctx, jc.UserGroups)
	if err != nil {
		common.WriteErrorWithExitCode(errorWriter, err, 1)
		return err
	}

	for _, group := range groups {
		groupBytes, err := json.Marshal(group)
		if err != nil {
			common.WriteErrorWithExitCode(errorWriter, err, 1)
			continue
		}

		var obj map[string]interface{}
		if err := json.Unmarshal(groupBytes, &obj); err != nil {
			common.WriteErrorWithExitCode(errorWriter, err, 1)
			continue
		}

		usersInGroup, err := f.jcc.ExpandUsersInGroup(ctx, group.ID, idLookup)
		if err != nil {
			common.WriteErrorWithExitCode(errorWriter, err, 1)
		}

		usersInGroupBytes, err := json.Marshal(usersInGroup)
		if err != nil {
			common.WriteErrorWithExitCode(errorWriter, err, 1)
		} else {
			var users []map[string]interface{}
			if err := json.Unmarshal(usersInGroupBytes, &users); err != nil {
				common.WriteErrorWithExitCode(errorWriter, err, 1)
			}

			obj["users"] = users
		}

		if err := writer.Write(obj); err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
		}
	}

	return nil
}
