package transform_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/aserto-dev/ds-load/sdk"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/common/msg"
	"github.com/aserto-dev/ds-load/sdk/transform"
	"github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestTransform(t *testing.T) {
	content, err := sdk.Assets().ReadFile("assets/peoplefinder.json")
	assert.NoError(t, err)
	input := make(map[string]interface{})
	err = json.Unmarshal(content, &input)
	assert.NoError(t, err)
	template, err := sdk.Assets().ReadFile("assets/test_template.tmpl")
	assert.NoError(t, err)

	transformer := transform.NewTransformer(1)

	output, err := transformer.TransformToTemplate(input, string(template))
	assert.NoError(t, err)
	t.Log(output)
	var out bytes.Buffer
	err = json.Indent(&out, []byte(output), "", "\t")
	assert.NoError(t, err)
	var directoryObject msg.Transform
	err = protojson.Unmarshal(out.Bytes(), &directoryObject)
	assert.NoError(t, err)
	objectCount := len(directoryObject.Objects)
	assert.Equal(t, objectCount, 5)
	relationCount := len(directoryObject.Relations)
	assert.Equal(t, relationCount, 2)
}

func TestTransformWithManyObjects(t *testing.T) {
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	numberOfObject := 100000

	t.Logf("before: %v Kb", bToKb(m1.TotalAlloc))
	content, err := sdk.Assets().ReadFile("assets/peoplefinder.json")
	assert.NoError(t, err)
	input := make(map[string]interface{})
	err = json.Unmarshal(content, &input)
	assert.NoError(t, err)
	template, err := sdk.Assets().ReadFile("assets/test_template.tmpl")
	assert.NoError(t, err)

	t.Logf("Attaching %d roles to peoplefinder", numberOfObject)
	// Increase number of objects
	for i := 0; i < numberOfObject; i++ {
		input["roles"] = append(input["roles"].([]interface{}), map[string]interface{}{
			"id":   fmt.Sprintf("%d", i),
			"name": fmt.Sprintf("role%d", i),
		})
	}

	transformer := transform.NewTransformer(1)

	output, err := transformer.TransformToTemplate(input, string(template))
	assert.NoError(t, err)
	var out bytes.Buffer
	err = json.Indent(&out, []byte(output), "", "\t")
	assert.NoError(t, err)
	var directoryObject msg.Transform
	err = protojson.Unmarshal(out.Bytes(), &directoryObject)
	assert.NoError(t, err)
	assert.Equal(t, len(directoryObject.Relations), 2)

	t.Log("Chunking")
	chunks := transformer.PrepareChunks(&directoryObject)
	t.Log("Object chunks", len(chunks))
	assert.NotNil(t, chunks)
	runtime.ReadMemStats(&m2)
	t.Logf("after: %v Kb", bToKb(m2.TotalAlloc))
	t.Logf("total diff: %v Kb", bToKb(m2.TotalAlloc-m1.TotalAlloc))
}

func bToKb(b uint64) uint64 {
	return b / 1024
}

func TestTransformChunking(t *testing.T) {
	var directoryObjects msg.Transform
	for i := 0; i < 30; i++ {
		key := fmt.Sprintf("%d", i)
		testType := "test_object" //nolint:goconst

		directoryObjects.Objects = append(directoryObjects.Objects, &common.Object{
			Key:  key,
			Type: testType,
		})

		directoryObjects.Relations = append(directoryObjects.Relations, &common.Relation{
			Subject: &common.ObjectIdentifier{
				Key:  &key,
				Type: &testType,
			},
			Object: &common.ObjectIdentifier{
				Key:  &key,
				Type: &testType,
			},
			Relation: "test",
		})
	}

	trans := transform.NewTransformer(10)

	chunks := trans.PrepareChunks(&directoryObjects)
	assert.Equal(t, len(chunks), 6)
}

func TestTransformWriter(t *testing.T) {
	var directoryObjects msg.Transform
	for i := 0; i < 30; i++ {
		key := fmt.Sprintf("%d", i)
		testType := "test_object"

		directoryObjects.Objects = append(directoryObjects.Objects, &common.Object{
			Key:  key,
			Type: testType,
		})

		directoryObjects.Relations = append(directoryObjects.Relations, &common.Relation{
			Subject: &common.ObjectIdentifier{
				Key:  &key,
				Type: &testType,
			},
			Object: &common.ObjectIdentifier{
				Key:  &key,
				Type: &testType,
			},
			Relation: "test",
		})
	}

	trans := transform.NewTransformer(5)

	chunks := trans.PrepareChunks(&directoryObjects)
	assert.Equal(t, len(chunks), 12)
	var output bytes.Buffer
	jsonWriter, err := js.NewJSONArrayWriter(&output)
	assert.NoError(t, err)
	err = trans.WriteChunks(jsonWriter, chunks)
	assert.NoError(t, err)
	t.Log(output.String())
}

func TestTransformEscapedChars(t *testing.T) {
	const auth0user string = `
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
	}`

	content := []byte(auth0user)
	input := make(map[string]interface{})
	err := json.Unmarshal(content, &input)
	assert.NoError(t, err)
	template, err := sdk.Assets().ReadFile("assets/test_template.tmpl")
	assert.NoError(t, err)

	transformer := transform.NewTransformer(1)

	output, err := transformer.TransformToTemplate(input, string(template))
	assert.NoError(t, err)
	t.Log(output)
	var out bytes.Buffer
	err = json.Indent(&out, []byte(output), "", "\t")
	assert.NoError(t, err)
	var directoryObject msg.Transform
	err = protojson.Unmarshal(out.Bytes(), &directoryObject)
	assert.NoError(t, err)
	objectCount := len(directoryObject.Objects)
	assert.Equal(t, objectCount, 2)
	relationCount := len(directoryObject.Relations)
	assert.Equal(t, relationCount, 2)

	userObject := directoryObject.Objects[0]
	assert.Equal(t, userObject.Type, "user")
	assert.Equal(t, userObject.DisplayName, "oana+test666")

	userEmail, ok := userObject.Properties.Fields["email"]
	assert.True(t, ok)
	assert.Equal(t, userEmail.GetStringValue(), "oana+test666@aserto.com")
}
