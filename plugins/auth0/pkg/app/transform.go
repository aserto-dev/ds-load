// get json from stdin and return json with v2 objects and v2 relations
package app

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"

	"text/template"

	"github.com/alecthomas/kong"
)

type TransformCmd struct {
	TemplateFile string `cmd:""`
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
	inputText, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	input := make(map[string]interface{})

	err = json.Unmarshal(inputText, &input)
	if err != nil {
		return err
	}
	output, err := Transform(input, string(template))
	if err != nil {
		return err
	}
	os.Stdout.WriteString(output)

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
