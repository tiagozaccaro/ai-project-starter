# AGENTS.md — ai-project-starter (Generator)

This repository contains the `ai-project-starter` generator itself: a Go CLI that
scaffolds standalone full-stack projects (Next.js web, FastAPI API, Flutter
mobile) from a validated `project.config.yaml`.

## Stack

- **Language**: Go 1.22+ (generics, `embed`, `io/fs`, `text/template`)
- **CLI**: `spf13/cobra`; interactive prompts: `charmbracelet/huh`
- **Config**: `gopkg.in/yaml.v3`; layering: `spf13/viper` (flags > env > config file)
- **Templates**: `text/template` + `embed.FS` (`internal/assets`, embedded via `//go:embed all:templates`)
- **Tests**: stdlib `testing` + `stretchr/testify`; golden comparison via `google/go-cmp`

## Development Rules

- Single authoritative config model: `internal/config/schema.go` (`ProjectConfig`). No duplicate models in the CLI or template layer.
- Feature selection lives only in `internal/dependencies/resolver.go`. Template files MUST NOT contain feature-selection branching that duplicates resolver logic.
- Every template is registered in `internal/templates/catalog.go` with an explicit condition expression evaluated against `ResolvedConfig`. Disabled features produce zero output (SC-005).
- Validation: `ProjectConfig.Validate()` reports ALL errors at once with error code, field reference, and suggested fix (FR-003).
- Auto-corrections (mongodb/none → migrations=false) are applied by the resolver, noted in the summary — never silent.
- Generated projects are standalone (Constitution II): never reference the generator at runtime.

## Testing Thresholds

- Generator test coverage ≥ 90% (SC-015): `make coverage` must show ≥90% total.
- `go test -race ./...` passes on macOS and Linux.
- Golden file tests (`tests/integration/golden_test.go`) must pass — they catch template regressions.
- Determinism: identical config → byte-identical output (modulo timestamps/lockfiles).

## Commit Conventions

- Conventional commits: `feat:`, `fix:`, `chore:`, `docs:`, `test:`.
- Scope optional: `feat(cli): add upgrade command`.
- One logical change per commit; commit after each task or logical group.

## Security

- Never commit secrets, tokens, or populated env files (Constitution VI).
- Generated `.gitignore` must exclude `.env*`, `*.pem`, `*.key`, `secrets.*`.
- Generated Dockerfiles use non-root `USER`.

## Directory Map

- `cmd/ai-project-starter/` — binary entry point
- `internal/cli/` — cobra commands (generate, validate, migrate-config, list-profiles, list-features, upgrade, completion)
- `internal/config/` — schema, loader, profiles, migration
- `internal/dependencies/` — resolver → `ResolvedConfig`
- `internal/templates/` — catalog, renderer, funcs
- `internal/filesystem/` — atomic writer, dry-run, destination safety
- `internal/manifests/` — per-app dependency manifests
- `internal/pipeline/` — generate orchestration + interactive wizard
- `internal/setup/` — post-generation integrity validation
- `internal/summary/` — output formatting
- `internal/assets/templates/` — embedded project templates
- `shadcn-test/` — **template validation bed**: a standalone Next.js app used to manually validate the shadcn/ui component and dashboard/login block templates under `internal/assets/templates/apps/web/` and `features/dashboard|login/` before they ship. It is intentionally excluded from the generator's template catalog and is not part of any generated project.

## Change Control

- Changes touching ≥ 3 files or modifying public APIs require spec/plan/tasks (Constitution VII).
- Do not silently change architecture, suppress failing tests, or introduce unnecessary dependencies (Constitution III).
