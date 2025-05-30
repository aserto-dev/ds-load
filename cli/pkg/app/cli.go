package app

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/cli/pkg/constants"
	"github.com/aserto-dev/ds-load/cli/pkg/plugin"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
	"github.com/aserto-dev/ds-load/sdk/common/version"
	"github.com/pkg/errors"
)

type CLI struct {
	Exec             ExecCmd             `cmd:"" help:"import data in directory by running fetch, transform and publish" default:"withargs"`
	Publish          PublishCmd          `cmd:"" help:"load data from stdin into directory"`
	GetPlugin        GetPluginCmd        `cmd:"" help:"download plugin"`
	SetDefaultPlugin SetDefaultPluginCmd `cmd:"" help:"sets a plugin as default"`
	ListPlugins      ListPluginsCmd      `cmd:"" help:"list available plugins"`
	Version          VersionCmd          `cmd:"" help:"version information"`

	Config    kong.ConfigFlag `short:"c" help:"Path to the config file. Any argument provided to the CLI will take precedence."`
	Verbosity int             `short:"v" type:"counter" help:"Use to increase output verbosity."`
}

type GetPluginCmd struct{}

func (getPlugin *GetPluginCmd) Run(c *cc.CommonCtx) error {
	fmt.Println("not implemented")
	return nil
}

type SetDefaultPluginCmd struct{}

func (defaultPlugin *SetDefaultPluginCmd) Run(c *cc.CommonCtx) error {
	fmt.Println("not implemented")
	return nil
}

type ListPluginsCmd struct{}

func (listPlugins *ListPluginsCmd) Run(c *cc.CommonCtx) error {
	find, err := plugin.NewHomeDirFinder(true)
	if err != nil {
		return err
	}

	plugins, err := find.Find()
	if err != nil {
		return err
	}

	for _, p := range plugins {
		if _, err := fmt.Fprintf(os.Stdout, "%-15s (%s)\n", p.Name, p.Path); err != nil {
			return errors.Wrapf(err, "failed to write plugin info for %q", p.Name)
		}
	}

	return nil
}

type VersionCmd struct{}

func (cmd *VersionCmd) Run(c *cc.CommonCtx) error {
	fmt.Printf("%s - %s\n",
		constants.AppName,
		version.GetInfo().String(),
	)

	return nil
}
