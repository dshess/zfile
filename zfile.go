package zfile

import (
	"io"
	"os"
)

func Create(path string) (io.WriteCloser, error) {
	inner, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return inner, nil
}
