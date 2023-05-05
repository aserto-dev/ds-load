package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
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

// AssetsDir returns the directory containing test assets.
func AssetsDir() string {
	_, filename, _, _ := runtime.Caller(0) //nolint: dogsled

	return filepath.Join(filepath.Dir(filename), "assets")
}

// AssetDefaultConfigExample returns the path of the default yaml config file.
func AssetDefaultPeoplefinderExample() string {
	return filepath.Join(AssetsDir(), "peoplefinder.json")
}

// AssetDefaultConfigExample returns the path of the default yaml config file.
func AssetDefaultTemplate() string {
	return filepath.Join(AssetsDir(), "transform_template.tmpl")
}
