package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteFilesMkdirFailure(t *testing.T) {
	dest := t.TempDir()
	// A file exists where a directory component is needed,
	// so MkdirAll fails.
	require.NoError(t, os.WriteFile(filepath.Join(dest, "blocker"), []byte("x"), 0o644))

	w := NewWriter(dest, false)
	files := []File{{TargetPath: "blocker/child.txt", Content: []byte("y")}}
	_, err := w.WriteFiles(files)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create dir")
}

func TestWriteFilesCreateTempFailure(t *testing.T) {
	dest := t.TempDir()
	// A read-only directory makes CreateTemp fail for
	// non-root users.
	sub := filepath.Join(dest, "ro")
	require.NoError(t, os.MkdirAll(sub, 0o755))
	require.NoError(t, os.Chmod(sub, 0o500))
	t.Cleanup(func() { _ = os.Chmod(sub, 0o755) })

	w := NewWriter(dest, false)
	files := []File{{TargetPath: "ro/x.txt", Content: []byte("y")}}
	_, err := w.WriteFiles(files)
	if err == nil {
		t.Skip("running as root; read-only directory does not block writes")
	}
	require.Error(t, err)
}

func TestWriteFilesRenameFailure(t *testing.T) {
	dest := t.TempDir()
	// A directory exists at the target path, so
	// renaming a file onto it fails.
	require.NoError(t, os.MkdirAll(filepath.Join(dest, "target"), 0o755))

	w := NewWriter(dest, false)
	files := []File{{TargetPath: "target", Content: []byte("y")}}
	_, err := w.WriteFiles(files)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rename")
}
