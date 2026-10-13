# ADR-002: Template Embedding via `embed.FS`

**Status**: Accepted
**Date**: 2026-10-12

## Context

The generator must render hundreds of project templates (apps, databases,
architectures, feature modules) with condition gating, while satisfying the
zero-runtime-dependency constraint: generation must not execute target-ecosystem
tooling (no `uvicorn`, no `npx`) and must not require network access (spec.md
Edge Cases, plan §2 Constraints).

## Decision

Bundle all project templates into the Go binary using the standard library
`embed` package (`//go:embed all:templates` in `internal/assets/assets.go`),
and render them with `text/template`. The template catalog
(`internal/templates/catalog.go`) is the single source of truth mapping each
embedded source path to an output target path plus a condition expression
evaluated against `ResolvedConfig`.

## Rationale

- Templates are static text; `text/template` provides sandboxed rendering with
  condition gating and custom functions (`toSnake`, `toKebab`, `toPascal`, …).
- `embed` compiles templates into the binary at build time: no template reads
  from disk, no network fetches, no missing-file failures in the field.
- Condition evaluation is centralized in the renderer: disabled features produce
  zero output (SC-005, FR-020).
- Cross-platform paths are normalized at write time (`filepath.Join`), keeping
  templates OS-agnostic.

## Consequences

- Adding a template requires: (1) create the `.tmpl` under
  `internal/assets/templates/`, (2) register a `TemplateManifest` entry with an
  explicit condition. The embed picks up new files automatically.
- Template content ships with the binary; template updates require a new release
  (compatibility window: 1 major version, plan §6.4).
- Binary size grows with the template library (~200–300 files, well under the
  50 MB budget).
