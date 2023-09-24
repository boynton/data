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

// Timestamp is a time.Time equivalent that marshals to/from JSON as a string.

package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Timestamp struct {
	time.Time
}

const RFC3339Milli = "%d-%02d-%02dT%02d:%02d:%02d.%03dZ"

func (ts Timestamp) String() string {
	if ts.IsZero() {
		return ""
	}
	return fmt.Sprintf(RFC3339Milli, ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), ts.Nanosecond()/1000000)
}

func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ts.String() + "\""), nil
}

func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err == nil {
		var tsp Timestamp
		tsp, err = ParseTimestamp(string(j))
		if err == nil {
			*ts = tsp
		}
	}
	return err
}

func ParseTimestamp(s string) (Timestamp, error) {
	layout := "2006-01-02T15:04:05.999Z" //derive this from the spec used for output?
	t, e := time.Parse(layout, s)
	if e != nil {
		if strings.HasSuffix(s, "+00:00") || strings.HasSuffix(s, "-00:00") {
			t, e = time.Parse(layout, s[:len(s)-6]+"Z")
		} else if strings.HasSuffix(s, "+0000") || strings.HasSuffix(s, "-0000") {
			t, e = time.Parse(layout, s[:len(s)-5]+"Z")
		}
		if e != nil {
			var ts Timestamp
			return ts, fmt.Errorf("Bad Timestamp: %q", s)
		}
	}
	return Timestamp{t}, nil
}

func (ts *Timestamp) ToRfc2616String() string {
	if ts.IsZero() {
		return ""
	}
	return ts.UTC().Format(http.TimeFormat)
}
