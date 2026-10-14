# Quickstart Validation Guide: project-forge

**Feature**: specs/001-project-forge/spec.md
**Date**: 2026-10-12

---

## Prerequisites

### System Requirements
- **Go 1.22+** (for building the generator from source)
- **Node.js ≥20** (for generated web projects)
- **Python ≥3.11** (for generated API projects)
- **Flutter ≥3.19** (for generated mobile projects)
- **Docker & Docker Compose** (optional, for Docker-enabled projects)
- **Git** (for version control)
- **GitHub CLI** (`gh`) (optional, for CI validation)

### Generator Installation

**Prerequisites**: Go 1.22+, Node.js ≥20 (for generated web projects), Python ≥3.11 + `uv` (for generated API projects), Flutter ≥3.19 (for generated mobile projects).

**From source**:
```bash
git clone https://github.com/tiagozaccaro/project-forge
cd project-forge
go build -o project-forge ./cmd/project-forge
./project-forge --version
```

**Via package managers** (after release):
```bash
# Binary download (recommended) - static binary, no runtime deps
curl -fsSL https://github.com/tiagozaccaro/project-forge/releases/latest/download/project-forge-linux-amd64.tar.gz | tar xz
sudo mv project-forge /usr/local/bin/

# Homebrew
brew tap tiagozaccaro/tap
brew install project-forge

# npm wrapper (downloads precompiled binary)
npm install -g @project-forge/cli
```

---

## Validation Scenarios

Each scenario validates a specific user story from the specification. Run in order.

---

### Scenario 1: Fullstack Profile Generation (P1)

**User Story**: Generate a new project with predefined `fullstack` profile.

**Prerequisites**: Generator installed, Docker running.

```bash
# Generate project
./project-forge generate --profile fullstack --name myapp --dest ./test-fullstack --non-interactive

# Verify structure
ls -la test-fullstack/
# Should show: apps/, docker-compose.yml, .github/, AGENTS.md, README.md, .env.example

# Verify web app
cd test-fullstack/apps/web
pnpm install --frozen-lockfile
pnpm run build          # Must succeed
pnpm run typecheck      # Must pass (TS strict)
pnpm run lint           # Must pass

# Verify API app
cd ../api
uv sync
uv run pytest           # Must pass, ≥80% coverage (tests use SQLite in-memory)
uv run pyright          # Must pass (strict)

# Verify Docker
cd ../..
docker compose build   # Must succeed
docker compose up -d
sleep 10
curl -f http://localhost:8000/health  # API health
curl -f http://localhost/           # Web (nginx)
docker compose down
```

**Expected Outcome**: All commands succeed. Project runs independently.

---

### Scenario 2: Interactive Custom Configuration (P1)

**User Story**: Customize project via interactive wizard.

**Prerequisites**: Generator installed.

```bash
# Run interactive (will prompt for each choice)
./project-forge generate --name custom-app --dest ./test-custom

# At prompts, select:
# - Profile: custom
# - Applications: web only
# - Database: postgresql, architecture: fastapi
# - Features: auth=google-oidc, dashboard=true, login=true, docker=true, github=true
# - All others: false

# Verify generated config
cat test-custom/project.config.yaml
# Should reflect all choices exactly

# Verify only web app generated
ls test-custom/apps/
# Should show ONLY: web/

# Verify features present
grep -r "google-oidc" test-custom/apps/web/
grep -r "dashboard-01" test-custom/apps/web/
grep -r "login-04" test-custom/apps/web/
ls test-custom/.github/workflows/
ls test-custom/docker-compose.yml
```

**Expected Outcome**: Config matches selections exactly. No api/ or mobile/ directories. Only selected features present.

---

### Scenario 3: Config File Generation (P2)

**User Story**: Generate from committed config file for repeatability.

**Prerequisites**: Generator installed, config file from Scenario 2.

```bash
# Copy config to new location
cp test-custom/project.config.yaml /tmp/my-config.yaml

# Generate in clean directory
./project-forge generate --config /tmp/my-config.yaml --dest ./test-repeat

# Compare outputs (excluding lockfiles and timestamps)
diff -r test-custom/ test-repeat/ --exclude='*',*' --exclude='node_modules' --exclude='.venv' --exclude='__pycache__' --exclude='*.pyc' --exclude='pnpm-lock.yaml' --exclude='uv.lock' --exclude='pubspec.lock'
# Should show NO differences (deterministic)

# Verify both build
cd test-repeat/apps/web && pnpm install --frozen-lockfile && pnpm run build
cd ../api && uv sync && uv run pytest
```

**Expected Outcome**: Zero diff (modulo lockfiles/timestamps). Both projects build and test identically.

---

### Scenario 4: Dry-Run & Destination Protection (P1)

**User Story**: Safe generation with dry-run and existing directory protection.

```bash
# Create non-empty directory
mkdir -p ./test-existing/existing-file.txt
echo "important" > ./test-existing/existing-file.txt

# Attempt generation without --force (should fail)
./project-forge generate --profile frontend --name test --dest ./test-existing --non-interactive
# Expected: Exit code 2, error message, existing-file.txt UNCHANGED

# Verify file untouched
cat ./test-existing/existing-file.txt
# Should still show "important"

# Dry-run (should show tree, write nothing)
./project-forge generate --profile frontend --name test --dest ./test-dryrun --dry-run
# Expected: Prints file tree, exit 0, NO files created

# Verify nothing created
ls ./test-dryrun 2>/dev/null || echo "Directory does not exist (correct)"
```

**Expected Outcome**: Non-empty directory protected. Dry-run shows tree without writing.

---

### Scenario 5: Configuration Validation (P1)

**User Story**: Invalid configurations caught before generation.

```bash
# Create invalid config (Supabase-first + API enabled)
cat > /tmp/invalid-config.yaml <<'EOF'
schemaVersion: 2
generatorVersion: "1.0.0"
profile: "custom"
projectName: "invalid-test2"
applications:
  web: true
  api: true
database:
  type: "postgresql"
  architecture: "supabase-first"
features:
  testing: true
EOF

# Validate (should fail)
./project-forge validate --config /tmp/invalid-config.yaml --strict
# Expected: Exit code 1, error about Supabase-first + API conflict

# Create config with MongoDB + migrations (auto-corrected, not error)
cat > /tmp/mongodb-config.yaml <<'EOF'
schemaVersion: 2
generatorVersion: "1.0.0"
profile: "custom"
projectName: "mongodb-test"
applications:
  api: true
database:
  type: "mongodb"
  architecture: "fastapi"
features:
  database.migrations: true  # Auto-forced to false by resolver (H2 decision)
  testing: true
EOF

./project-forge validate --config /tmp/mongodb-config.yaml
# Expected: Passes (resolves). Generation summary notes migrations was auto-forced false.
```

**Expected Outcome**: Clear, actionable error messages for each invalid combination.

---

### Scenario 6: Flutter Mobile Navigation (P2)

**User Story**: Generated Flutter app implements 9 navigation rules.

**Prerequisites**: Generator installed, Flutter ≥3.19.

```bash
# Generate with mobile
cat > /tmp/mobile-config.yaml <<'EOF'
schemaVersion: 2
generatorVersion: "1.0.0"
profile: "custom"
projectName: "mobile-test"
applications:
  mobile: true
database:
  type: "none"
  architecture: "fastapi"
features:
  testing: true
metadata:
  mobile:
    bundleId: "com.test.mobiletest"
    destinations:
      - id: "home"
        label: "Home"
        icon: "home"
        route: "/home"
      - id: "explore"
        label: "Explore"
        icon: "explore"
        route: "/explore"
      - id: "profile"
        label: "Profile"
        icon: "person"
        route: "/profile"
EOF

./project-forge generate --config /tmp/mobile-config.yaml --dest ./test-mobile

# Run widget tests
cd test-mobile/apps/mobile
flutter pub get --offline
flutter test test/navigation_test.dart
# Expected: All 9 navigation rule tests pass

# Manual verification (run app)
flutter run -d chrome  # or device
# Verify:
# 1. Bottom nav shows 3 tabs (Home, Explore, Profile)
# 2. Left drawer opens via menu button (top-left)
# 3. Right drawer opens via avatar tap (top-right)
# 4. Never both drawers open simultaneously
# 5. Drawers closed on startup
# 6. Tab switch preserves state
# 7. Android back button works
# 8. Adaptive layout on wide screen (NavigationRail)
# 9. No hardcoded labels (uses config)
```

**Expected Outcome**: All 9 navigation rules verified by automated tests and manual check.

---

### Scenario 7: Generated Project Independence (P1)

**User Story**: Generated project works without generator present.

**Prerequisites**: Clean machine/VM without generator, Docker.

```bash
# On clean machine, clone generated project (from Scenario 1)
git clone <generated-project-repo> test-independent
cd test-independent

# No generator installed - verify it's not in PATH
which project-forge || echo "Not found (correct)"

# Web app
cd apps/web
pnpm install --frozen-lockfile
pnpm run build
pnpm run test

# API app (if enabled)
cd ../api
uv sync
uv run pytest

# Mobile app (if enabled)
cd ../mobile
flutter pub get --offline
flutter test

# Full stack with Docker
cd ../..
docker compose up -d
sleep 15
curl -f http://localhost/health
docker compose down
```

**Expected Outcome**: All build, test, and run commands succeed without generator.

---

### Scenario 8: Profile Variations

Test each profile produces expected structure.

```bash
# Backend only
./project-forge generate --profile backend --name backend-only --dest ./test-backend --non-interactive
ls test-backend/apps/  # Should show ONLY: api/

# Frontend only
./project-forge generate --profile frontend --name frontend-only --dest ./test-frontend --non-interactive
ls test-frontend/apps/  # Should show ONLY: web/

# AI Fullstack
./project-forge generate --profile ai-fullstack --name ai-stack --dest ./test-ai --non-interactive
# Should have web + api + AI features
grep -r "ai.*gateway" test-ai/apps/api/
```

---

### Scenario 9: Database Architecture Variations

Each architecture produces different migration and backend structure. Use complete config files (not partial overrides) to avoid `--config`/`--profile` precedence issues (see CLI contract).

```bash
# SQLite (FastAPI) - write complete config file
cat > /tmp/sqlite-config.yaml <<'EOF'
schemaVersion: 2
generatorVersion: "1.0.0"
profile: "backend"
projectName: "sqlite-app"
applications:
  api: true
database:
  type: "sqlite"
  architecture: "fastapi"
features:
  testing: true
EOF
./project-forge generate --config /tmp/sqlite-config.yaml --dest ./test-sqlite --non-interactive

# Supabase-first
cat > /tmp/supabase-config.yaml <<'EOF'
schemaVersion: 2
generatorVersion: "1.0.0"
profile: "frontend"
projectName: "supabase-app"
applications:
  web: true
  api: false
database:
  type: "postgresql"
  architecture: "supabase-first"
features:
  auth:
    - "supabase"
  testing: true
EOF
./project-forge generate --config /tmp/supabase-config.yaml --dest ./test-supabase --non-interactive

# Hybrid
cat > /tmp/hybrid-config.yaml <<'EOF'
schemaVersion: 2
generatorVersion: "1.0.0"
profile: "fullstack"
projectName: "hybrid-app"
applications:
  web: true
  api: true
database:
  type: "postgresql"
  architecture: "hybrid"
features:
  auth:
    - "supabase"
  testing: true
EOF
./project-forge generate --config /tmp/hybrid-config.yaml --dest ./test-hybrid --non-interactive
```

**Verify each**:
- SQLite: `apps/api/alembic/` present, `pyproject.toml` has `aiosqlite` deps
- Supabase-first: `supabase/migrations/` present, NO `apps/api/`, web uses Supabase client
- Hybrid: BOTH `supabase/migrations/` AND `apps/api/alembic/` present

---

### Scenario 10: CI Pipeline Validation

```bash
# Use generated project from Scenario 1
cd test-fullstack

# Simulate GitHub Actions locally with act (optional)
# Or push to GitHub and watch Actions tab

# Verify workflow file exists and is valid YAML
cat .github/workflows/ci.yml | yq eval '.'

# Check matrix includes only enabled apps
# Check steps: lint, typecheck, test, build for each
```

---

## Validation Checklist Summary

| Scenario | User Story | Key Assertions |
|----------|------------|----------------|
| 1 | Fullstack profile | All apps build, test, run; Docker works |
| 2 | Interactive config | Config matches selections; no extra files |
| 3 | Config file repeatability | Zero diff between generations |
| 4 | Dry-run & protection | Non-empty dir protected; dry-run writes nothing |
| 5 | Config validation | Clear errors for invalid combos |
| 6 | Flutter navigation | 9 rules pass widget tests + manual |
| 7 | Project independence | Builds/tests/runs without generator |
| 8 | Profile variations | Each profile produces correct app set |
| 9 | DB architectures | Correct migration setup per architecture |
| 10 | CI pipeline | Workflow valid, matrix matches config |

---

## Troubleshooting

| Issue | Resolution |
|-------|------------|
| `go build` fails | Ensure Go 1.22+, run `go mod tidy && go build ./cmd/project-forge` |
| `pnpm install --frozen-lockfile` fails | Check Node ≥20, delete `node_modules` and `pnpm-lock.yaml`, retry |
| `uv sync` (generated API) fails | Ensure Python 3.11+, run `uv sync` from `apps/api/` |
| `flutter pub get --offline` fails | Run `flutter pub get` online once, then offline works |
| Docker build fails | Check Dockerfile syntax, ensure `docker compose build` context correct |
| Validation false positive | Check `project.config.yaml` against schema in `contracts/config-schema.md` |
| Generation too slow (>2min) | Check for network calls during template rendering; templates should be static text |

---

## Performance Benchmarks

Run after implementation to verify success criteria:

```bash
# SC-001: Fullstack generation < 2 min
time ./project-forge generate --profile fullstack --name bench --dest ./bench --non-interactive

# SC-002: Config file generation < 30 sec
time ./project-forge generate --config bench/project.config.yaml --dest ./bench2

# SC-015: Generator test coverage
cd <generator-repo>
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
# Must show ≥90%
```

---

## Next Steps

After all scenarios pass:
1. Run `/speckit.tasks` to generate implementation task list
2. Begin implementation following `tasks.md`
3. Re-run these validations after each major milestone