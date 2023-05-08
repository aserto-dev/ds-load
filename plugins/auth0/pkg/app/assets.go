package app

import (
	"path/filepath"
	"runtime"
)

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
