package data

import (
	"fmt"
)

type Any interface{}

type Number *Decimal

func From(val interface{}) Any {
	if val == nil {
		return nil
	}
	switch v := val.(type) {
	case string:
		return v
	case int:
		return Number(DecimalFromInt64(int64(v)))
	case int64:
		return Number(DecimalFromInt64(v))
	case float64:
		return Number(DecimalFromFloat64(v))
	case *Decimal:
		return Number(v)
	case bool:
		return v
	case *Object:
		return v
	case *Array:
		return v
	default:
		fmt.Printf("whoops, value type not implemented: %v\n", v)
		panic("here")
	}
}
