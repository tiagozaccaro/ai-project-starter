package pipeline

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYAMLUnmarshal(t *testing.T) {
	var m map[string]interface{}
	require.NoError(t, yamlUnmarshal([]byte("a: 1\nb: two\n"), &m))
	assert.Equal(t, 1, m["a"])
	assert.Equal(t, "two", m["b"])
}

func TestYAMLUnmarshalError(t *testing.T) {
	var m map[string]interface{}
	assert.Error(t, yamlUnmarshal([]byte(":\n- broken"), &m))
}

func TestYAMLMarshal(t *testing.T) {
	data, err := yamlMarshal(map[string]string{"a": "b"})
	require.NoError(t, err)
	assert.Contains(t, string(data), "a: b")
}

func TestIsTerminal(t *testing.T) {
	// A regular file is never a terminal.
	f, err := os.CreateTemp(t.TempDir(), "regular")
	require.NoError(t, err)
	defer f.Close()
	assert.False(t, isTerminal(f))
}

func TestIsTerminalClosedFile(t *testing.T) {
	// A closed file's Stat returns an error, exercising the
	// error-return path in isTerminal.
	f, err := os.CreateTemp(t.TempDir(), "regular")
	require.NoError(t, err)
	require.NoError(t, f.Close())
	assert.False(t, isTerminal(f))
}
