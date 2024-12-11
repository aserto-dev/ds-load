package oktaclient

import (
	"fmt"

	"github.com/okta/okta-sdk-golang/v5/okta"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OktaClient struct {
	User            okta.UserAPI
	Group           okta.GroupAPI
	RoleAssignments okta.RoleAssignmentAPI
}

func NewOktaClient(domain, token string, requestTimeout int64) (*OktaClient, error) {
	config, err := okta.NewConfiguration(
		okta.WithOrgUrl(fmt.Sprintf("https://%s", domain)),
		okta.WithToken(token),
		okta.WithConnectionTimeout(5),
		okta.WithRequestTimeout(requestTimeout),
		okta.WithRateLimitMaxBackOff(30),
		okta.WithRateLimitMaxRetries(3),
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
