package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/hubspot/pkg/app"
	"github.com/aserto-dev/ds-load/sdk/common/kongyaml"
)

func main() {
	cli := app.CLI{}

	defaultConfigPath := "~/.config/ds-load/cfg/hubspot.yaml"

	yamlLoader := kongyaml.NewYAMLResolver("hubspot")
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

	ctx := kong.Parse(&cli, options...)
	err := ctx.Run()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
