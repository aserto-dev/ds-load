package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/cli/pkg/app"
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

		kong.ConfigureHelp(kong.HelpOptions{
			NoAppSummary:        false,
			Summary:             false,
			Compact:             true,
			Tree:                true,
			FlagsLast:           true,
			Indenter:            kong.SpaceIndenter,
			NoExpandSubcommands: false,
		}),
	}

	ctx := kong.Parse(&cli, options...)
	context := &app.Context{
		Config:    cli.Config,
		Insecure:  cli.Insecure,
		Verbosity: cli.Verbosity,
	}

	err := ctx.Run(context)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

}
