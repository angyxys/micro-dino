package dino

type Meta map[string]any

func (m Meta) Get(key string) Value {
	if m == nil {
		return Value{}
	}
	return Value{raw: m[key]}
}

func (m Meta) Set(key string, v any) Meta {
	if m == nil {
		m = Meta{}
	}
	m[key] = v
	return m
}
