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

const JC_API_KEY string = "JC_API_KEY"

func TestMain(m *testing.M) {
	if os.Getenv(JC_API_KEY) == "" {
		fmt.Fprintf(os.Stderr, "env %q not set, tests skipped", JC_API_KEY)
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

	directories, err := jcc.ListDirectories()
	require.NoError(t, err)

	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(directories)
}

func TestListUsers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jcc, err := jc.NewJumpCloudClient(ctx, os.Getenv("JC_API_KEY"))
	require.NoError(t, err)
	assert.NoError(t, err)

	users, err := jcc.ListUsers()
	require.NoError(t, err)
	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(users)
}

func TestListGroups(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jcc, err := jc.NewJumpCloudClient(ctx, os.Getenv("JC_API_KEY"))
	require.NoError(t, err)
	assert.NoError(t, err)

	groups, err := jcc.ListGroups(jc.AllGroups)
	require.NoError(t, err)
	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(groups)
}

func TestGetSystemGroups(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jcc, err := jc.NewJumpCloudClient(ctx, os.Getenv("JC_API_KEY"))
	require.NoError(t, err)
	assert.NoError(t, err)

	groups, err := jcc.ListGroups(jc.SystemGroups)
	require.NoError(t, err)
	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(groups)
}

func TestGetUserGroups(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jcc, err := jc.NewJumpCloudClient(ctx, os.Getenv("JC_API_KEY"))
	require.NoError(t, err)
	assert.NoError(t, err)

	groups, err := jcc.ListGroups(jc.UserGroups)
	require.NoError(t, err)
	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(groups)
}

func TestGetMembersOfGroup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jcc, err := jc.NewJumpCloudClient(ctx, os.Getenv("JC_API_KEY"))
	require.NoError(t, err)
	assert.NoError(t, err)

	groups, err := jcc.ListGroups(jc.UserGroups)
	require.NoError(t, err)

	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	for _, group := range groups {
		users, err := jcc.GetUsersInGroup(group.ID)
		require.NoError(t, err)
		enc.Encode(users)
	}
}
