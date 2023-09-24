package data

import (
// "encoding/json"
// "fmt"
)

type Array struct {
	Elements []interface{}
}

func ArrayFromSlice(s []interface{}) *Array {
	return &Array{Elements: s}
}
