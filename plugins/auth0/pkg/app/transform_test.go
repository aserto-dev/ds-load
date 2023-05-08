package app

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/alecthomas/assert"
)

func TestTransform(t *testing.T) {
	content, err := os.ReadFile(AssetDefaultPeoplefinderExample())
	assert.NoError(t, err)
	input := make(map[string]interface{})
	err = json.Unmarshal(content, &input)
	assert.NoError(t, err)
	template, err := os.ReadFile(AssetDefaultTemplate())
	assert.NoError(t, err)
	output, err := Transform(input, string(template))
	assert.NoError(t, err)
	t.Log(output)
}
