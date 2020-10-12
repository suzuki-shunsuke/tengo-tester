package file

import (
	"io"
	"os"
)

type Reader struct{}

func (reader Reader) Open(p string) (io.ReadCloser, error) {
	f, err := os.Open(p)
	return f, err
}
