package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/app"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
	"github.com/aserto-dev/ds-load/sdk/common/kongyaml"
)

func main() {
	cli := app.CLI{}

	defaultConfigPath := "~/.config/ds-load/cfg/ldap.yaml"

	yamlLoader := kongyaml.NewYAMLResolver("ldap")
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
