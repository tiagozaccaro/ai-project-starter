# ProjectForge

A powerful, production-ready project scaffolding generator that forges consistent, maintainable, and AI-friendly software repositories from a single validated configuration. It eliminates boilerplate and enforces best practices from day one.

ProjectForge is a developer tool, not an application — it generates standalone, independently buildable projects you own completely. No runtime dependency on the generator.

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-blue.svg)](go.mod) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](#license) [![Test Coverage](https://img.shields.io/badge/Coverage-%E2%89%A590%25-brightgreen.svg)](#testing)

## Table of Contents

- [Why ProjectForge?](#why-projectforge)
- [What It Generates](#what-it-generates)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage Examples](#usage-examples)
- [Commands Reference](#commands-reference)
- [Configuration](#configuration)
- [Architecture Overview](#architecture-overview)
- [Generated Project Structure](#generated-project-structure)
- [Development](#development)
- [Testing](#testing)
- [Contributing](#contributing)
- [Troubleshooting](#troubleshooting)
- [License](#license)

## Why ProjectForge?

Setting up a modern full-stack project involves dozens of decisions: architecture, tooling, linting, testing, CI/CD, Docker, auth, and more. ProjectForge codifies these decisions into repeatable templates so you can:

- **Ship faster**: Generate production-ready boilerplate in seconds, not hours
- **Stay consistent**: Enforce architecture patterns across all your projects
- **Be AI-friendly**: Generated projects include AGENTS.md, specs, and clear structure so coding agents work effectively
- **Own your code**: Generated projects are completely standalone — no generator runtime dependency
- **Validate early**: Catch configuration errors before a single file is written
- **Stay deterministic**: Identical configuration produces byte-identical output

## What It Generates

ProjectForge scaffolds a monorepo containing independently buildable applications, all driven by a single `project.config.yaml` source of truth:

- **[Next.js (Web)](https://nextjs.org/)** — TypeScript, App Router, Tailwind CSS, shadcn/ui, ESLint, Prettier, Vitest/Jest ready
- **[FastAPI (API)](https://fastapi.tiangolo.com/)** — Python 3.11+, async-first, Pydantic v2, SQLAlchemy/Alembic, pytest, ruff, pyright, uv package manager
- **[Flutter (Mobile)](https://flutter.dev/)** — Cross-platform mobile app with clean architecture, Riverpod/BLoC ready, testing setup
- **Shared tooling** — Docker Compose, GitHub Actions CI/CD, editor configs (VS Code), agent instructions (AGENTS.md), and Speckit specifications

Each application builds, lints, type-checks, and tests independently.

## Features

| Feature | Description |
| --- | --- |
| **Preset Profiles** | Choose from `fullstack`, `backend`, `frontend`, `ai-fullstack`, or `custom` to jumpstart your setup |
| **Independent Apps** | Toggle web/api/mobile independently — generate exactly what you need |
| **Flexible Databases** | PostgreSQL, MySQL, SQLite, MongoDB, or none with proper migrations/ORM setup |
| **Architecture Modes** | Support for `fastapi`, `supabase-first`, or `hybrid` architectures |
| **Auth Providers** | Google OIDC or Supabase Auth with sensible defaults and security best practices |
| **Feature Modules** | Add dashboard, login, charts/graphs, AI capabilities, workers, observability, Docker, VS Code, GitHub Actions, Speckit, MCP |
| **Safe Generation** | Dry-run preview, destination protection with automatic backups, atomic writes, and conflict detection |
| **Rock-Solid Validation** | Full config validation with actionable, aggregated error messages — never fails fast on the first error |
| **Smart Auto-Corrections** | Intelligently resolves invalid combinations and reports non-blocking warnings transparently |
| **Deterministic Output** | Same config always produces identical files (timestamps/lockfiles aside) |
| **AI-First Design** | Includes AGENTS.md, agent configs, and Speckit specs for optimal LLM-assisted development |
| **Repeatable Scaffolding** | Commit your `project.config.yaml` and regenerate anytime with confidence |

## Prerequisites

To build ProjectForge from source: **Go 1.23+**

For generated projects, you'll need:

| Application | Requirements |
| --- | --- |
| **Web (Next.js)** | [Node.js](https://nodejs.org/) ≥ 20, [pnpm](https://pnpm.io/) recommended |
| **API (FastAPI)** | [Python](https://www.python.org/) ≥ 3.11, [uv](https://docs.astral.sh/uv/) (fast Python package manager) |
| **Mobile (Flutter)** | [Flutter](https://flutter.dev/) ≥ 3.19, [Dart](https://dart.dev/) SDK |

## Installation

### From Source (Recommended for Development)

```bash
git clone https://github.com/tiagozaccaro/project-forge.git
cd project-forge
go build -o projectforge ./cmd/project-forge
```

Verify the build:

```bash
./projectforge --version
```

### Add to PATH (Optional)

```bash
# Move to a directory in your PATH
sudo mv projectforge /usr/local/bin/
projectforge --version
```

### Cross-Platform Builds

```bash
# Linux amd64
GOOS=linux GOARCH=amd64 go build -o projectforge-linux-amd64 ./cmd/project-forge

# macOS arm64
GOOS=darwin GOARCH=arm64 go build -o projectforge-darwin-arm64 ./cmd/project-forge

# Windows amd64
GOOS=windows GOARCH=amd64 go build -o projectforge.exe ./cmd/project-forge
```

## Quick Start

### Generate a Full-Stack App (Non-Interactive)

Create a Next.js + FastAPI project in one command:

```bash
./projectforge generate --profile fullstack --name my-app --dest ./my-app --non-interactive
cd my-app
```

Spin up the web app:

```bash
cd apps/web
pnpm install
pnpm dev
```

Spin up the API:

```bash
cd apps/api
uv sync
uv run fastapi dev main.py
```

### Interactive Mode (Guided Wizard)

Prefer to configure step by step? Launch the interactive wizard:

```bash
./projectforge generate --name my-awesome-app
```

The wizard will walk you through profiles, apps, database, auth, features, and more.

### Preview Before Generating (Dry Run)

See exactly what will be created without writing files:

```bash
./projectforge generate --profile fullstack --name my-app --dry-run
```

Perfect for validating your choices or understanding the output structure.

## Usage Examples

### Backend-Only API

```bash
./projectforge generate \
  --profile backend \
  --name api-service \
  --dest ./api-service \
  --non-interactive
```

### Frontend-Only Next.js App

```bash
./projectforge generate \
  --profile frontend \
  --name web-app \
  --dest ./web-app \
  --non-interactive
```

### AI-Focused Full-Stack

```bash
./projectforge generate \
  --profile ai-fullstack \
  --name ai-app \
  --dest ./ai-app \
  --non-interactive
```

### Generate from a Saved Config

For reproducible scaffolding, commit your config:

```bash
# First, generate and save the config (it’s included in generated output)
# Or create it manually
cat > project.config.yaml << 'EOF'
schemaVersion: 2
generatorVersion: "1.0.0"
profile: fullstack
projectName: my-app
applications:
  web: true
  api: true
  mobile: false
database:
  type: postgresql
  architecture: fastapi
features:
  auth: ["google-oidc"]
  dashboard: true
  login: true
  database.migrations: true
  testing: true
EOF

./projectforge generate --config project.config.yaml --dest ./my-app
```

### Regenerate an Existing Project (Upgrade)

Safely update generated files while preserving your custom code:

```bash
cd my-existing-project
./projectforge upgrade --dest . --dry-run  # Preview changes
./projectforge upgrade --dest . --apply    # Apply safely
```

The upgrade command diffs against your destination and never touches lockfiles or user-edited files.

### Overwrite a Non-Empty Directory

If you need to regenerate into an existing directory:

```bash
./projectforge generate --profile fullstack --name my-app --dest ./my-app --force
```

This automatically creates a timestamped backup: `my-app.backup-<timestamp>/`.

## Commands Reference

| Command | Short | Description |
| --- | --- | --- |
| [`generate`](#generate) | — | Generate a new project from a profile, wizard, or config file |
| [`validate`](#validate) | — | Validate a `project.config.yaml` without generating |
| [`migrate-config`](#migrate-config) | — | Migrate a config file to the current schema version |
| [`list-profiles`](#list-profiles) | — | List available profiles with descriptions |
| [`list-features`](#list-features) | — | List configurable feature modules and dependencies |
| [`upgrade`](#upgrade) | — | Regenerate an existing project from its saved configuration |
| [`completion`](#completion) | — | Generate shell completion scripts (bash/zsh/fish/powershell) |

### `generate`

Generate a project with flexible input modes.

```bash
projectforge generate [flags]
```

| Flag | Short | Type | Default | Description |
| --- | --- | --- | --- | --- |
| `--profile` | `-p` | string | — | Profile: `fullstack`, `backend`, `frontend`, `ai-fullstack`, `custom` |
| `--name` | `-n` | string | — | Project name (kebab-case recommended) |
| `--dest` | `-d` | string | `./{name}` | Destination directory |
| `--config` | `-c` | string | — | Path to `project.config.yaml` |
| `--layout` | — | string | — | Project layout: `monorepo`, `flat` |
| `--interactive` | `-i` | bool | false | Force interactive wizard |
| `--non-interactive` | — | bool | false | Disable prompts; fail if incomplete |
| `--dry-run` | — | bool | false | Preview file tree without writing |
| `--force` | `-f` | bool | false | Overwrite non-empty destination (creates backup) |
| `--output` | — | string | — | Output format: `json`, `yaml` |
| `--verbose` | `-v` | bool | false | Increase verbosity |
| `--quiet` | `-q` | bool | false | Suppress non-error output |
| `--no-color` | — | bool | false | Disable colored output |

### `validate`

Validate configuration files with detailed error reporting.

```bash
projectforge validate --config project.config.yaml [flags]
```

| Flag | Short | Type | Default | Description |
| --- | --- | --- | --- | --- |
| `--config` | `-c` | string | — | Path to config file (**required**) |
| `--strict` | — | bool | false | Treat warnings as errors |
| `--output` | — | string | — | Output format: `json`, `yaml` |

**Exit codes:** `0` if valid, `1` if validation fails.

### `migrate-config`

Migrate older config schemas to the current version.

```bash
projectforge migrate-config --config old-config.yaml [flags]
```

| Flag | Short | Type | Default | Description |
| --- | --- | --- | --- | --- |
| `--config` | `-c` | string | — | Path to config file (**required**) |
| `--output` | `-o` | string | overwrite in-place | Output path for migrated config |
| `--no-backup` | — | bool | false | Skip creating `.bak` file |
| `--target-version` | — | int | latest | Target schema version |
| `--dry-run` | — | bool | false | Show diff without writing |

### `list-profiles`

Explore available starter profiles.

```bash
projectforge list-profiles [flags]
```

| Flag | Short | Type | Default | Description |
| --- | --- | --- | --- | --- |
| `--output` | `-o` | string | `table` | Output format: `table`, `json`, `yaml` |

### `list-features`

Discover all feature modules and their dependencies.

```bash
projectforge list-features [flags]
```

| Flag | Short | Type | Default | Description |
| --- | --- | --- | --- | --- |
| `--profile` | `-p` | string | — | Filter by specific profile |
| `--output` | `-o` | string | `table` | Output format: `table`, `json`, `yaml` |

### `upgrade`

Regenerate an existing project safely from its stored config.

```bash
projectforge upgrade [flags]
```

| Flag | Short | Type | Default | Description |
| --- | --- | --- | --- | --- |
| `--dest` | `-d` | string | `.` (current dir) | Destination directory containing `project.config.yaml` |
| `--apply` | — | bool | false | Write changes without prompting |
| `--dry-run` | — | bool | false | Preview changes without writing |

### `completion`

Generate shell completions for tab completion.

```bash
projectforge completion bash|zsh|fish|powershell
```

**Examples:**

```bash
# Bash
projectforge completion bash > ~/.bash_completion.d/projectforge
source ~/.bash_completion.d/projectforge

# Zsh
projectforge completion zsh > "${fpath[1]}/_projectforge"
autoload -U compinit && compinit

# Fish
projectforge completion fish > ~/.config/fish/completions/projectforge.fish

# PowerShell
projectforge completion powershell > projectforge.ps1
```

## Configuration

Every generated project includes a `project.config.yaml` file — the single source of truth for that project. This file captures your choices and enables reproducible regeneration/upgrade.

### Minimal Example (Full-Stack)

```yaml
schemaVersion: 2
generatorVersion: "1.0.0"
profile: fullstack
projectName: my-app
applications:
  web: true
  api: true
  mobile: false
database:
  type: postgresql
  architecture: fastapi
features:
  auth: ["google-oidc"]
  dashboard: true
  login: true
  database.migrations: true
  testing: true
  docker: true
  github.actions: true
```

### Backend-Only Example

```yaml
schemaVersion: 2
generatorVersion: "1.0.0"
profile: backend
projectName: api-service
applications:
  web: false
  api: true
  mobile: false
database:
  type: postgresql
  architecture: fastapi
features:
  database.migrations: true
  testing: true
  observability: true
```

### Key Concepts

- **Schema Versioning**: Configs are versioned (`schemaVersion`) and can be migrated with `migrate-config`
- **Validation First**: All constraints are enforced before generation — mutual exclusivity, prerequisites, and valid combinations
- **Auto-Corrections**: Invalid combinations are intelligently corrected (with warnings) rather than failing silently
- **Independent Applications**: `web`, `api`, `mobile` can be toggled independently based on your needs
- **Architecture Constraints**: `supabase-first`/`hybrid` require PostgreSQL; SQLite/MongoDB only support `fastapi`; `supabase-first` disables API generation

For the complete schema specification, see [contracts/config-schema.md](contracts/config-schema.md).

## Architecture Overview

ProjectForge follows a clean, layered architecture that separates concerns and is highly testable.

```text
CLI Flags / Config File
        │
        ▼
    ┌─────────────────┐
    │  ProjectConfig  │  ← Single source of truth (validated)
    │  (validation)   │
    └────────┬────────┘
             │
             ▼
    ┌───────────────────┐
    │  ResolvedConfig   │  ← Capabilities + auto-corrections + warnings
    │  (dependencies)   │
    └────────┬──────────┘
             │
             ▼
    ┌─────────────────┐
    │   Templates     │  ← Embedded via go:embed + text/template
    │   (rendered)    │
    └────────┬────────┘
             │
             ▼
    ┌─────────────────────┐
    │  Filesystem Writer  │  ← Atomic writes, dry-run, safe, protected
    └─────────────────────┘
```

### Core Packages

| Package | Purpose |
| --- | --- |
| `cmd/project-forge` | Application entry point |
| `internal/cli` | Cobra CLI wiring, flags, version handling |
| `internal/config` | Schema, validation, profiles, config migration |
| `internal/dependencies` | Resolves config → capabilities (enforces constraints & auto-corrections) |
| `internal/templates` | Embedded template catalog, condition evaluation, renderer |
| `internal/filesystem` | Atomic file writes, destination protection, backups, conflict detection |
| `internal/pipeline` | Generation orchestration (load → resolve → validate → render → write → summarize) |
| `internal/manifests` | Manifest/lockfile generation helpers |
| `internal/setup` | Post-generation integrity verification |
| `internal/summary` | Human-readable output formatting |
| `internal/assets/templates` | Embedded project template files (`go:embed all:templates`) |

### Key Design Principles

- **Single Source of Truth**: `ProjectConfig` in `internal/config/schema.go` — no duplication across layers
- **Feature Logic Centralized**: All feature dependencies live in `internal/dependencies/resolver.go`; templates never duplicate selection logic
- **Fail Aggregated, Never Fast**: Validation reports ALL errors at once with codes, field references, and suggested fixes
- **Template Safety**: Conditions are explicit per template in the catalog; disabled features produce zero output
- **Delimiter Isolation**: Templates use `[[`/`]]` delimiters to preserve `{{ }}` in generated code (React, GitHub Actions)
- **Testable by Design**: Renderer accepts `fs.FS` interface (not concrete embed.FS), enabling `fstest.MapFS` in unit tests

### Exit Codes

| Code | Meaning |
| --- | --- |
| `0` | Success |
| `1` | Validation error (invalid config) |
| `2` | Destination error (non-empty without `--force`) |
| `3` | Template/IO error |
| `4` | User cancelled (interactive mode) |

## Generated Project Structure

Here's a typical full-stack project generated by ProjectForge:

```text
my-app/
├── .github/workflows/        # CI/CD pipelines
├── .vscode/                  # Editor settings/debug configs
├── apps/
│   ├── web/                  # Next.js frontend (pnpm)
│   ├── api/                  # FastAPI backend (uv + pyproject.toml)
│   └── mobile/               # Flutter app (pub)
├── docker/                   # Dockerfiles & docker-compose
├── docs/                     # Project documentation
├── specs/                    # Speckit specifications
├── AGENTS.md                 # AI agent instructions
├── project.config.yaml       # Source of truth (critical - commit this!)
├── README.md                 # Generated project README
└── .gitignore
```

**Important**: Always commit `project.config.yaml`. It enables reproducible regeneration and the `upgrade` command.

## Development

### Build

```bash
go build -o projectforge ./cmd/project-forge
```

### Run Tests

```bash
# All tests
go test ./...

# With verbose output
go test -v ./...

# With race detection
go test -race ./...

# Specific package
go test -v ./internal/config/... ./internal/dependencies/...
```

### Code Quality

```bash
# Vet
go vet ./...

# Format
gofmt -l .        # Check formatting
gofmt -w .        # Auto-format

# Tidy modules
go mod tidy
```

### Coverage

ProjectForge maintains **≥90% test coverage**.

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View in terminal
go tool cover -func=coverage.out

# View HTML report
go tool cover -html=coverage.out
```

Target: **≥90% total coverage** (enforced by `make coverage` if using Makefile).

## Testing

### Test Strategy

- **Unit Tests** (`*_test.go`): Config validation, profiles, resolver logic, condition evaluation, filesystem safety, exit codes
- **Integration Tests**: Golden file tests (`tests/integration/golden_test.go`) catch template regressions
- **E2E Validation**: Generated projects are validated by building/running their native toolchains (pnpm build/typecheck/test for web, uv sync/pytest/ruff/pyright for API)

### Golden Tests

Golden tests compare generated output against committed expected fixtures. If templates change intentionally, update golden files:

```bash
# Update golden files (be careful - review diffs!)
go test -v ./tests/integration/... -update
```

### Determinism

ProjectForge guarantees deterministic output: identical `ProjectConfig` → byte-identical generated files (excluding timestamps and lockfiles like `pnpm-lock.yaml`, `uv.lock`, `pubspec.lock`).

## Contributing

Contributions are welcome! Please follow these guidelines:

### Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make your changes (follow conventions below)
4. Add/update tests (maintain ≥90% coverage)
5. Run tests: `go test -race ./...`
6. Format code: `gofmt -w .`
7. Commit using [Conventional Commits](https://www.conventionalcommits.org/)
8. Push and open a Pull Request

### Commit Conventions

Use Conventional Commits with optional scopes:

- `feat:` — New feature
- `fix:` — Bug fix
- `docs:` — Documentation changes
- `test:` — Test additions/changes
- `chore:` — Maintenance/refactoring
- `perf:` — Performance improvements

**Examples:**

```bash
feat(cli): add --output json support to list-profiles
fix(resolver): correct MongoDB migrations auto-correction
docs(readme): add troubleshooting section
test(config): add validation edge cases
```

### Guidelines

- **Single logical change per commit** — commit after completing each task
- **Architecture changes**: Changes touching ≥3 files or modifying public APIs require a spec/plan (see Constitution in CLAUDE.md/AGENTS.md)
- **No silent behavior changes**: Be explicit about auto-corrections and warnings
- **Test coverage**: New code must maintain or improve ≥90% coverage
- **Security first**: Never commit secrets, tokens, or populated `.env` files. Generated `.gitignore` must exclude sensitive patterns
- **Follow existing patterns**: Mimic code style, naming conventions, and architectural decisions in the codebase

### Change Control

Per project constitution: Do not silently change architecture, suppress failing tests, or introduce unnecessary dependencies without justification.

## Troubleshooting

### Common Issues

| Issue | Cause | Solution |
| --- | --- | --- |
| "destination is not empty" | Target directory has files | Add `--force` to create a backup and overwrite, or choose a different `--dest` |
| Validation errors on first attempt | Config has multiple issues | Fix reported errors (all shown at once with suggestions). Use `validate` command to debug |
| Templates produce unexpected output | Condition logic confusion | Check template conditions in `internal/templates/catalog.go` and resolver logic |
| Generated project won't build | Missing prerequisites | Verify Node/Python/Flutter versions match prerequisites |
| Dry-run shows nothing | All conditions false | Check your config matches expected schema and app/feature flags |
| Tests fail with golden mismatch | Templates changed | Review diffs, then update golden files with `-update` flag if intentional |

### Debugging

Increase verbosity for more details:

```bash
./projectforge generate --profile fullstack --name test --dest ./test -v
```

Validate configs before generating:

```bash
./projectforge validate --config project.config.yaml --strict
```

### Getting Help

- Run `projectforge <command> --help` for command-specific documentation
- Check [docs/architecture.md](docs/architecture.md) for deep architectural details
- Review [contracts/](contracts/) for CLI and schema specifications
- See [CLAUDE.md](CLAUDE.md) or [AGENTS.md](AGENTS.md) for agent/development instructions

## License

MIT License. See the [LICENSE](LICENSE) file for details.

---

**Built with Go. Forged for speed, consistency, and AI-assisted development.**