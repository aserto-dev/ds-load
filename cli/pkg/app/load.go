package app

import (
	"bufio"
	"os"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/cli/pkg/clients"
	"github.com/pkg/errors"
)

type LoadCmd struct {
	clients.Config
	dirClient clients.DirectoryClient
}

func (l *LoadCmd) Run(c *cc.CommonCtx) error {
	err := l.initializeDirectoryClient(c)
	if err != nil {
		return err
	}

	return l.readJSON()
}

func (l *LoadCmd) initializeDirectoryClient(c *cc.CommonCtx) error {
	cli, err := clients.NewDirectoryImportClient(c, &l.Config)
	if err != nil {
		return errors.Wrap(err, "Could not connect to the directory")
	}
	l.dirClient = cli

	return nil
}

func (l *LoadCmd) readJSON() error {
	reader := bufio.NewReader(os.Stdin)
	return l.dirClient.HandleMessages(reader)
}
