/*
Copyright 2021 Lee R. Boynton

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Map is a generic map of string key to parameterized value type that preserves key order.

package data

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Map[V any] struct {
	keys     []string
	bindings map[string]V
}

func NewMap[V any]() *Map[V] {
	return &Map[V]{
		bindings: make(map[string]V, 0),
	}
}

func (s *Map[V]) UnmarshalJSON(data []byte) error {
	keys, err := JsonKeysInOrder(data)
	if err != nil {
		return err
	}
	str := NewMap[V]()
	str.keys = keys
	err = json.Unmarshal(data, &str.bindings)
	if err != nil {
		return err
	}
	*s = *str
	return nil
}

func (s Map[V]) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	for i, key := range s.keys {
		value := s.bindings[key]
		if i > 0 {
			buffer.WriteString(",")
		}
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("%q:%s", key, string(jsonValue)))
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

func (s *Map[V]) String() string {
	return JsonEncode(s)
}

func (s *Map[V]) find(key string) int {
	for i, k := range s.keys {
		if k == key {
			return i
		}
	}
	return -1
}

func (s *Map[V]) Has(key string) bool {
	if s != nil {
		if _, ok := s.bindings[key]; ok {
			return true
		}
	}
	return false
}

func (s *Map[V]) Get(key string) V {
	return s.bindings[key]
}

func (s *Map[V]) Put(key string, val V) {
	if s == nil {
		*s = *NewMap[V]()
	}
	if _, ok := s.bindings[key]; !ok {
		s.keys = append(s.keys, key)
	}
	s.bindings[key] = val
}

func (s *Map[V]) Keys() []string {
	if s == nil {
		return nil
	}
	return s.keys
}

func (s *Map[V]) Length() int {
	if s == nil || s.keys == nil {
		return 0
	}
	return len(s.keys)
}
