package fetcher

import (
	"encoding/json"
	"iter"
)

func YieldError(err error) iter.Seq2[map[string]any, error] {
	return func(yield func(map[string]any, error) bool) {
		(yield(nil, err))
	}
}

func YieldMap[T any](objects []T, marshaller func(T) ([]byte, error)) iter.Seq2[map[string]any, error] {
	return func(yield func(map[string]any, error) bool) {
		for _, object := range objects {
			objectBytes, err := marshaller(object)
			if err != nil && !yield(nil, err) {
				return
			}

			var obj map[string]any
			err = json.Unmarshal(objectBytes, &obj)

			if !(yield(obj, err)) {
				return
			}
		}
	}
}
