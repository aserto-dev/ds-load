package transform

import (
	"bytes"
	"context"
	"io"
	"os"
	"text/template"

	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/common/msg"
	"github.com/aserto-dev/ds-load/sdk/plugin"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

type transform struct{}

func New() plugin.Transform {
	return &transform{}
}

func (t *transform) Transform(ctx context.Context, ioReader io.Reader, outputWriter, errorWriter io.Writer, transformTemplate []byte) error {
	jsonWriter, err := js.NewJSONArrayWriter(outputWriter)
	if err != nil {
		return err
	}
	defer jsonWriter.Close()
	reader, err := js.NewJSONArrayReader(ioReader)
	if err != nil {
		return err
	}

	for {
		var idpData map[string]interface{}
		err := reader.Read(&idpData)
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read idpData into map[string]interface{}")
		}
		err = t.doTransform(idpData, jsonWriter, transformTemplate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *transform) ExportDefaultTemplate(outputWriter io.Writer) {

}

func (t *transform) doTransform(idpData map[string]interface{}, jsonWriter *js.JSONArrayWriter, transformTemplate []byte) error {
	output, err := t.transformToTemplate(idpData, string(transformTemplate))
	if err != nil {
		return errors.Wrap(err, "transform transformTemplate execute failed")
	}
	if os.Getenv("DEBUG") != "" {
		os.Stdout.WriteString(output)
	}
	var directoryObject msg.Transform

	err = protojson.Unmarshal([]byte(output), &directoryObject)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal transformed data into directory objects and relations")
	}

	err = jsonWriter.WriteProtoMessage(&directoryObject)
	if err != nil {
		return errors.Wrap(err, "failed to write directory objects to output")
	}
	return nil
}

func (t *transform) transformToTemplate(input map[string]interface{}, templateString string) (string, error) {
	temp := template.New("transform")
	parsed, err := temp.Funcs(customFunctions()).Parse(templateString)
	if err != nil {
		return "", err
	}
	var filled bytes.Buffer
	err = parsed.Execute(&filled, input)
	if err != nil {
		return "", err
	}

	// pass input to templates
	return filled.String(), nil
}
