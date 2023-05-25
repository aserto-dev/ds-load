package finder

import (
	"os"
	"path/filepath"

	"github.com/aserto-dev/ds-load/cli/pkg/constants"
	"github.com/aserto-dev/ds-load/cli/pkg/plugin"
)

type PluginDir struct {
	dir string
}

func NewCustomDir(dir string) Finder {
	return &PluginDir{
		dir: dir,
	}
}

func NewHomeDir() (Finder, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &PluginDir{
		dir: homeDir,
	}, nil
}

func (path PluginDir) Find() ([]*plugin.Plugin, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	foundPlugins := []*plugin.Plugin{}

	files, err := filepath.Glob(filepath.Join(homeDir, ".ds-load", "plugins", constants.PluginPrefix+"*"))
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		for _, f := range files {
			foundPlugins = append(foundPlugins, plugin.New(f))
		}
	}

	return foundPlugins, nil
}
