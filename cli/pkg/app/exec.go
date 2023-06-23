package app

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/cli/pkg/clients"
	"github.com/aserto-dev/ds-load/cli/pkg/plugin"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type ExecCmd struct {
	clients.Config
	CommandArgs  []string `name:"command" passthrough:"" arg:"" help:"available commands are: ${plugins}"`
	Print        bool     `name:"print" short:"p" help:"print output to stdout"`
	PluginFolder string   `hidden:""`

	execPlugin *plugin.Plugin          `kong:"-"`
	pluginArgs []string                `kong:"-"`
	dirClient  clients.DirectoryClient `kong:"-"`
}

func (e *ExecCmd) Run(c *cc.CommonCtx) error {
	defaultPrintCmd := []string{"fetch", "version", "export-transform"}
	var err error
	var find *plugin.Finder
	if e.PluginFolder != "" {
		find = plugin.NewFinder(true, e.PluginFolder)
	} else {
		find, err = plugin.NewHomeDirFinder(true)
		if err != nil {
			return err
		}
	}
	pl := e.CommandArgs[0]

	plugins, err := find.Find()
	if err != nil {
		return err
	}
	for _, p := range plugins {
		if pl == p.Name {
			e.execPlugin = p
			break
		}
	}
	if e.execPlugin == nil {
		return errors.Errorf("plugin [%s] not found", pl)
	}

	e.pluginArgs = e.CommandArgs[1:]
	var pluginSubCommand string
	if len(e.CommandArgs) > 1 {
		pluginSubCommand = e.CommandArgs[1]
	}

	if slices.Contains(e.pluginArgs, "-h") || slices.Contains(e.pluginArgs, "--help") || slices.Contains(defaultPrintCmd, pluginSubCommand) {
		e.Print = true
	}

	if !e.Print {
		cli, err := clients.NewDirectoryImportClient(c, &e.Config)
		if err != nil {
			return errors.Wrap(err, "Could not connect to the directory")
		}
		e.dirClient = cli
	}
	return e.LaunchPlugin(c)
}

func (e *ExecCmd) LaunchPlugin(c *cc.CommonCtx) error {
	if (!slices.Contains(e.pluginArgs, "-c") || !slices.Contains(e.pluginArgs, "--config")) && c.ConfigPath != "" {
		e.pluginArgs = append(e.pluginArgs, "-c", c.ConfigPath)
	}
	pluginCmd := exec.Command(e.execPlugin.Path, e.pluginArgs...) //nolint:gosec
	var pStdout io.ReadCloser
	var wg sync.WaitGroup

	pStderr, err := pluginCmd.StderrPipe()
	if err != nil {
		return err
	}
	defer pStderr.Close()

	wg.Add(1)
	go listenOnStderr(c, &wg, pStderr)

	if e.Print {
		pluginCmd.Stdout = os.Stdout
	} else {
		pStdout, err = pluginCmd.StdoutPipe()
		if err != nil {
			return err
		}
		defer pStdout.Close()
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return err
	}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		pluginCmd.Stdin = os.Stdin
	}

	err = pluginCmd.Start()
	if err != nil {
		return err
	}

	if !e.Print {
		err = e.dirClient.HandleMessages(pStdout)
	}
	if err != nil {
		wg.Wait()
		return err
	}

	wg.Wait()
	return pluginCmd.Wait()
}

func listenOnStderr(c *cc.CommonCtx, wg *sync.WaitGroup, stderr io.ReadCloser) {
	scanner := bufio.NewReader(stderr)
	gotError := false

	for {
		line, err := scanner.ReadBytes('\n')
		os.Stderr.Write(line)

		if len(line) > 0 {
			gotError = true
		}

		// we have reached the end of the stream
		if err == io.EOF {
			if gotError {
				os.Exit(1)
			}

			break
		} else if err != nil {
			c.Log.Fatal().Err(err)
		}
	}
	wg.Done()
}
