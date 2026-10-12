# Architecture

This document describes the architecture of the `ai-project-starter` generator.

## Overview

`ai-project-starter` is a Go CLI code generator. It reads a configuration (`project.config.yaml` or CLI flags), resolves it into a set of capabilities, renders embedded templates, and writes a monorepo of independently buildable applications.

```
┌─────────────┐     ┌──────────────┐     ┌────────────────┐     ┌──────────────┐
│  CLI flags  │────>│  ProjectConfig│────>│  ResolvedConfig │────>│  Templates   │
│  / config   │     │  (validate)   │     │  (capabilities) │     │  (render)    │
└─────────────┘     └──────────────┘     └────────────────┘     └──────┬───────┘
                                                                         │
                                                                         v
                                                               ┌──────────────────┐
                                                               │  Filesystem write │
                                                               │  (atomic, safe)   │
                                                               └──────────────────┘
```

## Layering

| Layer | Package | Responsibility |
|-------|---------|----------------|
| Entry | `cmd/ai-project-starter`, `internal/cli` | Cobra command wiring, flags, version |
| Config | `internal/config` | `ProjectConfig` model, validation, profiles, migration |
| Resolution | `internal/dependencies` | `ProjectConfig` -> `ResolvedConfig` (auto-corrections, capability flags, warnings) |
| Templates | `internal/templates` | Embedded template catalog + `text/template` renderer |
| Filesystem | `internal/filesystem` | Atomic writes, dry-run, destination protection, conflict detection, backup |
| Pipeline | `internal/pipeline` | Orchestration: load -> resolve -> validate -> render -> write -> summary |
| Manifests | `internal/manifests` | Lockfile-generation script and env examples |
| Setup | `internal/setup` | Post-generation verification |
| Summary | `internal/summary` | Human-readable output formatting |

## Data flow

1. **Load**: CLI flags, `--config` file, or `DefaultProfileConfig(profile)` produce a `ProjectConfig`.
2. **Resolve**: `Resolver.Resolve(cfg)` produces a `ResolvedConfig`:
   - Auto-corrects invalid combinations (MongoDB/none disable Alembic migrations; supabase-first disables the FastAPI API).
   - Derives `GenerateWeb`, `GenerateAPI`, `GenerateMobile`, `GenerateAuth`, `AuthProvider`, and every feature flag.
   - Records non-blocking warnings (e.g. `LOGIN_WITHOUT_AUTH`).
3. **Validate**: `ProjectConfig.Validate()` returns a `ConfigValidationError` aggregating *all* violations (never fails fast).
4. **Render**: `Renderer.RenderAll` evaluates each catalog entry's condition against the `ResolvedConfig` and renders matching templates using `[[ ]]` delimiters.
5. **Write**: `filesystem.Writer` writes files atomically (temp file + rename), respecting `--dry-run` and destination protection.
6. **Summarize**: the summary formatter prints the generated tree, warnings, and next steps.

## Configuration model

`ProjectConfig` is the single source of truth and is serialized as `project.config.yaml` (schema version 2). Key invariants:

- **Applications are independent**: `web`, `api`, `mobile` are independently toggleable.
- **Database + architecture are coupled**: `supabase-first`/`hybrid` require PostgreSQL; SQLite/MongoDB only support `fastapi`; `supabase-first` cannot enable the API.
- **Auth is mutually exclusive**: `google-oidc` XOR `supabase`.
- **Features have prerequisites**: dashboard/login/charts/graphs require `web`; ai/workers require `api`.
- **Defaults**: `database.migrations` and `testing` default to `true` when unset.

## Template engine

- Templates are embedded via `//go:embed` in `internal/assets` and compiled into the binary.
- The catalog (`internal/templates/catalog.go`) is the registry: each entry has a `SourcePath`, `TargetPath`, and a `Condition` expression.
- Conditions are Go template expressions evaluated against `ResolvedConfig` (e.g. `{{ and .GenerateWeb .GenerateDashboard }}`).
- Delimiters are `[[` and `]]` so that generated JavaScript/JSX containing `{{ }}` (React inline styles, GitHub Actions `${{ }}`) is emitted verbatim.
- The `Renderer.FS` field is typed as the `fs.FS` interface (not the concrete `embed.FS`), which keeps the renderer unit-testable with `testing/fstest.MapFS`.

## Safety model

- **Dry-run** (`--dry-run`): prints the planned file tree, writes nothing.
- **Destination protection**: a non-empty destination is rejected (exit code 2) unless `--force` is passed; `--force` first renames the destination to a timestamped `.backup-<ts>` directory.
- **Atomic writes**: each file is written to a temp file in its target directory, chmod'd, then renamed into place.
- **Conflict detection**: `DetectConflicts` reports which planned paths already exist.

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Validation error |
| 2 | Destination error (non-empty, no `--force`) |
| 3 | Template/IO error |
| 4 | User cancelled (interactive) |

## Testing strategy

- **Unit tests** (`*_test.go`): config validation, profile copy-on-select, resolver auto-corrections/warnings, condition evaluation, atomic writes, destination protection, exit codes.
- **End-to-end**: generate to a temp dir, then run the generated project's own toolchain (`pnpm` for web, `uv` for api) to confirm the output builds, type-checks, lints, and tests green.

## Contracts

The public interface is defined in:

- [contracts/cli.md](../contracts/cli.md) — commands, flags, exit codes, output formats.
- [contracts/config-schema.md](../contracts/config-schema.md) — the `project.config.yaml` schema and validation rules.
