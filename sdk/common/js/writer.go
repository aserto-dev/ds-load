package js

import (
	"encoding/json"
	"io"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type JSONArrayWriter struct {
	writer           io.Writer
	addDelimiter     bool
	arrayInitialized bool
}

func NewJSONArrayWriter(w io.Writer) (*JSONArrayWriter, error) {
	return &JSONArrayWriter{
		writer:           w,
		addDelimiter:     false,
		arrayInitialized: false,
	}, nil
}

func (w *JSONArrayWriter) WriteProtoMessage(message protoreflect.ProtoMessage) error {
	err := w.writeDelimiters()
	if err != nil {
		return err
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
	err := w.writeDelimiters()
	if err != nil {
		return err
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
		if w.arrayInitialized {
			_, err := w.writer.Write([]byte{']', '\n'})
			if err != nil {
				return err
			}
		}
		w.addDelimiter = false
		w.writer = nil
	}
	return nil
}

func (w *JSONArrayWriter) writeDelimiters() error {
	if !w.arrayInitialized {
		_, err := w.writer.Write([]byte{'['})
		if err != nil {
			return err
		}
		w.arrayInitialized = true
	}
	if w.addDelimiter {
		_, err := w.writer.Write([]byte{','})
		if err != nil {
			return err
		}
	}
	return nil
}
