package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aserto-dev/ds-load/plugins/fusionauth/pkg/fusionauthclient"
	"github.com/aserto-dev/ds-load/sdk/common/js"
)

type Fetcher struct {
	fusionauthClient *fusionauthclient.FusionAuthClient
	groups           bool
	host             string
}

func New(client *fusionauthclient.FusionAuthClient) (*Fetcher, error) {
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

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	users, err := f.fusionauthClient.ListUsers(ctx)
	if err != nil {
		_, _ = errorWriter.Write([]byte(err.Error()))
	}

	for i := range users {
		user := &users[i]
		userBytes, err := json.Marshal(user)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			return err
		}
		var obj map[string]interface{}
		err = json.Unmarshal(userBytes, &obj)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			return err
		}
		if user.ImageUrl != "" {
			obj["picture"] = fmt.Sprintf("%s%s", f.host, user.ImageUrl)
		}

		err = writer.Write(obj)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
		}
	}

	if f.groups {
		groups, err := f.fusionauthClient.ListGroups(ctx)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
		}

		for _, group := range groups {
			err = writer.Write(group)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
			}
		}
	}

	return nil
}
