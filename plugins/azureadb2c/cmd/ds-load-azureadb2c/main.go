package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/azureadb2c/pkg/app"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/kongyaml"
)

func main() {
	cli := app.CLI{}

	defaultConfigPath := "~/.config/ds-load/cfg/azureadb2c.yaml"

	yamlLoader := kongyaml.NewYAMLResolver("azureadb2c")
	options := []kong.Option{
		kong.Name(app.AppName),
		kong.Description(app.AppDescription),
		kong.UsageOnError(),
		kong.Configuration(yamlLoader.Loader, defaultConfigPath),
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

	ctx := cc.NewCommonContext(cli.Verbosity, string(cli.Config))
	kongCtx := kong.Parse(&cli, options...)
	if err := kongCtx.Run(ctx); err != nil {
		kongCtx.FatalIfErrorf(err)
	}

	os.Exit(common.GetExitCode())
}