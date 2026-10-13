package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitInitGitNotFound(t *testing.T) {
	// With no git on PATH, GitInit is a best-effort no-op.
	t.Setenv("PATH", "/nonexistent-path-for-test")
	pg := NewPostGeneration(t.TempDir())
	assert.NoError(t, pg.GitInit())
}
