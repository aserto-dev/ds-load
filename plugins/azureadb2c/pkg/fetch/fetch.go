package fetch

import (
	"context"
	"errors"
	"io"
	"iter"

	"github.com/aserto-dev/ds-load/plugins/azureadb2c/pkg/azureclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/fetcher"
	"github.com/aserto-dev/msgraph-sdk-go/models"
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
		return fetcher.YieldError(err)
	}

	return fetcher.YieldMap(aadUsers, func(object any) ([]byte, error) {
		user, ok := object.(models.Userable)
		if !ok {
			return nil, errors.ErrUnsupported
		}

		writer := kiota.NewJsonSerializationWriter()

		err := user.Serialize(writer)
		if err != nil {
			return nil, err
		}

		objBytes, err := writer.GetSerializedContent()
		if err != nil {
			return nil, err
		}

		objString := "{" + string(objBytes) + "}"

		return []byte(objString), nil
	})
}

func (f *Fetcher) fetchGroups(ctx context.Context) iter.Seq2[map[string]any, error] {
	aadGroups, err := f.azureClient.ListGroups(ctx, f.groupProps)
	if err != nil {
		return fetcher.YieldError(err)
	}

	return fetcher.YieldMap(aadGroups, func(object any) ([]byte, error) {
		group, ok := object.(models.Groupable)
		if !ok {
			return nil, errors.ErrUnsupported
		}

		writer := kiota.NewJsonSerializationWriter()

		if err := group.Serialize(writer); err != nil {
			return nil, err
		}

		objBytes, err := writer.GetSerializedContent()
		if err != nil {
			return nil, err
		}

		objString := "{" + string(objBytes) + "}"

		return []byte(objString), nil
	})
}
