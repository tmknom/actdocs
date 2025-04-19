package cli

import (
	"io"
	"os"
)

func ReadSource(filename string) (raw []byte, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) { err = file.Close() }(file)

	return io.ReadAll(file)
}
