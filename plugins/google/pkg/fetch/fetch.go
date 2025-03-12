package fetch

import (
	"context"
	"encoding/json"
	"io"

	"github.com/aserto-dev/ds-load/plugins/google/pkg/googleclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
)

type Fetcher struct {
	gClient *googleclient.GoogleClient
	Groups  bool
}

func New(client *googleclient.GoogleClient) (*Fetcher, error) {
	return &Fetcher{
		gClient: client,
	}, nil
}

func (f *Fetcher) WithGroups(groups bool) *Fetcher {
	f.Groups = groups
	return f
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	users, err := f.gClient.ListUsers()
	if err != nil {
		common.WriteErrorWithExitCode(errorWriter, err, 1)
		return err
	}

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
	}

	if f.Groups {
		groups, err := f.gClient.ListGroups()
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

			usersInGroup, err := f.gClient.GetUsersInGroup(group.Id)
			if err != nil {
				common.WriteErrorWithExitCode(errorWriter, err, 1)
			} else {
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
			}

			if err := writer.Write(obj); err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
			}
		}
	}

	return nil
}
