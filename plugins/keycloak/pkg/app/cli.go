package app

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/sdk/common/version"
)

type CLI struct {
	Version         VersionCmd         `cmd:"" help:"version information"`
	Fetch           FetchCmd           `cmd:"" help:"fetch keycloak directory data"`
	Transform       TransformCmd       `cmd:"" help:"transform keycloak directory data"`
	ExportTransform ExportTransformCmd `cmd:"" help:"export default transform template"`
	Exec            ExecCmd            `cmd:"" help:"fetch and transform keycloak directory" default:"withargs"`
	Verify          VerifyCmd          `cmd:"" help:"verify fetcher configuration and credentials"`
	Config          kong.ConfigFlag    `flag:"" short:"c" help:"Configuration file path" `
	Verbosity       int                `flag:"" short:"v" type:"counter" help:"Use to increase output verbosity."`
}

type VersionCmd struct{}

func (cmd *VersionCmd) Run() error {
	fmt.Printf("%s - %s\n",
		AppName,
		version.GetInfo().String(),
	)

	return nil
}
