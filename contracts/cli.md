# CLI Contract: ai-project-starter

**Version**: 1.0.0
**Generated**: 2026-10-12

---

## Global Options

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--config` | `-c` | string | `""` | Path to `project.config.yaml` |
| `--non-interactive` | | bool | `false` | Disable interactive prompts; fail if config incomplete |
| `--dry-run` | | bool | `false` | Show file tree without writing |
| `--force` | `-f` | bool | `false` | Overwrite non-empty destination |
| `--verbose` | `-v` | count | `0` | Increase verbosity (`-v`, `-vv`) |
| `--quiet` | `-q` | bool | `false` | Suppress non-error output |
| `--no-color` | | bool | `false` | Disable colored output |
| `--help` | `-h` | bool | `false` | Show help |
| `--version` | | bool | `false` | Print version and exit |

---

## Commands

### `generate`

Generate a new project from profile, interactive config, or config file.

```
ai-project-starter generate [flags]
```

**Flags**:

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--profile` | `-p` | string | `""` | Profile: `fullstack`, `backend`, `frontend`, `ai-fullstack`, `custom` |
| `--name` | `-n` | string | `""` | Project name (kebab-case) |
| `--dest` | `-d` | string | `"./{name}"` | Destination directory |
| `--interactive` | `-i` | bool | `true`* | Enable interactive configuration wizard |

*Default `true` when no `--config` and no `--profile`; `false` when `--config` or `--profile` provided.

**Behavior**:
1. Load config from `--config` if provided
2. Else if `--profile`: start with profile defaults
3. Else: start with empty config
4. If `--interactive` or no config source: run interactive wizard
5. Validate complete config
6. If `--dry-run`: print file tree, exit 0
7. Check destination (error if non-empty without `--force`)
8. Generate files atomically
9. Print summary

**Exit Codes**:
- `0` — Success
- `1` — Validation error
- `2` — Destination error (non-empty, no `--force`)
- `3` — Template/IO error
- `4` — User cancelled (interactive)

---

### `validate`

Validate a `project.config.yaml` file without generating.

```
ai-project-starter validate [flags]
```

**Flags**:

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--config` | `-c` | string | **required** | Path to config file |
| `--strict` | | bool | `false` | Treat warnings as errors |

**Output**: Human-readable validation result (or JSON with `--output json`).

**Exit Codes**:
- `0` — Valid
- `1` — Invalid (errors found)
- `2` — File not found / parse error

---

### `migrate-config`

Migrate a config file to the current schema version.

```
ai-project-starter migrate-config [flags]
```

**Flags**:

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--config` | `-c` | string | **required** | Path to config file |
| `--output` | `-o` | string | `""` | Output path (default: overwrite in place with backup) |
| `--no-backup` | | bool | `false` | Skip creating `.bak` file |
| `--target-version` | | int | `latest` | Target schema version |
| `--dry-run` | | bool | `false` | Show diff without writing |

**Behavior**:
1. Load config, detect current `schemaVersion`
2. Apply migrations sequentially to target version
3. Show diff (added/changed/removed fields)
4. Write migrated config (with backup unless `--no-backup`)

**Exit Codes**:
- `0` — Success (migrated or already current)
- `1` — Migration error
- `2` — Config parse error
- `3` — No migration path to target version

---

### `list-profiles`

List available profiles with descriptions.

```
ai-project-starter list-profiles [flags]
```

**Flags**:

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `table` | Output format: `table`, `json`, `yaml` |

**Output Example** (table):
```
NAME           DESCRIPTION
fullstack      Next.js web + FastAPI backend
backend        FastAPI backend only
frontend       Next.js web only
ai-fullstack   Fullstack with AI integration
custom         Fully customizable
```

---

### `list-features`

List all configurable feature modules with dependencies.

```
ai-project-starter list-features [flags]
```

**Flags**:

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--profile` | `-p` | string | `""` | Show features for specific profile |
| `--output` | `-o` | string | `table` | Output format: `table`, `json`, `yaml` |

**Output Example** (table):
```
FEATURE              REQUIRES              CONFLICTS
auth:google-oidc     api                   auth:supabase, arch:supabase-first
auth:supabase        arch:supabase-first|hybrid  auth:google-oidc
dashboard            web                   -
login                web                   -
charts               web                   -
graphs               web                   -
ai                   api                   -
workers              api                   -
observability        web|api               -
docker               -                     -
vscode               -                     -
github               -                     -
speckit              -                     -
mcp                  -                     -
testing              -                     -
database:migrations  -                     db:mongodb (unless explicit)
```

---

### `upgrade`

Upgrade a generated project to the latest generator version.

```
ai-project-starter upgrade [flags]
```

**Flags**:

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--config` | `-c` | string | **required** | Path to `project.config.yaml` in generated project |
| `--generator-version` | | string | `latest` | Target generator version |
| `--dry-run` | | bool | `false` | Show diff without applying |
| `--backup` | | bool | `true` | Create backup before upgrading |

**Behavior**:
1. Read config + detect current `generatorVersion`
2. Fetch/check target generator version
3. Run config migrations
4. Regenerate in temp directory
5. Diff against current project
6. Prompt to apply (or auto-apply with `--yes`)

---

### `completion`

Generate shell completion scripts.

```
ai-project-starter completion [shell]
```

**Shells**: `bash`, `zsh`, `fish`, `powershell`

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `AI_PROJECT_STARTER_CONFIG` | Default config file path |
| `AI_PROJECT_STARTER_TEMPLATES_DIR` | Override embedded templates (dev only) |
| `NO_COLOR` | Disable colored output (standard) |
| `TERM` | Terminal type (for interactive UI) |

---

## Output Formats

### Summary (generate success)

```
✓ Project "myapp" generated at ./myapp

Profile:        fullstack
Applications:   web (Next.js), api (FastAPI)
Database:       PostgreSQL (FastAPI architecture)
Features:       auth:[google-oidc], dashboard, login, docker, github, testing
                (12 enabled, 3 disabled)

Structure:
  myapp/
  ├── apps/
  │   ├── web/          # Next.js App Router
  │   └── api/          # FastAPI + SQLAlchemy
  ├── docker-compose.yml
  ├── .github/workflows/ci.yml
  └── AGENTS.md

Next steps:
  cd myapp
  docker compose up -d    # Start database
  cd apps/web && npm run dev
  cd apps/api && uvicorn app.main:app --reload
```

### Dry-Run Tree

```
DRY RUN - would create 147 files:

myapp/
├── .gitignore
├── .env.example
├── README.md
├── AGENTS.md
├── docker-compose.yml
├── .github/
│   └── workflows/
│       └── ci.yml
├── apps/
│   ├── web/
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   ├── next.config.js
│   │   ├── tailwind.config.ts
│   │   ├── src/
│   │   │   ├── app/
│   │   │   ├── components/
│   │   │   └── lib/api.ts
│   │   └── ...
│   └── api/
│       ├── pyproject.toml
│       ├── requirements.lock
│       ├── alembic.ini
│       ├── app/
│       │   ├── models/
│       │   ├── api/
│       │   └── main.py
│       └── tests/
└── ...
```

### Validation Error

```
✗ Configuration validation failed (3 errors):

[ERROR] INVALID_ARCH_APP: Supabase-first architecture cannot have FastAPI backend enabled
  → features.database.architecture=supabase-first conflicts with applications.api=true
  → Fix: Set applications.api=false or change architecture to "fastapi" or "hybrid"

[ERROR] FEATURE_REQUIRES_WEB: dashboard requires web application to be enabled
  → features.dashboard=true but applications.web=false
  → Fix: Enable web application or disable dashboard

[ERROR] INVALID_DESTINATIONS: Mobile bottom navigation requires 3-5 destinations
  → mobile.destinations has 2 items
  → Fix: Add 1-3 more destinations

See 'ai-project-starter validate --config project.config.yaml --strict' for full report.
```