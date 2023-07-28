package fetch

import (
	"context"
	"encoding/json"
	"github.com/aserto-dev/ds-load/plugins/google/pkg/googleclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"io"
)

type Fetcher struct {
	gClient *googleclient.GoogleClient
	Groups  bool
}

func New(ctx context.Context, clientID, clientSecret, refrestToken, customer string) (*Fetcher, error) {
	gClent, err := googleclient.NewGoogleClient(
		ctx,
		clientID,
		clientSecret,
		refrestToken,
		customer)

	if err != nil {
		return nil, err
	}

	return &Fetcher{
		gClient: gClent,
	}, nil
}

func (f *Fetcher) WithGroups(groups bool) *Fetcher {
	f.Groups = groups
	return f
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	writer, err := js.NewJSONArrayWriter(outputWriter)
	if err != nil {
		return err
	}
	defer writer.Close()

	users, err := f.gClient.ListUsers()
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

		writer.Write(obj)
	}

	if f.Groups {
		groups, err := f.gClient.ListGroups()
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

			usersInGroup, err := f.gClient.GetUsersInGroup(group.Id)
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
			writer.Write(obj)
		}
	}

	return nil
}
