package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/aserto-dev/ds-load/common/msg"
	"github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestTransform(t *testing.T) {
	content, err := Assets().ReadFile("assets/peoplefinder.json")
	assert.NoError(t, err)
	input := make(map[string]interface{})
	err = json.Unmarshal(content, &input)
	assert.NoError(t, err)
	template, err := Assets().ReadFile("assets/transform_template.tmpl")
	assert.NoError(t, err)
	output, err := Transform(input, string(template))
	assert.NoError(t, err)
	t.Log(output)
	var out bytes.Buffer
	err = json.Indent(&out, []byte(output), "", "\t")
	assert.NoError(t, err)
	var directoryObject msg.Transform
	err = protojson.Unmarshal(out.Bytes(), &directoryObject)
	assert.NoError(t, err)
	assert.Equal(t, len(directoryObject.Objects), 5)
	assert.Equal(t, len(directoryObject.Relations), 2)
}

func TestTransformChunking(t *testing.T) {
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
	trans := TransformCmd{}
	trans.MaxChunkSize = 10
	objectChunks, relationChunks := trans.prepareChunks(&directoryObjects)
	assert.Equal(t, len(objectChunks), 3)
	assert.Equal(t, len(relationChunks), 3)
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
	trans := TransformCmd{}
	trans.MaxChunkSize = 5
	objectChunks, relationChunks := trans.prepareChunks(&directoryObjects)
	assert.Equal(t, len(objectChunks), 6)
	assert.Equal(t, len(relationChunks), 6)
	var output bytes.Buffer
	err := writeResponse(&output, objectChunks, relationChunks)
	assert.NoError(t, err)
	t.Log(output.String())
}
