package transform

import (
	"bytes"
	"html/template"
	"io"
	"reflect"

	"github.com/aserto-dev/ds-load/common/msg"
	v2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

type Tranformer struct {
	MaxChunkSize int
}

func NewTransformer(chunkSize int) *Tranformer {
	return &Tranformer{MaxChunkSize: chunkSize}
}

func (t *Tranformer) WriteChunks(writer io.Writer, objectChunks [][]*v2.Object, relationChunks [][]*v2.Relation) error {
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
			err := t.writeProtoMessage(writer, &start)
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
			err := t.writeProtoMessage(writer, &message)
			if err != nil {
				return err
			}
		}
		if t.MaxChunkSize > 1 {
			err := t.writeProtoMessage(writer, &end)
			if err != nil {
				return err
			}
		}
	}

	for _, chunk := range relationChunks {
		if t.MaxChunkSize > 1 {
			err := t.writeProtoMessage(writer, &start)
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
			err := t.writeProtoMessage(writer, &message)
			if err != nil {
				return err
			}
		}
		if t.MaxChunkSize > 1 {
			err := t.writeProtoMessage(writer, &end)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Tranformer) writeProtoMessage(writer io.Writer, message *msg.PluginMessage) error {
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

func (t *Tranformer) PrepareChunks(directoryObject *msg.Transform) ([][]*v2.Object, [][]*v2.Relation) {
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

func (t *Tranformer) TransformToTemplate(input map[string]interface{}, templateString string) (string, error) {
	temp := template.New("transform")
	parsed, err := temp.Funcs(fns).Parse(templateString)
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
