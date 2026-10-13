package pipeline

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tiagozaccaro/ai-project-starter/internal/config"
	"github.com/tiagozaccaro/ai-project-starter/internal/dependencies"
	"github.com/tiagozaccaro/ai-project-starter/internal/filesystem"
	"github.com/tiagozaccaro/ai-project-starter/internal/templates"
)

// UpgradeOptions holds the options for the upgrade command.
type UpgradeOptions struct {
	Dest   string
	Apply  bool
	DryRun bool
}

// change describes a single file difference between the
// rendered tree and the destination.
type change struct {
	Action string // "add" or "modify"
	Path   string
}

// Upgrade regenerates an existing project from its saved
// project.config.yaml: load → validate → resolve → render
// to a temp directory → diff → report (or apply).
//
// The generator never runs package managers or touches
// lockfiles (zero-runtime-dependency constraint); only
// template-rendered files are compared and written.
func Upgrade(opts UpgradeOptions) error {
	out := os.Stdout

	destination := opts.Dest
	if destination == "" {
		destination = "."
	}
	configPath := filepath.Join(destination, "project.config.yaml")

	// 1. Load the saved configuration.
	result, err := config.LoadFromFile(configPath)
	if err != nil {
		return fmt.Errorf(
			"load project.config.yaml (run from the project root or pass --dest): %w",
			err,
		)
	}
	cfg := result.Config

	// 2. Validate the complete configuration (FR-003).
	if verr := cfg.Validate(); verr != nil {
		return ValidationError(verr.Error())
	}

	// 3. Resolve dependencies and auto-corrections.
	resolver := dependencies.NewResolver()
	rc, err := resolver.Resolve(cfg)
	if err != nil {
		return ValidationError(err.Error())
	}

	// 4. Render templates into a temporary directory.
	renderer := templates.NewRenderer(templates.AssetsFS())
	manifest := templates.Catalog()
	files, err := renderer.RenderAll(manifest, rc)
	if err != nil {
		return TemplateError(err.Error())
	}

	tmpDir, err := os.MkdirTemp("", "ai-project-starter-upgrade-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	writer := filesystem.NewWriter(tmpDir, false)
	writeFiles := make([]filesystem.File, len(files))
	for i, f := range files {
		writeFiles[i] = filesystem.File{
			TargetPath: f.TargetPath,
			Content:    f.Content,
			Executable: f.Executable,
		}
	}
	if _, err := writer.WriteFiles(writeFiles); err != nil {
		return TemplateError(err.Error())
	}

	// 5. Diff the rendered tree against the destination.
	changes, err := diffTrees(tmpDir, destination)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "Upgrading %q (generator %s)\n",
		rc.ProjectName, cfg.GeneratorVersion)
	if len(changes) == 0 {
		fmt.Fprintln(out, "Project is up to date — no template changes detected.")
		return nil
	}

	fmt.Fprintf(out, "%d file(s) would change:\n", len(changes))
	for _, c := range changes {
		fmt.Fprintf(out, "  %-6s %s\n", c.Action, c.Path)
	}

	if opts.DryRun {
		return nil
	}

	// 6. Prompt to apply unless --apply was given.
	if !opts.Apply {
		fmt.Fprint(out, "\nApply these changes? [y/N] ")
		var answer string
		if _, err := fmt.Fscanln(os.Stdin, &answer); err != nil {
			// Non-interactive stdin: treat as no.
			fmt.Fprintln(out, "\nAborted (re-run with --apply to write changes).")
			return nil
		}
		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer != "y" && answer != "yes" {
			fmt.Fprintln(out, "Aborted.")
			return nil
		}
	}

	// 7. Apply changes.
	for _, c := range changes {
		src := filepath.Join(tmpDir, filepath.FromSlash(c.Path))
		dst := filepath.Join(destination, filepath.FromSlash(c.Path))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return fmt.Errorf("create dirs for %s: %w", c.Path, err)
		}
		data, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("read rendered %s: %w", c.Path, err)
		}
		mode := os.FileMode(0o644)
		if info, err := os.Stat(src); err == nil {
			mode = info.Mode()
		}
		if err := os.WriteFile(dst, data, mode); err != nil {
			return fmt.Errorf("write %s: %w", c.Path, err)
		}
	}
	fmt.Fprintf(out, "Applied %d change(s).\n", len(changes))
	return nil
}

// diffTrees compares the rendered tree (src) against the
// destination (dst), reporting added and modified files.
// Removals are not applied: regenerating from scratch is
// the supported path for removed templates.
func diffTrees(src, dst string) ([]change, error) {
	var changes []change

	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		target := filepath.Join(dst, filepath.FromSlash(rel))
		existing, err := os.ReadFile(target)
		if err != nil {
			if os.IsNotExist(err) {
				changes = append(changes, change{Action: "add", Path: rel})
				return nil
			}
			return err
		}
		rendered, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !bytes.Equal(existing, rendered) {
			changes = append(changes, change{Action: "modify", Path: rel})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("diff trees: %w", err)
	}

	// Stable, readable ordering.
	sortChanges(changes)
	return changes, nil
}

func sortChanges(changes []change) {
	for i := 1; i < len(changes); i++ {
		for j := i; j > 0 && changes[j].Path < changes[j-1].Path; j-- {
			changes[j], changes[j-1] = changes[j-1], changes[j]
		}
	}
}
