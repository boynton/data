/*
Copyright 2021 Lee Boynton

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

	"github.com/shopspring/decimal"
)

type Decimal struct {
	Value decimal.Decimal
}

// Encode as a JSON number. The JSON spec allows for arbitrary precision, so this is the correct thing to do.
func (d Decimal) MarshalJSON() ([]byte, error) {
	repr := d.String()
	return []byte(repr), nil
}

func (d *Decimal) UnmarshalJSON(b []byte) error {
	stringRepr := string(b)
	num, err := DecimalFromString(stringRepr)
	if err == nil {
		*d = *num
		return nil
	}
	return fmt.Errorf("Bad number: %s", string(b))
}

func DecimalFromString(text string) (*Decimal, error) {
	num, err := decimal.NewFromString(text)
	if err != nil {
		return nil, fmt.Errorf("Bad number: %s", text)
	}
	return &Decimal{Value: num}, nil
}

func DecimalFromFloat64(val float64) *Decimal {
	return &Decimal{Value: decimal.NewFromFloat(val)}
}

func DecimalFromInt64(val int64) *Decimal {
	return &Decimal{Value: decimal.NewFromInt(val)}
}

func (d *Decimal) String() string {
	return d.Value.String()
}

func (d *Decimal) AsInt() int {
	n := d.Value.InexactFloat64()
	return int(n)
}

func (d *Decimal) AsInt32() int32 {
	n := d.Value.InexactFloat64()
	return int32(n)
}

func (d *Decimal) AsInt64() int64 {
	i := d.Value.InexactFloat64() //fixme. Maybe convert to big.Int first?
	return int64(i)
}

func (d *Decimal) AsFloat64() float64 {
	f := d.Value.InexactFloat64()
	return f
}

func (d *Decimal) AsBigFloat() *big.Float {
	return d.Value.BigFloat()
}

func DecimalValue(val *Decimal, defval interface{}) *Decimal {
	if val != nil {
		return val
	}
	if defval != nil {
		switch n := defval.(type) {
		case *Decimal:
			return n
		case int64:
			d, _ := DecimalFromString(fmt.Sprint(n))
			return d
		case int32:
			return DecimalFromInt64(int64(n))
		case int16:
			return DecimalFromInt64(int64(n))
		case int8:
			return DecimalFromInt64(int64(n))
		case float32:
			return DecimalFromFloat64(float64(n))
		case float64:
			return DecimalFromFloat64(n)
		}
	}
	return nil
}
