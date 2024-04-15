package auth0client

import (
	"context"

	"github.com/auth0/go-auth0/management"
)

type Auth0Client struct {
	Mgmt *management.Management
}

func New(ctx context.Context, clientID, clientSecret, domain string) (*Auth0Client, error) {
	options := []management.Option{
		management.WithClientCredentials(
			ctx,
			clientID,
			clientSecret,
		),
	}

	mgmt, err := management.New(
		domain,
		options...,
	)
	if err != nil {
		return nil, err
	}

	return &Auth0Client{
		Mgmt: mgmt,
	}, nil
}
