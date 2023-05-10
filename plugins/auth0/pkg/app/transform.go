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
	v2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
)

type TransformCmd struct {
	TemplateFile string `cmd:""`
	MaxChunkSize int    `cmd:""`
}

type transformObject struct {
	Objects   []v2.Object   `json:"objects"`
	Relations []v2.Relation `json:"relations"`
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
			return err
		}
		output, err := Transform(input, string(template))
		if err != nil {
			return err
		}
		var directoryObject transformObject
		err = json.Unmarshal([]byte(output), &directoryObject)
		if err != nil {
			return err
		}

		err = t.writeObjects(os.Stdout, directoryObject)
		if err != nil {
			return err
		}

	}

	return nil
}

func (t *TransformCmd) writeObjects(writer io.Writer, directoryObject transformObject) error {
	if len(directoryObject.Objects) > t.MaxChunkSize {
		for i := 0; i < len(directoryObject.Objects); i += t.MaxChunkSize {
			end := i + t.MaxChunkSize
			if end > len(directoryObject.Objects) {
				end = len(directoryObject.Objects)
			}
			chunk := directoryObject.Objects[i:end]
			chunkBytes, err := json.Marshal(chunk)
			if err != nil {
				return err
			}
			writer.Write(chunkBytes)
			writer.Write([]byte("\n"))
		}
	} else {
		chunk, err := json.Marshal(directoryObject.Objects)
		if err != nil {
			return err
		}
		writer.Write(chunk)
		writer.Write([]byte("\n"))

	}
	if len(directoryObject.Relations) > t.MaxChunkSize {
		for i := 0; i < len(directoryObject.Relations); i += t.MaxChunkSize {
			end := i + t.MaxChunkSize
			if end > len(directoryObject.Relations) {
				end = len(directoryObject.Relations)
			}
			chunk := directoryObject.Relations[i:end]
			chunkBytes, err := json.Marshal(chunk)
			if err != nil {
				return err
			}
			writer.Write(chunkBytes)
			writer.Write([]byte("\n"))
		}
	} else {
		chunk, err := json.Marshal(directoryObject.Relations)
		if err != nil {
			return err
		}
		writer.Write(chunk)
		writer.Write([]byte("\n"))
	}
	return nil
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
