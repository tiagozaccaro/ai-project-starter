# AI Project Starter Constitution

## Core Principles

### I. Configuration-Driven Generation
A single validated `project.config.yaml` is the source of truth for each generated project. Every optional technology, feature, integration, and service is controlled by configuration. The generator must never produce unused dependencies, empty modules, placeholder features, or services that the project does not need.

**Compliance Requirements:**
- `project.config.yaml` schema validation passes for every generation run.
- Zero dependencies, modules, or files in the generated output that are not explicitly enabled in the configuration.
- Configuration changes produce idempotent, deterministic output.

### II. Standalone Generated Projects
Every generated project must work independently of the generator. No runtime dependency on the generator is permitted after generation. Generated projects include their own documentation, configuration, lockfiles, tests, and development instructions.

**Compliance Requirements:**
- Generated project builds, tests, and runs in a clean environment without the generator present.
- All required lockfiles (`package-lock.json`, `pnpm-lock.yaml`, `requirements.lock`, `pubspec.lock`, etc.) are present and valid.
- `README.md`, `AGENTS.md`, `CONTRIBUTING.md`, and architecture decision records (ADRs) are generated.

### III. AI-First Development
Support GitHub Spec Kit and Kilo Code. Provide a root `AGENTS.md` and appropriately scoped instructions. Give AI agents explicit rules for security, architecture, typing, testing, dependency management, and change control. Agents must not silently change architecture, suppress failing tests, fabricate successful results, or introduce unnecessary dependencies.

**Compliance Requirements:**
- `AGENTS.md` exists at the repository root and in each major component directory.
- Agent instructions cover: security boundaries, architectural invariants, type-safety rules, testing requirements, dependency policies, and change-control procedures.
- CI includes a step that validates agent instruction compliance (e.g., no `any` without justification, no skipped tests).

### IV. Type Safety and Quality
TypeScript strict mode; avoid `any` unless justified with a documented exception. Python type annotations, Ruff, Pyright, pytest, and coverage. Automated formatting, linting, static analysis, and tests. Test generated projects from clean environments.

**Compliance Requirements:**
- TypeScript: `"strict": true`, `noImplicitAny`, `noUncheckedIndexedAccess` enabled; zero `any` without `// @ts-expect-error` + justification comment.
- Python: `pyproject.toml` configures Ruff, Pyright, pytest with coverage ≥ 80%; `mypy` or `pyright` passes in strict mode.
- CI runs format, lint, type-check, and test suites on every PR; failure blocks merge.
- Generated project test suites execute successfully in a fresh container/VM.

### V. Free and Open-Source First
Prefer free, open-source components and tools. Do not require paid UI components or paid hosted services to develop locally. Clearly distinguish open-source software from optional hosted services and paid AI providers.

**Compliance Requirements:**
- Zero mandatory paid dependencies in generated `package.json`, `requirements.txt`, `pubspec.yaml`, or `Cargo.toml`.
- Any optional paid service integration is behind a configuration flag and documented as such.
- License compatibility check (MIT, Apache-2.0, BSD-3-Clause, etc.) passes for all transitive dependencies.

### VI. Security by Default
Never commit secrets, access tokens, private keys, or populated environment files. Use least-privilege OAuth scopes and secure session/token handling. Run containers as non-root where practical. Require explicit authorization for destructive operations, production access, publishing, merging, or force-pushing.

**Compliance Requirements:**
- `.gitignore` excludes `.env*`, `*.pem`, `*.key`, `secrets.*`, `*.kubeconfig`.
- Pre-commit hook runs `trufflehog`/`gitleaks`/`detect-secrets`; failures block commit.
- Dockerfiles use non-root `USER`; `docker-compose.yml` drops capabilities.
- CI/CD pipelines require manual approval for production deploy, publish, merge-to-main, and force-push.

### VII. Incremental Implementation
Use specifications, plans, tasks, acceptance criteria, and tests for substantial changes. Permit lightweight changes without unnecessary full specifications. Never claim completion without running and reporting applicable verification.

**Compliance Requirements:**
- Changes touching ≥ 3 files or modifying public APIs require a spec (`specs/<feature>.md`), plan (`plans/<feature>.md`), and task list (`tasks/<feature>.md`).
- Lightweight changes (< 3 files, no public API change) may proceed with a checklist only.
- Every PR description includes a verification report: test results, lint/type-check status, and manual validation steps performed.

### VIII. Reproducibility and Maintainability
Pin dependencies through supported lockfiles. Validate compatible dependency versions and current official integration instructions. Support macOS and Linux. Document architecture decisions and intentional trade-offs.

**Compliance Requirements:**
- All lockfiles committed; `npm ci` / `pip install --require-hashes` / `flutter pub get --offline` succeed offline.
- Dependabot/Renovate configured with auto-merge for patch/minor updates passing CI.
- ADR directory (`docs/adr/`) exists; each significant decision has an ADR (MADR format).
- CI matrix includes macOS-latest and ubuntu-latest.

### IX. Configurable Platform Parity
When web and mobile are enabled, define their shared business capabilities and acceptance criteria. Differences in presentation are acceptable; missing required functionality is not.

**Compliance Requirements:**
- Shared capability matrix documented in `docs/platform-parity.md` with acceptance criteria per capability.
- E2E test suite covers each capability on both platforms; CI runs both.
- Feature flags in `project.config.yaml` control platform-specific toggles, not core capability omission.

### X. Generated Mobile Navigation Standard
When Flutter mobile is enabled, the bottom navigation bar is the primary navigation for 3–5 configurable destinations. A dedicated app-menu control opens the left drawer for application functionality. Tapping the profile/avatar icon opens the right drawer for user and application settings. Never use the profile drawer as the main application navigation. The drawers are independent, mutually exclusive, and closed initially. Preserve navigation state where appropriate and test back-button/gesture behavior.

**Compliance Requirements:**
- Generated Flutter app implements `Scaffold` with `bottomNavigationBar` (3–5 items from config).
- Left drawer (`Drawer`/`NavigationDrawer`) opened only by app-menu `IconButton`; right drawer (`endDrawer`) opened only by profile `CircleAvatar`/`IconButton`.
- Drawers never open simultaneously; initial state: both closed.
- Integration tests verify: back-button pops routes correctly, gesture navigation preserves state, drawer toggle exclusivity.

## Governance

This constitution supersedes all other practices. Amendments require: (1) a written proposal with rationale, (2) approval by project maintainers, (3) a migration plan for existing generated projects, and (4) version bump in this document.

All PRs and reviews must verify compliance with applicable principles. Complexity must be justified against the simplest working alternative. Use `AGENTS.md` for runtime development guidance.

**Version**: 1.0.0 | **Ratified**: 2026-10-12 | **Last Amended**: 2026-10-12