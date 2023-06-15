package transform

import (
	"bytes"
	v2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
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
	var chunkIndex = 0
	chunks = append(chunks, &msg.Transform{Objects: []*v2.Object{}, Relations: []*v2.Relation{}})

	for _, obj := range directoryObject.Objects {
		if len(chunks[chunkIndex].Objects) >= t.MaxChunkSize {
			chunkIndex += 1
			chunks = append(chunks, &msg.Transform{Objects: []*v2.Object{}, Relations: []*v2.Relation{}})
		}

		chunks[chunkIndex].Objects = append(chunks[chunkIndex].Objects, obj)
	}

	for _, rel := range directoryObject.Relations {
		if len(chunks[chunkIndex].Objects)+len(chunks[chunkIndex].Relations) >= t.MaxChunkSize {
			chunkIndex += 1
			chunks = append(chunks, &msg.Transform{Objects: []*v2.Object{}, Relations: []*v2.Relation{}})
		}
		chunks[chunkIndex].Relations = append(chunks[chunkIndex].Relations, rel)
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
