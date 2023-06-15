package transform

import (
	"bytes"
	"html/template"
	"reflect"
	"strings"

	"github.com/aserto-dev/ds-load/common/js"
	"github.com/aserto-dev/ds-load/common/msg"
)

type Tranformer struct {
	MaxChunkSize int
}

func NewTransformer(chunkSize int) *Tranformer {
	return &Tranformer{
		MaxChunkSize: chunkSize,
	}
}

func (t *Tranformer) WriteChunks(writer *js.JSONWriter, chunks []*msg.Transform) error {
	for _, chunk := range chunks {
		err := t.writeProtoMessage(writer, chunk)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Tranformer) writeProtoMessage(writer *js.JSONWriter, message *msg.Transform) error {
	// messageBytes, err := protojson.Marshal(message)
	// if err != nil {
	// 	return err
	// }
	err := writer.WriteProtoMessage(message)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tranformer) PrepareChunks(directoryObject *msg.Transform) []*msg.Transform {
	var chunks []*msg.Transform

	if len(directoryObject.Objects) > t.MaxChunkSize {
		for i := 0; i < len(directoryObject.Objects); i += t.MaxChunkSize {
			end := i + t.MaxChunkSize
			if end > len(directoryObject.Objects) {
				end = len(directoryObject.Objects)
			}
			chunks = append(chunks, &msg.Transform{
				Objects: directoryObject.Objects[i:end],
			})
		}
	} else {
		chunks = append(chunks, &msg.Transform{
			Objects: directoryObject.Objects,
		})
	}

	remaining := t.MaxChunkSize - len(chunks[len(chunks)-1].Objects)
	if remaining > 0 {
		if remaining < len(directoryObject.Relations) {
			chunks[len(chunks)-1].Relations = directoryObject.Relations[:remaining-1]
		} else {
			chunks[len(chunks)-1].Relations = directoryObject.Relations
			return chunks
		}
	}

	if len(directoryObject.Relations)-remaining > t.MaxChunkSize {
		for i := remaining; i < len(directoryObject.Relations); i += t.MaxChunkSize {
			end := i + t.MaxChunkSize
			if end > len(directoryObject.Relations) {
				end = len(directoryObject.Relations)
			}
			chunks = append(chunks, &msg.Transform{
				Relations: directoryObject.Relations[i:end],
			})
		}
	} else {
		chunks = append(chunks, &msg.Transform{
			Relations: directoryObject.Relations[remaining:],
		})
	}
	return chunks
}

var fns = template.FuncMap{
	"last": func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len()-1
	},
	"contains": strings.Contains,
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
