//molint:dupl
package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/aserto-dev/ds-load/plugins/keycloak/pkg/client"
	"github.com/stretchr/testify/require"
)

const (
	EnvKeyCloakClientID     string = `KEYCLOAK_CLIENT_ID`     //
	EnvKeycloakClientSecret string = `KEYCLOAK_CLIENT_SECRET` //
	EnvKeyCloakTokenURL     string = `KEYCLOAK_TOKEN_URL`     //nolint: gosec // not a secret
)

func TestMain(m *testing.M) {
	if os.Getenv(EnvKeycloakClientSecret) == "" {
		fmt.Fprintf(os.Stderr, "env %q not set, tests skipped", EnvKeycloakClientSecret)
		return
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func keycloakClientConfig() *client.KeycloakConfig {
	return &client.KeycloakConfig{
		ClientID:     os.Getenv(EnvKeyCloakClientID),
		ClientSecret: os.Getenv(EnvKeycloakClientSecret),
		TokenURL:     os.Getenv(EnvKeyCloakTokenURL),
	}
}

func TestListUsers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kcc, err := client.NewKeycloakClient(ctx, keycloakClientConfig())
	require.NoError(t, err)

	users, err := kcc.ListUsers(ctx)
	require.NoError(t, err)

	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	if err := enc.Encode(users); err != nil {
		require.NoError(t, err)
	}
}

func TestListGroups(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kcc, err := client.NewKeycloakClient(ctx, keycloakClientConfig())
	require.NoError(t, err)

	groups, err := kcc.ListGroups(ctx)
	require.NoError(t, err)

	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	if err := enc.Encode(groups); err != nil {
		require.NoError(t, err)
	}
}

func TestListRoles(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kcc, err := client.NewKeycloakClient(ctx, keycloakClientConfig())
	require.NoError(t, err)

	roles, err := kcc.ListRoles(ctx)
	require.NoError(t, err)

	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	if err := enc.Encode(roles); err != nil {
		require.NoError(t, err)
	}
}

func TestGetRolesOfUser(t *testing.T) { //nolint:dupl
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kcc, err := client.NewKeycloakClient(ctx, keycloakClientConfig())
	require.NoError(t, err)

	users, err := kcc.ListUsers(ctx)
	require.NoError(t, err)
	require.NotNil(t, users)

	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	for _, user := range users {
		realmMappings, err := kcc.GetRolesOfUser(ctx, user.ID)
		require.NoError(t, err)
		require.NotNil(t, realmMappings)

		if err := enc.Encode(user); err != nil {
			require.NoError(t, err)
		}

		for _, role := range realmMappings.RealmMappings {
			if err := enc.Encode(role); err != nil {
				require.NoError(t, err)
			}
		}
	}
}

func TestGetRolesOfGroup(t *testing.T) { //nolint:dupl
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kcc, err := client.NewKeycloakClient(ctx, keycloakClientConfig())
	require.NoError(t, err)

	groups, err := kcc.ListGroups(ctx)
	require.NoError(t, err)
	require.NotNil(t, groups)

	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	for _, group := range groups {
		realmMappings, err := kcc.GetRolesOfGroup(ctx, group.ID)
		require.NoError(t, err)
		require.NotNil(t, realmMappings)

		if err := enc.Encode(group); err != nil {
			require.NoError(t, err)
		}

		for _, role := range realmMappings.RealmMappings {
			if err := enc.Encode(role); err != nil {
				require.NoError(t, err)
			}
		}
	}
}

func TestGetUsersOfRole(t *testing.T) { //nolint:dupl
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kcc, err := client.NewKeycloakClient(ctx, keycloakClientConfig())
	require.NoError(t, err)

	roles, err := kcc.ListRoles(ctx)
	require.NoError(t, err)
	require.NotNil(t, roles)

	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	for _, role := range roles {
		users, err := kcc.GetUsersOfRole(ctx, role.Name)
		require.NoError(t, err)
		require.NotNil(t, users)

		if err := enc.Encode(role); err != nil {
			require.NoError(t, err)
		}

		for _, user := range users {
			if err := enc.Encode(user); err != nil {
				require.NoError(t, err)
			}
		}
	}
}

func TestGetUsersOfGroup(t *testing.T) { //nolint:dupl
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kcc, err := client.NewKeycloakClient(ctx, keycloakClientConfig())
	require.NoError(t, err)

	groups, err := kcc.ListGroups(ctx)
	require.NoError(t, err)
	require.NotNil(t, groups)

	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	for _, group := range groups {
		users, err := kcc.GetUsersOfGroup(ctx, group.ID)
		require.NoError(t, err)
		require.NotNil(t, users)

		if err := enc.Encode(group); err != nil {
			require.NoError(t, err)
		}

		for _, user := range users {
			if err := enc.Encode(user); err != nil {
				require.NoError(t, err)
			}
		}
	}
}

func TestGetGroupsOfUser(t *testing.T) { //nolint:dupl
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kcc, err := client.NewKeycloakClient(ctx, keycloakClientConfig())
	require.NoError(t, err)

	users, err := kcc.ListUsers(ctx)
	require.NoError(t, err)
	require.NotNil(t, users)

	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	for _, user := range users {
		groups, err := kcc.GetGroupsOfUser(ctx, user.ID)
		require.NoError(t, err)
		require.NotNil(t, users)

		if err := enc.Encode(user); err != nil {
			require.NoError(t, err)
		}

		for _, group := range groups {
			if err := enc.Encode(group); err != nil {
				require.NoError(t, err)
			}
		}
	}
}
