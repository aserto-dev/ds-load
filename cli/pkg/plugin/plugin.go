package plugin

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aserto-dev/ds-load/cli/pkg/constants"
)

type Plugin struct {
	Name string
	Path string
}

func NewPlugin(path string) *Plugin {
	return &Plugin{
		Name: pluginName(path),
		Path: path,
	}
}

func pluginName(path string) string {
	file := filepath.Base(path)
	name := strings.TrimPrefix(file, constants.PluginPrefix)
	if runtime.GOOS == "windows" {
		name = strings.TrimSuffix(name, ".exe")
	}
	return name
}
