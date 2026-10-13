package pipeline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// generateProject generates a valid fullstack project
// into dest and returns the path.
func generateProject(t *testing.T, dest string) {
	t.Helper()
	path := writeConfig(t, "project.config.yaml", validConfig)
	var err error
	captureStdout(t, func() {
		err = Generate(GenerateOptions{ConfigPath: path, Dest: dest})
	})
	require.NoError(t, err)
}

func TestUpgradeUpToDate(t *testing.T) {
	dest := t.TempDir()
	generateProject(t, dest)

	var err error
	out := captureStdout(t, func() {
		err = Upgrade(UpgradeOptions{Dest: dest})
	})
	require.NoError(t, err)
	assert.Contains(t, out, "no template changes detected")
}

func TestUpgradeDryRunDetectsModification(t *testing.T) {
	dest := t.TempDir()
	generateProject(t, dest)

	// Modify a generated file so the diff detects it.
	target := filepath.Join(dest, "README.md")
	require.NoError(t, os.WriteFile(target, []byte("changed\n"), 0o644))

	var err error
	out := captureStdout(t, func() {
		err = Upgrade(UpgradeOptions{Dest: dest, DryRun: true})
	})
	require.NoError(t, err)
	assert.Contains(t, out, "would change")
	assert.Contains(t, out, "modify")
	assert.Contains(t, out, "README.md")
}

func TestUpgradeApplyModifiesFile(t *testing.T) {
	dest := t.TempDir()
	generateProject(t, dest)

	target := filepath.Join(dest, "README.md")
	require.NoError(t, os.WriteFile(target, []byte("changed\n"), 0o644))

	var err error
	out := captureStdout(t, func() {
		err = Upgrade(UpgradeOptions{Dest: dest, Apply: true})
	})
	require.NoError(t, err)
	assert.Contains(t, out, "Applied 1 change(s).")

	// The file is restored to the rendered content.
	data, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.NotEqual(t, "changed\n", string(data))
}

func TestUpgradePromptsDecline(t *testing.T) {
	dest := t.TempDir()
	generateProject(t, dest)

	target := filepath.Join(dest, "README.md")
	require.NoError(t, os.WriteFile(target, []byte("changed\n"), 0o644))

	// Non-interactive stdin (empty) declines the prompt.
	var err error
	out := captureStdout(t, func() {
		err = Upgrade(UpgradeOptions{Dest: dest})
	})
	require.NoError(t, err)
	assert.Contains(t, out, "Aborted")
}

func TestUpgradeMissingConfig(t *testing.T) {
	dest := t.TempDir()
	err := Upgrade(UpgradeOptions{Dest: dest})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "project.config.yaml")
}

func TestUpgradeDefaultDest(t *testing.T) {
	// Empty dest defaults to the current directory,
	// which has no project.config.yaml in tests.
	err := Upgrade(UpgradeOptions{})
	require.Error(t, err)
}

func TestUpgradeDryRunNoChanges(t *testing.T) {
	dest := t.TempDir()
	generateProject(t, dest)

	// No modifications, dry run should report up to date
	var err error
	out := captureStdout(t, func() {
		err = Upgrade(UpgradeOptions{Dest: dest, DryRun: true})
	})
	require.NoError(t, err)
	assert.Contains(t, out, "no template changes detected")
}

func TestUpgradeApplyNoChanges(t *testing.T) {
	dest := t.TempDir()
	generateProject(t, dest)

	// No modifications, apply should report no changes needed
	var err error
	out := captureStdout(t, func() {
		err = Upgrade(UpgradeOptions{Dest: dest, Apply: true})
	})
	require.NoError(t, err)
	assert.Contains(t, out, "Project is up to date")
}

func TestDiffTreesAddAndModify(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	// Same file in both (no change).
	require.NoError(t, os.WriteFile(filepath.Join(src, "same.txt"), []byte("x"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dst, "same.txt"), []byte("x"), 0o644))
	// New file only in src (add).
	require.NoError(t, os.WriteFile(filepath.Join(src, "new.txt"), []byte("x"), 0o644))
	// Modified file.
	require.NoError(t, os.WriteFile(filepath.Join(src, "mod.txt"), []byte("new"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dst, "mod.txt"), []byte("old"), 0o644))

	changes, err := diffTrees(src, dst)
	require.NoError(t, err)
	require.Len(t, changes, 2)

	byPath := make(map[string]string)
	for _, c := range changes {
		byPath[c.Path] = c.Action
	}
	assert.Equal(t, "add", byPath["new.txt"])
	assert.Equal(t, "modify", byPath["mod.txt"])
}

func TestDiffTreesNoChanges(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(src, "a.txt"), []byte("x"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dst, "a.txt"), []byte("x"), 0o644))

	changes, err := diffTrees(src, dst)
	require.NoError(t, err)
	assert.Empty(t, changes)
}

func TestDiffTreesMissingSrc(t *testing.T) {
	_, err := diffTrees(filepath.Join(t.TempDir(), "missing"), t.TempDir())
	require.Error(t, err)
}

func TestSortChanges(t *testing.T) {
	changes := []change{
		{Path: "c.txt"}, {Path: "a.txt"}, {Path: "b.txt"},
	}
	sortChanges(changes)
	assert.Equal(t, "a.txt", changes[0].Path)
	assert.Equal(t, "b.txt", changes[1].Path)
	assert.Equal(t, "c.txt", changes[2].Path)
}

func TestDiffTreesModifyDifferentFile(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(src, "diff.txt"), []byte("rendered"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dst, "diff.txt"), []byte("existing"), 0o644))

	changes, err := diffTrees(src, dst)
	require.NoError(t, err)
	require.Len(t, changes, 1)
	assert.Equal(t, "modify", changes[0].Action)
	assert.Equal(t, "diff.txt", changes[0].Path)
}