package app_test

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type testSetup struct {
	assert *require.Assertions
}

func setupTest(t *testing.T) *testSetup {
	assert := require.New(t)

	return &testSetup{
		assert: assert,
	}
}

func TestExec(t *testing.T) {

}
