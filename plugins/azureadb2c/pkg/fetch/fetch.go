package fetch

import (
	"context"
	"encoding/json"
	"io"
	"iter"

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
		for obj, err := range f.fetchGroups(ctx) {
			if err != nil {
				common.WriteErrorWithExitCode(errorWriter, err, 1)
			}

			if err := jsonWriter.Write(obj); err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
			}
		}
	}

	for user, err := range f.fetchUsers(ctx) {
		if err != nil {
			common.WriteErrorWithExitCode(errorWriter, err, 1)
		}

		if err := jsonWriter.Write(user); err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
		}
	}

	return nil
}

func (f *Fetcher) fetchUsers(ctx context.Context) iter.Seq2[map[string]any, error] {
	aadUsers, err := f.azureClient.ListUsers(ctx, f.Groups, f.userProps)
	if err != nil {
		return func(yield func(map[string]any, error) bool) {
			if !yield(nil, err) {
				return
			}
		}
	}

	return func(yield func(map[string]any, error) bool) {
		for _, user := range aadUsers {
			writer := kiota.NewJsonSerializationWriter()

			err := user.Serialize(writer)
			if err != nil {
				if !yield(nil, err) {
					return
				}
			}

			objBytes, err := writer.GetSerializedContent()
			if err != nil {
				if !(yield(nil, err)) {
					return
				}
			}

			objString := "{" + string(objBytes) + "}"

			var obj map[string]any
			if err := json.Unmarshal([]byte(objString), &obj); err != nil {
				if !(yield(obj, err)) {
					return
				}
			}
		}
	}
}

func (f *Fetcher) fetchGroups(ctx context.Context) iter.Seq2[map[string]any, error] {
	aadGroups, err := f.azureClient.ListGroups(ctx, f.groupProps)
	if err != nil {
		return func(yield func(map[string]any, error) bool) {
			if !(yield(nil, err)) {
				return
			}
		}
	}

	return func(yield func(map[string]any, error) bool) {
		for _, group := range aadGroups {
			writer := kiota.NewJsonSerializationWriter()

			if err := group.Serialize(writer); err != nil {
				if !yield(nil, err) {
					return
				}
			}

			objBytes, err := writer.GetSerializedContent()
			if err != nil {
				if !(yield(nil, err)) {
					return
				}
			}

			objString := "{" + string(objBytes) + "}"

			var obj map[string]any
			if err := json.Unmarshal([]byte(objString), &obj); err != nil {
				if !(yield(obj, err)) {
					return
				}
			}
		}
	}
}
