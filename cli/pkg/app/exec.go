package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"

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
	e.LaunchCommands()
	return nil
}

func (e *ExecCmd) LaunchCommands() {
	fetchCmd := exec.Command(getPluginExecPath(e.Plugin), "fetch", "-c", e.PluginConfig)
	transformCmd := exec.Command(getPluginExecPath(e.Plugin), "transform", "-c", e.PluginConfig)

	r, w, _ := os.Pipe()
	transformCmd.Stdin = w
	fetchCmd.Stdout = r

	stdout, err := transformCmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}

	stderr, err := transformCmd.StderrPipe()
	if err != nil {
		log.Println(err)
	}
	go listenOnStderr(stderr)

	err = fetchCmd.Run()
	if err != nil {
		log.Println(err)
	}

	err = transformCmd.Run()
	if err != nil {
		log.Println(err)
	}

	listenOnStdout(stdout)
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
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	return path.Join(DefaultPluginLocation, "ds-load-"+pluginName+ext)
}
