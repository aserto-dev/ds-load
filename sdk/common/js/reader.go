package js

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type JSONArrayReader struct {
	decoder *json.Decoder
}

var ErrInvalidInput = errors.New("invalid input")

func NewJSONArrayReader(r io.Reader) (*JSONArrayReader, error) {
	decoder := json.NewDecoder(r)

	tk, err := decoder.Token()
	if err != nil {
		return nil, errors.Wrap(err, "could not read json output")
	}

	delim, ok := tk.(json.Delim)
	if !ok {
		return nil, errors.Wrap(ErrInvalidInput, "first token not a json delimiter")
	}

	if delim != json.Delim('[') {
		return nil, errors.Wrapf(ErrInvalidInput, "unexpected delimiter. expected '['. found '%s'", delim)
	}

	return &JSONArrayReader{
		decoder: decoder,
	}, nil
}

// reads next json object as proto message
// returns io.EOF at the end of the input stream.
func (r *JSONArrayReader) ReadProtoMessage(message proto.Message) error {
	more, err := r.more()
	switch {
	case err != nil:
		return err
	case more:
		return UnmarshalNext(r.decoder, message)
	default:
		return io.EOF
	}
}

func (r *JSONArrayReader) Read(message any) error {
	more, err := r.more()
	if err != nil {
		return err
	}
	if more {
		if err := r.decoder.Decode(&message); err != nil {
			return err
		}
	}
	return nil
}

func (r *JSONArrayReader) more() (bool, error) {
	if r.decoder.More() {
		return true, nil
	}

	// if no more data in array read ] character at end of array
	tok, err := r.decoder.Token()
	if err != nil {
		return false, err
	}
	if delim, ok := tok.(json.Delim); !ok && delim.String() != "]" {
		return false, errors.Errorf("file does not contain a JSON array")
	}

	return false, io.EOF
}

func UnmarshalNext(d *json.Decoder, m proto.Message) error {
	var b json.RawMessage
	if err := d.Decode(&b); err != nil {
		return err
	}

	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}.Unmarshal(b, m)
}
