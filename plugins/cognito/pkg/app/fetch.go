package app

import (
	"context"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/cognito/pkg/fetch"
)

type FetchCmd struct {
	AccessKey  string `short:"k" help:"AWS Access Key" env:"AWS_ACCESS_KEY" required:""`
	SecretKey  string `short:"s" help:"AWS Secret Key" env:"AWS_SECRET_KEY" required:""`
	UserPoolID string `short:"p" help:"AWS Cognito User Pool ID" env:"AWS_COGNITO_USER_POOL_ID" required:""`
	Region     string `short:"r" help:"AWS Region" env:"AWS_REGION" required:""`
	Groups     bool   `short:"g" help:"Retrieve Cognito groups" env:"AWS_COGNITO_GROUPS" default:"false" negatable:""`
}

func (cmd *FetchCmd) Run(kongCtx *kong.Context) error {
	fetcher, err := fetch.New(cmd.AccessKey, cmd.SecretKey, cmd.UserPoolID, cmd.Region)
	if err != nil {
		return err
	}
	fetcher = fetcher.WithGroups(cmd.Groups)

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	return fetcher.Fetch(timeoutCtx, os.Stdout, os.Stderr)
}
