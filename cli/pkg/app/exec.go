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
	"github.com/aserto-dev/ds-load/common/js"
	"github.com/aserto-dev/ds-load/common/msg"
	dsi "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
)

type ExecCmd struct {
	clients.Config
	CommandArgs  []string `name:"command" passthrough:"" arg:"" help:"available commands are: ${plugins}"`
	Print        bool     `name:"print" short:"p" help:"print output to stdout"`
	PluginFolder string   `hidden:""`

	execPlugin *plugin.Plugin     `kong:"-"`
	pluginArgs []string           `kong:"-"`
	dirClient  dsi.ImporterClient `kong:"-"`
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
		err = e.handleMessages(c, pStdout)
	}
	if err != nil {
		wg.Wait()
		return err
	}

	wg.Wait()
	return pluginCmd.Wait()
}

func (e *ExecCmd) handleMessages(c *cc.CommonCtx, stdout io.ReadCloser) error {
	reader, err := js.NewJSONArrayReader(stdout)
	if err != nil {
		return err
	}

	for {
		var message msg.Transform
		err := reader.ReadProtoMessage(&message)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = e.importToDirectory(c, &message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *ExecCmd) importToDirectory(c *cc.CommonCtx, message *msg.Transform) error {
	var sErr error
	errGroup, iCtx := errgroup.WithContext(c.Context)
	stream, err := e.dirClient.Import(iCtx)
	if err != nil {
		return err
	}
	errGroup.Go(receiver(stream))

	// import objects
	for _, object := range message.Objects {
		c.UI.Note().Msgf("object: [%s] type [%s]", object.Key, object.Type)
		sErr = stream.Send(&dsi.ImportRequest{
			Msg: &dsi.ImportRequest_Object{
				Object: object,
			},
		})
	}

	for _, relation := range message.Relations {
		c.UI.Note().Msgf("relation: [%s] obj: [%s] subj [%s]", relation.Relation, *relation.Object.Key, *relation.Subject.Key)
		sErr = stream.Send(&dsi.ImportRequest{
			Msg: &dsi.ImportRequest_Relation{
				Relation: relation,
			},
		})
	}

	err = stream.CloseSend()
	if err != nil {
		return err
	}

	err = errGroup.Wait()
	if err != nil {
		return err
	}

	// TODO handle stream errors
	if sErr != nil {
		c.Log.Err(sErr)
	}

	return nil
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

func receiver(stream dsi.Importer_ImportClient) func() error {
	return func() error {
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}
		}
	}
}
