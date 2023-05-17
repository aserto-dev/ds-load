package app

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/cli/pkg/clients"
	"github.com/aserto-dev/ds-load/common/msg"
	dsi "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
)

type ExecCmd struct {
	Plugin       string `cmd:""`
	PluginConfig string `cmd:""`
	MaxChunkSize int    `cmd:""`
	clients.Config
	dirClient dsi.ImporterClient
}

func (e *ExecCmd) Run(c *cc.CommonCtx) error {
	// TODO improve plugin support check
	if e.Plugin != "auth0" {
		return errors.New("plugin not supported")
	}
	cli, err := clients.NewDirectoryImportClient(c, &e.Config)
	if err != nil {
		return err
	}
	e.dirClient = cli
	return e.LaunchCommands(c)
}

func (e *ExecCmd) LaunchCommands(c *cc.CommonCtx) error {
	fetchCmd := exec.Command(getPluginExecPath(e.Plugin), "fetch", "-c", e.PluginConfig)         //nolint:gosec
	transformCmd := exec.Command(getPluginExecPath(e.Plugin), "transform", "-c", e.PluginConfig) //nolint:gosec

	var wg sync.WaitGroup

	fStdout, err := fetchCmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer fStdout.Close()

	fStderr, err := fetchCmd.StderrPipe()
	if err != nil {
		return err
	}
	defer fStderr.Close()

	tStdin, err := transformCmd.StdinPipe()
	if err != nil {
		return err
	}
	defer tStdin.Close()

	tStdout, err := transformCmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := transformCmd.StderrPipe()
	if err != nil {
		return err
	}

	go listenOnStderr(stderr)
	go listenOnStderr(fStderr)

	err = fetchCmd.Start()
	if err != nil {
		return err
	}

	err = transformCmd.Start()
	if err != nil {
		return err
	}
	wg.Add(1)

	go func() {
		scanner := bufio.NewScanner(fStdout)
		for scanner.Scan() {
			line := scanner.Bytes()
			_, err = tStdin.Write(line)
			if err != nil {
				c.UI.Problem().Msg(err.Error())
			}
			_, err = tStdin.Write([]byte("\n"))
			if err != nil {
				c.UI.Problem().Msg(err.Error())
			}
		}

		wg.Done()
		err = tStdin.Close()
		if err != nil {
			c.UI.Problem().Msg(err.Error())
		}
	}()

	err = e.handleMessages(c, tStdout)
	if err != nil {
		return err
	}

	wg.Wait()
	err = fetchCmd.Wait()
	if err != nil {
		c.UI.Problem().Msg(err.Error())
	}
	return transformCmd.Wait()
}

func (e *ExecCmd) handleMessages(c *cc.CommonCtx, stdout io.ReadCloser) error {
	scanner := bufio.NewReader(stdout)

	var stream dsi.Importer_ImportClient
	var errGroup *errgroup.Group
	var iCtx context.Context
	streamOpen := false

	for {
		line, err := scanner.ReadBytes('\n')
		if err == io.EOF {
			// we have reached the end of the stream
			break
		} else if err != nil {
			return err
		}

		var sErr error
		protoMsg := convertToProto(line)
		switch protoMsg.Data.(type) {
		case *msg.PluginMessage_Batch:
			batch := protoMsg.GetBatch()
			switch batch.Type {
			case msg.BatchType_BEGIN:
				if streamOpen {
					return errors.New("received batch begin on already open stream")
				}
				errGroup, iCtx = errgroup.WithContext(c.Context)
				stream, err = e.dirClient.Import(iCtx)
				if err != nil {
					return err
				}
				errGroup.Go(receiver(stream))
				streamOpen = true
			case msg.BatchType_END:
				if !streamOpen {
					return errors.New("received batch end on already closed stream")
				}
				err = stream.CloseSend()
				if err != nil {
					return err
				}
				streamOpen = false
				err = errGroup.Wait()
				if err != nil {
					return err
				}
			case msg.BatchType_NONE:
				return errors.New("received unexpected batch type none")
			}
		case *msg.PluginMessage_Object:
			object := protoMsg.GetObject()
			c.UI.Note().Msgf("object: [%s] type [%s]", object.Key, object.Type)
			sErr = stream.Send(&dsi.ImportRequest{
				Msg: &dsi.ImportRequest_Object{
					Object: object,
				},
			})
		case *msg.PluginMessage_Relation:
			relation := protoMsg.GetRelation()
			c.UI.Note().Msgf("relation: [%s] obj: [%s] subj [%s]", relation.Relation, *relation.Object.Key, *relation.Subject.Key)
			sErr = stream.Send(&dsi.ImportRequest{
				Msg: &dsi.ImportRequest_Relation{
					Relation: relation,
				},
			})
		}
		// TODO handle stream errors
		if sErr != nil {
			return sErr
		}
	}

	return nil
}

func convertToProto(line []byte) *msg.PluginMessage {
	pluginMsg := &msg.PluginMessage{}
	err := protojson.Unmarshal(line, pluginMsg)
	if err != nil {
		log.Println(err)
	}

	return pluginMsg
}

func listenOnStderr(stderr io.ReadCloser) {
	scanner := bufio.NewReader(stderr)

	for {
		line, err := scanner.ReadBytes('\n')
		if err == io.EOF {
			// we have reached the end of the stream
			break
		} else if err != nil {
			log.Fatal(err)
		}

		os.Stderr.WriteString(string(line))
		os.Stderr.WriteString("\n")
	}
}

func getPluginExecPath(pluginName string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	return filepath.Join(homeDir, ".ds-load", "plugins", "ds-load-"+pluginName+ext)
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
