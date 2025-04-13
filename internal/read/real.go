package read

import (
	"io"
	"os"
)

type SourceReader struct{}

func (r *SourceReader) Read(filename string) (raw []byte, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) { err = file.Close() }(file)

	return io.ReadAll(file)
}
