package data

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Object - a map-like structure that maintains the order of its keys, and provides
// some utilities for converting from generic types.
// This represents a JSON object, and marshals to/from JSON while maintaining key order.
// It also handles arbitrary precision decimal numbers while doing so.
type Object struct {
	bindings []*Binding
	index    map[string]interface{}
}

type Binding struct {
	Key   string
	Value interface{}
}

// when false, use the Decimal type to represent numbers without precision loss. Set to true to use float64 for numbers.
var UseFloatNumbers = false

// Create a new empty object
func NewObject() *Object {
	return &Object{
		bindings: nil,
		index:    make(map[string]interface{}, 0),
	}
}

var end = fmt.Errorf("invalid end of array or object")

func parseValue(d *json.Decoder) (interface{}, error) {
	t, err := d.Token()
	if err != nil {
		return nil, err
	}
	switch t {
	case json.Delim('['):
		var ary []interface{}
		for {
			if item, err := parseValue(d); err != nil {
				if err == end {
					break
				}
				return nil, err
			} else {
				ary = append(ary, item)
			}
		}
		return ary, nil
	case json.Delim('{'):
		obj := NewObject()
		for {
			t, err = d.Token()
			if err != nil {
				return nil, err
			}
			if t == json.Delim('}') {
				break
			}
			key := t.(string)
			var val interface{}
			if val, err = parseValue(d); err != nil {
				if err == end {
					break
				}
				return nil, err
			}
			obj.Put(key, val)
		}
		return obj, nil
	case json.Delim(']'), json.Delim('}'):
		return nil, end
	default:
		switch val := t.(type) {
		case string:
			return val, nil
		case bool:
			return val, nil
		case json.Number: //only when UseNumber on the decoder is used
			return DecimalFromString(val.String())
		case float64: //only when UseNumber on the decoder is *not* used
			return val, nil
		default:
			if val != nil {
				return nil, fmt.Errorf("Malformed JSON at %v\n", val)
			}
		}
	}
	return nil, nil
}

func (s *Object) UnmarshalJSON(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	if !UseFloatNumbers {
		d.UseNumber()
	}
	val, err := parseValue(d)
	if err != nil {
		return err
	}
	if ss, ok := val.(*Object); ok {
		*s = *ss
		return nil
	}
	return fmt.Errorf("JSON Unmarshal: not an object when one was expected")
}

// note: the map has already lost key order, this is just for convenience/debugging
func ObjectFromMap(m map[string]interface{}) *Object {
	o := NewObject()
	for k, v := range m {
		o.Put(k, v)
	}
	return o
}

func (s *Object) String() string {
	return Pretty(s)
}

func (s Object) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	for i, b := range s.bindings {
		if i > 0 {
			buffer.WriteString(",")
		}
		jsonValue, err := json.Marshal(b.Value)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("%q:%s", b.Key, string(jsonValue)))
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

func (s *Object) Bindings() []*Binding {
	if s == nil {
		return nil
	}
	return s.bindings
}

func (s *Object) Keys() []string {
	var result []string
	for _, b := range s.bindings {
		result = append(result, b.Key)
	}
	return result
}

func (s *Object) Has(key string) bool {
	if s != nil {
		if _, ok := s.index[key]; ok {
			return true
		}
	}
	return false
}

func (s *Object) Put(key string, val interface{}) {
	if s == nil {
		*s = *NewObject()
	}
	if _, ok := s.index[key]; !ok {
		s.bindings = append(s.bindings, &Binding{Key: key, Value: val})
	}
	s.index[key] = val
}

func (s *Object) Get(key string) interface{} {
	if s == nil {
		return nil
	}
	return s.index[key]
}

func (s *Object) GetString(key string) string {
	return AsString(s.Get(key))
}

func (s *Object) GetBool(key string) bool {
	return AsBool(s.Get(key))
}

func (s *Object) GetInt(key string) int {
	return AsInt(s.Get(key))
}

func (s *Object) GetInt64(key string) int64 {
	return AsInt64(s.Get(key))
}

func (s *Object) GetDecimal(key string) *Decimal {
	return AsDecimal(s.Get(key))
}

func (s *Object) GetObject(key string) *Object {
	return AsObject(s.Get(key))
}

func (s *Object) GetSlice(key string) []interface{} {
	return AsSlice(s.Get(key))
}

func (s *Object) GetStringSlice(key string) []string {
	return AsStringSlice(s.Get(key))
}

func (s *Object) GetMap(key string) map[string]interface{} {
	return AsMap(s.Get(key))
}
