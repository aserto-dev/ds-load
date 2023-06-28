package plugin

import (
	"io"
	"log"
	"os"

	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/common/msg"
	"github.com/aserto-dev/ds-load/sdk/transform"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

type DSPlugin struct {
	template     []byte
	outWriter    io.Writer
	errWriter    io.Writer
	maxChunkSize int
	transformer  *transform.Transformer
}

func NewDSPlugin(options ...PluginOption) *DSPlugin {
	plugin := &DSPlugin{}
	for _, o := range options {
		o(plugin)
	}

	if plugin.outWriter == nil {
		plugin.outWriter = os.Stdout
	}

	if plugin.errWriter == nil {
		plugin.errWriter = os.Stderr
	}

	if plugin.maxChunkSize == 0 {
		plugin.maxChunkSize = 20
	}

	plugin.transformer = transform.NewTransformer(plugin.maxChunkSize)

	return plugin
}

// json encodes results and prints to plugin writer.
func (plugin *DSPlugin) WriteFetchOutput(results chan map[string]interface{}, errCh chan error, transformMessage bool) error {
	go func() {
		for err := range errCh {
			_, wErr := plugin.errWriter.Write([]byte(err.Error() + "\n"))
			if wErr != nil {
				log.Fatalf("cannot write to output: %s", wErr.Error())
			}
		}
	}()

	writer, err := js.NewJSONArrayWriter(plugin.outWriter)
	if err != nil {
		return err
	}
	defer writer.Close()
	for result := range results {
		err := writer.Write(result)
		if err != nil {
			return err
		}
		if transformMessage {
			err = plugin.doTransform(result, writer)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (plugin *DSPlugin) Transform() error {
	jsonWriter, err := js.NewJSONArrayWriter(plugin.outWriter)
	if err != nil {
		return err
	}
	defer jsonWriter.Close()
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
		err = plugin.doTransform(input, jsonWriter)
		if err != nil {
			return err
		}
	}

	return nil
}

func (plugin *DSPlugin) ExportTransform() error {
	_, err := plugin.outWriter.Write(plugin.template)
	if err != nil {
		log.Fatalf("cannot write to output: %s", err.Error())
	}

	return nil
}

func (plugin *DSPlugin) doTransform(input map[string]interface{}, jsonWriter *js.JSONArrayWriter) error {
	output, err := plugin.transformer.TransformToTemplate(input, string(plugin.template))
	if err != nil {
		return errors.Wrap(err, "transform template execute failed")
	}
	if os.Getenv("DEBUG") != "" {
		os.Stdout.WriteString(output)
	}
	var directoryObject msg.Transform

	err = protojson.Unmarshal([]byte(output), &directoryObject)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal transformed data into directory objects and relations")
	}

	chunks := plugin.transformer.PrepareChunks(&directoryObject)
	err = plugin.transformer.WriteChunks(jsonWriter, chunks)
	if err != nil {
		return errors.Wrap(err, "failed to write chunks to output")
	}
	return nil
}
