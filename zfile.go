package zfile

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func Open(path string) (io.ReadCloser, error) {
	inner, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if filepath.Ext(path) == ".gz" {
		outer, err := gzip.NewReader(inner)
		if err != nil {
			inner.Close()
			return nil, err
		}
		return &wrappedReadCloser{
			wrappedCloser: inner,
			readCloser:    outer,
		}, nil
	}

	return inner, nil
}

func Create(path string) (io.WriteCloser, error) {
	inner, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	if filepath.Ext(path) == ".gz" {
		outer, err := gzip.NewWriterLevel(inner, gzip.BestCompression)
		if err != nil {
			inner.Close()
			return nil, err
		}
		return &wrappedWriteCloser{
			wrappedCloser: inner,
			writeCloser:   outer,
		}, nil
	}
	return inner, nil
}

type wrappedReadCloser struct {
	wrappedCloser io.Closer
	readCloser    io.ReadCloser
}

func (rc *wrappedReadCloser) Read(p []byte) (n int, err error) {
	return rc.readCloser.Read(p)
}

func (rc *wrappedReadCloser) Close() error {
	outerErr := rc.readCloser.Close()
	innerErr := rc.wrappedCloser.Close()
	if outerErr != nil {
		return outerErr
	}
	return innerErr
}

type wrappedWriteCloser struct {
	wrappedCloser io.Closer
	writeCloser   io.WriteCloser
}

func (wc *wrappedWriteCloser) Write(p []byte) (n int, err error) {
	return wc.writeCloser.Write(p)
}

func (wc *wrappedWriteCloser) Close() error {
	outerErr := wc.writeCloser.Close()
	innerErr := wc.wrappedCloser.Close()
	if outerErr != nil {
		return outerErr
	}
	return innerErr
}
