// get json from stdin and return json with v2 objects and v2 relations
package app

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strings"

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

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		inputText := scanner.Bytes()

		input := make(map[string]interface{})

		err = json.Unmarshal(inputText, &input)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal input into map[string]interface{}")
		}
		output, err := Transform(input, string(templateContent))
		if err != nil {
			return errors.Wrap(err, "transform template execute failed")
		}
		var directoryObject msg.Transform

		err = protojson.Unmarshal([]byte(output), &directoryObject)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal transformed data into directory objects and relations")
		}

		objectChunks, relationChunks := t.prepareChunks(&directoryObject)

		err = t.writeResponse(os.Stdout, objectChunks, relationChunks)
		if err != nil {
			return errors.Wrap(err, "failed to write chunks to output")
		}
	}

	return nil
}

func (t *TransformCmd) writeResponse(writer io.Writer, objectChunks [][]*v2.Object, relationChunks [][]*v2.Relation) error {
	start := msg.PluginMessage{
		Data: &msg.PluginMessage_Batch{
			Batch: &msg.Batch{
				Type: msg.BatchType_BEGIN,
			},
		},
	}
	end := msg.PluginMessage{
		Data: &msg.PluginMessage_Batch{
			Batch: &msg.Batch{
				Type: msg.BatchType_END,
			},
		},
	}

	for _, chunk := range objectChunks {
		if t.MaxChunkSize > 1 {
			err := writeProtoMessage(writer, &start)
			if err != nil {
				return err
			}
		}
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
		if t.MaxChunkSize > 1 {
			err := writeProtoMessage(writer, &end)
			if err != nil {
				return err
			}
		}
	}

	for _, chunk := range relationChunks {
		if t.MaxChunkSize > 1 {
			err := writeProtoMessage(writer, &start)
			if err != nil {
				return err
			}
		}
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
		if t.MaxChunkSize > 1 {
			err := writeProtoMessage(writer, &end)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func writeProtoMessage(writer io.Writer, message *msg.PluginMessage) error {
	messageBytes, err := protojson.Marshal(message)
	if err != nil {
		return err
	}
	_, err = writer.Write(messageBytes)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte("\n"))
	if err != nil {
		return err
	}
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
	"marshal": func(array []interface{}) string {
		builder := strings.Builder{}
		for i := range array {
			content, err := json.MarshalIndent(array[i], "", "  ")
			if err != nil {
				return "error marshaling object"
			}
			builder.WriteString(string(content))
			if i < len(array)-1 {
				builder.WriteString(",\n")
			}
		}
		return builder.String()
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
