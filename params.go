package dino

type Params map[string]any

func (p Params) Get(key string) Value {
	if p == nil {
		return Value{}
	}
	return Value{raw: p[key]}
}

func (p Params) Has(key string) bool {
	if p == nil {
		return false
	}
	_, ok := p[key]
	return ok
}

func (p Params) Set(key string, value any) Params {
	if p == nil {
		p = Params{}
	}
	p[key] = value
	return p
}

func (p Params) Merge(other Params) Params {
	out := Params{}
	for k, v := range p {
		out[k] = v
	}
	for k, v := range other {
		out[k] = v
	}
	return out
}
