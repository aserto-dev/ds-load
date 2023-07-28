package fetch

import (
	"context"
	"encoding/json"
	"io"

	"github.com/aserto-dev/ds-load/plugins/cognito/pkg/cognitoclient"
	"github.com/aserto-dev/ds-load/sdk/common/js"
)

type Fetcher struct {
	cognitoClient *cognitoclient.CognitoClient
	groups        bool
}

func New(accessKey, secretKey, userPoolID, region string) (*Fetcher, error) {
	cognitoClient, err := cognitoclient.NewCognitoClient(
		accessKey,
		secretKey,
		userPoolID,
		region)
	if err != nil {
		return nil, err
	}

	return &Fetcher{
		cognitoClient: cognitoClient,
	}, nil
}

func (f *Fetcher) WithGroups(groups bool) *Fetcher {
	f.groups = groups
	return f
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	writer, err := js.NewJSONArrayWriter(outputWriter)
	if err != nil {
		return err
	}
	defer writer.Close()

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
			_, _ = errorWriter.Write([]byte(err.Error()))
			return err
		}
		var obj map[string]interface{}
		err = json.Unmarshal(userBytes, &obj)
		if err != nil {
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
			err = json.Unmarshal(groupBytes, &grps)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
			}
			obj["Groups"] = grps
		}

		err = writer.Write(obj)
		_, _ = errorWriter.Write([]byte(err.Error()))
	}

	return nil
}
