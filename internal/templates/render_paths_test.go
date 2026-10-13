package templates

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderTemplateParseError(t *testing.T) {
	fsys := fstest.MapFS{
		"bad.tmpl": &fstest.MapFile{Data: []byte("{% if .GenerateWeb %}")},
	}
	manifest := []TemplateManifest{
		{SourcePath: "bad.tmpl", TargetPath: "out.txt", Condition: ``},
	}
	r := &Renderer{FS: fsys}
	_, err := r.RenderAll(manifest, resolvedConfig())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse template")
}

func TestRenderTemplateExecutionError(t *testing.T) {
	// ProjectName is a string; accessing a field on it
	// fails at execution time.
	fsys := fstest.MapFS{
		"bad.tmpl": &fstest.MapFile{Data: []byte("{% .ProjectName.Bogus %}")},
	}
	manifest := []TemplateManifest{
		{SourcePath: "bad.tmpl", TargetPath: "out.txt", Condition: ``},
	}
	r := &Renderer{FS: fsys}
	_, err := r.RenderAll(manifest, resolvedConfig())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "render template")
}

func TestPlannedFilesTargetPathError(t *testing.T) {
	manifest := []TemplateManifest{
		{SourcePath: "x.tmpl", TargetPath: "{% if .GenerateWeb %}", Condition: ``},
	}
	r := &Renderer{FS: fstest.MapFS{}}
	_, err := r.PlannedFiles(manifest, resolvedConfig())
	require.Error(t, err)
}

func TestEvaluateConditionParseError(t *testing.T) {
	_, err := EvaluateCondition("{{ if }}", resolvedConfig())
	assert.Error(t, err)
}

func TestEvaluateTargetPathEmpty(t *testing.T) {
	got, err := EvaluateTargetPath("", resolvedConfig())
	require.NoError(t, err)
	assert.Empty(t, got)
}
