// Package manifests computes per-application dependency
// sets from a ResolvedConfig and produces the
// lockfile-generation workflow.
package manifests

import (
	"fmt"
	"strings"

	"github.com/tiagozaccaro/ai-project-starter/internal/dependencies"
)

// DependencySet is the computed set of dependencies for
// one application.
type DependencySet struct {
	App     string
	Runtime []string
	Dev     []string
}

// Generator computes dependency manifests.
type Generator struct{}

// NewGenerator returns a configured Generator.
func NewGenerator() *Generator { return &Generator{} }

// WebDependencies returns the web app dependency set.
func (g *Generator) WebDependencies(rc *dependencies.ResolvedConfig) DependencySet {
	if !rc.GenerateWeb {
		return DependencySet{}
	}
	ds := DependencySet{App: "web"}
	ds.Runtime = []string{
		"next",
		"react",
		"react-dom",
	}
	ds.Dev = []string{
		"typescript",
		"@types/react",
		"@types/node",
		"eslint",
		"@typescript-eslint/parser",
		"@typescript-eslint/eslint-plugin",
		"prettier",
	}
	// Tailwind + shadcn/ui are used by the design system.
	ds.Dev = append(ds.Dev, "tailwindcss", "postcss", "autoprefixer")
	if rc.GenerateDashboard || rc.GenerateLogin || rc.GenerateCharts || rc.GenerateGraphs {
		ds.Dev = append(ds.Dev, "lucide-react")
		ds.Dev = append(ds.Dev, "class-variance-authority", "clsx", "tailwind-merge")
	}
	if rc.GenerateDashboard {
		ds.Dev = append(ds.Dev, "@radix-ui/react-slot")
	}
	if rc.GenerateCharts {
		ds.Runtime = append(ds.Runtime, "recharts")
	}
	if rc.GenerateGraphs {
		ds.Runtime = append(ds.Runtime, "@xyflow/react")
	}
	if rc.GenerateAPI {
		ds.Runtime = append(ds.Runtime, "@tanstack/react-query")
	}
	return ds
}

// APIDependencies returns the API app dependency set.
func (g *Generator) APIDependencies(rc *dependencies.ResolvedConfig) DependencySet {
	if !rc.GenerateAPI {
		return DependencySet{}
	}
	ds := DependencySet{App: "api"}
	ds.Runtime = []string{
		"fastapi",
		"uvicorn[standard]",
		"pydantic",
		"pydantic-settings",
		"sqlalchemy[asyncio]",
	}
	switch rc.Database.Type {
	case "postgresql":
		ds.Runtime = append(ds.Runtime, "psycopg[binary]", "aiosqlite")
	case "sqlite":
		ds.Runtime = append(ds.Runtime, "aiosqlite")
	case "mysql":
		ds.Runtime = append(ds.Runtime, "aiomysql", "aiosqlite")
	case "mongodb":
		ds.Runtime = append(ds.Runtime, "motor", "beanie")
	}
	if rc.GenerateAuth {
		if rc.AuthProvider == "google-oidc" {
			ds.Runtime = append(ds.Runtime, "authlib", "python-jose[cryptography]", "google-auth", "httpx")
		} else if rc.AuthProvider == "supabase" {
			ds.Runtime = append(ds.Runtime, "supabase")
		}
	}
	if rc.GenerateAI {
		ds.Runtime = append(ds.Runtime, "openai")
	}
	if rc.GenerateWorkers {
		ds.Runtime = append(ds.Runtime, "celery", "redis")
	}
	if rc.GenerateObservability {
		ds.Runtime = append(ds.Runtime, "opentelemetry-api", "opentelemetry-sdk")
	}
	ds.Dev = []string{
		"pytest",
		"pytest-asyncio",
		"pytest-cov",
		"httpx",
		"ruff",
		"pyright",
	}
	if rc.GenerateAlembic {
		ds.Dev = append(ds.Dev, "alembic")
		ds.Runtime = append(ds.Runtime, "alembic")
	}
	return ds
}

// LockfileWorkflow returns the shell script that generates
// lockfiles for all enabled applications. The generator
// never runs package managers at generation time (zero
// runtime dependency constraint); instead it emits this
// workflow for the user to run once after generation.
func (g *Generator) LockfileWorkflow(rc *dependencies.ResolvedConfig) string {
	var b strings.Builder
	b.WriteString("#!/usr/bin/env bash\n")
	b.WriteString("# Lockfile generation workflow for " + rc.ProjectName + ".\n")
	b.WriteString("# Run once after generation to produce committed lockfiles.\n")
	b.WriteString("set -euo pipefail\n\n")
	if rc.GenerateWeb {
		b.WriteString("# Web: pnpm-lock.yaml\n")
		b.WriteString("cd \"$(dirname \"$0\")/../apps/web\"\n")
		b.WriteString("pnpm install --lockfile-only\n")
		b.WriteString("cd - >/dev/null\n\n")
	}
	if rc.GenerateAPI {
		b.WriteString("# API: uv.lock\n")
		b.WriteString("cd \"$(dirname \"$0\")/../apps/api\"\n")
		b.WriteString("uv lock\n")
		b.WriteString("cd - >/dev/null\n\n")
	}
	if rc.GenerateMobile {
		b.WriteString("# Mobile: pubspec.lock\n")
		b.WriteString("cd \"$(dirname \"$0\")/../apps/mobile\"\n")
		b.WriteString("flutter pub get\n")
		b.WriteString("cd - >/dev/null\n\n")
	}
	b.WriteString("echo \"Lockfiles generated. Commit them to the repository.\"\n")
	return b.String()
}

// Summary returns a one-line dependency summary for diagnostics.
func (g *Generator) Summary(rc *dependencies.ResolvedConfig) string {
	var parts []string
	if rc.GenerateWeb {
		parts = append(parts, fmt.Sprintf("web:%d", len(g.WebDependencies(rc).Runtime)))
	}
	if rc.GenerateAPI {
		parts = append(parts, fmt.Sprintf("api:%d", len(g.APIDependencies(rc).Runtime)))
	}
	return strings.Join(parts, " ")
}
