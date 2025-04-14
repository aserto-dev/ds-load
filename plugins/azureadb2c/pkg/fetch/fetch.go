package fetch

import (
	"context"
	"encoding/json"
	"io"

	"github.com/aserto-dev/ds-load/plugins/azureadb2c/pkg/azureclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	kiota "github.com/microsoft/kiota-serialization-json-go"
)

type Fetcher struct {
	azureClient *azureclient.AzureADClient
	Groups      bool
	userProps   []string
	groupProps  []string
}

func New(ctx context.Context, client *azureclient.AzureADClient, userProps, groupProps []string) (*Fetcher, error) {
	return &Fetcher{
		azureClient: client,
		userProps:   userProps,
		groupProps:  groupProps,
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
		if err := f.fetchGroups(ctx, jsonWriter, errorWriter); err != nil {
			return err
		}
	}

	aadUsers, err := f.azureClient.ListUsers(ctx, f.Groups, f.userProps)
	if err != nil {
		common.WriteErrorWithExitCode(errorWriter, err, 1)
	}

	for _, user := range aadUsers {
		writer := kiota.NewJsonSerializationWriter()

		err := user.Serialize(writer)
		if err != nil {
			common.WriteErrorWithExitCode(errorWriter, err, 1)
			return err
		}

		if err := writeObject(jsonWriter, writer, errorWriter); err != nil {
			return err
		}
	}

	return nil
}

func (f *Fetcher) fetchGroups(ctx context.Context, jsonWriter *js.JSONArrayWriter, errorWriter io.Writer) error {
	aadGroups, err := f.azureClient.ListGroups(ctx, f.groupProps)
	if err != nil {
		common.WriteErrorWithExitCode(errorWriter, err, 1)
	}

	for _, group := range aadGroups {
		writer := kiota.NewJsonSerializationWriter()

		err := group.Serialize(writer)
		if err != nil {
			common.WriteErrorWithExitCode(errorWriter, err, 1)
			return err
		}

		if err := writeObject(jsonWriter, writer, errorWriter); err != nil {
			return err
		}
	}

	return nil
}

func writeObject(jsonWriter *js.JSONArrayWriter, writer *kiota.JsonSerializationWriter, errorWriter io.Writer) error {
	objBytes, err := writer.GetSerializedContent()
	if err != nil {
		common.WriteErrorWithExitCode(errorWriter, err, 1)
		return err
	}

	objString := "{" + string(objBytes) + "}"

	var obj map[string]any
	if err := json.Unmarshal([]byte(objString), &obj); err != nil {
		common.WriteErrorWithExitCode(errorWriter, err, 1)
		return err
	}

	if err := jsonWriter.Write(obj); err != nil {
		_, _ = errorWriter.Write([]byte(err.Error()))
	}

	return nil
}
