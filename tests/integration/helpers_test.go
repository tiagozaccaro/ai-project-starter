package integration

import (
	"errors"
	"path/filepath"
	"sort"
	"testing"

	"github.com/tiagozaccaro/ai-project-starter/internal/config"
	"github.com/tiagozaccaro/ai-project-starter/internal/dependencies"
	"github.com/tiagozaccaro/ai-project-starter/internal/pipeline"
	"github.com/tiagozaccaro/ai-project-starter/internal/templates"
)

// fixturesDir is the path to the test fixtures relative to
// the tests/integration package.
const fixturesDir = ".." + string(filepath.Separator) + "fixtures"

// fixturePath returns the absolute path of a config fixture.
func fixturePath(name string) string {
	return filepath.Join(fixturesDir, "configs", name)
}

// goldenPath returns the absolute path of a golden manifest.
func goldenPath(profile string) string {
	return filepath.Join(fixturesDir, "golden", profile, "manifest.txt")
}

// loadResolved loads a config fixture and resolves it against
// the dependency resolver. The destination is pinned to a fixed
// value so rendered output is a pure function of the config and
// therefore comparable against golden manifests.
func loadResolved(t *testing.T, fixture string) *dependencies.ResolvedConfig {
	t.Helper()
	result, err := config.LoadFromFile(fixturePath(fixture))
	if err != nil {
		t.Fatalf("load fixture %s: %v", fixture, err)
	}
	cfg := result.Config
	cfg.Destination = "myproject"
	if err := cfg.Validate(); err != nil {
		t.Fatalf("validate fixture %s: %v", fixture, err)
	}
	resolver := dependencies.NewResolver()
	rc, err := resolver.Resolve(cfg)
	if err != nil {
		t.Fatalf("resolve fixture %s: %v", fixture, err)
	}
	return rc
}

// renderAll renders every catalog template for rc and returns
// a map of target path to rendered content.
func renderAll(t *testing.T, rc *dependencies.ResolvedConfig) map[string][]byte {
	t.Helper()
	renderer := templates.NewRenderer(templates.AssetsFS())
	files, err := renderer.RenderAll(templates.Catalog(), rc)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	out := make(map[string][]byte, len(files))
	for _, f := range files {
		out[f.TargetPath] = f.Content
	}
	return out
}

// buildManifest renders rc and returns the sorted golden
// manifest body: one "sha256  path" line per rendered file.
func buildManifest(t *testing.T, rc *dependencies.ResolvedConfig) []byte {
	t.Helper()
	files := renderAll(t, rc)
	paths := make([]string, 0, len(files))
	for p := range files {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	var buf []byte
	for _, p := range paths {
		sum := sha256Hex(files[p])
		buf = append(buf, []byte(sum+"  "+p+"\n")...)
	}
	return buf
}

// profiles returns the built-in profile names in golden order.
func profiles() []string {
	return []string{"fullstack", "backend", "frontend", "ai-fullstack", "custom"}
}

// generateTo runs the full generation pipeline for a config
// fixture into dest. It returns the resulting CLI error, or
// nil on success. Option functions mutate the options before
// the pipeline runs.
func generateTo(t *testing.T, fixture, dest string, opts ...func(*pipeline.GenerateOptions)) *pipeline.CLIError {
	t.Helper()
	options := pipeline.GenerateOptions{
		ConfigPath:     fixturePath(fixture),
		Dest:           dest,
		NonInteractive: true,
	}
	for _, opt := range opts {
		opt(&options)
	}
	err := pipeline.Generate(options)
	if err == nil {
		return nil
	}
	var cliErr *pipeline.CLIError
	if errors.As(err, &cliErr) {
		return cliErr
	}
	t.Fatalf("unexpected error type %T: %v", err, err)
	return nil
}
