// get json from stdin and return json with v2 objects and v2 relations
package app

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"

	"text/template"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/common/msg"
	v2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

type TransformCmd struct {
	TemplateFile string `cmd:""`
	MaxChunkSize int    `cmd:""`
}

func (t *TransformCmd) Run(context *kong.Context) error {
	var template []byte
	var err error

	if t.TemplateFile == "" {
		template, err = Assets().ReadFile("assets/transform_template.tmpl")
		if err != nil {
			return err
		}
	} else {
		template, err = os.ReadFile(t.TemplateFile)
		if err != nil {
			return err
		}
	}

	if t.MaxChunkSize == 0 {
		t.MaxChunkSize = 10
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		inputText := scanner.Bytes()

		input := make(map[string]interface{})

		err = json.Unmarshal(inputText, &input)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal input into map[string]interface{}")
		}
		output, err := Transform(input, string(template))
		if err != nil {
			return errors.Wrap(err, "transform template execute failed")
		}
		var directoryObject msg.Transform

		err = protojson.Unmarshal([]byte(output), &directoryObject)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal transformed data into directory objects and relations")
		}

		objectChunks, relationChunks := t.prepareChunks(&directoryObject)

		err = writeResponse(os.Stdout, objectChunks, relationChunks)
		if err != nil {
			return errors.Wrap(err, "failed to write chunks to output")
		}
	}

	return nil
}

func writeResponse(writer io.Writer, objectChunks [][]*v2.Object, relationChunks [][]*v2.Relation) error {
	start := msg.PluginMessage{
		Data: &msg.PluginMessage_Batch{
			Batch: &msg.Batch{
				Data: &msg.Batch_Begin{Begin: true},
			},
		},
	}
	end := msg.PluginMessage{
		Data: &msg.PluginMessage_Batch{
			Batch: &msg.Batch{
				Data: &msg.Batch_End{End: true},
			},
		},
	}

	for _, chunk := range objectChunks {
		writeProtoMessage(writer, &start)
		for index := range chunk {
			message := msg.PluginMessage{
				Data: &msg.PluginMessage_Object{
					Object: chunk[index],
				},
			}
			err := writeProtoMessage(writer, &message)
			if err != nil {
				return err
			}
		}
		writeProtoMessage(writer, &end)
	}

	for _, chunk := range relationChunks {
		writeProtoMessage(writer, &start)
		for index := range chunk {
			message := msg.PluginMessage{
				Data: &msg.PluginMessage_Relation{
					Relation: chunk[index],
				},
			}
			err := writeProtoMessage(writer, &message)
			if err != nil {
				return err
			}
		}
		writeProtoMessage(writer, &end)
	}

	return nil
}

func writeProtoMessage(writer io.Writer, message *msg.PluginMessage) error {
	messageBytes, err := protojson.Marshal(message)
	if err != nil {
		return err
	}
	writer.Write(messageBytes)
	writer.Write([]byte("\n"))
	return nil
}

func (t *TransformCmd) prepareChunks(directoryObject *msg.Transform) ([][]*v2.Object, [][]*v2.Relation) {
	var objectChunks [][]*v2.Object
	var relationChunks [][]*v2.Relation
	if len(directoryObject.Objects) > t.MaxChunkSize {
		for i := 0; i < len(directoryObject.Objects); i += t.MaxChunkSize {
			end := i + t.MaxChunkSize
			if end > len(directoryObject.Objects) {
				end = len(directoryObject.Objects)
			}
			objectChunks = append(objectChunks, directoryObject.Objects[i:end])
		}
	} else {
		objectChunks = append(objectChunks, directoryObject.Objects)
	}
	if len(directoryObject.Relations) > t.MaxChunkSize {
		for i := 0; i < len(directoryObject.Relations); i += t.MaxChunkSize {
			end := i + t.MaxChunkSize
			if end > len(directoryObject.Relations) {
				end = len(directoryObject.Relations)
			}
			relationChunks = append(relationChunks, directoryObject.Relations[i:end])
		}
	} else {
		relationChunks = append(relationChunks, directoryObject.Relations)
	}
	return objectChunks, relationChunks
}

var fns = template.FuncMap{
	"last": func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len()-1
	},
}

func Transform(input map[string]interface{}, templateString string) (string, error) {
	t := template.New("transform")
	parsed, err := t.Funcs(fns).Parse(templateString)
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
