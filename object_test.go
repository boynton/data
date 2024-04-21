package data

import (
	//	"encoding/json"
	"fmt"
	"testing"
)

const jsonStream = `{
  "country"     : "United States",
  "countryCode" : "US",
  "region"      : "CA",
  "regionName"  : "California",
  "city"        : "Mountain View",
  "zip"         : "94043",
  "lat"         : 37.4192,
  "lon"         : -122.0574,
  "timezone"    : "America/Los_Angeles",
  "isp"         : "Google Cloud",
  "org"         : "Google Cloud",
  "as"          : "AS15169 Google Inc.",
  "mobile"      : true,
  "proxy"       : false,
  "n"           : null,
  "pi"          : 3.1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679,
  "ary"         : [1, 2, "three"],
  "obj"         : {"x": 23, "y": [1, 2, 3]},
  "query"       : "35.192.xx.xxx"
}`

func TestStruct(test *testing.T) {
	s := NewObject()
	fmt.Println("s(1):", s)

	s.Put("x", 23)
	fmt.Println("s(2):", s)

	d := JsonDecode(jsonStream)
	fmt.Println("d:", d)
}
