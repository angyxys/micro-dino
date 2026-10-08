package dino

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Value struct {
	raw any
}

func (v Value) Any() any { return v.raw }

func (v Value) IsNil() bool { return v.raw == nil }

func (v Value) String() string {
	switch x := v.raw.(type) {
	case nil:
		return ""
	case string:
		return x
	case fmt.Stringer:
		return x.String()
	case []byte:
		return string(x)
	default:
		return fmt.Sprintf("%v", x)
	}
}

func (v Value) Int() int {
	n, _ := v.Int64()
	return int(n)
}

func (v Value) Int64() (int64, error) {
	switch x := v.raw.(type) {
	case nil:
		return 0, nil
	case int:
		return int64(x), nil
	case int64:
		return x, nil
	case float64:
		return int64(x), nil
	case string:
		return strconv.ParseInt(x, 10, 64)
	case json.Number:
		return x.Int64()
	default:
		return 0, fmt.Errorf("micro-dino: cannot convert %T to int64", x)
	}
}

func (v Value) Float64() (float64, error) {
	switch x := v.raw.(type) {
	case nil:
		return 0, nil
	case float64:
		return x, nil
	case int:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case string:
		return strconv.ParseFloat(x, 64)
	default:
		return 0, fmt.Errorf("micro-dino: cannot convert %T to float64", x)
	}
}

func (v Value) Bool() bool {
	switch x := v.raw.(type) {
	case bool:
		return x
	case string:
		b, _ := strconv.ParseBool(x)
		return b
	case int, int64, float64:
		return v.Int() != 0
	}
	return false
}

func (v Value) Bind(target any) error {
	b, err := json.Marshal(v.raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}
