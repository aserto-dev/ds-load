package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/cli/pkg/app"
	"github.com/aserto-dev/ds-load/cli/pkg/constants"
	"github.com/aserto-dev/ds-load/cli/pkg/plugin"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
	"github.com/aserto-dev/ds-load/sdk/common/kongyaml"
)

func main() {
	pluginEnum := ""

	pluginFinder, err := plugin.NewHomeDirFinder(true)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	plugins, err := pluginFinder.Find()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	for _, p := range plugins {
		pluginEnum += (p.Name + "|")
	}

	pluginEnum = strings.TrimSuffix(pluginEnum, "|")

	yamlLoader := kongyaml.NewYAMLResolver("")

	cli := app.CLI{}
	options := []kong.Option{
		kong.Name(constants.AppName),
		kong.Vars{"plugins": pluginEnum},
		kong.Description(constants.AppDescription),
		kong.UsageOnError(),
		kong.Configuration(yamlLoader.Loader, constants.DefaultConfigPath),
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

	ctx := cc.NewCommonContext(cli.Verbosity, string(cli.Config))

	if err := kongCtx.Run(ctx); err != nil {
		kongCtx.FatalIfErrorf(err)
	}

	os.Exit(common.GetExitCode())
}
