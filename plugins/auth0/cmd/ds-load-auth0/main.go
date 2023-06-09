package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/common/kongyaml"
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/app"
)

func main() {
	cli := app.CLI{}

	defaultConfigPath := "~/.config/ds-load/cfg/auth0.yaml"

	options := []kong.Option{
		kong.Name(app.AppName),
		kong.Exit(func(exitCode int) {
			os.Exit(exitCode)
		}),
		kong.Description(app.AppDescription),
		kong.UsageOnError(),
		kong.Configuration(kongyaml.Loader, defaultConfigPath),
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
