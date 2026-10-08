package transport

import "fmt"

type Codec interface {
	Name() string

	Marshal(m *Message) ([]byte, error)

	Unmarshal(data []byte, m *Message) error
}

var codecRegistry = map[string]Codec{}

func RegisterCodec(c Codec) {
	if c == nil {
		return
	}
	codecRegistry[c.Name()] = c
}

func GetCodec(name string) (Codec, bool) {
	c, ok := codecRegistry[name]
	return c, ok
}

func MustCodec(name string) Codec {
	c, ok := GetCodec(name)
	if !ok {
		panic(fmt.Sprintf("micro-dino: codec %q not registered", name))
	}
	return c
}

func RegisteredCodecs() []string {
	out := make([]string, 0, len(codecRegistry))
	for name := range codecRegistry {
		out = append(out, name)
	}
	return out
}

func init() {
	RegisterCodec(NewJSONCodec())
}
