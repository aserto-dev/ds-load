// get json from stdin and return json with v2 objects and v2 relations
package app

import (
	"bytes"
	"errors"
	"reflect"

	"text/template"

	"github.com/alecthomas/kong"
)

type TransformCmd struct {
}

func (t *TransformCmd) Run(context *kong.Context) error {

	return errors.New("not implemented")
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
