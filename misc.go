package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func Pretty(obj interface{}) string {
	indentSize := "  "
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", indentSize)
	if err := enc.Encode(&obj); err != nil {
		return fmt.Sprint(obj)
	}
	s := buf.String()
	return string(s)
}

func JsonEncode(obj interface{}) string {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(&obj); err != nil {
		return fmt.Sprint(obj)
	}
	return strings.TrimRight(string(buf.String()), " \t\n\v\f\r")
}

func JsonDecode(j string) *Object {
	var v Object
	err := json.Unmarshal(([]byte)(j), &v)
	if err != nil {
		return nil
	}
	return &v
}

func JsonDecodeAs[T any](j string, target *T) error {
	return json.Unmarshal([]byte(j), target)
}

func AsMap(v interface{}) map[string]interface{} {
	if v != nil {
		switch m := v.(type) {
		case map[string]interface{}:
			return m
			//		case *Struct:
			//			return m.value.bindings
		default:
			return nil
		}
	}
	return nil
}

func AsArray(o interface{}) *Array {
	switch v := o.(type) {
	case *Array:
		return v
	case []interface{}:
		return ArrayFromSlice(v)
	}
	return nil
}

func AsObject(v interface{}) *Object {
	if v != nil {
		switch o := v.(type) {
		case *Object:
			return o
		case map[string]interface{}:
			return ObjectFromMap(o)
		}
	}
	return nil
}

func AsSlice(v interface{}) []interface{} {
	if v != nil {
		switch a := v.(type) {
		case []interface{}:
			return a
		case *Array:
			return a.Elements
		default:
			return nil
		}
	}
	return nil
}

func AsStringSlice(v interface{}) []string {
	var sa []string
	a := AsSlice(v)
	if a != nil {
		for _, i := range a {
			switch s := i.(type) {
			case *Array:
				var lst []string
				for _, v := range s.Elements {
					lst = append(lst, AsString(v))
				}
				return lst
			case *string:
				sa = append(sa, *s)
			case string:
				sa = append(sa, s)
			default:
				return nil
			}
		}
	}
	return sa
}

func TypeOf(v interface{}) string {
	return fmt.Sprint(reflect.TypeOf(v))
}

func AsString(v interface{}) string {
	if v != nil {
		switch s := v.(type) {
		case string:
			return s
		case *string:
			return *s
		default:
			fmt.Println("not a string or *string:", TypeOf(v))
			panic("what?!")
		}
	}
	return ""
}

func AsBool(v interface{}) bool {
	if v != nil {
		if b, isBool := v.(bool); isBool {
			return b
		}
		return true
	}
	return false
}

func AsInt(v interface{}) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int32:
		return int(n)
	case int64:
		return int(n)
	case int:
		return n
	case *Integer:
		return n.AsInt()
	case *Decimal:
		return n.AsInt()
	}
	return 0
}

func AsInt64(v interface{}) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case float64:
		return int64(n)
	case int32:
		return int64(n)
	case int:
		return int64(n)
	case *Integer:
		return n.AsInt64()
	case *Decimal:
		return n.AsInt64()
	}
	return 0
}

func AsFloat64(v interface{}) float64 {
	if n, ok := v.(float64); ok {
		return n
	}
	return 0
}

func AsDecimal(v interface{}) *Decimal {
	switch n := v.(type) {
	case Decimal:
		return &n
	case *Decimal:
		return n
	case float64:
		return DecimalFromFloat64(n)
	default:
		return nil
	}
}

func Get(m map[string]interface{}, key string) interface{} {
	if m != nil {
		if v, ok := m[key]; ok {
			return v
		}
	}
	return nil
}

func GetString(m map[string]interface{}, key string) string {
	return AsString(Get(m, key))
}
func GetStringSlice(m map[string]interface{}, key string) []string {
	return AsStringSlice(Get(m, key))
}
func GetBool(m map[string]interface{}, key string) bool {
	return AsBool(Get(m, key))
}
func GetInt(m map[string]interface{}, key string) int {
	return AsInt(Get(m, key))
}
func GetInt64(m map[string]interface{}, key string) int64 {
	return AsInt64(Get(m, key))
}
func GetSlice(m map[string]interface{}, key string) []interface{} {
	return AsSlice(Get(m, key))
}
func GetMap(m map[string]interface{}, key string) map[string]interface{} {
	return AsMap(Get(m, key))
}
func GetDecimal(m map[string]interface{}, key string) *Decimal {
	return AsDecimal(Get(m, key))
}

var NOT_FOUND_ERROR = fmt.Errorf("Not found")

// methods don't do generics
func ValueTo[T any](v interface{}, storage *T) error {
	//this could be improved. It is to leverage json's object mapping ;-)
	return JsonDecodeAs(JsonEncode(v), storage)
}

func ValueOf(v interface{}) *Object {
	val := JsonDecode(JsonEncode(v))
	return AsObject(val)
}
