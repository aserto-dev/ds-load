package finder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aserto-dev/ds-load/cli/pkg/constants"
	"github.com/aserto-dev/ds-load/cli/pkg/plugin"
)

type Environment struct {
}

func NewEnvironment() Finder {
	return &Environment{}
}

func (path Environment) Find() ([]*plugin.Plugin, error) {
	pathEnv := os.Getenv("PATH")
	dirs := strings.Split(pathEnv, ":")
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	dirs = append(dirs, pwd)

	foundPlugins := []*plugin.Plugin{}
	for _, dir := range dirs {
		files, err := filepath.Glob(filepath.Join(dir, constants.PluginPrefix+"*"))
		if err != nil {
			return nil, err
		}
		if len(files) > 0 {
			for _, f := range files {
				foundPlugins = append(foundPlugins, plugin.New(f))
			}
		}
	}

	return foundPlugins, nil
}
