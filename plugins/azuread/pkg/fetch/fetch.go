package fetch

import (
	"context"
	"encoding/json"
	"io"

	"github.com/aserto-dev/ds-load/plugins/azuread/pkg/azureclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	kiota "github.com/microsoft/kiota-serialization-json-go"
)

type Fetcher struct {
	azureClient *azureclient.AzureADClient
	Groups      bool
}

func New(ctx context.Context, client *azureclient.AzureADClient) (*Fetcher, error) {
	return &Fetcher{
		azureClient: client,
	}, nil
}

func (f *Fetcher) WithGroups(groups bool) *Fetcher {
	f.Groups = groups
	return f
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	jsonWriter := js.NewJSONArrayWriter(outputWriter)
	defer jsonWriter.Close()

	if f.Groups {
		aadGroups, err := f.azureClient.ListGroups(ctx)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
		}

		for _, group := range aadGroups {
			writer := kiota.NewJsonSerializationWriter()
			err := group.Serialize(writer)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
				return err
			}
			groupBytes, err := writer.GetSerializedContent()
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
				return err
			}

			groupString := "{" + string(groupBytes) + "}"
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(groupString), &obj)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
				return err
			}
			err = jsonWriter.Write(obj)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
			}
		}
	}

	aadUsers, err := f.azureClient.ListUsers(ctx, f.Groups)
	if err != nil {
		_, _ = errorWriter.Write([]byte(err.Error()))
		common.SetExitCode(1)
	}

	for _, user := range aadUsers {
		writer := kiota.NewJsonSerializationWriter()
		err := user.Serialize(writer)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
			return err
		}
		userBytes, err := writer.GetSerializedContent()
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
			return err
		}

		userString := "{" + string(userBytes) + "}"
		var obj map[string]interface{}
		err = json.Unmarshal([]byte(userString), &obj)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
			return err
		}
		err = jsonWriter.Write(obj)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
		}
	}

	return nil
}
