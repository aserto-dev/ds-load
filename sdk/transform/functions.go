package transform

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/dongri/phonenumber"
)

func customFunctions() map[string]any {
	return template.FuncMap{
		"last": func(x int, a interface{}) bool {
			return x == reflect.ValueOf(a).Len()-1
		},
		"contains":     strings.Contains,
		"separator":    separator,
		"marshal":      marshal,
		"fromEnv":      fromEnv,
		"phoneIso3166": phoneIso3166,
		"add": func(a int, b int) int {
			return a + b
		},
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

func marshal(v interface{}) string {
	a, _ := json.Marshal(v)
	return string(a)
}

func fromEnv(key, envName string) string {
	value := os.Getenv(envName)
	strValue, _ := json.Marshal(value)
	return fmt.Sprintf("%q:%s", key, string(strValue))
}

func phoneIso3166(phone string) string {
	country := phonenumber.GetISO3166ByNumber(phone, true)
	return phonenumber.ParseWithLandLine(phone, country.Alpha2)
}
