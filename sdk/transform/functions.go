package transform

import (
	"reflect"
	"strings"
	"text/template"
)

func customFunctions() map[string]any {
	return template.FuncMap{
		"last": func(x int, a interface{}) bool {
			return x == reflect.ValueOf(a).Len()-1
		},
		"contains":  strings.Contains,
		"separator": separator,
	}
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
