package app

import (
	"fmt"
	"sync/atomic"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/sdk/common/version"
)

var (
	exitCode int32
)

func GetExitCode() int {
	return int(atomic.LoadInt32(&exitCode))
}

func SetExitCode(code int) {
	atomic.StoreInt32(&exitCode, int32(code))
}

type CLI struct {
	Config          kong.ConfigFlag    `help:"Configuration file path" short:"c"`
	Version         VersionCmd         `cmd:"" help:"version information"`
	Fetch           FetchCmd           `cmd:"" help:"fetch okta data"`
	Transform       TransformCmd       `cmd:"" help:"transform okta data"`
	ExportTransform ExportTransportCmd `cmd:"" help:"export default transform template"`
	Exec            ExecCmd            `cmd:"" help:"fetch and transform auth0 data" default:"withargs"`
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
