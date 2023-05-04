package app

import (
	"fmt"
	"github.com/aserto-dev/ds-load/common/version"
)

type Context struct {
	Config   string
	LogLevel string
	Insecure bool
}

type CLI struct {
	Exec             ExecCmd             `cmd:"" help:"import data in directory"`
	GetPlugin        GetPluginCmd        `cmd:"" help:"download plugin"`
	SetDefaultPlugin SetDefaultPluginCmd `cmd:"" help:"sets a plugin as default"`
	ListPlugins      ListPluginsCmd      `cmd:"" help:"list available plugins"`
	Version          VersionCmd          `cmd:"" help:"version information"`

	Config   string `short:"c" help:"Path to the config file. Any argument provided to the CLI will take precedence."`
	LogLevel string `short:"l" help:"Specify log level."`
	Insecure bool   `short:"i" help:"Disable TLS verification"`
}

type ExecCmd struct{}

func (exec *ExecCmd) Run(ctx *Context) error {
	fmt.Println("not implemented")
	return nil
}

type GetPluginCmd struct{}

func (getPlugin *GetPluginCmd) Run() error {
	fmt.Println("not implemented")
	return nil
}

type SetDefaultPluginCmd struct{}

func (defaultPlugin *SetDefaultPluginCmd) Run() error {
	fmt.Println("not implemented")
	return nil
}

type ListPluginsCmd struct{}

func (listPlugins *ListPluginsCmd) Run() error {
	fmt.Println("not implemented")
	return nil
}

type VersionCmd struct{}

func (cmd *VersionCmd) Run() error {
	fmt.Printf("%s - %s\n",
		AppName,
		version.GetInfo().String(),
	)
	return nil
}
