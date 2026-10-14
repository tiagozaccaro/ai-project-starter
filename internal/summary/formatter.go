// Package summary formats generation results, dry-run
// trees, and diagnostics.
package summary

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/tiagozaccaro/project-forge/internal/dependencies"
)

// Formatter renders human-readable and machine-readable
// summaries.
type Formatter struct {
	Out io.Writer
}

// NewFormatter returns a Formatter writing to out.
func NewFormatter(out io.Writer) *Formatter {
	return &Formatter{Out: out}
}

// PrintSuccess prints the post-generation summary.
func (f *Formatter) PrintSuccess(rc *dependencies.ResolvedConfig, written []string) {
	fmt.Fprintf(f.Out, "Project %q generated at %s\n\n", rc.ProjectName, rc.Destination())
	fmt.Fprintf(f.Out, "Profile:      %s\n", rc.Profile)
	fmt.Fprintf(f.Out, "Applications: %s\n", formatApps(rc))
	fmt.Fprintf(f.Out, "Database:     %s (%s architecture)\n", rc.Database.Type, rc.Database.Architecture)
	fmt.Fprintf(f.Out, "Features:     %s\n", formatFeatures(rc))
	fmt.Fprintf(f.Out, "Files:        %d\n\n", len(written))

	if len(rc.AutoCorrections) > 0 {
		fmt.Fprintln(f.Out, "Auto-corrections:")
		for _, ac := range rc.AutoCorrections {
			fmt.Fprintf(f.Out, "  - %s: %s -> %s (%s)\n", ac.Field, ac.From, ac.To, ac.Reason)
		}
		fmt.Fprintln(f.Out)
	}
	if len(rc.Warnings) > 0 {
		fmt.Fprintln(f.Out, "Warnings:")
		for _, w := range rc.Warnings {
			fmt.Fprintf(f.Out, "  - [%s] %s\n", w.Code, w.Message)
		}
		fmt.Fprintln(f.Out)
	}

	fmt.Fprintln(f.Out, "Next steps:")
	fmt.Fprintf(f.Out, "  cd %s\n", rc.Destination())
	if rc.GenerateAPI {
		fmt.Fprintf(f.Out, "  cd %s && uv sync && uv run pytest\n", rc.APIDir())
	}
	if rc.GenerateWeb {
		fmt.Fprintf(f.Out, "  cd %s && pnpm install && pnpm dev\n", rc.WebDir())
	}
}

// PrintDryRun prints the dry-run file tree.
func (f *Formatter) PrintDryRun(rc *dependencies.ResolvedConfig, planned []string) {
	fmt.Fprintf(f.Out, "DRY RUN - would create %d files:\n\n", len(planned))
	fmt.Fprintln(f.Out, formatTree(rc.ProjectName, planned))
}

// PrintDiagnostics prints verbose diagnostics.
func (f *Formatter) PrintDiagnostics(rc *dependencies.ResolvedConfig) {
	fmt.Fprintf(f.Out, "Resolved config: %s\n", rc.String())
	fmt.Fprintf(f.Out, "Destination: %s\n", rc.Destination())
	fmt.Fprintf(f.Out, "Web: %v  API: %v  Mobile: %v\n", rc.GenerateWeb, rc.GenerateAPI, rc.GenerateMobile)
	fmt.Fprintf(f.Out, "Auth: %s  Migrations: %v  Testing: %v\n", rc.AuthProvider, rc.MigrationsEnabled, rc.TestingEnabled)
}

// JSONSummary is the machine-readable summary output.
type JSONSummary struct {
	ProjectName     string                        `json:"projectName"`
	Profile         string                        `json:"profile"`
	Destination     string                        `json:"destination"`
	Applications    []string                      `json:"applications"`
	DatabaseType    string                        `json:"databaseType"`
	Architecture    string                        `json:"architecture"`
	AuthProvider    string                        `json:"authProvider"`
	FileCount       int                           `json:"fileCount"`
	Files           []string                      `json:"files"`
	AutoCorrections []dependencies.AutoCorrection `json:"autoCorrections,omitempty"`
	Warnings        []dependencies.Warning        `json:"warnings,omitempty"`
}

// PrintJSON prints the summary as JSON.
func (f *Formatter) PrintJSON(rc *dependencies.ResolvedConfig, files []string) error {
	s := JSONSummary{
		ProjectName:     rc.ProjectName,
		Profile:         rc.Profile,
		Destination:     rc.Destination(),
		Applications:    formatAppsList(rc),
		DatabaseType:    rc.Database.Type,
		Architecture:    rc.Database.Architecture,
		AuthProvider:    rc.AuthProvider,
		FileCount:       len(files),
		Files:           files,
		AutoCorrections: rc.AutoCorrections,
		Warnings:        rc.Warnings,
	}
	enc := json.NewEncoder(f.Out)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

func formatApps(rc *dependencies.ResolvedConfig) string {
	return strings.Join(formatAppsList(rc), ", ")
}

func formatAppsList(rc *dependencies.ResolvedConfig) []string {
	var apps []string
	if rc.GenerateWeb {
		apps = append(apps, "web (Next.js)")
	}
	if rc.GenerateAPI {
		apps = append(apps, "api (FastAPI)")
	}
	if rc.GenerateMobile {
		apps = append(apps, "mobile (Flutter)")
	}
	return apps
}

func formatFeatures(rc *dependencies.ResolvedConfig) string {
	var feats []string
	if rc.GenerateAuth {
		feats = append(feats, "auth:"+rc.AuthProvider)
	}
	if rc.MigrationsEnabled {
		feats = append(feats, "database.migrations")
	}
	if rc.GenerateDashboard {
		feats = append(feats, "dashboard")
	}
	if rc.GenerateLogin {
		feats = append(feats, "login")
	}
	if rc.GenerateCharts {
		feats = append(feats, "charts")
	}
	if rc.GenerateGraphs {
		feats = append(feats, "graphs")
	}
	if rc.GenerateAI {
		feats = append(feats, "ai")
	}
	if rc.GenerateWorkers {
		feats = append(feats, "workers")
	}
	if rc.GenerateObservability {
		feats = append(feats, "observability")
	}
	if rc.GenerateDocker {
		feats = append(feats, "docker")
	}
	if rc.GenerateVSCode {
		feats = append(feats, "vscode")
	}
	if rc.GenerateGitHub {
		feats = append(feats, "github")
	}
	if rc.GenerateSpeckit {
		feats = append(feats, "speckit")
	}
	if rc.GenerateMCP {
		feats = append(feats, "mcp")
	}
	if rc.TestingEnabled {
		feats = append(feats, "testing")
	}
	return strings.Join(feats, ", ")
}

// treeNode is an internal tree node for rendering dry-run trees.
type treeNode struct {
	name     string
	children map[string]*treeNode
	isFile   bool
}

// formatTree renders a simple tree view of the planned files.
func formatTree(root string, files []string) string {
	sorted := make([]string, len(files))
	copy(sorted, files)
	sort.Strings(sorted)

	rootNode := &treeNode{name: root, children: map[string]*treeNode{}}
	for _, f := range sorted {
		parts := strings.Split(f, "/")
		cur := rootNode
		for i, p := range parts {
			isFile := i == len(parts)-1
			child, ok := cur.children[p]
			if !ok {
				child = &treeNode{name: p, children: map[string]*treeNode{}, isFile: isFile}
				cur.children[p] = child
			}
			cur = child
		}
	}

	var b strings.Builder
	renderNode(&b, rootNode, "", true)
	return b.String()
}

func renderNode(b *strings.Builder, n *treeNode, prefix string, isLast bool) {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	if n != nil && n.name != "" {
		b.WriteString(prefix + connector + n.name + "\n")
	}
	newPrefix := prefix
	if isLast {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}
	// Stable ordering: directories first, then files, alphabetical.
	names := make([]string, 0, len(n.children))
	for name := range n.children {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		ci := n.children[names[i]]
		cj := n.children[names[j]]
		if ci.isFile != cj.isFile {
			return !ci.isFile
		}
		return names[i] < names[j]
	})
	for i, name := range names {
		renderNode(b, n.children[name], newPrefix, i == len(names)-1)
	}
}
