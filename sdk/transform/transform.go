package transform

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/common/msg"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

type GoTemplateTransform struct {
	template []byte
}

func NewGoTemplateTransform(transformTemplate []byte) *GoTemplateTransform {
	return &GoTemplateTransform{
		template: transformTemplate,
	}
}

func (t *GoTemplateTransform) ExportTransform(outputWriter io.Writer) error {
	_, err := outputWriter.Write(t.template)
	if err != nil {
		log.Fatalf("cannot write to output: %s", err.Error())
		return err
	}

	return nil
}

func (t *GoTemplateTransform) Transform(
	ctx context.Context,
	ioReader io.Reader,
	outputWriter,
	errorWriter io.Writer,
) error {
	jsonWriter := js.NewJSONArrayWriter(outputWriter)
	defer jsonWriter.Close()

	reader, err := js.NewJSONArrayReader(ioReader)
	if err != nil {
		return err
	}

	for {
		var idpData map[string]any

		err := reader.Read(&idpData)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return errors.Wrap(err, "failed to read idpData into map[string]any")
		}

		if err := t.doTransform(idpData, jsonWriter); err != nil {
			return err
		}
	}

	return nil
}

func (t *GoTemplateTransform) doTransform(idpData map[string]any, jsonWriter *js.JSONArrayWriter) error {
	dirV3msg, err := t.TransformObject(idpData)
	if err != nil {
		return errors.Wrap(err, "failed to transform idpData into directory objects and relations")
	}

	if err := jsonWriter.WriteProtoMessage(dirV3msg); err != nil {
		return errors.Wrap(err, "failed to write directory objects to output")
	}

	return nil
}

func (t *GoTemplateTransform) TransformObject(idpData map[string]any) (*msg.Transform, error) {
	output, err := t.transformToTemplate(idpData, string(t.template))
	if err != nil {
		return nil, errors.Wrap(err, "GoTemplateTransform transformTemplate execute failed")
	}

	if os.Getenv("DEBUG") != "" {
		if _, err := os.Stdout.WriteString(output); err != nil {
			return nil, errors.Wrap(err, "failed to write to stdout")
		}
	}

	var dirV3msg msg.Transform

	opts := protojson.UnmarshalOptions{
		AllowPartial:   false,
		DiscardUnknown: false,
	}

	if err := opts.Unmarshal([]byte(output), &dirV3msg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal transformed data into directory v3 objects and relations")
	}

	return &dirV3msg, nil
}

func (t *GoTemplateTransform) transformToTemplate(input map[string]any, templateString string) (string, error) {
	temp := template.New("GoTemplateTransform")

	parsed, err := temp.Funcs(customFunctions()).Parse(templateString)
	if err != nil {
		return "", err
	}

	var filled bytes.Buffer

	if err := parsed.Execute(&filled, input); err != nil {
		return "", err
	}

	// pass input to templates
	return filled.String(), nil
}
