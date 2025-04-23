package fetcher

import (
	"encoding/json"
	"iter"
)

func YieldError(err error) iter.Seq2[map[string]any, error] {
	return func(yield func(map[string]any, error) bool) {
		if !(yield(nil, err)) {
			return
		}
	}
}

func YieldMap[T any](marshaller func(any) ([]byte, error), objects []T) iter.Seq2[map[string]any, error] {
	return func(yield func(map[string]any, error) bool) {
		for _, object := range objects {
			groupBytes, err := marshaller(object)
			if err != nil && !yield(nil, err) {
				return
			}

			var obj map[string]any
			err = json.Unmarshal(groupBytes, &obj)

			if !(yield(obj, err)) {
				return
			}
		}
	}
}
