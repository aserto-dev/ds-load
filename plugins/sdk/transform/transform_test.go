package transform_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/aserto-dev/ds-load/common/js"
	"github.com/aserto-dev/ds-load/common/msg"
	"github.com/aserto-dev/ds-load/plugins/sdk"
	"github.com/aserto-dev/ds-load/plugins/sdk/transform"
	"github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestTransform(t *testing.T) {
	content, err := sdk.Assets().ReadFile("assets/peoplefinder.json")
	assert.NoError(t, err)
	input := make(map[string]interface{})
	err = json.Unmarshal(content, &input)
	assert.NoError(t, err)
	template, err := sdk.Assets().ReadFile("assets/transform_template.tmpl")
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
	assert.Equal(t, len(directoryObject.Relations), 1)
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
	template, err := sdk.Assets().ReadFile("assets/transform_template.tmpl")
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
	assert.Equal(t, len(directoryObject.Relations), 1)

	t.Log("Chunking")
	chunks := transformer.PrepareChunks(&directoryObject)
	t.Log("Object chunks", len(chunks))
	assert.NotNil(t, chunks)
	// assert.NotNil(t, relationChunks)
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
	assert.Equal(t, len(chunks), 3)
	// assert.Equal(t, len(relationChunks), 3)
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
	assert.Equal(t, len(chunks), 6)
	// assert.Equal(t, len(relationChunks), 6)
	var output bytes.Buffer
	jsonWriter := js.NewJSONArrayWriter(&output)
	err := trans.WriteChunks(jsonWriter, chunks)
	assert.NoError(t, err)
	t.Log(output.String())
}
