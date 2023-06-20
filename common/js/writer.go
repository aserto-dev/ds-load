package js

import (
	"encoding/json"
	"io"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type JSONArrayWriter struct {
	writer       io.Writer
	addDelimiter bool
}

func NewJSONArrayWriter(w io.Writer) (*JSONArrayWriter, error) {
	_, err := w.Write([]byte{'['})
	if err != nil {
		return nil, err
	}
	return &JSONArrayWriter{
		writer:       w,
		addDelimiter: false,
	}, nil
}

func (w *JSONArrayWriter) WriteProtoMessage(message protoreflect.ProtoMessage) error {
	if w.addDelimiter {
		_, _ = w.writer.Write([]byte{','})
	}
	jsonObj, err := protojson.Marshal(message)
	if err != nil {
		return err
	}
	_, err = w.writer.Write(jsonObj)
	if !w.addDelimiter {
		w.addDelimiter = true
	}
	return err
}

func (w *JSONArrayWriter) Write(message any) error {
	if w.addDelimiter {
		_, _ = w.writer.Write([]byte{','})
	}
	jsonObj, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = w.writer.Write(jsonObj)
	if !w.addDelimiter {
		w.addDelimiter = true
	}
	return err
}

func (w *JSONArrayWriter) Close() error {
	if w.writer != nil {
		_, _ = w.writer.Write([]byte{']'})
		w.addDelimiter = false
		w.writer = nil
	}
	return nil
}
