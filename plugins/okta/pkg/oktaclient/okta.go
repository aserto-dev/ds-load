package oktaclient

import (
	"github.com/okta/okta-sdk-golang/v5/okta"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OktaClient struct {
	User            okta.UserAPI
	Group           okta.GroupAPI
	RoleAssignments okta.RoleAssignmentAPI
}

const (
	defaultConnectionTimeout = 5
	defaultRateLimitBackoff  = 30
	defaultRateLimitRetries  = 3
)

func NewOktaClient(domain, token string, requestTimeout int64) (*OktaClient, error) {
	config, err := okta.NewConfiguration(
		okta.WithOrgUrl("https://"+domain),
		okta.WithToken(token),
		okta.WithConnectionTimeout(defaultConnectionTimeout),
		okta.WithRequestTimeout(requestTimeout),
		okta.WithRateLimitMaxBackOff(defaultRateLimitBackoff),
		okta.WithRateLimitMaxRetries(defaultRateLimitRetries),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create Okta configuration: %s", err.Error())
	}

	client := okta.NewAPIClient(config)

	return &OktaClient{
		User:            client.UserAPI,
		Group:           client.GroupAPI,
		RoleAssignments: client.RoleAssignmentAPI,
	}, nil
}
