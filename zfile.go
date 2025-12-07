package zfile

import (
	"io"
	"os"
)

func Open(path string) (io.ReadCloser, error) {
	inner, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return inner, nil
}

func Create(path string) (io.WriteCloser, error) {
	inner, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return inner, nil
}
