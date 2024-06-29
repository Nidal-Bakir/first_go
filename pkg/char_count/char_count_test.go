package charcount_test

import (
	"fmt"
	"os"

	"testing"

	charcount "github.com/Nidal-Bakir/first_go/pkg/char_count"
	"github.com/stretchr/testify/require"
)

const (
	charCountExpected int = 2000
)

func TestCountInFile(t *testing.T) {

	buffSize := 5

	file, err := os.Open("testdata/data.txt")
	if err != nil {
		t.Fatalf("error while creating the test file: %v", err)
	}

	actualCount, err := charcount.CountInFile(file.Name(), buffSize)
	if err != nil {
		t.Fatalf("error while counting chars in the file: %v", err)
	}

	require.Equal(t, charCountExpected, actualCount)
}

func BenchmarkCharCount(b *testing.B) {
	file, err := os.Open("testdata/data.txt")
	if err != nil {
		b.Fatalf("error while creating the test file: %v", err)
	}
	defer file.Close()

	for _, v := range []int{100, 200, 300, 500, 1000, 1500, 2000,5000,10000,50000,100000} {
		b.Run(fmt.Sprintf("CountInFile buff=%d", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := charcount.CountInFile(file.Name(), v)
				if err != nil {
					b.Fatalf("error while counting chars in the file: %v", err)
				}
			}
		})
	}
}
 