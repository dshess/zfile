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

// Return a compression type based on extension.
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

// Like Open() but allows explicit compression type selection.
func OpenType(path string, ctype Compressor) (io.ReadCloser, error) {
	if ctype == CSuffix {
		ctype = CDerive(path)
	}

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

// Opens the named file for reading like os.Open() with automatic
// decompression based on suffix.  For instance, Open("file.gz") reads a
// gzip-compressed file.
func Open(path string) (io.ReadCloser, error) {
	return OpenType(path, CSuffix)
}

// Like ReadFile() but allows explicit compression type selection.
func ReadFileType(path string, ctype Compressor) ([]byte, error) {
	f, err := OpenType(path, ctype)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}

// Reads the contents of the file like os.ReadFile(), with automatic
// decompression based on suffix.  For instance, ReadFile("file.gz") reads
// a gzip-compressed file.
func ReadFile(path string) ([]byte, error) {
	return ReadFileType(path, CSuffix)
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

func innerCreateType(path string, perm os.FileMode, ctype Compressor) (io.WriteCloser, error) {
	if ctype == CSuffix {
		ctype = CDerive(path)
	}

	inner, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
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

// Like Create(), but allows passing an explicit compression type.  See
// WriteFileType().
func CreateType(path string, ctype Compressor) (io.WriteCloser, error) {
	return innerCreateType(path, 0666, ctype)
}

// Opens the named file for writing like os.Create() with automatic
// compression based on suffix.  For instance, Create("file.gz") writes a
// gzip-compressed file.
func Create(path string) (io.WriteCloser, error) {
	return CreateType(path, CSuffix)
}

// Like WriteFile(), but allows passing an explicit compression type.  This
// is useful when writing a temporary file which will later be renamed into
// place.  For instance:
//
//	tmpPath := path+".tmp"
//	err := zfile.WriteFileType(tmpPath, data, perm, zfile.CDerive(path))
//	if err != nil { return err }
//	return os.Rename(tmpPath, path)
func WriteFileType(path string, data []byte, perm os.FileMode, ctype Compressor) error {
	f, err := innerCreateType(path, perm, ctype)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return f.Close()
}

// Writes the contents of the file like os.WriteFile(), with automatic
// compression based on suffix.  For instance, WriteFile("file.gz") writes
// a gzip-compressed file.
func WriteFile(path string, data []byte, perm os.FileMode) error {
	return WriteFileType(path, data, perm, CSuffix)
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
