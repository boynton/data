/*
Copyright 2023 Lee Boynton

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

// Decimal is a big.Float equivalent that marshals to/from JSON.

package data

import (
	"fmt"
	"math/big"
)

type Integer struct {
	big.Int
}

func (i Integer) MarshalJSON() ([]byte, error) {
	repr := i.String()
	return []byte(repr), nil
}

func (i *Integer) UnmarshalJSON(b []byte) error {
	text := string(b)
	//	var i2 Integer
	var _, ok = i.SetString(text, 0)
	if !ok {
		return fmt.Errorf("Bad integer: %s", text)
	}
	return nil
}

func (i Integer) String() string {
	return fmt.Sprintf("%v", i.Int)
}

func (i Integer) AsInt() int {
	return int(i.Int64())
}

func (i Integer) AsInt64() int64 {
	return i.Int64()
}
