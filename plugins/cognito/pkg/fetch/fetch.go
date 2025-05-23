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
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
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

func (f *Fetcher) Fetch(ctx context.Context, outputWriter io.Writer, errorWriter common.ErrorWriter) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	if f.groups {
		for obj, err := range f.fetchGroups() {
			errorWriter.Error(err, common.WithExitCode)

			err := writer.Write(obj)
			errorWriter.Error(err)
		}
	}

	users, err := f.cognitoClient.ListUsers(ctx)
	errorWriter.Error(err, common.WithExitCode)

	for _, user := range users {
		attributes := make(map[string]string)
		for _, attr := range user.Attributes {
			attributes[*attr.Name] = *attr.Value
		}

		userBytes, err := json.Marshal(user)
		if err != nil {
			errorWriter.Error(err, common.WithExitCode)
			return err
		}

		var obj map[string]any
		err = json.Unmarshal(userBytes, &obj)
		errorWriter.Error(err)

		obj["Attributes"] = attributes

		if f.groups {
			groups, err := f.cognitoClient.GetGroupsForUser(*user.Username)
			if err != nil {
				errorWriter.Error(err)
				continue
			}

			groupBytes, err := json.Marshal(groups.Groups)
			if err != nil {
				errorWriter.Error(err)
				return err
			}

			var grps []map[string]string
			err = json.Unmarshal(groupBytes, &grps)
			errorWriter.Error(err)

			obj["Groups"] = grps
		}

		err = writer.Write(obj)
		errorWriter.Error(err)
	}

	return nil
}

func (f *Fetcher) fetchGroups() iter.Seq2[map[string]any, error] {
	groups, err := f.cognitoClient.ListGroups()
	if err != nil {
		return fetcher.YieldError(err)
	}

	return fetcher.YieldMap(groups, func(group *cognitoidentityprovider.GroupType) ([]byte, error) {
		return json.Marshal(group)
	})
}
