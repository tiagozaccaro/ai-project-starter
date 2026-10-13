package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafetyErrorMessage(t *testing.T) {
	err := &SafetyError{Code: "X", Message: "unsafe"}
	assert.Equal(t, "unsafe", err.Error())
}

func TestCheckDestinationNotADirectory(t *testing.T) {
	// A regular file as destination makes inspection fail.
	file := filepath.Join(t.TempDir(), "file.txt")
	require.NoError(t, os.WriteFile(file, []byte("x"), 0o644))

	err := CheckDestination(file, false)
	require.Error(t, err)
	var sErr *SafetyError
	require.ErrorAs(t, err, &sErr)
	assert.Equal(t, "DESTINATION_ERROR", sErr.Code)
	assert.Contains(t, sErr.Message, "cannot inspect destination")
}

func TestCheckDestinationMissing(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	assert.NoError(t, CheckDestination(missing, false))
}

func TestDetectConflictsMixed(t *testing.T) {
	dest := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dest, "exists.txt"), []byte("x"), 0o644))

	conflicts, err := DetectConflicts(dest, []string{"exists.txt", "missing.txt"})
	require.NoError(t, err)
	assert.Equal(t, []string{"exists.txt"}, conflicts)
}

func TestCreateBackupMissing(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	_, err := CreateBackup(missing)
	assert.Error(t, err)
}
