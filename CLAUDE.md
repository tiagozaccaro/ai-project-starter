# AGENTS.md — ai-project-starter (generator)

Instructions for AI coding agents working on the `ai-project-starter` generator itself.

## What this is

A Go CLI that scaffolds new software projects. It is a **code generator**, not a runtime application. The output is a monorepo of independently buildable apps (`apps/web`, `apps/api`, `apps/mobile`) driven by a validated `project.config.yaml`.

## Commands

```bash
go build -o ai-project-starter ./cmd/ai-project-starter   # build
go test ./...                                              # run all tests
go vet ./...                                               # static analysis
gofmt -l .                                                 # check formatting
```

There is no `golangci-lint` in this environment; `go vet` + `gofmt` are the linters.

## Architecture

- **Entry point**: `cmd/ai-project-starter/main.go` -> `internal/cli`.
- **`internal/config`**: `ProjectConfig` is the authoritative model. `Validate()` returns *all* violations at once (`ConfigValidationError`). Profiles are presets in `profiles.go` (copy-on-select: `DefaultProfileConfig` returns a deep copy, never a live reference).
- **`internal/dependencies`**: `Resolver.Resolve()` turns a `ProjectConfig` into a `ResolvedConfig` — applies auto-corrections (e.g. MongoDB/none disable migrations), derives capability flags (`GenerateWeb`, `GenerateAPI`, ...), and records warnings.
- **`internal/templates`**: templates are embedded via `//go:embed` (`internal/assets`). The catalog (`catalog.go`) maps each template to a target path and a condition expression. The renderer (`renderer.go`) uses `text/template` with `[[ ]]` delimiters (so generated JS/JSX containing `{{ }}` is not misparsed).
- **`internal/filesystem`**: atomic write-then-rename, dry-run collection, destination protection (`CheckDestination`), conflict detection, backup creation.
- **`internal/pipeline`**: orchestrates load -> resolve -> validate -> render -> write -> summary. Exit codes are defined in `errors.go` (`CLIError`).

## Conventions

- **Templates use `[[ ]]` delimiters**, never `{{ }}`. Generated code that needs `{{ }}` (e.g. React JSX, GitHub Actions `${{ }}`) is written literally inside templates.
- **Go template syntax only**: `if`/`else if`/`else`/`end`, `and`/`or`/`not`, `eq`/`ne`. There is no `||` or `&&`.
- **Always close actions with `]]`**: an `else if`/`else` branch must end with `]]`, even when the following text starts with `[` (e.g. `[[ else if .X ]]["web"]`). A single `]` is a parse error.
- **Validation is aggregate**: add every violation via `errs.Add(...)`; do not return early on the first error.
- **Profiles are immutable presets**: mutate the copy returned by `DefaultProfileConfig`, never the package-level data.
- **Generated apps are independent**: no cross-`apps/*` runtime dependencies.

## Testing

- Unit tests live next to the code (`*_test.go`).
- End-to-end validation: generate to a temp dir, then run the generated app's own toolchain (`pnpm install/build/typecheck/test`, `uv sync/pytest/ruff/pyright`).
- When adding a template, register it in `internal/templates/catalog.go` with a condition, and verify it renders for every combination that enables it.

## Do not

- Do not add dependencies without justification.
- Do not introduce runtime network calls during template rendering (templates are static text).
- Do not commit generated `node_modules/`, `.venv/`, or lockfiles of the *generator* repo itself.
