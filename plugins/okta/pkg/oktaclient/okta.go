package oktaclient

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OktaClient struct {
	User  *okta.UserResource
	Group *okta.GroupResource
}

func NewOktaClient(ctx context.Context, domain, token string, requestTimeout int64) (*OktaClient, error) {
	_, client, err := okta.NewClient(
		ctx,
		okta.WithOrgUrl(fmt.Sprintf("https://%s", domain)),
		okta.WithToken(token),
		okta.WithConnectionTimeout(5),
		okta.WithRequestTimeout(requestTimeout),
		okta.WithRateLimitMaxBackOff(30),
		okta.WithRateLimitMaxRetries(3),
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to Okta: %s", err.Error())
	}

	return &OktaClient{
		User:  client.User,
		Group: client.Group,
	}, nil
}
