package fusionauthclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
)

type FusionAuthClient struct {
	fusionauthClient *fusionauth.FusionAuthClient
	host             string
}

const defaultTimeout = 10 * time.Second

func NewFusionAuthClient(host, apiKey string) (*FusionAuthClient, error) {
	c := &FusionAuthClient{}

	httpClient := &http.Client{
		Timeout: defaultTimeout,
	}

	baseURL, _ := url.Parse(host)

	c.fusionauthClient = fusionauth.NewClient(httpClient, baseURL, apiKey)
	c.host = host

	return c, nil
}

func (c *FusionAuthClient) ListUsers(ctx context.Context) ([]fusionauth.User, error) {
	users := make([]fusionauth.User, 0)
	pageSize := 100
	startRow := 0

	for {
		searchRequest := fusionauth.SearchRequest{
			Search: fusionauth.UserSearchCriteria{
				BaseElasticSearchCriteria: fusionauth.BaseElasticSearchCriteria{
					QueryString: "*",
					BaseSearchCriteria: fusionauth.BaseSearchCriteria{
						NumberOfResults: pageSize,
						StartRow:        startRow,
					},
				},
			},
		}

		response, faErrs, err := c.fusionauthClient.SearchUsersByQueryWithContext(ctx, searchRequest)
		if err != nil {
			fmt.Println("Failed to list users:", err)
			return nil, err
		}

		if faErrs != nil {
			fmt.Println("Failed to list users:", faErrs)
			return nil, faErrs
		}

		users = append(users, response.Users...)
		startRow += pageSize

		if int64(startRow) >= response.Total {
			break
		}
	}

	return users, nil
}

func (c *FusionAuthClient) ListGroups(ctx context.Context) ([]fusionauth.Group, error) {
	groupResponse, err := c.fusionauthClient.RetrieveGroupsWithContext(ctx)
	return groupResponse.Groups, err
}

func (c *FusionAuthClient) GetUserByID(ctx context.Context, id string) (fusionauth.User, error) {
	userResponse, _, err := c.fusionauthClient.RetrieveUserWithContext(ctx, id)
	return userResponse.User, err
}

func (c *FusionAuthClient) GetUserByEmail(ctx context.Context, email string) (fusionauth.User, error) {
	userResponse, _, err := c.fusionauthClient.RetrieveUserByEmailWithContext(ctx, email)
	return userResponse.User, err
}
