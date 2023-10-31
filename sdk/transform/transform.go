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
	convert "github.com/aserto-dev/go-directory/pkg/convert"
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

func (t *GoTemplateTransform) Transform(ctx context.Context, ioReader io.Reader, outputWriter, errorWriter io.Writer) error {
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
		err = t.doTransform(idpData, jsonWriter, t.template, errorWriter)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *GoTemplateTransform) doTransform(idpData map[string]interface{}, jsonWriter *js.JSONArrayWriter, transformTemplate []byte, errorWriter io.Writer) error {
	output, err := t.transformToTemplate(idpData, string(transformTemplate))
	if err != nil {
		return errors.Wrap(err, "GoTemplateTransform transformTemplate execute failed")
	}
	if os.Getenv("DEBUG") != "" {
		os.Stdout.WriteString(output)
	}
	var dirV3msg msg.Transform

	opts := protojson.UnmarshalOptions{
		AllowPartial:   false,
		DiscardUnknown: false,
	}
	err = opts.Unmarshal([]byte(output), &dirV3msg)
	if err != nil {
		errorWriter.Write([]byte(err.Error()))
		var dirV2msg msg.TransformV2
		err = opts.Unmarshal([]byte(output), &dirV2msg)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal transformed data into directory objects and relations")
		}
		dirV3msg.Objects = convert.ObjectArrayToV3(dirV2msg.Objects)
		dirV3msg.Relations = convert.RelationArrayToV3(dirV2msg.Relations)
	}

	err = jsonWriter.WriteProtoMessage(&dirV3msg)
	if err != nil {
		return errors.Wrap(err, "failed to write directory objects to output")
	}
	return nil
}

func (t *GoTemplateTransform) transformToTemplate(input map[string]interface{}, templateString string) (string, error) {
	temp := template.New("GoTemplateTransform")
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
