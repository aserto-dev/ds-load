package js

import (
	"encoding/json"
	"io"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type JSONArrayWriter struct {
	writer io.Writer
	first  bool
}

func NewJSONArrayWriter(w io.Writer) (*JSONArrayWriter, error) {
	_, err := w.Write([]byte{'['})
	if err != nil {
		return nil, err
	}
	return &JSONArrayWriter{
		writer: w,
		first:  false,
	}, nil
}

func (w *JSONArrayWriter) WriteProtoMessage(m protoreflect.ProtoMessage) error {
	if w.first {
		_, _ = w.writer.Write([]byte{','})
	}
	jsonObj, err := protojson.Marshal(m)
	if err != nil {
		return err
	}
	_, err = w.writer.Write(jsonObj)
	if !w.first {
		w.first = true
	}
	return err
}

func (w *JSONArrayWriter) Write(m any) error {
	if w.first {
		_, _ = w.writer.Write([]byte{','})
	}
	jsonObj, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = w.writer.Write(jsonObj)
	if !w.first {
		w.first = true
	}
	return err
}

func (w *JSONArrayWriter) Close() error {
	if w.writer != nil {
		_, _ = w.writer.Write([]byte{']'})
		w.first = false
		w.writer = nil
	}
	return nil
}
