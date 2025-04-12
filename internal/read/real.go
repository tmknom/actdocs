package read

import (
	"io"
	"os"
)

type YamlReader struct {
	Filename string
}

func (r *YamlReader) Read() (raw []byte, err error) {
	file, err := os.Open(r.Filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) { err = file.Close() }(file)

	return io.ReadAll(file)
}
