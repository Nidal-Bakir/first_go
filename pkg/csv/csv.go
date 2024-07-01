package csv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func Marshal(data interface{}) ([][]string, error) {
	dataVal := reflect.ValueOf(data)
	if dataVal.Kind() == reflect.Ptr {
		dataVal = dataVal.Elem()
	}
	if dataVal.Kind() != reflect.Slice {
		return nil, errors.New("the data should be of kind Slice")
	}
	dataElemType := dataVal.Type().Elem()
	if dataElemType.Kind() != reflect.Struct {
		return nil, errors.New("the slice should be of type struct")
	}

	out := make([][]string, dataVal.Len()+1)

	headers, err := extractHeadersFromStructTags(dataElemType)
	if err != nil {
		return nil, err
	}
	out[0] = headers

	for i := 0; i < dataVal.Len(); i++ {
		row := make([]string, len(headers))
		err = marshalOne(dataVal.Index(i), row)
		if err != nil {
			return nil, err
		}
		out[i+1] = row
	}

	return out, nil
}

func extractHeadersFromStructTags(rv reflect.Type) ([]string, error) {
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
		case reflect.Float32, reflect.Float64:
			var bitSize int
			if fieldVal.Kind() == reflect.Float64 {
				bitSize = 64
			} else {
				bitSize = 32
			}
			strValue = strconv.FormatFloat(fieldVal.Float(), 'f', 10, bitSize)
		case reflect.Bool:
			strValue = strconv.FormatBool(fieldVal.Bool())
		default:
			return fmt.Errorf("unsupported field type %v", fieldVal.Kind())
		}
		out[i] = strValue
	}
	return nil
}

func Unmarshal(data [][]string, v interface{}) error {

	sliceValPtr := reflect.ValueOf(v)
	if sliceValPtr.Kind() != reflect.Ptr {
		return fmt.Errorf("the [v interface] value should be of kind pointer to slice of struct, but %T was given", v)
	}
	sliceVal := sliceValPtr.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("the [v interface] value should be of kind pointer to slice of struct, but %T was given", v)
	}

	sliceStructType := sliceVal.Type().Elem()
	if sliceStructType.Kind() != reflect.Struct {
		return fmt.Errorf("the [v interface] value should be of kind pointer to slice of struct, but %T was given", v)
	}

	if len(data) == 0 {
		sliceVal.Set(reflect.ValueOf([]interface{}{}))
		return nil
	}

	headersIndexes := mapHeadersToIndexes(data[0])
	for _, row := range data[1:] {
		val, err := unmarshalOne(row, headersIndexes, sliceStructType)
		if err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, val.Elem()))
	}

	return nil
}

func mapHeadersToIndexes(headers []string) map[string]int {
	headersMapFelids := make(map[string]int, len(headers))
	for index, header := range headers {
		headersMapFelids[header] = index
	}
	return headersMapFelids
}

func unmarshalOne(row []string, headersIndexes map[string]int, t reflect.Type) (reflect.Value, error) {
	ptr := reflect.New(t)
	ptrVal := ptr.Elem()
	for i := range t.NumField() {
		header, ok := t.Field(i).Tag.Lookup("csv")
		if !ok {
			continue
		}
		indexInRow, ok := headersIndexes[header]
		if !ok {
			continue
		}

		field := ptrVal.Field(i)
		rowVal := row[indexInRow]
		switch field.Kind() {

		case reflect.Int:
			v, err := strconv.Atoi(rowVal)
			if err != nil {
				return reflect.ValueOf(nil), err
			}
			field.SetInt(int64(v))

		case reflect.Float64, reflect.Float32:
			var bitSize int
			if field.Kind() == reflect.Float64 {
				bitSize = 64
			} else {
				bitSize = 32
			}
			v, err := strconv.ParseFloat(rowVal, bitSize)
			if err != nil {
				return reflect.ValueOf(nil), err
			}
			field.SetFloat(v)

		case reflect.Bool:
			v, err := strconv.ParseBool(rowVal)
			if err != nil {
				return reflect.ValueOf(nil), err
			}
			field.SetBool(v)

		case reflect.String:
			field.SetString(rowVal)
		}
	}
	return ptr, nil
}
