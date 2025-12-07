package zfile

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
)

type Compressor int

const (
	CSuffix Compressor = iota
	CGzip
	CZstd
	CXz
	CNone
	CError
)

func CDerive(path string) Compressor {
	switch filepath.Ext(path) {
	case ".gz":
		return CGzip
	case ".zst":
		return CZstd
	case ".xz":
		return CXz
	default:
		return CNone
	}
}

func cDecoder(inner io.ReadCloser, ctype Compressor) (io.ReadCloser, error) {
	switch ctype {
	case CGzip:
		return gzip.NewReader(inner)
	case CZstd:
		decoder, err := zstd.NewReader(inner)
		if err != nil {
			return nil, err
		}
		return decoder.IOReadCloser(), nil
	case CXz:
		decoder, err := xz.NewReader(inner)
		if err != nil {
			return nil, err
		}
		return io.NopCloser(decoder), nil
	}
	return nil, nil
}

func OpenType(path string, ctype Compressor) (io.ReadCloser, error) {
	inner, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// NOTE: Take care that inner does not escape.

	outer, err := cDecoder(inner, ctype)
	if err != nil {
		inner.Close()
		return nil, err
	}

	if outer != nil {
		return &wrappedReadCloser{
			wrappedCloser: inner,
			readCloser:    outer,
		}, nil
	}

	return inner, nil
}

func Open(path string) (io.ReadCloser, error) {
	return OpenType(path, CDerive(path))
}

func cEncoder(inner io.WriteCloser, ctype Compressor) (io.WriteCloser, error) {
	switch ctype {
	case CGzip:
		return gzip.NewWriterLevel(inner, gzip.BestCompression)
	case CZstd:
		return zstd.NewWriter(
			inner,
			zstd.WithEncoderLevel(zstd.SpeedBestCompression),
		)
	case CXz:
		return xz.NewWriter(inner)
	}
	return nil, nil
}

func CreateType(path string, ctype Compressor) (io.WriteCloser, error) {
	inner, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	// NOTE: Take care that inner does not escape.

	outer, err := cEncoder(inner, ctype)
	if err != nil {
		inner.Close()
		return nil, err
	}

	if outer != nil {
		return &wrappedWriteCloser{
			wrappedCloser: inner,
			writeCloser:   outer,
		}, nil
	}

	return inner, nil
}

func Create(path string) (io.WriteCloser, error) {
	return CreateType(path, CDerive(path))
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
