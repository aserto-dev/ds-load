package main

import (
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/cli/pkg/app"
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/cli/pkg/constants"
	"github.com/aserto-dev/ds-load/cli/pkg/plugin"
	"github.com/aserto-dev/ds-load/common/kongyaml"
)

func main() {
	pluginEnum := ""
	find, err := plugin.NewHomeDirFinder(true)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	plugins, err := find.Find()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	for _, p := range plugins {
		pluginEnum += (p.Name + "|")
	}
	pluginEnum = strings.TrimSuffix(pluginEnum, "|")

	cli := app.CLI{}
	options := []kong.Option{
		kong.Name(constants.AppName),
		kong.Exit(func(exitCode int) {
			os.Exit(exitCode)
		}),
		kong.Vars{"plugins": pluginEnum},
		kong.Description(constants.AppDescription),
		kong.UsageOnError(),
		kong.Configuration(kongyaml.Loader, constants.DefaultConfigPath),
		kong.ConfigureHelp(kong.HelpOptions{
			NoAppSummary:        false,
			Summary:             false,
			Compact:             true,
			Tree:                false,
			FlagsLast:           true,
			Indenter:            kong.SpaceIndenter,
			NoExpandSubcommands: false,
		}),
	}

	kongCtx := kong.Parse(&cli, options...)
	ctx := cc.NewCommonContext(cli.Verbosity)
	err = kongCtx.Run(ctx)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
