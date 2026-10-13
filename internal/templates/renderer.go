package templates

import (
	"bytes"
	"fmt"
	"io/fs"
	"strings"
	"text/template"

	"github.com/tiagozaccaro/ai-project-starter/internal/dependencies"
)

// TemplateManifest describes a single template file: its source
// path in the embedded FS, its output path relative to the
// project root, and an optional condition expression evaluated
// against the ResolvedConfig.
type TemplateManifest struct {
	SourcePath string
	TargetPath string
	Condition  string
	Executable bool
}

// RenderedFile is a fully rendered output file.
type RenderedFile struct {
	TargetPath string
	Content    []byte
	Executable bool
}

// Renderer renders embedded templates against a ResolvedConfig.
type Renderer struct {
	FS fs.FS
}

// NewRenderer returns a Renderer backed by the given filesystem.
func NewRenderer(fsys fs.FS) *Renderer {
	return &Renderer{FS: fsys}
}

// EvaluateCondition evaluates a condition expression against the
// resolved config. An empty condition is always true. The
// expression is a Go template whose output must be "true" or
// "false".
func EvaluateCondition(condition string, rc *dependencies.ResolvedConfig) (bool, error) {
	if strings.TrimSpace(condition) == "" {
		return true, nil
	}
	tmpl, err := template.New("condition").Funcs(Funcs()).Parse(condition)
	if err != nil {
		return false, fmt.Errorf("parse condition %q: %w", condition, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, rc); err != nil {
		return false, fmt.Errorf("evaluate condition %q: %w", condition, err)
	}
	result := strings.TrimSpace(buf.String())
	switch result {
	case "true", "1":
		return true, nil
	case "false", "0", "":
		return false, nil
	default:
		return false, fmt.Errorf("condition %q evaluated to non-boolean %q", condition, result)
	}
}

// EvaluateTargetPath evaluates a manifest target path against
// the resolved config. Target paths may be Go template
// expressions (using {% %} delimiters) so they can adapt to
// the resolved layout (e.g. flat vs monorepo). A plain path
// with no delimiters is returned unchanged.
func EvaluateTargetPath(targetPath string, rc *dependencies.ResolvedConfig) (string, error) {
	tmpl, err := template.New("targetPath").Funcs(Funcs()).Delims("{%", "%}").Parse(targetPath)
	if err != nil {
		return "", fmt.Errorf("parse target path %q: %w", targetPath, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, rc); err != nil {
		return "", fmt.Errorf("evaluate target path %q: %w", targetPath, err)
	}
	return buf.String(), nil
}

// RenderAll iterates the given manifest, evaluates each
// condition, and renders every matching template.
func (r *Renderer) RenderAll(manifest []TemplateManifest, rc *dependencies.ResolvedConfig) ([]RenderedFile, error) {
	var out []RenderedFile
	for _, entry := range manifest {
		match, err := EvaluateCondition(entry.Condition, rc)
		if err != nil {
			return nil, fmt.Errorf("template %s: %w", entry.SourcePath, err)
		}
		if !match {
			continue
		}
		content, err := r.renderTemplate(entry.SourcePath, rc)
		if err != nil {
			return nil, err
		}
		targetPath, err := EvaluateTargetPath(entry.TargetPath, rc)
		if err != nil {
			return nil, fmt.Errorf("template %s: %w", entry.SourcePath, err)
		}
		out = append(out, RenderedFile{
			TargetPath: targetPath,
			Content:    content,
			Executable: entry.Executable,
		})
	}
	return out, nil
}

// renderTemplate reads a single template from the embedded FS
// and renders it against the resolved config. Templates use
// {% %} delimiters so generated JavaScript/JSX that contains
// {{ }} (e.g. React inline styles) is not misinterpreted by
// the template engine.
func (r *Renderer) renderTemplate(sourcePath string, rc *dependencies.ResolvedConfig) ([]byte, error) {
	data, err := fs.ReadFile(r.FS, sourcePath)
	if err != nil {
		return nil, fmt.Errorf("read template %s: %w", sourcePath, err)
	}
	tmpl, err := template.New(sourcePath).Funcs(Funcs()).Delims("{%", "%}").Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse template %s: %w", sourcePath, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, rc); err != nil {
		return nil, fmt.Errorf("render template %s: %w", sourcePath, err)
	}
	return buf.Bytes(), nil
}

// PlannedFiles returns the list of target paths that would be
// generated for the given resolved config, without rendering
// contents. Used by dry-run mode.
func (r *Renderer) PlannedFiles(manifest []TemplateManifest, rc *dependencies.ResolvedConfig) ([]string, error) {
	var out []string
	for _, entry := range manifest {
		match, err := EvaluateCondition(entry.Condition, rc)
		if err != nil {
			return nil, fmt.Errorf("template %s: %w", entry.SourcePath, err)
		}
		if !match {
			continue
		}
		targetPath, err := EvaluateTargetPath(entry.TargetPath, rc)
		if err != nil {
			return nil, fmt.Errorf("template %s: %w", entry.SourcePath, err)
		}
		out = append(out, targetPath)
	}
	return out, nil
}

// CatalogFS returns the embedded filesystem used by the renderer.
func (r *Renderer) CatalogFS() fs.FS { return r.FS }
