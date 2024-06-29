package csv

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_extractHeader(t *testing.T) {
	type Book struct {
		Author    string  `csv:"author"`
		ReadCount int     `csv:"read_count"`
		Price     float64 `csv:"price"`
	}

	headers, err := extractHeader(reflect.TypeOf(Book{}))
	require.Equal(t, []string{"author", "read_count", "price"}, headers)
	require.Nil(t, err)

	headers, err = extractHeader(reflect.TypeOf(""))
	require.Nil(t, headers)
	require.NotNil(t, err)
}

func Test_marshalOne(t *testing.T) {
	type Book struct {
		Author    string  `csv:"author"`
		ReadCount int     `csv:"read_count"`
		Price     float64 `csv:"price"`
	}

	out := make([]string, 3)
	err := marshalOne(reflect.ValueOf(Book{
		Author:    "nidal",
		ReadCount: 3,
		Price:     55.222,
	}), out)

	require.Nil(t, err)
	require.Equal(t, []string{"nidal", "3", "55.2220000000"}, out)
}

func TestMarshal(t *testing.T) {
	type Book struct {
		Author    string  `csv:"author"`
		ReadCount int     `csv:"read_count"`
		Price     float64 `csv:"price"`
	}

	arr := []Book{
		{
			Author:    "nidal",
			ReadCount: 1,
			Price:     1.111,
		},
		{
			Author:    "bakir",
			ReadCount: 2,
			Price:     22.222,
		},
	}

	result, err := Marshal(arr)

	ex := [][]string{
		{"author", "read_count", "price"},
		{"nidal", "1", "1.1110000000"},
		{"bakir", "2", "22.2220000000"},
	}

	require.Nil(t, err)
	require.Equal(t, ex, result)
}
