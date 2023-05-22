// get json from stdin and return json with v2 objects and v2 relations
package app

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/common/msg"
	"github.com/aserto-dev/ds-load/plugins/sdk/transform"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

type TransformCmd struct {
	TemplateFile string `cmd:""`
	MaxChunkSize int    `cmd:""`
}

func (t *TransformCmd) Run(context *kong.Context) error {
	var templateContent []byte
	var err error

	if t.TemplateFile == "" {
		templateContent, err = Assets().ReadFile("assets/transform_template.tmpl")
		if err != nil {
			return err
		}
	} else {
		templateContent, err = os.ReadFile(t.TemplateFile)
		if err != nil {
			return err
		}
	}

	if t.MaxChunkSize == 0 {
		t.MaxChunkSize = 1 // By default do not write begin and end batches
	}

	tranformer := transform.NewTransformer(t.MaxChunkSize)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		inputText := scanner.Bytes()

		input := make(map[string]interface{})

		err = json.Unmarshal(inputText, &input)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal input into map[string]interface{}")
		}
		output, err := tranformer.TransformToTemplate(input, string(templateContent))
		if err != nil {
			return errors.Wrap(err, "transform template execute failed")
		}
		var directoryObject msg.Transform
		if os.Getenv("DEBUG") != "" {
			os.Stdout.WriteString(output)
		}
		err = protojson.Unmarshal([]byte(output), &directoryObject)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal transformed data into directory objects and relations")
		}

		objectChunks, relationChunks := tranformer.PrepareChunks(&directoryObject)

		err = tranformer.WriteChunks(os.Stdout, objectChunks, relationChunks)
		if err != nil {
			return errors.Wrap(err, "failed to write chunks to output")
		}
	}

	return nil
}
