package transform

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/dongri/phonenumber"
	"github.com/rs/zerolog/log"
)

func customFunctions() template.FuncMap {
	f := sprig.GenericFuncMap()
	delete(f, "last")

	extra := template.FuncMap{
		"last": func(x int, a interface{}) bool {
			return x == reflect.ValueOf(a).Len()-1
		},
		"separator":    separator,
		"marshal":      marshal,
		"fromEnv":      fromEnv,
		"phoneIso3166": phoneIso3166,
		"array_contains": func(a []interface{}, b string) bool {
			for _, x := range a {
				stringX, ok := x.(string)
				if !ok {
					return false
				}

				if stringX == b {
					return true
				}
			}
			return false
		},
	}

	for k, v := range extra {
		if _, ok := f[k]; !ok {
			f[k] = v
		} else {
			fmt.Fprintf(os.Stderr, "duplicate template func [%s]\n", k)
		}
	}

	return f
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

func marshal(v any) string {
	a, err := json.Marshal(v)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal any")
	}

	return string(a)
}

func fromEnv(key, envName string) string {
	value := os.Getenv(envName)
	strValue, err := json.Marshal(value)

	if err != nil {
		log.Error().Err(err).Msg("failed to marshal value")
	}

	return fmt.Sprintf("%q:%s", key, string(strValue))
}

func phoneIso3166(phone string) string {
	country := phonenumber.GetISO3166ByNumber(phone, true)
	return phonenumber.ParseWithLandLine(phone, country.Alpha2)
}
