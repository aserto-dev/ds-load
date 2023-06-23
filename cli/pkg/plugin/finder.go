package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aserto-dev/ds/cli/pkg/constants"
	"golang.org/x/exp/slices"
)

type Finder struct {
	dirs []string
	env  bool
}

func NewFinder(env bool, dirs ...string) *Finder {
	return &Finder{
		dirs: dirs,
		env:  env,
	}
}

func NewHomeDirFinder(env bool) (*Finder, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &Finder{
		dirs: []string{filepath.Join(homeDir, ".ds", "plugins")},
		env:  env,
	}, nil
}

func (f Finder) Find() ([]*Plugin, error) {
	addedPlugins := []string{}
	dirs := f.dirs
	if f.env {
		pathEnv := os.Getenv("PATH")
		dirs = append(dirs, strings.Split(pathEnv, ":")...)
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}

		dirs = append(dirs, pwd)
	}

	foundPlugins := []*Plugin{}

	for _, dir := range dirs {
		files, err := filepath.Glob(filepath.Join(dir, constants.PluginPrefix+"*"))
		if err != nil {
			return nil, err
		}
		if len(files) > 0 {
			for _, f := range files {
				p := NewPlugin(f)
				if !slices.Contains(addedPlugins, p.Name) {
					foundPlugins = append(foundPlugins, p)
					addedPlugins = append(addedPlugins, p.Name)
				}
			}
		}
	}

	return foundPlugins, nil
}
