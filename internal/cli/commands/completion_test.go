package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCompletionCommand(t *testing.T) {
	cmd := NewCompletionCommand()
	assert.Equal(t, "completion", cmd.Name())
	assert.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)
	assert.NotNil(t, cmd.Flags())

	// Test bash completion
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"bash"})
	require.NoError(t, cmd.Execute())
	assert.Contains(t, out.String(), "bash completion")

	// Reset
	cmd = NewCompletionCommand()
	out = &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"zsh"})
	require.NoError(t, cmd.Execute())
	assert.Contains(t, out.String(), "compdef")

	// Reset
	cmd = NewCompletionCommand()
	out = &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"fish"})
	require.NoError(t, cmd.Execute())
	assert.Contains(t, out.String(), "fish completion")

	// Reset
	cmd = NewCompletionCommand()
	out = &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"powershell"})
	require.NoError(t, cmd.Execute())
	assert.Contains(t, out.String(), "powershell completion")
}

func TestNewCompletionCommandInvalid(t *testing.T) {
	cmd := NewCompletionCommand()
	cmd.SetArgs([]string{"bad"})
	err := cmd.Execute()
	require.Error(t, err)
}
