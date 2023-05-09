package app

import (
	"bufio"
	"fmt"
	"github.com/aserto-dev/ds-load/common/msg"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"log"
	"os/exec"
)

func LaunchCommands() {
	cmd := exec.Command("sleep", "60")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}
	go listenOnStdout(stdout)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
	}
	go listenOnStderr(stderr)

	err = cmd.Run()

	if err != nil {
		log.Println(err)
	}
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

}
