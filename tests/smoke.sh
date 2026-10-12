#!/usr/bin/env bash
# Generated-app smoke tests for ai-project-starter.
#
# Builds the generator, generates each profile into a temporary
# directory, and runs the generated project's own toolchain
# (build, typecheck, test, lint for web; sync, pytest, ruff,
# pyright for api).
#
# Prerequisites: go, pnpm, uv (Node and Python toolchains are
# provisioned by pnpm/uv).
#
# Usage: ./tests/smoke.sh [profile ...]   # default: all profiles
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BINARY="$REPO_ROOT/ai-project-starter"
WORKDIR="$(mktemp -d)"
trap 'rm -rf "$WORKDIR"' EXIT

PROFILES=("$@")
if [ ${#PROFILES[@]} -eq 0 ]; then
  PROFILES=(fullstack backend frontend)
fi

# Tool availability (checked per-app below).
have() { command -v "$1" >/dev/null 2>&1; }

pass=0
fail=0
skipped=0

log()  { printf '\n\033[1m==> %s\033[0m\n' "$*"; }
ok()   { printf '   \033[32m✓ %s\033[0m\n' "$*"; pass=$((pass+1)); }
bad()  { printf '   \033[31m✗ %s\033[0m\n' "$*"; fail=$((fail+1)); }
skip() { printf '   \033[33m- %s (skipped)\033[0m\n' "$*"; skipped=$((skipped+1)); }

# Build the generator.
log "Building generator"
( cd "$REPO_ROOT" && go build -o "$BINARY" ./cmd/ai-project-starter )
ok "generator built"

run_web_checks() {
  local dir="$1"
  # The generator does not commit lockfiles; the lockfile
  # workflow produces them (pnpm-lock.yaml, uv.lock).
  ( cd "$dir" && ./scripts/generate-lockfiles.sh >/dev/null 2>&1 ) \
    && ok "web: generate lockfiles" || { bad "web: generate lockfiles"; return 1; }
  ( cd "$dir/apps/web" && pnpm install --frozen-lockfile >/dev/null 2>&1 ) \
    && ok "web: pnpm install" || { bad "web: pnpm install"; return 1; }
  local failed=0
  ( cd "$dir/apps/web" && pnpm run build >/dev/null 2>&1 )      && ok "web: build"      || { bad "web: build"; failed=1; }
  ( cd "$dir/apps/web" && pnpm run typecheck >/dev/null 2>&1 )  && ok "web: typecheck"  || { bad "web: typecheck"; failed=1; }
  ( cd "$dir/apps/web" && pnpm run lint >/dev/null 2>&1 )       && ok "web: lint"       || { bad "web: lint"; failed=1; }
  ( cd "$dir/apps/web" && pnpm run test >/dev/null 2>&1 )       && ok "web: test"       || { bad "web: test"; failed=1; }
  return $failed
}

run_api_checks() {
  local dir="$1"
  ( cd "$dir/apps/api" && uv sync >/dev/null 2>&1 ) \
    && ok "api: uv sync" || { bad "api: uv sync"; return 1; }
  local failed=0
  ( cd "$dir/apps/api" && uv run pytest >/dev/null 2>&1 )   && ok "api: pytest"  || { bad "api: pytest"; failed=1; }
  ( cd "$dir/apps/api" && uv run ruff check . >/dev/null 2>&1 ) && ok "api: ruff" || { bad "api: ruff"; failed=1; }
  ( cd "$dir/apps/api" && uv run pyright >/dev/null 2>&1 )  && ok "api: pyright" || { bad "api: pyright"; failed=1; }
  return $failed
}

for profile in "${PROFILES[@]}"; do
  log "Profile: $profile"
  dest="$WORKDIR/$profile"
  if ! "$BINARY" generate --profile "$profile" --name "smoke-$profile" --dest "$dest" --non-interactive >/dev/null 2>&1; then
    bad "$profile: generation failed"
    continue
  fi
  ok "$profile: generated"

  if [ -d "$dest/apps/web" ]; then
    if have pnpm; then run_web_checks "$dest" || true; else skip "web checks (pnpm not installed)"; fi
  fi
  if [ -d "$dest/apps/api" ]; then
    if have uv; then run_api_checks "$dest" || true; else skip "api checks (uv not installed)"; fi
  fi
done

printf '\n\033[1mResults: %d passed, %d failed, %d skipped\033[0m\n' "$pass" "$fail" "$skipped"
[ "$fail" -eq 0 ]
