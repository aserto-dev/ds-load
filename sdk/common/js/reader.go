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

func NewJSONArrayReader(r io.Reader) (*JSONArrayReader, error) {
	decoder := json.NewDecoder(r)

	tk, err := decoder.Token()
	if err != nil {
		return nil, errors.Wrap(err, "could not read json output")
	}

	delim, ok := tk.(json.Delim)
	if !ok {
		return nil, errors.New("first token not a delimiter")
	}

	if delim != json.Delim('[') {
		return nil, errors.New("first token not a [")
	}

	return &JSONArrayReader{
		decoder: decoder,
	}, nil
}

// reads next json object as proto message
// returns io.EOF at the end of the input stream.
func (r *JSONArrayReader) ReadProtoMessage(message proto.Message) error {
	more, err := r.more()
	if err != nil {
		return err
	}
	if more {
		if err := UnmarshalNext(r.decoder, message); err != nil {
			return err
		}
	}
	return nil
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
	// check for next token in case of multiple json sources
	tok, err = r.decoder.Token()
	if err != nil {
		return false, err
	}

	delim, ok := tok.(json.Delim)
	if !ok {
		return false, errors.Wrap(err, "first token not a delimiter")
	}

	if delim == json.Delim('[') {
		return true, nil
	} else {
		return false, errors.Wrap(err, "first token not a [")
	}
}

// func (r *JSONReader) Read(results chan any) error {
// 	for {
// 		tk, err := r.decoder.Token()
// 		if err != nil {
// 			return errors.Wrap(err, "could not read")
// 		}

// 		delim, ok := tk.(json.Delim)
// 		if !ok {
// 			return errors.Wrap(err, "first token not a delimiter")
// 		}

// 		if delim == json.Delim('{') {
// 			results <- BatchBegin{}
// 			for r.decoder.More() {
// 				t, err := r.decoder.Token()
// 				if err != nil {
// 					return err
// 				}

// 				keyStr := ""
// 				if key, ok := t.(string); ok {
// 					keyStr = key
// 				}

// 				tok, err := r.decoder.Token()
// 				if err != nil {
// 					return err
// 				}
// 				if delim, ok := tok.(json.Delim); !ok && delim.String() != "[" {
// 					return errors.Errorf("%s does not contain a JSON array", keyStr)
// 				}

// 				for r.decoder.More() {
// 					switch keyStr {
// 					case "objects":
// 						var obj v2.Object
// 						err := UnmarshalNext(r.decoder, &obj)
// 						if err != nil {
// 							return err
// 						}
// 						results <- obj
// 					case "relations":
// 						var rel v2.Relation
// 						err := UnmarshalNext(r.decoder, &rel)
// 						if err != nil {
// 							return err
// 						}
// 						results <- rel
// 					}
// 				}

// 				tok, err = r.decoder.Token()
// 				if err != nil {
// 					return err
// 				}
// 				if delim, ok := tok.(json.Delim); !ok && delim.String() != "]" {
// 					return errors.Errorf("%s does not contain a JSON array", keyStr)
// 				}
// 			}
// 		}
// 		if delim == json.Delim('}') {
// 			results <- BatchEnd{}

// 			tk, err := r.decoder.Token()
// 			if err != nil {
// 				return errors.Wrap(err, "could not read")
// 			}

// 			delim, ok := tk.(json.Delim)
// 			if !ok {
// 				return errors.Wrap(err, "first token not a delimiter")
// 			}
// 			if delim == json.Delim('}') {
// 				// reached end of json
// 				break
// 			}
// 		}
// 	}

// 	return nil
// }

func UnmarshalNext(d *json.Decoder, m proto.Message) error {
	var b json.RawMessage
	if err := d.Decode(&b); err != nil {
		return err
	}

	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}.Unmarshal(b, m)
}
