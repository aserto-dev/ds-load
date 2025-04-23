package fetch

import (
	"context"
	"encoding/json"
	"io"
	"iter"

	"github.com/aserto-dev/ds-load/plugins/cognito/pkg/cognitoclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/fetcher"
)

type Fetcher struct {
	cognitoClient *cognitoclient.CognitoClient
	groups        bool
}

func New(client *cognitoclient.CognitoClient) (*Fetcher, error) {
	return &Fetcher{
		cognitoClient: client,
	}, nil
}

func (f *Fetcher) WithGroups(groups bool) *Fetcher {
	f.groups = groups
	return f
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	if f.groups {
		for obj, err := range f.fetchGroups() {
			if err != nil {
				common.WriteErrorWithExitCode(errorWriter, err, 1)
			}

			if err := writer.Write(obj); err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
			}
		}
	}

	users, err := f.cognitoClient.ListUsers(ctx)
	if err != nil {
		_, _ = errorWriter.Write([]byte(err.Error()))
	}

	for _, user := range users {
		attributes := make(map[string]string)
		for _, attr := range user.Attributes {
			attributes[*attr.Name] = *attr.Value
		}

		userBytes, err := json.Marshal(user)
		if err != nil {
			common.WriteErrorWithExitCode(errorWriter, err, 1)
			return err
		}

		var obj map[string]interface{}
		if err := json.Unmarshal(userBytes, &obj); err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
		}

		obj["Attributes"] = attributes

		if f.groups {
			groups, err := f.cognitoClient.GetGroupsForUser(*user.Username)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				continue
			}

			groupBytes, err := json.Marshal(groups.Groups)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				return err
			}

			var grps []map[string]string
			if err := json.Unmarshal(groupBytes, &grps); err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
			}

			obj["Groups"] = grps
		}

		if err := writer.Write(obj); err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
		}
	}

	return nil
}

func (f *Fetcher) fetchGroups() iter.Seq2[map[string]any, error] {
	groups, err := f.cognitoClient.ListGroups()
	if err != nil {
		return fetcher.YieldError(err)
	}

	return fetcher.YieldMap(json.Marshal, groups)
}
