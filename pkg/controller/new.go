package controller

import (
	"io"

	"github.com/suzuki-shunsuke/tengo-tester/pkg/config"
)

type Controller struct {
	Config     config.Config
	FileReader FileReader
}

type FileReader interface {
	Open(path string) (io.ReadCloser, error)
}
