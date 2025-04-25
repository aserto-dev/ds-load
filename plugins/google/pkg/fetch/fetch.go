package fetch

import (
	"context"
	"encoding/json"
	"io"
	"iter"

	"github.com/aserto-dev/ds-load/plugins/google/pkg/googleclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/fetcher"

	admin "google.golang.org/api/admin/directory/v1"
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

func (f *Fetcher) Fetch(ctx context.Context, outputWriter io.Writer, errorWriter common.ErrorWriter) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	for user, err := range f.fetchUsers() {
		if err != nil {
			errorWriter.Error(err)
			continue
		}

		err := writer.Write(user)
		errorWriter.ErrorNoExitCode(err)
	}

	if f.Groups {
		for group, err := range f.fetchGroups() {
			if err != nil {
				errorWriter.Error(err)
				continue
			}

			err := writer.Write(group)
			errorWriter.ErrorNoExitCode(err)
		}
	}

	return nil
}

func (f *Fetcher) fetchUsers() iter.Seq2[map[string]any, error] {
	users, err := f.gClient.ListUsers()
	if err != nil {
		return fetcher.YieldError(err)
	}

	return fetcher.YieldMap(users, func(user *admin.User) ([]byte, error) {
		return json.Marshal(user)
	})
}

func (f *Fetcher) fetchGroups() iter.Seq2[map[string]any, error] {
	groups, err := f.gClient.ListGroups()
	if err != nil {
		return fetcher.YieldError(err)
	}

	return func(yield func(map[string]any, error) bool) {
		for _, group := range groups {
			groupBytes, err := json.Marshal(group)
			if err != nil {
				if !yield(nil, err) {
					return
				}
			}

			var obj map[string]any
			if err := json.Unmarshal(groupBytes, &obj); err != nil {
				if !yield(nil, err) {
					return
				}
			}

			users, err := f.fetchUsersInGroup(group.Id)
			if err != nil {
				if !yield(nil, err) {
					return
				}
			}

			obj["users"] = users
			if !(yield(obj, nil)) {
				return
			}
		}
	}
}

func (f *Fetcher) fetchUsersInGroup(groupId string) ([]map[string]any, error) {
	usersInGroup, err := f.gClient.GetUsersInGroup(groupId)
	if err != nil {
		return nil, err
	}

	usersInGroupBytes, err := json.Marshal(usersInGroup)
	if err != nil {
		return nil, err
	}

	var users []map[string]any
	if err := json.Unmarshal(usersInGroupBytes, &users); err != nil {
		return nil, err
	}

	return users, nil
}
