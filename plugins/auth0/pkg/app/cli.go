package app

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/common/version"
)

type CLI struct {
	Config  kong.ConfigFlag `help:"Configuration file path"`
	Info    InfoCmd         `cmd:"" help:""`
	Version VersionCmd      `cmd:"" help:"version information"`
	Fetch   FetchCmd        `cmd:"" help:"fetch auth0 data"`
}

type VersionCmd struct {
}

func (cmd *VersionCmd) Run() error {
	fmt.Printf("%s - %s\n",
		AppName,
		version.GetInfo().String(),
	)
	return nil
}
