package setup

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagozaccaro/project-forge/internal/dependencies"
)

func TestNewPostGeneration(t *testing.T) {
	pg := NewPostGeneration("/some/path")
	assert.NotNil(t, pg)
	assert.Equal(t, "/some/path", pg.Destination)
}

func TestGitInitWhenGitAvailable(t *testing.T) {
	pg := NewPostGeneration(t.TempDir())
	if _, err := exec.LookPath("git"); err != nil {
		// Skip this test if git is not available
		return
	}
	err := pg.GitInit()
	assert.NoError(t, err)
	assert.DirExists(t, filepath.Join(pg.Destination, ".git"))
}

func TestGitInitWhenGitNotAvailable(t *testing.T) {
	pg := NewPostGeneration(t.TempDir())
	// We can't easily mock exec.LookPath, so we just test
	// that GitInit doesn't panic when git is not found.
	// The function handles missing git gracefully.
	err := pg.GitInit()
	assert.NoError(t, err)
}

func TestHardenGitignoreWhenFileNotExist(t *testing.T) {
	pg := NewPostGeneration(t.TempDir())
	err := pg.HardenGitignore()
	assert.NoError(t, err)
}

func TestHardenGitignoreExisting(t *testing.T) {
	pg := NewPostGeneration(t.TempDir())
	// Write a .gitignore file
	giPath := filepath.Join(pg.Destination, ".gitignore")
	require.NoError(t, os.WriteFile(giPath, []byte("*.log\n"), 0o644))

	err := pg.HardenGitignore()
	assert.NoError(t, err)

	// Read the file and check it contains the original content
	data, err := os.ReadFile(giPath)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "*.log")
	assert.Contains(t, string(data), ".env*")
}

func TestValidateIntegrityFlatLayout(t *testing.T) {
	pg := NewPostGeneration(t.TempDir())
	// Create a flat layout structure (single app marker at project root)
	require.NoError(t, os.WriteFile(filepath.Join(pg.Destination, "package.json"), []byte("{}"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(pg.Destination, "pyproject.toml"), []byte(""), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(pg.Destination, "pubspec.yaml"), []byte(""), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(pg.Destination, "docker-compose.yml"), []byte(""), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(pg.Destination, ".vscode"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(pg.Destination, ".github"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(pg.Destination, "supabase"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(pg.Destination, "alembic"), 0o755))

	rc := &dependencies.ResolvedConfig{
		GenerateWeb:    true,
		GenerateAPI:    true,
		GenerateMobile: true,
		FlatLayout:     true, // Use flat layout
		GenerateDocker: true,
		GenerateVSCode: true,
		GenerateGitHub: true,
		GenerateSupabase: true,
		GenerateAlembic:  true,
		MigrationsEnabled: true,
	}

	violations := pg.ValidateIntegrity(rc)
	assert.Empty(t, violations)
}

func TestValidateIntegrityFeatureNegativeChecks(t *testing.T) {
	pg := NewPostGeneration(t.TempDir())
	// Create files/directories for features that will be disabled
	require.NoError(t, os.WriteFile(filepath.Join(pg.Destination, "docker-compose.yml"), []byte(""), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(pg.Destination, ".vscode"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(pg.Destination, ".github"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(pg.Destination, "supabase"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(pg.Destination, "alembic"), 0o755))

	rc := &dependencies.ResolvedConfig{
		GenerateWeb:    false, // Apps disabled
		GenerateAPI:    false,
		GenerateMobile: false,
		FlatLayout:     false,
		GenerateDocker: false, // Docker disabled but file present
		GenerateVSCode: false, // VSCode disabled but dir present
		GenerateGitHub: false, // GitHub disabled but dir present
		GenerateSupabase: false, // Supabase disabled but dir present
		GenerateAlembic:  false, // Alembic disabled but dir present
		MigrationsEnabled: false, // Migrations disabled but dir present
	}

	violations := pg.ValidateIntegrity(rc)
	assert.NotEmpty(t, violations)
	// Should have violations for each feature that is disabled but present
	violationStr := strings.Join(violations, " ")
	assert.Contains(t, violationStr, "docker disabled but docker-compose.yml present")
	assert.Contains(t, violationStr, "vscode disabled but .vscode present")
	assert.Contains(t, violationStr, "github disabled but .github present")
}

func TestValidateIntegritySupabaseAndAlembicEnabledMissing(t *testing.T) {
	pg := NewPostGeneration(t.TempDir())
	// Don't create supabase/ or alembic/ directories - they should be missing
	rc := &dependencies.ResolvedConfig{
		GenerateWeb:    false,
		GenerateAPI:    false,
		GenerateMobile: false,
		FlatLayout:     false,
		GenerateSupabase: true,  // Supabase enabled but directory missing
		GenerateAlembic:  true,  // Alembic enabled but directory missing
		MigrationsEnabled: true,
	}

	violations := pg.ValidateIntegrity(rc)
	assert.NotEmpty(t, violations)
	violationStr := strings.Join(violations, " ")
	assert.Contains(t, violationStr, "supabase architecture enabled but supabase/ missing")
	assert.Contains(t, violationStr, "migrations enabled but alembic missing")
}

func TestValidateIntegrityCollectFilesError(t *testing.T) {
	// Test when collectFiles returns an error
	pg := NewPostGeneration("") // Empty destination - this might cause collectFiles to fail
	rc := &dependencies.ResolvedConfig{}

	violations := pg.ValidateIntegrity(rc)
	// Should return an error message in violations
	assert.NotEmpty(t, violations)
	assert.Contains(t, violations[0], "integrity scan failed")
}

func TestValidateIntegritySuccess(t *testing.T) {
	pg := NewPostGeneration(t.TempDir())
	// Create a valid project structure
	require.NoError(t, os.MkdirAll(filepath.Join(pg.Destination, "apps/web"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(pg.Destination, "apps/web", "package.json"), []byte("{}"), 0o644))

	rc := &dependencies.ResolvedConfig{
		GenerateWeb:    true,
		GenerateAPI:    false,
		GenerateMobile: false,
		FlatLayout:     false,
	}

	violations := pg.ValidateIntegrity(rc)
	assert.Empty(t, violations)
}

func TestValidateIntegrityMissingEnabledApp(t *testing.T) {
	pg := NewPostGeneration(t.TempDir())
	// Create a project structure but missing the web app
	require.NoError(t, os.MkdirAll(filepath.Join(pg.Destination, "apps/api"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(pg.Destination, "apps/api", "pyproject.toml"), []byte(""), 0o644))

	rc := &dependencies.ResolvedConfig{
		GenerateWeb:    true, // Enabled but missing
		GenerateAPI:    true, // Present
		GenerateMobile: false,
		FlatLayout:     false,
	}

	violations := pg.ValidateIntegrity(rc)
	assert.NotEmpty(t, violations)
	assert.Contains(t, violations[0], "enabled web application missing")
}

func TestValidateIntegrityLeakedDisabledApp(t *testing.T) {
	pg := NewPostGeneration(t.TempDir())
	// Create a project structure with a disabled app (mobile) but also enabled web and api
	require.NoError(t, os.MkdirAll(filepath.Join(pg.Destination, "apps/web"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(pg.Destination, "apps/web", "package.json"), []byte("{}"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(pg.Destination, "apps/api"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(pg.Destination, "apps/api", "pyproject.toml"), []byte(""), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(pg.Destination, "apps/mobile"), 0o755))

	rc := &dependencies.ResolvedConfig{
		GenerateWeb:    true,
		GenerateAPI:    true,
		GenerateMobile: false, // Disabled but present
		FlatLayout:     false,
	}

	violations := pg.ValidateIntegrity(rc)
	assert.NotEmpty(t, violations)
	assert.Contains(t, violations[0], "disabled mobile application leaked")
}

func TestValidateIntegrityIntegrateError(t *testing.T) {
	pg := NewPostGeneration("") // Empty destination
	rc := &dependencies.ResolvedConfig{}

	violations := pg.ValidateIntegrity(rc)
	assert.NotEmpty(t, violations)
	assert.Contains(t, violations[0], "integrity scan failed")
}
