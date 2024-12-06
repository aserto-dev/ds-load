package app

import (
	"bufio"
	"os"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/cli/pkg/clients"
	"github.com/aserto-dev/ds-load/cli/pkg/publish"
	"github.com/aserto-dev/ds-load/sdk/plugin"
	"github.com/pkg/errors"
)

type PublishCmd struct {
	clients.Config
}

func (l *PublishCmd) Run(commonCtx *cc.CommonCtx) error {
	var publisher publish.Publisher

	dirClient, err := clients.NewDirectoryV3ImportClient(commonCtx.Context, &l.Config)
	if err != nil {
		return errors.Wrap(err, "Could not connect to the directory")
	}
	publisher = publish.NewDirectoryPublisher(commonCtx, dirClient)

	return l.processMessagesFromStdIn(commonCtx, publisher)
}

func (l *PublishCmd) processMessagesFromStdIn(commonCtx *cc.CommonCtx, publisher plugin.Publisher) error {
	reader := bufio.NewReader(os.Stdin)
	return publisher.Publish(commonCtx.Context, reader)
}
