package main

import (
	"os"

	"github.com/aserto-dev/ds-load/plugins/openapi/pkg/app"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
	"github.com/aserto-dev/ds-load/sdk/common/kongyaml"

	"github.com/alecthomas/kong"
)

func main() {
	cli := app.CLI{}

	defaultConfigPath := "~/.config/ds-load/cfg/openapi.yaml"

	yamlLoader := kongyaml.NewYAMLResolver("openapi")
	options := []kong.Option{
		kong.Name(app.AppName),
		kong.Exit(func(exitCode int) {
			os.Exit(exitCode)
		}),
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
