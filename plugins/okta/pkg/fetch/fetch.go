package fetch

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

type Fetcher struct {
	oktaClient oktaclient.OktaClient
	Groups     bool
	Roles      bool
}

func New(ctx context.Context, client oktaclient.OktaClient) (*Fetcher, error) {

	return &Fetcher{
		oktaClient: client,
	}, nil
}

func (fetcher *Fetcher) WithGroups(groups bool) *Fetcher {
	fetcher.Groups = groups
	return fetcher
}

func (fetcher *Fetcher) WithRoles(roles bool) *Fetcher {
	fetcher.Roles = roles
	return fetcher
}

func (fetcher *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	var response *okta.Response
	var users []*okta.User
	var err error

	writer, err := js.NewJSONArrayWriter(outputWriter)
	if err != nil {
		return err
	}
	defer writer.Close()

	users, response, err = fetcher.oktaClient.ListUsers(ctx, nil)
	if err != nil {
		_, _ = errorWriter.Write([]byte(err.Error()))
		common.SetExitCode(1)
		return err
	}

	for {
		logIfRateLimitExceeded(response, errorWriter)

		for _, user := range users {
			userResult, err := fetcher.processUser(ctx, user, errorWriter)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
			}
			err = writer.Write(userResult)
			_, _ = errorWriter.Write([]byte(err.Error()))
		}

		if response != nil && response.HasNextPage() {
			response, err = response.Next(ctx, &users)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
			}
		} else {
			break
		}
	}

	return nil
}

func (fetcher *Fetcher) processUser(ctx context.Context, user *okta.User, errorWriter io.Writer) (map[string]interface{}, error) {
	userBytes, err := json.Marshal(user)
	if err != nil {
		common.SetExitCode(1)
		return nil, err
	}
	var userResult map[string]interface{}
	err = json.Unmarshal(userBytes, &userResult)
	if err != nil {
		common.SetExitCode(1)
		return nil, err
	}

	if fetcher.Groups {
		// Write all user groups
		groups, err := fetcher.getGroups(ctx, user.Id, errorWriter)
		if err != nil {
			common.SetExitCode(1)
			return nil, err
		}
		userResult["groups"] = groups
	}

	if fetcher.Roles {
		// Write all user roles
		roles, err := fetcher.getRoles(ctx, user.Id, errorWriter)
		if err != nil {
			common.SetExitCode(1)
			return nil, err
		}
		userResult["roles"] = roles
	}
	return userResult, nil
}

func (fetcher *Fetcher) getGroups(ctx context.Context, userID string, errorWriter io.Writer) ([]map[string]interface{}, error) {
	var response *okta.Response
	var result []map[string]interface{}
	var groups []*okta.Group
	var err error

	groups, response, err = fetcher.oktaClient.ListUserGroups(ctx, userID)
	if err != nil {
		return nil, err
	}

	for {
		logIfRateLimitExceeded(response, errorWriter)

		for _, group := range groups {
			groupBytes, err := json.Marshal(group)
			if err != nil {
				return nil, err
			}
			var obj map[string]interface{}
			err = json.Unmarshal(groupBytes, &obj)
			if err != nil {
				return nil, err
			}
			result = append(result, obj)
		}

		if response != nil && response.HasNextPage() {
			response, err = response.Next(ctx, &groups)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return result, nil
}

func (fetcher *Fetcher) getRoles(ctx context.Context, userID string, errorWriter io.Writer) ([]map[string]interface{}, error) {
	var response *okta.Response
	var result []map[string]interface{}
	var roles []*okta.Role
	var err error

	roles, response, err = fetcher.oktaClient.ListAssignedRolesForUser(ctx, userID, nil)
	if err != nil {
		return nil, err
	}

	for {
		logIfRateLimitExceeded(response, errorWriter)

		for _, role := range roles {
			roleBytes, err := json.Marshal(role)
			if err != nil {
				return nil, err
			}
			var obj map[string]interface{}
			err = json.Unmarshal(roleBytes, &obj)
			if err != nil {
				return nil, err
			}
			result = append(result, obj)
		}

		if response != nil && response.HasNextPage() {
			response, err = response.Next(ctx, &roles)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return result, nil
}

func logIfRateLimitExceeded(resp *okta.Response, errorWriter io.Writer) {
	if resp.Response != nil && resp.StatusCode == http.StatusTooManyRequests {
		_, _ = errorWriter.Write([]byte("Rate limit exceeded"))
	}
}
