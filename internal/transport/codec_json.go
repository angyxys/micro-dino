package transport

import (
	"bytes"
	"encoding/json"
	"errors"
)

const CodecJSONName = "json"

type JSONCodec struct{}

func NewJSONCodec() *JSONCodec { return &JSONCodec{} }

func (c *JSONCodec) Name() string { return CodecJSONName }

func (c *JSONCodec) Marshal(m *Message) ([]byte, error) {
	if m == nil {
		return nil, errors.New("micro-dino: cannot marshal nil message")
	}
	return json.Marshal(m)
}

func (c *JSONCodec) Unmarshal(data []byte, m *Message) error {
	if m == nil {
		return errors.New("micro-dino: cannot unmarshal into nil message")
	}
	if len(data) == 0 {
		return errors.New("micro-dino: empty payload")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	return dec.Decode(m)
}
