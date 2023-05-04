package app

import (
	"fmt"

	"github.com/aserto-dev/ds-load/common/version"
)

type CLI struct {
	Info    InfoCmd    `cmd:"" help:""`
	Version VersionCmd `cmd:"" help:"version information"`
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
