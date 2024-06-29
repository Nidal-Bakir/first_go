package charcount

import (
	"errors"
	"io"
	"os"
)

func CountInFile(path string, buffSize int) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	count := 0
	buff := make([]byte, buffSize)

	for {
		n, err := file.Read(buff)
		count += n
		if err != nil {
			if errors.Is(err, io.EOF) {
				return count, nil
			}
			return 0, err
		}

	}
}
