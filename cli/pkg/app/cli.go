package app

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/common/version"
)

type Context struct {
	Config    string
	Verbosity int
	Insecure  bool
}

type CLI struct {
	Exec             ExecCmd             `cmd:"" help:"import data in directory"`
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
	fmt.Println("not implemented")
	return nil
}

type VersionCmd struct{}

func (cmd *VersionCmd) Run(c *cc.CommonCtx) error {
	fmt.Printf("%s - %s\n",
		AppName,
		version.GetInfo().String(),
	)
	return nil
}
