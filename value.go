package tslogs

import (
	"strconv"
)

type Value struct {
	val interface{}
}

func (self *Value) Float64() (float64, error) {
	return strconv.ParseFloat(self.val.(string), 64)
}

func (self *Value) String() (string, error) {
	return self.val.(string), nil
}