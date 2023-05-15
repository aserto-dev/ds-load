package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/aserto-dev/ds-load/common/msg"
	"google.golang.org/protobuf/encoding/protojson"
)

type ExecCmd struct {
	Plugin       string `cmd:""`
	Directory    string `cmd:""`
	APIKey       string `cmd:""`
	PluginConfig string `cmd:""`
	Insecure     bool   `short:"i" help:"Disable TLS verification"`
}

func (e *ExecCmd) Run(ctx *Context) error {
	// TODO improve plugin support check
	if e.Plugin != "auth0" {
		return errors.New("plugin not supported")
	}
	return e.LaunchCommands()
}

func (e *ExecCmd) LaunchCommands() error {
	fetchCmd := exec.Command(getPluginExecPath(e.Plugin), "fetch", "-c", e.PluginConfig)
	transformCmd := exec.Command(getPluginExecPath(e.Plugin), "transform", "-c", e.PluginConfig)

	var wg sync.WaitGroup

	fStdout, err := fetchCmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer fStdout.Close()

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

	err = fetchCmd.Start()
	if err != nil {
		return err
	}

	err = transformCmd.Start()
	if err != nil {
		return err
	}
	wg.Add(2)

	go func() {
		scanner := bufio.NewScanner(fStdout)
		for scanner.Scan() {
			line := scanner.Bytes()
			tStdin.Write(line)
			tStdin.Write([]byte("\n"))
		}

		wg.Done()
		tStdin.Close()
	}()

	go func() {
		listenOnStdout(tStdout)
		wg.Done()
	}()

	wg.Wait()
	fetchCmd.Wait()
	transformCmd.Wait()

	return nil
}

func listenOnStdout(stdout io.ReadCloser) {
	scanner := bufio.NewReader(stdout)

	for {
		line, err := scanner.ReadBytes('\n')
		if err == io.EOF {
			// we have reached the end of the stream
			break
		} else if err != nil {
			log.Fatal(err)
		}

		protoMsg := convertToProto(line)
		writeToDirectory(protoMsg)
	}
}

func convertToProto(line []byte) *msg.PluginMessage {
	pluginMsg := &msg.PluginMessage{}
	err := protojson.Unmarshal(line, pluginMsg)
	if err != nil {
		log.Println(err)
	}

	return pluginMsg
}

func writeToDirectory(protoMsg *msg.PluginMessage) {
	fmt.Println(protoMsg)
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
