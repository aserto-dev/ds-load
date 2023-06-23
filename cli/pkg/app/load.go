package app

import (
	"bufio"
	"os"

	"github.com/aserto-dev/ds/cli/pkg/cc"
	"github.com/aserto-dev/ds/cli/pkg/clients"
	"github.com/pkg/errors"
)

type LoadCmd struct {
	clients.Config
}

func (l *LoadCmd) Run(c *cc.CommonCtx) error {
	dirClient, err := clients.NewDirectoryImportClient(c, &l.Config)
	if err != nil {
		return errors.Wrap(err, "Could not connect to the directory")
	}

	return l.processMessagesFromStdIn(dirClient)
}

func (l *LoadCmd) processMessagesFromStdIn(dirClient clients.DirectoryClient) error {
	reader := bufio.NewReader(os.Stdin)
	return dirClient.HandleMessages(reader)
}
