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

	Tenant       string
	ClientID     string
	ClientSecret string
	RefreshToken string
}

func New(ctx context.Context, tenant, clientID, clientSecret, refreshToken string) (*Fetcher, error) {
	azureClient, err := createAzureAdClient(ctx, tenant, clientID, clientSecret, refreshToken)
	if err != nil {
		return nil, err
	}

	return &Fetcher{
		azureClient: azureClient,
	}, nil
}

func createAzureAdClient(ctx context.Context, tenant, clientID, clientSecret, refreshToken string) (azureClient *azureclient.AzureADClient, err error) {
	if refreshToken != "" {
		return azureclient.NewAzureADClientWithRefreshToken(
			ctx,
			tenant,
			clientID,
			clientSecret,
			refreshToken)
	}

	return azureclient.NewAzureADClient(
		ctx,
		tenant,
		clientID,
		clientSecret)
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	jsonWriter, err := js.NewJSONArrayWriter(outputWriter)
	if err != nil {
		return err
	}
	defer jsonWriter.Close()

	aadUsers, err := f.azureClient.ListUsers(ctx)
	if err != nil {
		_, _ = errorWriter.Write([]byte(err.Error()))
		common.SetExitCode(1)
	}

	for _, user := range aadUsers.GetValue() {
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
		_, _ = errorWriter.Write([]byte(err.Error()))
	}

	return nil
}
