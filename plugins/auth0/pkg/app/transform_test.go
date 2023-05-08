package app

import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/assert"
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
}
