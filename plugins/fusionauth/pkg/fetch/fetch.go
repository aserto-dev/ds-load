package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aserto-dev/ds-load/plugins/fusionauth/pkg/client"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
)

type Fetcher struct {
	fusionauthClient *client.FusionAuthClient
	groups           bool
	host             string
}

func New(client *client.FusionAuthClient) (*Fetcher, error) {
	return &Fetcher{
		fusionauthClient: client,
	}, nil
}

func (f *Fetcher) WithGroups(groups bool) *Fetcher {
	f.groups = groups
	return f
}

func (f *Fetcher) WithHost(host string) *Fetcher {
	f.host = host
	return f
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter io.Writer, errorWriter common.ErrorWriter) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	users, err := f.fusionauthClient.ListUsers(ctx)
	errorWriter.Error(err)

	for i := range users {
		user := &users[i]

		userBytes, err := json.Marshal(user)
		if err != nil {
			errorWriter.Error(err)
			return err
		}

		var obj map[string]any
		if err := json.Unmarshal(userBytes, &obj); err != nil {
			errorWriter.Error(err)
			return err
		}

		if user.ImageUrl != "" {
			obj["picture"] = fmt.Sprintf("%s%s", f.host, user.ImageUrl)
		}

		err = writer.Write(obj)
		errorWriter.Error(err)
	}

	if f.groups {
		groups, err := f.fusionauthClient.ListGroups(ctx)
		errorWriter.Error(err)

		for i := range groups {
			group := &groups[i]
			err := writer.Write(group)
			errorWriter.Error(err)
		}
	}

	return nil
}
