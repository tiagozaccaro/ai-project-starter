# Feature Specification: ai-project-starter

**Feature Branch**: `001-ai-project-starter`

**Created**: 2026-10-12

**Status**: Draft

**Input**: User description: "Build a reusable, configurable project generator named `ai-project-starter`. Its purpose is to generate consistent, maintainable, AI-friendly software repositories for new applications. It is a developer tool, not an application-specific business product."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Generate a new project with a predefined profile (Priority: P1)

A developer wants to quickly scaffold a new fullstack application using the `fullstack` profile (Next.js + FastAPI) with minimal configuration.

**Why this priority**: This is the primary use case - developers need to generate working projects quickly using sensible defaults.

**Independent Test**: Can be fully tested by running the CLI with `--profile fullstack --name myapp --non-interactive` and verifying the generated project structure, dependencies, and that it builds/runs.

**Acceptance Scenarios**:

1. **Given** the CLI is installed, **When** the user runs `ai-project-starter generate --profile fullstack --name myapp --dest ./myapp --non-interactive`, **Then** a complete fullstack project is created at `./myapp` with Next.js web app and FastAPI backend.
2. **Given** a valid configuration, **When** the generator runs, **Then** it produces a summary of enabled technologies and generated modules.

### User Story 2 - Customize a project using interactive configuration (Priority: P1)

A developer wants to create a project with specific choices: frontend only, PostgreSQL database, Google OAuth, Docker, and GitHub Actions - but no mobile or AI features.

**Why this priority**: Many developers need fine-grained control over their stack. Interactive mode guides them through choices.

**Independent Test**: Can be tested by running the CLI in interactive mode, making selections, and verifying the output matches exactly what was selected.

**Acceptance Scenarios**:

1. **Given** the CLI is run in interactive mode, **When** the user selects `frontend` profile, enables PostgreSQL, Google OAuth, Docker, and GitHub Actions, **Then** the generated project contains only the selected components.
2. **Given** interactive configuration, **When** the user completes the wizard, **Then** a `project.config.yaml` is saved that can be reused for repeatable generation.

### User Story 3 - Generate project from existing configuration file (Priority: P2)

A team wants to reproduce the exact same project structure across multiple environments (development, staging, CI) using a committed `project.config.yaml`.

**Why this priority**: Repeatability and version-controlled configuration are essential for team workflows.

**Independent Test**: Can be tested by generating a project, committing the config, then regenerating in a clean directory and comparing outputs.

**Acceptance Scenarios**:

1. **Given** a valid `project.config.yaml`, **When** the user runs `ai-project-starter generate --config project.config.yaml --dest ./new-project`, **Then** the generated project matches the original exactly.
2. **Given** the same configuration run twice, **Then** the output is deterministic (identical file contents, modulo timestamps).

### User Story 4 - Safe generation with dry-run and existing destination protection (Priority: P1)

A developer accidentally points the generator at an existing directory and wants to preview changes without data loss.

**Why this priority**: Preventing accidental data loss is critical for a developer tool.

**Independent Test**: Can be tested by running generation against an existing non-empty directory and verifying it fails safely; running with `--dry-run` shows what would be created without writing files.

**Acceptance Scenarios**:

1. **Given** a non-empty destination directory, **When** the user runs generation without `--force`, **Then** the CLI exits with an error and does not modify the directory.
2. **Given** any generation command, **When** the user adds `--dry-run`, **Then** the CLI outputs a tree of files that would be created without writing anything.

### User Story 5 - Validate configuration before generation (Priority: P1)

A developer provides invalid configuration (e.g., MongoDB with SQLAlchemy selected, or Supabase-first with FastAPI backend enabled).

**Why this priority**: Early validation prevents wasted time and confusing errors during generation.

**Independent Test**: Can be tested by providing various invalid configurations and verifying clear, actionable error messages.

**Acceptance Scenarios**:

1. **Given** an invalid configuration (mutually exclusive options), **When** the user runs generation, **Then** the CLI validates and reports specific conflicts with suggested fixes.
2. **Given** a valid configuration, **When** validation runs, **Then** it passes silently and generation proceeds.

### User Story 6 - Generate Flutter mobile app with proper navigation (Priority: P2)

A developer enables the mobile app and needs the specified navigation model: bottom navigation (3-5 tabs), left drawer for features, right drawer for profile/settings, proper Android back behavior.

**Why this priority**: Mobile is a first-class application type with specific UX requirements.

**Independent Test**: Can be tested by generating a mobile-enabled project, running the Flutter app, and verifying navigation behavior matches all 9 requirements.

**Acceptance Scenarios**:

1. **Given** mobile is enabled with 4 configured destinations, **When** the app launches, **Then** bottom navigation shows 4 tabs, left drawer opens via menu button, right drawer opens via avatar tap, drawers never open simultaneously, drawers closed on startup.
2. **Given** the generated Flutter project, **When** widget tests run, **Then** all navigation rules are verified by automated tests.

### User Story 7 - Generated project works independently (Priority: P1)

A generated project must build, test, and run without the generator present.

**Why this priority**: This is a core constitutional requirement - generated projects must be standalone.

**Independent Test**: Can be tested by copying a generated project to a clean machine without the generator and running install/build/test commands.

**Acceptance Scenarios**:

1. **Given** a generated fullstack project, **When** `npm ci && npm run build` runs in the web app, **Then** it succeeds with zero dependencies on the generator.
2. **Given** a generated backend project, **When** `pip install --require-hashes -r requirements.lock && pytest` runs, **Then** it succeeds with all tests passing.
3. **Given** a generated mobile project, **When** `flutter pub get --offline && flutter test` runs, **Then** it succeeds.

---

### Edge Cases

- What happens when the destination directory exists but is empty? (Should proceed with warning)
- What happens when configuration references a non-existent profile? (Clear error with list of valid profiles)
- What happens when database selection conflicts with architecture (e.g., MongoDB with SQLAlchemy)? (Validation error)
- What happens when Supabase-first is selected but FastAPI backend is also enabled? (Validation error - mutually exclusive)
- What happens when user cancels interactive mode mid-way? (Clean exit, no partial generation)
- What happens when a required system dependency is missing (e.g., Docker not installed but Docker enabled)? (Clear error with installation hint)
- What happens when network is unavailable during generation (for fetching templates)? (Graceful error with offline guidance)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The CLI MUST accept a project name, destination path, and configuration (via profile, interactive mode, or config file).
- **FR-002**: The CLI MUST support interactive configuration prompting and non-interactive generation from a config file.
- **FR-003**: The CLI MUST validate configuration before creating any files, reporting all errors at once with actionable messages.
- **FR-004**: The CLI MUST support repeatable generation from the same `project.config.yaml` producing deterministic output.
- **FR-005**: The CLI MUST produce a clear summary of enabled technologies, applications, and generated modules after generation.
- **FR-006**: The CLI MUST never silently overwrite an existing non-empty destination directory without explicit `--force` flag.
- **FR-007**: The CLI MUST support a `--dry-run` mode that outputs the file tree without writing files.
- **FR-008**: The CLI MUST run on macOS and Linux (x86_64 and ARM64).
- **FR-009**: The generator MUST provide these starting profiles: `fullstack`, `backend`, `frontend`, `ai-fullstack`, `custom`.
- **FR-010**: Every profile MUST be customizable - users can override any default setting.
- **FR-011**: The generator MUST allow enabling applications independently: `apps/web` (Next.js), `apps/api` (FastAPI), `apps/mobile` (Flutter).
- **FR-012**: Disabled applications MUST NOT generate any files, dependencies, or configuration.
- **FR-013**: The generator MUST support these database/backend architectures: FastAPI with selected database, Supabase-first, Hybrid FastAPI + Supabase.
- **FR-014**: The generator MUST support these databases: PostgreSQL, SQLite, MySQL, MongoDB, or no database.
- **FR-015**: Database-specific dependencies MUST be correctly chosen: SQLAlchemy 2 + Alembic + Psycopg for PostgreSQL; SQLite-specific config for SQLite; appropriate drivers for MySQL; suitable document implementation for MongoDB.
- **FR-016**: Supabase-first architecture MUST use Supabase-managed PostgreSQL and MUST NOT introduce FastAPI unnecessarily.
- **FR-017**: Hybrid mode MUST clearly define which operations belong to FastAPI and which use Supabase services.
- **FR-018**: The generator MUST NOT install relational ORM or migration tools for MongoDB unless explicitly required by a selected component.
- **FR-019**: The generator MUST make these modules independently configurable: Google OAuth 2.0/OIDC login, database and migrations, shadcn/ui dashboard and login blocks, charts (shadcn/Recharts), interactive node graphs (@xyflow/react), AI provider abstraction and AI gateway, background workers and queue infrastructure, structured logging/Sentry/OpenTelemetry, Docker and Docker Compose, VS Code debugging and tasks, GitHub Actions/issue templates/PR templates/Dependabot, GitHub Spec Kit and Kilo Code instructions, optional MCP integrations, frontend/backend/mobile testing infrastructure.
- **FR-020**: Disabled features MUST NOT leave behind dependencies, environment variables, configuration, or documentation implying they are enabled.
- **FR-021**: The web application MUST use a responsive design system built from free shadcn/ui components.
- **FR-022**: When dashboard/login blocks are enabled, the generator MUST add the official `dashboard-01` and `login-04` blocks as starting scaffolds.
- **FR-023**: The login block MUST adapt to the selected authentication provider (not introduce an unrelated auth mechanism).
- **FR-024**: When charts are enabled, a chart component MUST be included.
- **FR-025**: When FastAPI backend and web frontend are both enabled, TypeScript API clients MUST be generated from FastAPI's OpenAPI schema.
- **FR-026**: When Flutter is enabled, the generator MUST implement the required navigation model: bottom navigation (3-5 configurable primary destinations), left drawer for features (opened via menu control), right drawer for profile/settings (opened via avatar tap), never both drawers simultaneously, drawers closed on startup, tab switching preserves state, Android back behavior and swipe gestures follow platform conventions, adaptive navigation on larger screens, no hardcoded illustrative destinations.
- **FR-027**: The generator MUST create reusable Flutter navigation components and automated widget tests for navigation rules.
- **FR-028**: Generated projects MUST include: root documentation and architecture overview, scoped `AGENTS.md` files, Git configuration and safe branch conventions, Spec Kit integration when enabled, VS Code launch configurations for web/backend/tests/full-stack, Docker development support when enabled, CI quality gates appropriate to the stack, `.env.example`, appropriate `.gitignore` files, dependency lockfiles, unit and integration tests.
- **FR-029**: The CLI MUST provide useful error messages for all failure modes.
- **FR-030**: The generator MUST be installable via standard package managers (npm, pipx, homebrew) and runnable as a standalone binary.

### Key Entities

- **Project Configuration**: The `project.config.yaml` file that serves as the single source of truth for a generated project. Contains profile, enabled applications, database choice, architecture mode, and all optional feature flags.
- **Profile**: A named preset (`fullstack`, `backend`, `frontend`, `ai-fullstack`, `custom`) that provides sensible defaults for a project type.
- **Application**: An independently enableable component: `web` (Next.js), `api` (FastAPI), `mobile` (Flutter).
- **Database Configuration**: The selected database (PostgreSQL, SQLite, MySQL, MongoDB, none) and architecture mode (FastAPI, Supabase-first, Hybrid).
- **Feature Module**: An independently configurable optional capability (auth, database, UI blocks, charts, AI, workers, observability, Docker, VS Code, GitHub, Spec Kit, MCP, testing).
- **Generated Project**: The output repository structure containing all selected applications, infrastructure, and tooling, ready to build and run independently.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can generate a working fullstack project (Next.js + FastAPI) in under 2 minutes from CLI invocation to running application.
- **SC-002**: Users can generate a project from a saved configuration file in under 30 seconds.
- **SC-003**: 100% of generated projects build successfully on first attempt in a clean environment (macOS and Linux).
- **SC-004**: 100% of generated projects pass their test suites on first run in a clean environment.
- **SC-005**: Zero dependencies, files, or configuration artifacts exist in generated output that are not explicitly enabled in the configuration.
- **SC-006**: Configuration validation catches 100% of invalid combinations (mutually exclusive options, missing required dependencies) before file generation.
- **SC-007**: Dry-run mode accurately predicts 100% of files that would be generated.
- **SC-008**: Existing non-empty destination directories are never modified without explicit `--force` (0 data loss incidents).
- **SC-009**: Generated Flutter apps pass 100% of navigation widget tests covering all 9 navigation rules.
- **SC-010**: Generated projects achieve ≥80% test coverage (backend Python, frontend TypeScript, mobile Dart).
- **SC-011**: Generated projects include valid lockfiles for all package managers used (npm/pnpm, pip/uv, flutter pub).
- **SC-012**: CI pipelines in generated projects pass on both macOS-latest and ubuntu-latest GitHub runners.
- **SC-013**: Generated web code passes strict type checking with zero unsafe types without justification.
- **SC-014**: Generated backend code passes strict static type analysis.
- **SC-015**: The generator CLI itself has ≥90% test coverage and passes all tests on macOS and Linux.

## Assumptions

- Target users are software developers familiar with CLI tools, Git, and modern web/mobile development stacks.
- Users have Node.js ≥20, Python ≥3.11, and Flutter ≥3.19 installed for the respective stacks they generate.
- Users have Git installed and configured.
- Users have Docker installed if they enable Docker features (generator will validate and error clearly if not).
- Network access is available during generation for fetching template assets (with graceful offline error handling).
- The generator is implemented in **Go** and distributed as a standalone binary via npm, pipx, and Homebrew - users don't need to clone the generator source.
- Generated projects target modern evergreen browsers (last 2 versions) and mobile OS versions (iOS 16+, Android 12+).
- Free/open-source tools are preferred; paid services (Supabase hosted, Sentry, AI providers) are optional and behind config flags.
- The initial release focuses on a reliable generator with well-tested templates, not a hosted SaaS or GUI configuration tool.

## Clarifications

### Session 2026-10-12

- Q: What programming language should the `ai-project-starter` generator itself be implemented in, given it must be distributable as a standalone binary via npm, pipx, and Homebrew? → A: Go
- Q: How should the `project.config.yaml` schema handle versioning and migrations when the generator releases new versions with changed/removed/added configuration options? → A: Explicit version field + automatic migration with opt-out
- Q: When a user selects a profile (e.g., `fullstack`) and then overrides individual features in interactive mode or config file, should the profile be treated as a **starting preset** (copied values that can be independently changed) or a **live constraint** (changing the profile later resets overridden features to profile defaults)? → A: Starting preset (copy-on-select)
- Q: For projects enabling both web (`apps/web`) and API (`apps/api`) with Google OAuth: should authentication be centralized in **FastAPI** (backend issues JWTs, web uses HTTP-only cookies) or **Supabase Auth** (Supabase manages auth, both web and API validate Supabase JWTs), given these are mutually exclusive architectures? → A: Exclusive per project
- Q: How should database migrations be handled for each architecture, specifically: (a) FastAPI + PostgreSQL uses Alembic, but what generates the initial migration? (b) Supabase-first: Supabase CLI manages migrations - does the generator create Supabase migration files or expect users to use Supabase dashboard? (c) Hybrid: which layer owns migrations for shared tables? (d) MongoDB: no relational migrations, but schema validation/indexes needed? → A: Architecture-specific starters