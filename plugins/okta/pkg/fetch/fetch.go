//nolint:dupl // similar code in getting group details
package fetch

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

type Fetcher struct {
	oktaClient *oktaclient.OktaClient
	Groups     bool
	Roles      bool
}

func New(ctx context.Context, client *oktaclient.OktaClient) (*Fetcher, error) {
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

func (fetcher *Fetcher) Fetch(ctx context.Context, outputWriter io.Writer, errorWriter common.ErrorWriter) error {
	writer := js.NewJSONArrayWriter(outputWriter)
	defer writer.Close()

	if fetcher.Roles {
		err := fetcher.fetchGroups(ctx, writer, errorWriter)
		if err != nil {
			return err
		}
	}

	return fetcher.fetchUsers(ctx, writer, errorWriter)
}

func (fetcher *Fetcher) fetchUsers(ctx context.Context, writer *js.JSONArrayWriter, errorWriter common.ErrorWriter) error {
	users, response, err := fetcher.oktaClient.User.ListUsers(ctx).Execute()
	if err != nil {
		errorWriter.Error(err)
		return err
	}

	for {
		logIfRateLimitExceeded(response, errorWriter)

		for i := range users {
			user := &users[i]

			userResult, err := fetcher.processUser(ctx, user, errorWriter)
			errorWriter.Error(err)

			err = writer.Write(userResult)
			errorWriter.Error(err)
		}

		if response != nil && response.HasNextPage() {
			response, err = response.Next(&users)
			errorWriter.Error(err)
		} else {
			break
		}
	}

	return nil
}

func (fetcher *Fetcher) fetchGroups(ctx context.Context, writer *js.JSONArrayWriter, errorWriter common.ErrorWriter) error {
	groups, response, err := fetcher.oktaClient.Group.ListGroups(ctx).Execute()
	if err != nil {
		errorWriter.Error(err)
		return err
	}

	for {
		logIfRateLimitExceeded(response, errorWriter)

		for _, group := range groups {
			groupResult, err := fetcher.processGroup(ctx, &group, errorWriter)
			errorWriter.Error(err, common.WithExitCode)

			err = writer.Write(groupResult)
			errorWriter.Error(err)
		}

		if response != nil && response.HasNextPage() {
			response, err = response.Next(&groups)
			errorWriter.Error(err)
		} else {
			break
		}
	}

	return nil
}

func (fetcher *Fetcher) processUser(ctx context.Context, user *okta.User, errorWriter common.ErrorWriter) (map[string]any, error) {
	userBytes, err := json.Marshal(user)
	if err != nil {
		common.SetExitCode(1)
		return nil, err
	}

	var userResult map[string]any

	if err := json.Unmarshal(userBytes, &userResult); err != nil {
		common.SetExitCode(1)
		return nil, err
	}

	if fetcher.Groups {
		// Write all user groups
		groups, err := fetcher.getGroups(ctx, *user.Id, errorWriter)
		if err != nil {
			common.SetExitCode(1)
			return nil, err
		}

		userResult["groups"] = groups
	}

	if fetcher.Roles {
		// Write all user roles
		roles, err := fetcher.getUserRoles(ctx, *user.Id, errorWriter)
		if err != nil {
			common.SetExitCode(1)
			return nil, err
		}

		userResult["roles"] = roles
	}

	return userResult, nil
}

func (fetcher *Fetcher) processGroup(ctx context.Context, group *okta.Group, errorWriter common.ErrorWriter) (map[string]any, error) {
	userBytes, err := json.Marshal(group)
	if err != nil {
		common.SetExitCode(1)
		return nil, err
	}

	var groupResult map[string]any

	if err := json.Unmarshal(userBytes, &groupResult); err != nil {
		common.SetExitCode(1)
		return nil, err
	}

	if fetcher.Roles {
		// Write all group roles
		roles, err := fetcher.getGroupRoles(ctx, *group.Id, errorWriter)
		if err != nil {
			common.SetExitCode(1)
			return nil, err
		}

		groupResult["roles"] = roles
	}

	return groupResult, nil
}

func (fetcher *Fetcher) getGroups(ctx context.Context, userID string, errorWriter common.ErrorWriter) ([]map[string]any, error) {
	var (
		response *okta.APIResponse
		result   []map[string]any
		groups   []okta.Group
		err      error
	)

	groups, response, err = fetcher.oktaClient.User.ListUserGroups(ctx, userID).Execute()
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

			var obj map[string]any
			if err := json.Unmarshal(groupBytes, &obj); err != nil {
				return nil, err
			}

			result = append(result, obj)
		}

		if response != nil && response.HasNextPage() {
			response, err = response.Next(&groups)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return result, nil
}

func (fetcher *Fetcher) getUserRoles(ctx context.Context, userID string, errorWriter common.ErrorWriter) ([]map[string]any, error) {
	var (
		response *okta.APIResponse
		result   []map[string]any
		roles    []okta.Role
		err      error
	)

	roles, response, err = fetcher.oktaClient.RoleAssignments.ListAssignedRolesForUser(ctx, userID).Execute()
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

			var obj map[string]any

			if err := json.Unmarshal(roleBytes, &obj); err != nil {
				return nil, err
			}

			result = append(result, obj)
		}

		if response != nil && response.HasNextPage() {
			response, err = response.Next(&roles)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return result, nil
}

func (fetcher *Fetcher) getGroupRoles(ctx context.Context, groupID string, errorWriter common.ErrorWriter) ([]map[string]any, error) {
	var (
		response *okta.APIResponse
		result   []map[string]any
		roles    []okta.Role
		err      error
	)

	roles, response, err = fetcher.oktaClient.RoleAssignments.ListGroupAssignedRoles(ctx, groupID).Execute()
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

			var obj map[string]any
			if err := json.Unmarshal(roleBytes, &obj); err != nil {
				return nil, err
			}

			result = append(result, obj)
		}

		if response != nil && response.HasNextPage() {
			response, err = response.Next(&roles)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return result, nil
}

func logIfRateLimitExceeded(resp *okta.APIResponse, errorWriter common.ErrorWriter) {
	if resp.Response != nil && resp.StatusCode == http.StatusTooManyRequests {
		_, _ = errorWriter.Write([]byte("Rate limit exceeded"))
	}
}
