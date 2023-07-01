package app

import (
	"context"
	"encoding/json"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/cognito/pkg/cognitoclient"
	"github.com/aserto-dev/ds-load/sdk/plugin"
)

type FetchCmd struct {
	AccessKey  string `short:"k" help:"AWS Access Key" env:"AWS_ACCESS_KEY" required:""`
	SecretKey  string `short:"s" help:"AWS Secret Key" env:"AWS_SECRET_KEY" required:""`
	UserPoolID string `short:"p" help:"AWS Cognito User Pool ID" env:"AWS_COGNITO_USER_POOL_ID" required:""`
	Region     string `short:"r" help:"AWS Region" env:"AWS_REGION" required:""`
}

func (cmd *FetchCmd) Run(ctx *kong.Context) error {
	cognitoClient, err := createCognitoClient(cmd.AccessKey, cmd.SecretKey, cmd.UserPoolID, cmd.Region)
	if err != nil {
		return err
	}

	results := make(chan map[string]interface{}, 1)
	errors := make(chan error, 1)
	go func() {
		Fetch(cognitoClient, results, errors)
		close(results)
		close(errors)
	}()
	return plugin.NewDSPlugin().WriteFetchOutput(results, errors, false)
}

func Fetch(cognitoClient *cognitoclient.CognitoClient, results chan map[string]interface{}, errors chan error) {
	users, err := cognitoClient.ListUsers()
	if err != nil {
		errors <- err
	}

	for _, user := range users.Users {
		attributes := make(map[string]string)
		for _, attr := range user.Attributes {
			attributes[*attr.Name] = *attr.Value
		}

		userBytes, err := json.Marshal(user)
		if err != nil {
			errors <- err
			return
		}
		var obj map[string]interface{}
		err = json.Unmarshal(userBytes, &obj)
		if err != nil {
			errors <- err
		}
		obj["Attributes"] = attributes
		results <- obj
	}
}

func createCognitoClient(accessKey, accessSecretKey, userPoolID, region string) (cognitoClient *cognitoclient.CognitoClient, err error) {
	return cognitoclient.NewCognitoClient(
		context.Background(),
		accessKey,
		accessSecretKey,
		userPoolID,
		region)
}
