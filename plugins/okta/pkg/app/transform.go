// get json from stdin and return json with v2 objects and v2 relations
package app

import (
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/common/js"
	"github.com/aserto-dev/ds-load/common/msg"
	"github.com/aserto-dev/ds-load/plugins/sdk/transform"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

type TransformCmd struct {
	TemplateFile string `name:"template-file" env:"DS_TEMPLATE_FILE" help:"transformation template file path" type:"path" optional:""`
	MaxChunkSize int    `name:"max-chunk-size" env:"DS_MAX_CHUNK_SIZE" help:"maximum chunk size" default:"20" optional:""`
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

	jsonWriter, err := js.NewJSONArrayWriter(os.Stdout)
	if err != nil {
		return err
	}
	defer jsonWriter.Close()

	tranformer := transform.NewTransformer(t.MaxChunkSize)
	reader, err := js.NewJSONArrayReader(os.Stdin)
	if err != nil {
		return err
	}

	for {
		var input map[string]interface{}
		err := reader.Read(&input)
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read input into map[string]interface{}")
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

		chunks := tranformer.PrepareChunks(&directoryObject)
		err = tranformer.WriteChunks(jsonWriter, chunks)
		if err != nil {
			return errors.Wrap(err, "failed to write chunks to output")
		}
	}

	return nil
}
