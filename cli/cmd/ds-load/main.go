package main

import (
	"os"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/aserto-dev/ds-load/cli/pkg/app"
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
)

func main() {

	cli := app.CLI{}
	options := []kong.Option{
		kong.Name(app.AppName),
		kong.Exit(func(exitCode int) {
			os.Exit(exitCode)
		}),
		kong.Description(app.AppDescription),
		kong.UsageOnError(),
		kong.Configuration(kongyaml.Loader, app.DefaultConfigPath),
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
	err := kongCtx.Run(ctx)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
