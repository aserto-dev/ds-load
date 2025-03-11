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
		_, _ = errorWriter.Write([]byte(err.Error()))
		common.SetExitCode(1)
		return err
	}

	for _, user := range users {
		userBytes, err := json.Marshal(user)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
			continue
		}
		var obj map[string]interface{}
		err = json.Unmarshal(userBytes, &obj)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
			continue
		}

		err = writer.Write(obj)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
		}
	}

	if f.Groups {
		groups, err := f.jcc.ListGroups(ctx, jc.UserGroups)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
			return err
		}

		for _, group := range groups {
			groupBytes, err := json.Marshal(group)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
				continue
			}
			var obj map[string]interface{}
			err = json.Unmarshal(groupBytes, &obj)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
				continue
			}

			usersInGroup, err := f.jcc.GetUsersInGroup(ctx, group.ID)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
			} else {
				usersInGroupBytes, err := json.Marshal(usersInGroup)
				if err != nil {
					_, _ = errorWriter.Write([]byte(err.Error()))
					common.SetExitCode(1)
				} else {
					var users []map[string]interface{}
					err = json.Unmarshal(usersInGroupBytes, &users)
					if err != nil {
						_, _ = errorWriter.Write([]byte(err.Error()))
						common.SetExitCode(1)
					}
					obj["users"] = users
				}
			}
			err = writer.Write(obj)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
			}
		}
	}

	return nil
}
