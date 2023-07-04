package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/azuread/pkg/app"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/kongyaml"
)

func main() {
	cli := app.CLI{}

	defaultConfigPath := "~/.config/ds-load/cfg/azuread.yaml"

	yamlLoader := kongyaml.NewYAMLResolver("azuread")
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

	ctx := kong.Parse(&cli, options...)
	if err := ctx.Run(); err != nil {
		ctx.FatalIfErrorf(err)
	}
	os.Exit(common.GetExitCode())
}
