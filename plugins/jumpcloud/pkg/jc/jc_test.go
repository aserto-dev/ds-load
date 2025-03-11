package jc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aserto-dev/ds-load/plugins/jumpcloud/pkg/jc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const JcAPIKey string = "JC_API_KEY" // nolint: gosec // no hardcoded credentials.

func TestMain(m *testing.M) {
	os.Setenv(JcAPIKey, "jca_8R6RtyhQcFuUfVdwzjzouUfXaVKhCwtqS3av")

	if os.Getenv(JcAPIKey) == "" {
		fmt.Fprintf(os.Stderr, "env %q not set, tests skipped", JcAPIKey)
		return
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestListDirectories(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jcc, err := jc.NewJumpCloudClient(ctx, os.Getenv("JC_API_KEY"))
	require.NoError(t, err)
	assert.NoError(t, err)

	directories, err := jcc.ListDirectories(ctx)
	require.NoError(t, err)

	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(directories); err != nil {
		require.NoError(t, err)
	}
}

func TestListUsers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jcc, err := jc.NewJumpCloudClient(ctx, os.Getenv("JC_API_KEY"))
	require.NoError(t, err)
	assert.NoError(t, err)

	users, err := jcc.ListUsers(ctx)
	require.NoError(t, err)
	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(users); err != nil {
		require.NoError(t, err)
	}
}

func TestListGroups(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jcc, err := jc.NewJumpCloudClient(ctx, os.Getenv("JC_API_KEY"))
	require.NoError(t, err)
	assert.NoError(t, err)

	groups, err := jcc.ListGroups(ctx, jc.AllGroups)
	require.NoError(t, err)
	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(groups); err != nil {
		require.NoError(t, err)
	}
}

func TestGetSystemGroups(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jcc, err := jc.NewJumpCloudClient(ctx, os.Getenv("JC_API_KEY"))
	require.NoError(t, err)
	assert.NoError(t, err)

	groups, err := jcc.ListGroups(ctx, jc.SystemGroups)
	require.NoError(t, err)
	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(groups); err != nil {
		require.NoError(t, err)
	}
}

func TestGetUserGroups(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jcc, err := jc.NewJumpCloudClient(ctx, os.Getenv("JC_API_KEY"))
	require.NoError(t, err)
	assert.NoError(t, err)

	groups, err := jcc.ListGroups(ctx, jc.UserGroups)
	require.NoError(t, err)
	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(groups); err != nil {
		require.NoError(t, err)
	}
}

func TestGetMembersOfGroup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	jcc, err := jc.NewJumpCloudClient(ctx, os.Getenv("JC_API_KEY"))
	require.NoError(t, err)
	assert.NoError(t, err)

	groups, err := jcc.ListGroups(ctx, jc.UserGroups)
	require.NoError(t, err)

	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	for _, group := range groups {
		users, err := jcc.GetUsersInGroup(ctx, group.ID)
		require.NoError(t, err)
		if err := enc.Encode(users); err != nil {
			require.NoError(t, err)
		}
	}
}
