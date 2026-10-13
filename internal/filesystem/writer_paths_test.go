package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDestinationEmptyNotADirectory(t *testing.T) {
	file := filepath.Join(t.TempDir(), "file.txt")
	require.NoError(t, os.WriteFile(file, []byte("x"), 0o644))
	_, err := DestinationEmpty(file)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a directory")
}

func TestListDestinationFilesMissing(t *testing.T) {
	_, err := ListDestinationFiles(filepath.Join(t.TempDir(), "missing"))
	assert.Error(t, err)
}

func TestWriteFilesNestedDirs(t *testing.T) {
	dest := t.TempDir()
	w := NewWriter(dest, false)
	files := []File{
		{TargetPath: "a/b/c/deep.txt", Content: []byte("deep")},
	}
	result, err := w.WriteFiles(files)
	require.NoError(t, err)
	assert.Equal(t, 1, result.FileCount)

	data, err := os.ReadFile(filepath.Join(dest, "a/b/c/deep.txt"))
	require.NoError(t, err)
	assert.Equal(t, "deep", string(data))
}

func TestWriteFileNonExecutableMode(t *testing.T) {
	dest := t.TempDir()
	w := NewWriter(dest, false)
	files := []File{
		{TargetPath: "data.txt", Content: []byte("x")},
	}
	_, err := w.WriteFiles(files)
	require.NoError(t, err)

	info, err := os.Stat(filepath.Join(dest, "data.txt"))
	require.NoError(t, err)
	assert.Zero(t, info.Mode()&0o111, "regular file must not be executable")
}
