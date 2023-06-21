package app

import (
	"context"
	"encoding/json"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/azuread/pkg/azureclient"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	kiota "github.com/microsoft/kiota-serialization-json-go"
)

type FetchCmd struct {
	Tenant       string `short:"t" help:"AzureAD tenant" env:"AZUREAD_TENANT" required:""`
	ClientID     string `short:"i" help:"AzureAD Client ID" env:"AZUREAD_CLIENT_ID" required:""`
	ClientSecret string `short:"s" help:"AzureAD Client Secret" env:"AZUREAD_CLIENT_SECRET" required:""`
	RefreshToken string `short:"r" help:"AzureAD Refresh Token" env:"AZUREAD_REFRESH_TOKEN" optional:""`
}

func (cmd *FetchCmd) Run(ctx *kong.Context) error {
	azureClient, err := createAzureAdClient(cmd.Tenant, cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken)
	if err != nil {
		return err
	}

	results := make(chan map[string]interface{}, 1)
	errors := make(chan error, 1)
	go func() {
		Fetch(azureClient, results, errors)
		close(results)
		close(errors)
	}()
	if err != nil {
		return err
	}

	go printErrors(errors)

	writer, err := js.NewJSONArrayWriter(os.Stdout)
	if err != nil {
		return err
	}
	defer writer.Close()
	for result := range results {
		err := writer.Write(result)
		if err != nil {
			return err
		}
	}
	return nil
}

func Fetch(azureClient *azureclient.AzureADClient, results chan map[string]interface{}, errors chan error) {
	aadUsers, err := azureClient.ListUsers()
	if err != nil {
		errors <- err
	}

	for _, user := range aadUsers.GetValue() {
		writer := kiota.NewJsonSerializationWriter()
		err := user.Serialize(writer)
		if err != nil {
			errors <- err
			return
		}
		userBytes, err := writer.GetSerializedContent()
		if err != nil {
			errors <- err
			return
		}

		userString := "{" + string(userBytes) + "}"
		var obj map[string]interface{}
		err = json.Unmarshal([]byte(userString), &obj)
		if err != nil {
			errors <- err
			return
		}
		results <- obj
	}
}

func createAzureAdClient(tenant, clientID, clientSecret, refreshToken string) (azureClient *azureclient.AzureADClient, err error) {
	if refreshToken != "" {
		return azureclient.NewAzureADClientWithRefreshToken(
			context.Background(),
			tenant,
			clientID,
			clientSecret,
			refreshToken)
	}

	return azureclient.NewAzureADClient(
		context.Background(),
		tenant,
		clientID,
		clientSecret)
}

func printErrors(errors chan error) {
	for err := range errors {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
	}
}
