package internal

import (
	"io"
	"os"
)

var (
	OpenReader = func(name string) (io.ReadCloser, error) {
		return os.Open(name)
	}
	CreateWriter = func(name string) (io.WriteCloser, error) {
		return os.Create(name)
	}
)
