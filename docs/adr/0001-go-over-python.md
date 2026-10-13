# ADR-001: Go for the Generator

**Status**: Accepted
**Date**: 2026-10-12
**Deciders**: Project maintainers (per user decision during `/speckit.analyze`)

## Context

The generator (`ai-project-starter`) must be distributable as a standalone binary
via npm, pipx, and Homebrew (spec.md FR-030). An early draft of the plan chose
Python 3.11+ with `uv`; during analysis the user explicitly reverted that
deviation and selected Go, matching spec.md §188 and Clarifications §197.

## Decision

Implement the generator in **Go 1.22+** using the standard library plus:
`spf13/cobra` (CLI), `spf13/viper` (config layering), `gopkg.in/yaml.v3` (YAML),
`text/template` (rendering), `charmbracelet/huh` (interactive prompts),
`stretchr/testify` and `google/go-cmp` (tests).

## Rationale

- Go compiles to a single static binary with zero runtime dependencies — natively
  satisfying FR-030 without PyInstaller-style packaging workarounds.
- Cross-compilation: one toolchain produces macOS (amd64/arm64) and Linux
  (amd64/arm64) binaries.
- `cobra` provides subcommands, flags, and exit-code management matching
  `contracts/cli.md`.
- `embed.FS` bundles ~200–300 template files into the binary; no network or
  filesystem reads at generation time.
- The Kilo Code / Spec Kit workflow already supports Go generators via
  `go build` / `go run` in CI.

## Consequences

- Generated projects are unchanged: web stays Next.js/pnpm, API stays
  FastAPI/uv, mobile stays Flutter. Only the generator is Go.
- Config validation is implemented as explicit Go validator functions in
  `internal/config/schema.go` (not pydantic model_validators).
- Template embedding uses `embed` (ADR-002), not `importlib.resources`.
- Python offers richer runtime schema validation; mitigated by golden file tests
  that enforce template correctness.
