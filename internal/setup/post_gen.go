// Package setup performs post-generation tasks:
// git initialization, .gitignore hardening, and
// output integrity validation.
package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tiagozaccaro/ai-project-starter/internal/dependencies"
)

// PostGeneration performs git init and integrity checks.
type PostGeneration struct {
	Destination string
}

// NewPostGeneration returns a PostGeneration for destination.
func NewPostGeneration(destination string) *PostGeneration {
	return &PostGeneration{Destination: destination}
}

// GitInit initializes a git repository in the destination
// when git is available. It is best-effort: a missing git
// binary is not fatal.
func (p *PostGeneration) GitInit() error {
	if _, err := exec.LookPath("git"); err != nil {
		return nil
	}
	cmd := exec.Command("git", "init")
	cmd.Dir = p.Destination
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git init: %w", err)
	}
	return nil
}

// HardenGitignore ensures the generated .gitignore excludes
// secrets and environment files (Constitution VI).
func (p *PostGeneration) HardenGitignore() error {
	giPath := filepath.Join(p.Destination, ".gitignore")
	data, err := os.ReadFile(giPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	required := []string{".env*", "*.pem", "*.key", "secrets.*", "*.kubeconfig"}
	content := string(data)
	var missing []string
	for _, pattern := range required {
		if !strings.Contains(content, pattern) {
			missing = append(missing, pattern)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	f, err := os.OpenFile(giPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, pattern := range missing {
		fmt.Fprintln(f, pattern)
	}
	return nil
}

// ValidateIntegrity verifies that all enabled features are
// present and no disabled features leaked into the output.
// It returns a list of integrity violations.
func (p *PostGeneration) ValidateIntegrity(rc *dependencies.ResolvedConfig) []string {
	var violations []string
	present, err := collectFiles(p.Destination)
	if err != nil {
		return []string{fmt.Sprintf("integrity scan failed: %v", err)}
	}

	// Expected directories for enabled applications. In monorepo
	// mode each application lives under apps/<name>/; in flat
	// mode the single application's marker file lives at the
	// project root.
	appPresent := func(dir, marker string) bool {
		if rc.FlatLayout {
			return hasFile(present, marker)
		}
		return hasDir(present, dir)
	}
	if rc.GenerateWeb {
		if !appPresent("apps/web", "package.json") {
			violations = append(violations, "enabled web application missing")
		}
	}
	if rc.GenerateAPI {
		if !appPresent("apps/api", "pyproject.toml") {
			violations = append(violations, "enabled api application missing")
		}
	}
	if rc.GenerateMobile {
		if !appPresent("apps/mobile", "pubspec.yaml") {
			violations = append(violations, "enabled mobile application missing")
		}
	}

	// Disabled applications must not be present.
	if !rc.GenerateWeb && hasDir(present, "apps/web") {
		violations = append(violations, "disabled web application leaked apps/web")
	}
	if !rc.GenerateAPI && hasDir(present, "apps/api") {
		violations = append(violations, "disabled api application leaked apps/api")
	}
	if !rc.GenerateMobile && hasDir(present, "apps/mobile") {
		violations = append(violations, "disabled mobile application leaked apps/mobile")
	}

	// Feature-specific presence checks.
	if rc.GenerateDocker {
		if !hasFile(present, "docker-compose.yml") {
			violations = append(violations, "docker enabled but docker-compose.yml missing")
		}
	} else {
		if hasFile(present, "docker-compose.yml") {
			violations = append(violations, "docker disabled but docker-compose.yml present")
		}
	}
	if rc.GenerateVSCode {
		if !hasDir(present, ".vscode") {
			violations = append(violations, "vscode enabled but .vscode missing")
		}
	} else {
		if hasDir(present, ".vscode") {
			violations = append(violations, "vscode disabled but .vscode present")
		}
	}
	if rc.GenerateGitHub {
		if !hasDir(present, ".github") {
			violations = append(violations, "github enabled but .github missing")
		}
	} else {
		if hasDir(present, ".github") {
			violations = append(violations, "github disabled but .github present")
		}
	}
	if rc.GenerateSupabase {
		if !hasDir(present, "supabase") {
			violations = append(violations, "supabase architecture enabled but supabase/ missing")
		}
	}
	if rc.GenerateAlembic {
		alembicDir := filepath.Join(rc.APIDir(), "alembic")
		if !hasDir(present, alembicDir) {
			violations = append(violations, "migrations enabled but alembic missing")
		}
	}
	if !rc.MigrationsEnabled {
		alembicDir := filepath.Join(rc.APIDir(), "alembic")
		if hasDir(present, alembicDir) {
			violations = append(violations, "migrations disabled but alembic present")
		}
	}

	return violations
}

func collectFiles(root string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		out = append(out, rel)
		return nil
	})
	return out, err
}

func hasDir(files []string, dir string) bool {
	for _, f := range files {
		if f == dir || strings.HasPrefix(f, dir+string(filepath.Separator)) || strings.HasPrefix(f, dir+"/") {
			return true
		}
	}
	return false
}

func hasFile(files []string, name string) bool {
	for _, f := range files {
		if f == name {
			return true
		}
	}
	return false
}
