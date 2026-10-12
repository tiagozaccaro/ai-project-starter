# ai-project-starter

A reusable, configurable project generator that scaffolds consistent, maintainable, AI-friendly software repositories. It is a developer tool, not an application-specific product.

`ai-project-starter` produces a monorepo with independently buildable applications (`apps/web`, `apps/api`, `apps/mobile`), a versioned configuration file (`project.config.yaml`), lockfiles, CI workflows, editor debugging configs, and agent instructions — all driven by a single, validated configuration model.

## Features

- **Profiles**: `fullstack`, `backend`, `frontend`, `ai-fullstack`, `custom` presets.
- **Independent applications**: Next.js (web), FastAPI (api), Flutter (mobile) — each builds and tests on its own.
- **Database & architecture**: PostgreSQL, SQLite, MySQL, MongoDB, or none; `fastapi`, `supabase-first`, or `hybrid` architecture.
- **Feature modules**: auth (Google OIDC / Supabase), dashboard, login, charts, graphs, AI, workers, observability, Docker, VS Code, GitHub Actions, Speckit, MCP.
- **Safe generation**: dry-run mode, non-empty destination protection with backup, atomic file writes, conflict detection.
- **Repeatability**: generate from a committed `project.config.yaml`; deterministic output.
- **Validation**: full config validation with actionable, aggregated error messages before any file is written.
- **AI-friendly**: generated `AGENTS.md` files and Speckit specs so coding agents understand the project.

## Installation

**Prerequisites**: Go 1.22+ (to build from source). Generated projects need Node.js ≥20 (web), Python ≥3.11 + [uv](https://docs.astral.sh/uv/) (api), Flutter ≥3.19 (mobile).

**From source**:

```bash
git clone https://github.com/your-org/ai-project-starter
cd ai-project-starter
go build -o ai-project-starter ./cmd/ai-project-starter
./ai-project-starter --version
```

## Quickstart

Generate a fullstack project (Next.js + FastAPI) non-interactively:

```bash
./ai-project-starter generate --profile fullstack --name my-app --dest ./my-app --non-interactive
```

Then build and test each app independently:

```bash
# Web (Next.js)
cd my-app/apps/web
pnpm install && pnpm run build && pnpm run typecheck && pnpm test

# API (FastAPI)
cd ../api
uv sync && uv run pytest && uv run ruff check . && uv run pyright
```

### Interactive mode

Run without flags to launch the guided wizard:

```bash
./ai-project-starter generate --name my-app
```

### Generate from a config file

```bash
./ai-project-starter generate --config project.config.yaml --dest ./my-app
```

### Dry run

Preview the file tree without writing anything:

```bash
./ai-project-starter generate --profile fullstack --name my-app --dry-run
```

## Commands

| Command | Description |
|---------|-------------|
| `generate` | Generate a project from a profile, wizard, or config file |
| `validate` | Validate a `project.config.yaml` without generating |
| `migrate-config` | Migrate a config file to the current schema version |
| `list-profiles` | List available profiles |
| `list-features` | List configurable feature modules and their dependencies |

Run `ai-project-starter <command> --help` for full flag documentation. See [contracts/cli.md](contracts/cli.md) for the complete CLI contract and [contracts/config-schema.md](contracts/config-schema.md) for the configuration schema.

## Configuration

Every generated project contains a `project.config.yaml` — the single source of truth. Example:

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
```

## Development

```bash
# Build
go build -o ai-project-starter ./cmd/ai-project-starter

# Test
go test ./...

# Lint / vet / format
go vet ./...
gofmt -l .
```

See [CLAUDE.md](CLAUDE.md) for agent instructions and [docs/architecture.md](docs/architecture.md) for the architecture overview.

## Project structure (generator)

```
cmd/ai-project-starter/   # Entry point
internal/
  cli/                    # Cobra command wiring
  config/                 # Configuration model, validation, profiles, migration
  dependencies/           # Dependency resolution (config -> resolved capabilities)
  templates/              # Embedded template catalog + renderer
  filesystem/             # Safe atomic writes, destination protection
  pipeline/               # Generation pipeline orchestration
  manifests/              # Lockfile/manifest generation
  setup/                  # Post-generation verification
  summary/                # Output formatting
internal/assets/          # Embedded templates (//go:embed)
contracts/                # CLI and config-schema contracts
specs/                    # Speckit feature specifications
```

## License

MIT
