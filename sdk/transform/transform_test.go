package transform_test

import (
	"bufio"
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/aserto-dev/ds-load/sdk"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/common/msg"
	"github.com/aserto-dev/ds-load/sdk/transform"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransform(t *testing.T) {
	// Arrange
	content, err := sdk.Assets().ReadFile("assets/peoplefinder.json")
	require.NoError(t, err)

	contentReader := strings.NewReader(string(content))

	template, err := sdk.Assets().ReadFile("assets/test_template.tmpl")
	require.NoError(t, err)

	var transformBuffer bytes.Buffer
	writer := bufio.NewWriter(&transformBuffer)

	transformer := transform.NewGoTemplateTransform(template)
	ctx := context.Background()

	// Act
	err = transformer.Transform(ctx, contentReader, writer, nil)
	require.NoError(t, err)
	err = writer.Flush()
	require.NoError(t, err)

	// Assert
	bufLen := transformBuffer.Len()
	transformOutput := make([]byte, bufLen)

	reader := bufio.NewReader(&transformBuffer)
	_, err = reader.Read(transformOutput)
	require.NoError(t, err)

	arrayReader, err := js.NewJSONArrayReader(bytes.NewReader(transformOutput))
	require.NoError(t, err)

	var directoryObj msg.Transform
	err = arrayReader.ReadProtoMessage(&directoryObj)
	require.NoError(t, err)

	objectCount := len(directoryObj.GetObjects())
	assert.Equal(t, 5, objectCount)

	relationCount := len(directoryObj.GetRelations())
	assert.Equal(t, 2, relationCount)
}

func TestTransformEscapedChars(t *testing.T) {
	// Arrange
	//nolint:lll
	const auth0user string = `
	[
	  {
		"created_at": "2023-06-19T10:18:13.702Z",
		"email": "oana+test666@aserto.com",
		"email_verified": true,
		"identities": [
		  {
			"connection": "Username-Password-Authentication",
			"provider": "auth0",
			"user_id": "64902b655c2e91cb3dee85a4",
			"isSocial": false
		  }
		],
		"name": "oana+test666@aserto.com",
		"nickname": "oana+test666",
		"picture": "https://s.gravatar.com/avatar/de191b7ce00efcc0cd07690f793c5186?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Foa.png", 
		"updated_at": "2023-06-30T12:47:52.762Z",
		"user_id": "auth0|64902b655c2e91cb3dee85a4",
		"user_metadata": {
		  "aserto-allow-tenant-creation": 5
		},
		"username": "oanatest1231",
		"last_password_reset": "2023-06-19T10:19:13.475Z",
		"last_ip": "109.99.219.89",
		"last_login": "2023-06-30T12:47:52.762Z",
		"logins_count": 9,
		"blocked_for": [],
		"guardian_authenticators": []
	  }
	]`

	content := []byte(auth0user)
	transformTemplate, err := sdk.Assets().ReadFile("assets/test_template.tmpl")
	require.NoError(t, err)

	contentReader := strings.NewReader(string(content))

	var transformBuffer bytes.Buffer

	writer := bufio.NewWriter(&transformBuffer)

	transformer := transform.NewGoTemplateTransform(transformTemplate)
	ctx := context.Background()

	// Act
	err = transformer.Transform(ctx, contentReader, writer, nil)
	require.NoError(t, err)
	err = writer.Flush()
	require.NoError(t, err)

	// Assert
	bufLen := transformBuffer.Len()
	transformOutput := make([]byte, bufLen)
	reader := bufio.NewReader(&transformBuffer)
	_, err = reader.Read(transformOutput)
	require.NoError(t, err)

	t.Log(transformOutput)

	arrayReader, err := js.NewJSONArrayReader(bytes.NewReader(transformOutput))
	require.NoError(t, err)

	var directoryObject msg.Transform
	err = arrayReader.ReadProtoMessage(&directoryObject)
	require.NoError(t, err)

	objectCount := len(directoryObject.GetObjects())
	assert.Equal(t, 2, objectCount)

	relationCount := len(directoryObject.GetRelations())
	assert.Equal(t, 2, relationCount)

	userObject := directoryObject.GetObjects()[0]
	assert.Equal(t, "user", userObject.GetType())
	assert.Equal(t, "oana+test666", userObject.GetDisplayName())

	userEmail, ok := userObject.GetProperties().GetFields()["email"]
	assert.True(t, ok)
	assert.Equal(t, "oana+test666@aserto.com", userEmail.GetStringValue())
}
