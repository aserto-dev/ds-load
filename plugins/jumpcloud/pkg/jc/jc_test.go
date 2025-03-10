package jc_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/aserto-dev/ds-load/plugins/jumpcloud/pkg/jc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Setenv("JC_API_KEY", "jca_7F8FDtW94HCXnVA97WkrdXDd7v7P47XkALqq")
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

func TestGetUserByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jcc, err := jc.NewJumpCloudClient(ctx, os.Getenv("JC_API_KEY"))
	require.NoError(t, err)
	assert.NoError(t, err)
	_ = jcc

	// req, err := jcc.NewRequest(http.MethodGet, "/api/Systemusers/67cb39b4c1189814618d9554", nil)
	// require.NoError(t, err)

	// var body any
	// resp, err := jcc.Do(req, &body)
	// require.NoError(t, err)

	// enc := json.NewEncoder(os.Stderr)
	// enc.SetEscapeHTML(false)
	// enc.SetIndent("", "  ")
	// enc.Encode(body)

	// fmt.Fprintf(os.Stderr, "status code: %d (%s)\n", resp.StatusCode, resp.Status)
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
