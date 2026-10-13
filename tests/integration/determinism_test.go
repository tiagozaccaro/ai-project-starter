package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagozaccaro/ai-project-starter/internal/pipeline"
)

// TestRenderDeterminism verifies rendering the same
// resolved configuration twice produces byte-identical
// output for every file (SC-007, SC-008).
func TestRenderDeterminism(t *testing.T) {
	for _, profile := range profiles() {
		t.Run(profile, func(t *testing.T) {
			rc := loadResolved(t, profile+".yaml")
			first := renderAll(t, rc)
			second := renderAll(t, rc)

			require.Equal(t, len(first), len(second))
			for path, content := range first {
				other, ok := second[path]
				require.True(t, ok, "file %s missing in second render", path)
				assert.Equal(t, string(content), string(other),
					"rendered file %s is not byte-identical across runs", path)
			}
		})
	}
}

// TestGenerateDeterminism verifies two full generations
// into the same destination produce byte-identical files
// on disk (modulo the .git directory, which is created by
// post-generation git init).
func TestGenerateDeterminism(t *testing.T) {
	dest := t.TempDir()
	require.Nil(t, generateTo(t, "fullstack.yaml", dest))
	first := snapshotTree(t, dest)

	// Regenerate into the same (now non-empty) destination
	// with --force to exercise the full pipeline twice.
	require.Nil(t, generateTo(t, "fullstack.yaml", dest, forceOpt))
	second := snapshotTree(t, dest)

	require.Equal(t, len(first), len(second))
	for path, content := range first {
		other, ok := second[path]
		require.True(t, ok, "file %s missing in second generation", path)
		assert.Equal(t, string(content), string(other),
			"generated file %s is not byte-identical across generations", path)
	}
}

// TestGenerateDeterminismAcrossDestinations verifies that
// generating the same config into two different empty
// destinations yields identical content for every file
// except project.config.yaml, which embeds the destination
// path by design.
func TestGenerateDeterminismAcrossDestinations(t *testing.T) {
	destA := t.TempDir()
	destB := t.TempDir()
	require.Nil(t, generateTo(t, "frontend.yaml", destA))
	require.Nil(t, generateTo(t, "frontend.yaml", destB))

	first := snapshotTree(t, destA)
	second := snapshotTree(t, destB)

	require.Equal(t, len(first), len(second))
	for path, content := range first {
		if path == "project.config.yaml" {
			continue
		}
		other, ok := second[path]
		require.True(t, ok, "file %s missing in second destination", path)
		assert.Equal(t, string(content), string(other),
			"generated file %s differs across destinations", path)
	}
}

// forceOpt enables --force for a generation.
func forceOpt(o *pipeline.GenerateOptions) {
	o.Force = true
}

// snapshotTree returns a map of relative path to file
// content for every file under root, excluding the .git
// directory.
func snapshotTree(t *testing.T, root string) map[string]string {
	t.Helper()
	snapshot := make(map[string]string)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		snapshot[rel] = string(data)
		return nil
	})
	require.NoError(t, err)
	return snapshot
}
