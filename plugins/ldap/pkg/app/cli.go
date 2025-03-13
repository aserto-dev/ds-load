package app

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/sdk/common/version"
)

type CLI struct {
	Config          kong.ConfigFlag    `help:"Configuration file path" short:"c"`
	Version         VersionCmd         `cmd:"" help:"version information"`
	Fetch           FetchCmd           `cmd:"" help:"fetch ldap data"`
	Transform       TransformCmd       `cmd:"" help:"transform ldap data"`
	ExportTransform ExportTransformCmd `cmd:"" help:"export default transform template"`
	Exec            ExecCmd            `cmd:"" help:"fetch and transform ldap data" default:"withargs"`
	Verbosity       int                `short:"v" type:"counter" help:"Use to increase output verbosity."`
	Verify          VerifyCmd          `cmd:"verify" help:"verify fetcher configuration and credentials"`
}

type VersionCmd struct{}

func (cmd *VersionCmd) Run() error {
	fmt.Printf("%s - %s\n",
		AppName,
		version.GetInfo().String(),
	)

	return nil
}
