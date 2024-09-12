// Package jsonwrap contains json wrappers for jsoniter.
package jsonwrap

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

// Decoder is a JSON decoder.
type Decoder struct {
	dec *jsoniter.Decoder
}

// Decode decode JSON into interface{}.
func (d *Decoder) Decode(obj interface{}) error {
	return handleDecodingError(d.dec.Decode(obj))
}

// Marshal returns the JSON encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

// NewDecoder returns a new decoder that reads from reader.
func NewDecoder(reader io.Reader) *Decoder {
	return newDecoder(reader, false)
}

// NewDecoderStrict returns a new decoder that reads from reader and returns an error if an unexpected field is decoded.
func NewDecoderStrict(reader io.Reader) *Decoder {
	return newDecoder(reader, true)
}

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	return handleDecodingError(jsoniter.Unmarshal(data, v))
}

func newDecoder(reader io.Reader, strict bool) *Decoder {
	dec := jsoniter.NewDecoder(reader)

	if strict {
		dec.DisallowUnknownFields()
	}

	return &Decoder{
		dec: dec,
	}
}
