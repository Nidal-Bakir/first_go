package csv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// func Unmarshal(data [][]string, i interface{}) error {

// }

func Marshal(data interface{}) ([][]string, error) {
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Slice {
		return nil, errors.New("the data should be of kind Slice")
	}
	dataElemType := dataType.Elem()
	if dataElemType.Kind() != reflect.Struct {
		return nil, errors.New("the slice should be of type struct")
	}

	dataVal := reflect.ValueOf(data)

	out := make([][]string, dataVal.Len()+1)

	headers, err := extractHeader(dataElemType)
	if err != nil {
		return nil, err
	}

	out[0] = headers

	for i := 0; i < dataVal.Len(); i++ {
		row := make([]string, len(headers))
		err = marshalOne(dataVal.Index(i), row)
		if err != nil {
			return nil, err
		} else {
			out[i+1] = row
		}
	}

	return out, nil
}

func extractHeader(rv reflect.Type) ([]string, error) {
	if rv.Kind() != reflect.Struct {
		return nil, errors.New("the Type value should be of kind struct")
	}

	headers := make([]string, 0)
	for i := 0; i < rv.NumField(); i++ {
		if key, ok := rv.Field(i).Tag.Lookup("csv"); ok {
			headers = append(headers, key)
		}
	}
	return headers, nil
}

func marshalOne(dataVal reflect.Value, out []string) error {
	dataType := dataVal.Type()

	if dataType.Kind() != reflect.Struct {
		return errors.New("the Type value should be of kind struct")
	}

	for i := 0; i < dataType.NumField(); i++ {
		if _, ok := dataType.Field(i).Tag.Lookup("csv"); !ok {
			continue
		}

		fieldVal := dataVal.Field(i)
		var strValue string
		switch fieldVal.Kind() {
		case reflect.String:
			strValue = fieldVal.String()
		case reflect.Int:
			strValue = strconv.FormatInt(fieldVal.Int(), 10)
		case reflect.Float64:
			strValue = strconv.FormatFloat(fieldVal.Float(), 'f', 10, 64)
		case reflect.Float32:
			strValue = strconv.FormatFloat(fieldVal.Float(), 'f', 10, 32)
		case reflect.Bool:
			strValue = strconv.FormatBool(fieldVal.Bool())
		default:
			return fmt.Errorf("unsupported field type %v", fieldVal.Kind())
		}
		out[i] = strValue
	}
	return nil
}
