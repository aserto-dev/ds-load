package transform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/common/msg"
	"github.com/dongri/phonenumber"

	v2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
)

type Transformer struct {
	MaxChunkSize int
}

func NewTransformer(chunkSize int) *Transformer {
	return &Transformer{
		MaxChunkSize: chunkSize,
	}
}

func (t *Transformer) WriteChunks(writer *js.JSONArrayWriter, chunks []*msg.Transform) error {
	for _, chunk := range chunks {
		err := t.writeProtoMessage(writer, chunk)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Transformer) writeProtoMessage(writer *js.JSONArrayWriter, message *msg.Transform) error {
	err := writer.WriteProtoMessage(message)
	if err != nil {
		return err
	}

	return nil
}

func (t *Transformer) PrepareChunks(directoryObject *msg.Transform) []*msg.Transform {
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

func (t *Transformer) nextFreeChunk(chunks []*msg.Transform) (*msg.Transform, []*msg.Transform) {
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

func (t *Transformer) isRoomInChunk(chunk *msg.Transform) bool {
	return len(chunk.Objects)+len(chunk.Relations) < t.MaxChunkSize
}

func (t *Transformer) addEmptyChunk(chunks []*msg.Transform) []*msg.Transform {
	return append(chunks, &msg.Transform{Objects: []*v2.Object{}, Relations: []*v2.Relation{}})
}

var fns = template.FuncMap{
	"last": func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len()-1
	},
	"contains":  strings.Contains,
	"separator": separator,
	"marshal": func(v interface{}) string {
		a, _ := json.Marshal(v)
		return string(a)
	},
	"fromEnv": func(key, envName string) string {
		value := os.Getenv(envName)
		strValue, _ := json.Marshal(value)
		return fmt.Sprintf("%q:%s", key, string(strValue))
	},
	"phoneIso3166": func(phone string) string {
		country := phonenumber.GetISO3166ByNumber(phone, true)
		return phonenumber.ParseWithLandLine(phone, country.Alpha2)
	},
}

func separator(s string) func() string {
	i := -1
	return func() string {
		i++
		if i == 0 {
			return ""
		}
		return s
	}
}

func (t *Transformer) TransformToTemplate(input map[string]interface{}, templateString string) (string, error) {
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
