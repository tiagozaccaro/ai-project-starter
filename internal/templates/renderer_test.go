package templates

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagozaccaro/ai-project-starter/internal/config"
	"github.com/tiagozaccaro/ai-project-starter/internal/dependencies"
)

func resolvedConfig() *dependencies.ResolvedConfig {
	return &dependencies.ResolvedConfig{
		ProjectConfig: &config.ProjectConfig{
			ProjectName: "my-app",
		},
		GenerateWeb: true,
		GenerateAPI: true,
	}
}

func TestEvaluateConditionEmpty(t *testing.T) {
	ok, err := EvaluateCondition("", resolvedConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("empty condition should match")
	}
}

func TestEvaluateConditionTrue(t *testing.T) {
	ok, err := EvaluateCondition("{{ .GenerateWeb }}", resolvedConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true condition to match")
	}
}

func TestEvaluateConditionFalse(t *testing.T) {
	ok, err := EvaluateCondition("{{ .GenerateMobile }}", resolvedConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false condition to not match")
	}
}

func TestEvaluateConditionAnd(t *testing.T) {
	ok, err := EvaluateCondition("{{ and .GenerateWeb .GenerateAPI }}", resolvedConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected and condition to match")
	}
}

func TestEvaluateConditionNonBoolean(t *testing.T) {
	_, err := EvaluateCondition("{{ .ProjectName }}", resolvedConfig())
	if err == nil {
		t.Fatal("expected error for non-boolean condition")
	}
}

func TestEvaluateConditionInvalidSyntax(t *testing.T) {
	_, err := EvaluateCondition("{{ if .GenerateWeb }}", resolvedConfig())
	if err == nil {
		t.Fatal("expected error for unterminated if")
	}
}

func TestRenderAllSelectsMatchingTemplates(t *testing.T) {
	fsys := fstest.MapFS{
		"web.tmpl":    &fstest.MapFile{Data: []byte("web: {% .ProjectName %}")},
		"api.tmpl":    &fstest.MapFile{Data: []byte("api: {% .ProjectName %}")},
		"mobile.tmpl": &fstest.MapFile{Data: []byte("mobile: {% .ProjectName %}")},
	}
	manifest := []TemplateManifest{
		{SourcePath: "web.tmpl", TargetPath: "out/web.txt", Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "api.tmpl", TargetPath: "out/api.txt", Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "mobile.tmpl", TargetPath: "out/mobile.txt", Condition: `{{ .GenerateMobile }}`},
	}
	r := &Renderer{FS: fsys}
	files, err := r.RenderAll(manifest, resolvedConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if string(files[0].Content) != "web: my-app" {
		t.Errorf("unexpected web content: %q", string(files[0].Content))
	}
	if string(files[1].Content) != "api: my-app" {
		t.Errorf("unexpected api content: %q", string(files[1].Content))
	}
}

func TestRenderAllDelims(t *testing.T) {
	fsys := fstest.MapFS{
		"t.tmpl": &fstest.MapFile{Data: []byte("{% .ProjectName %} {{ notATemplate }}")},
	}
	manifest := []TemplateManifest{
		{SourcePath: "t.tmpl", TargetPath: "out.txt", Condition: ``},
	}
	r := &Renderer{FS: fsys}
	files, err := r.RenderAll(manifest, resolvedConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := string(files[0].Content)
	expected := "my-app {{ notATemplate }}"
	if got != expected {
		t.Errorf("expected delimiters rendered and {{ }} preserved, got %q", got)
	}
}

func TestRenderAllMissingTemplate(t *testing.T) {
	fsys := fstest.MapFS{}
	manifest := []TemplateManifest{
		{SourcePath: "missing.tmpl", TargetPath: "out.txt", Condition: ``},
	}
	r := &Renderer{FS: fsys}
	_, err := r.RenderAll(manifest, resolvedConfig())
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}

func TestPlannedFiles(t *testing.T) {
	manifest := []TemplateManifest{
		{SourcePath: "web.tmpl", TargetPath: "out/web.txt", Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "mobile.tmpl", TargetPath: "out/mobile.txt", Condition: `{{ .GenerateMobile }}`},
	}
	r := &Renderer{FS: fstest.MapFS{}}
	paths, err := r.PlannedFiles(manifest, resolvedConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 1 || paths[0] != "out/web.txt" {
		t.Errorf("expected [out/web.txt], got %v", paths)
	}
}

func TestEvaluateTargetPathParseError(t *testing.T) {
	_, err := EvaluateTargetPath("{% invalid syntax here", resolvedConfig())
	if err == nil {
		t.Fatal("expected error for malformed target path")
	}
}

func TestEvaluateTargetPathExecError(t *testing.T) {
	_, err := EvaluateTargetPath("{% .UndefinedField %}", resolvedConfig())
	if err == nil {
		t.Fatal("expected error for undefined field in target path")
	}
}

func TestRenderAllConditionError(t *testing.T) {
	fsys := fstest.MapFS{
		"t.tmpl": &fstest.MapFile{Data: []byte("content")},
	}
	manifest := []TemplateManifest{
		{SourcePath: "t.tmpl", TargetPath: "out.txt", Condition: "{{ if true }}"},
	}
	r := &Renderer{FS: fsys}
	_, err := r.RenderAll(manifest, resolvedConfig())
	if err == nil {
		t.Fatal("expected error for unterminated condition")
	}
}

func TestRenderAllTargetPathError(t *testing.T) {
	fsys := fstest.MapFS{
		"t.tmpl": &fstest.MapFile{Data: []byte("content")},
	}
	manifest := []TemplateManifest{
		{SourcePath: "t.tmpl", TargetPath: "{% undefined %}/out.txt", Condition: ""},
	}
	r := &Renderer{FS: fsys}
	_, err := r.RenderAll(manifest, resolvedConfig())
	if err == nil {
		t.Fatal("expected error for malformed target path")
	}
}

func TestRenderAllExecError(t *testing.T) {
	fsys := fstest.MapFS{
		"t.tmpl": &fstest.MapFile{Data: []byte("{% .MissingField %}")},
	}
	manifest := []TemplateManifest{
		{SourcePath: "t.tmpl", TargetPath: "out.txt", Condition: ""},
	}
	r := &Renderer{FS: fsys}
	_, err := r.RenderAll(manifest, resolvedConfig())
	if err == nil {
		t.Fatal("expected error for missing field execution")
	}
}

func TestNewRenderer(t *testing.T) {
	fs := fstest.MapFS{}
	r := NewRenderer(fs)
	require.NotNil(t, r)
	require.Equal(t, fs, r.FS)
}

func TestCatalog(t *testing.T) {
	catalog := Catalog()
	require.NotEmpty(t, catalog)
	// Should have many templates
	assert.GreaterOrEqual(t, len(catalog), 10)
}

func TestAssetsFS(t *testing.T) {
	fs := AssetsFS()
	require.NotNil(t, fs)
}
