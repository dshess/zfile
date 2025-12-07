package zfile

import (
	"io"
	"os"
	"path/filepath"
	"testing"

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
}
