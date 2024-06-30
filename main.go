package main

import (
	"reflect"
)

func main() {

}

func IsValidInterface(i interface{}) bool {
	iv := reflect.ValueOf(i)
	return iv.IsValid() && !iv.IsNil()
}
