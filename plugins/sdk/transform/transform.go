package transform

import (
	"bytes"
	"html/template"
	"reflect"
	"strings"

	"github.com/aserto-dev/ds-load/common/js"
	"github.com/aserto-dev/ds-load/common/msg"

	v2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
)

type Tranformer struct {
	MaxChunkSize int
}

func NewTransformer(chunkSize int) *Tranformer {
	return &Tranformer{
		MaxChunkSize: chunkSize,
	}
}

func (t *Tranformer) WriteChunks(writer *js.JSONArrayWriter, chunks []*msg.Transform) error {
	for _, chunk := range chunks {
		err := t.writeProtoMessage(writer, chunk)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Tranformer) writeProtoMessage(writer *js.JSONArrayWriter, message *msg.Transform) error {
	err := writer.WriteProtoMessage(message)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tranformer) PrepareChunks(directoryObject *msg.Transform) []*msg.Transform {
	var chunks []*msg.Transform
	var freeChunk *msg.Transform

	for _, obj := range directoryObject.Objects {
		freeChunk, chunks = t.nextFreeChunk(chunks)
		freeChunk.Objects = append(freeChunk.Objects, obj)
	}

	for _, rel := range directoryObject.Relations {
		freeChunk, chunks = t.nextFreeChunk(chunks)
		freeChunk.Relations = append(freeChunk.Relations, rel)
	}

	return chunks
}

func (t *Tranformer) nextFreeChunk(chunks []*msg.Transform) (*msg.Transform, []*msg.Transform) {
	if len(chunks) == 0 {
		chunks = t.addEmptyChunk(chunks)
	}

	lastChunk := chunks[len(chunks)-1]
	if t.isRoomInChunk(lastChunk) {
		return lastChunk, chunks
	}

	chunks = t.addEmptyChunk(chunks)

	return chunks[len(chunks)-1], chunks
}

func (t *Tranformer) isRoomInChunk(chunk *msg.Transform) bool {
	return len(chunk.Objects)+len(chunk.Relations) < t.MaxChunkSize
}

func (t *Tranformer) addEmptyChunk(chunks []*msg.Transform) []*msg.Transform {
	return append(chunks, &msg.Transform{Objects: []*v2.Object{}, Relations: []*v2.Relation{}})
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
