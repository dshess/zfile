package zfile

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZfileRead(t *testing.T) {
	tmpDir := t.TempDir()

	const testdata = "This is a test"

	t.Run("plain", func(t *testing.T) {
		path := filepath.Join(tmpDir, "file")

		err := os.WriteFile(path, []byte(testdata), 0666)
		require.Nil(t, err)

		in, err := Open(path)
		require.Nil(t, err)
		defer in.Close()

		data, err := io.ReadAll(in)
		require.Nil(t, err)

		assert.Equal(t, testdata, string(data))
	})

	t.Run("gzip", func(t *testing.T) {
		path := filepath.Join(tmpDir, "file.gz")

		out, err := os.Create(path)
		require.Nil(t, err)
		require.NotNil(t, out)

		writer := gzip.NewWriter(out)
		require.NotNil(t, writer)

		_, err = writer.Write([]byte(testdata))
		require.Nil(t, err)
		require.Nil(t, writer.Close())
		require.Nil(t, out.Close())

		in, err := Open(path)
		require.Nil(t, err)
		defer in.Close()

		data, err := io.ReadAll(in)
		require.Nil(t, err)

		assert.Equal(t, testdata, string(data))
	})

	t.Run("zstd", func(t *testing.T) {
		path := filepath.Join(tmpDir, "file.zst")

		out, err := os.Create(path)
		require.Nil(t, err)
		require.NotNil(t, out)

		writer, err := zstd.NewWriter(out)
		require.Nil(t, err)
		require.NotNil(t, writer)

		_, err = writer.Write([]byte(testdata))
		require.Nil(t, err)
		require.Nil(t, writer.Close())
		require.Nil(t, out.Close())

		in, err := Open(path)
		require.Nil(t, err)
		defer in.Close()

		data, err := io.ReadAll(in)
		require.Nil(t, err)

		assert.Equal(t, testdata, string(data))
	})

	t.Run("xz", func(t *testing.T) {
		path := filepath.Join(tmpDir, "file.xz")

		out, err := os.Create(path)
		require.Nil(t, err)
		require.NotNil(t, out)

		writer, err := xz.NewWriter(out)
		require.Nil(t, err)
		require.NotNil(t, writer)

		_, err = writer.Write([]byte(testdata))
		require.Nil(t, err)
		require.Nil(t, writer.Close())
		require.Nil(t, out.Close())

		in, err := Open(path)
		require.Nil(t, err)
		defer in.Close()

		data, err := io.ReadAll(in)
		require.Nil(t, err)

		assert.Equal(t, testdata, string(data))
	})
}

func TestZfileWrite(t *testing.T) {
	tmpDir := t.TempDir()

	const testdata = "This is a test"

	t.Run("plain", func(t *testing.T) {
		path := filepath.Join(tmpDir, "file")

		out, err := Create(path)
		require.Nil(t, err)
		defer out.Close()

		_, err = io.WriteString(out, testdata)
		require.Nil(t, err)

		err = out.Close()
		require.Nil(t, err)

		data, err := os.ReadFile(path)
		require.Nil(t, err)

		assert.Equal(t, testdata, string(data))
	})

	t.Run("gzip", func(t *testing.T) {
		path := filepath.Join(tmpDir, "file.gz")

		out, err := Create(path)
		require.Nil(t, err)
		defer out.Close()

		_, err = io.WriteString(out, testdata)
		require.Nil(t, err)

		err = out.Close()
		require.Nil(t, err)

		in, err := os.Open(path)
		require.Nil(t, err)
		defer in.Close()

		decoder, err := gzip.NewReader(in)
		require.Nil(t, err)

		data, err := io.ReadAll(decoder)
		require.Nil(t, err)

		require.Nil(t, decoder.Close())
		require.Nil(t, in.Close())

		assert.Equal(t, testdata, string(data))
	})

	t.Run("zstd", func(t *testing.T) {
		path := filepath.Join(tmpDir, "file.zst")

		out, err := Create(path)
		require.Nil(t, err)
		defer out.Close()

		_, err = io.WriteString(out, testdata)
		require.Nil(t, err)

		err = out.Close()
		require.Nil(t, err)

		in, err := os.Open(path)
		require.Nil(t, err)
		defer in.Close()

		decoder, err := zstd.NewReader(in)
		require.Nil(t, err)

		data, err := io.ReadAll(decoder)
		require.Nil(t, err)

		decoder.Close()
		require.Nil(t, in.Close())

		assert.Equal(t, testdata, string(data))
	})

	t.Run("xz", func(t *testing.T) {
		path := filepath.Join(tmpDir, "file.xz")

		out, err := Create(path)
		require.Nil(t, err)
		defer out.Close()

		_, err = io.WriteString(out, testdata)
		require.Nil(t, err)

		err = out.Close()
		require.Nil(t, err)

		in, err := os.Open(path)
		require.Nil(t, err)
		defer in.Close()

		decoder, err := xz.NewReader(in)
		require.Nil(t, err)

		data, err := io.ReadAll(decoder)
		require.Nil(t, err)

		require.Nil(t, in.Close())

		assert.Equal(t, testdata, string(data))
	})
}
