package js

import (
	"io"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type JSONWriter struct {
	writer io.Writer
	first  bool
}

func NewJSONArrayWriter(w io.Writer) *JSONWriter {
	w.Write([]byte{'['})
	return &JSONWriter{
		writer: w,
		first:  false,
	}
}

func (w *JSONWriter) WriteProtoMessage(m protoreflect.ProtoMessage) error {
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

func (w *JSONWriter) Close() error {
	if w.writer != nil {
		_, _ = w.writer.Write([]byte{']'})
		w.first = false
		w.writer = nil
	}
	return nil
}
