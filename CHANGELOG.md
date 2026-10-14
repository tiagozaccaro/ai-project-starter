# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Platform parity templates for various languages and frameworks (#245)
- shadcn UI components and AI chat integrations (#244)
- Completion command and converge fixes (#236)
- Git-flow commands for kilo workflow (#212)
- GitHub extension for Specify with task-to-issue sync (#211)

### Changed
- Git fetch --prune added to git-flow cleanup step (#213)

### Fixed
- Various converge fixes (#236)

## [1.0.0] - 2026-10-14

### Added
- Scaffold AI project starter CLI with template pipeline and specs (#214)
- Initial project structure with Next.js, FastAPI, and Flutter templates
- Configuration validation and migration system
- Dependency resolver for feature selection
- Atomic file writer with dry-run support
- Post-generation integrity validation
- Interactive wizard for project generation
- Multiple profile support (web, api, mobile, fullstack)
- Feature flags for optional components
- Template catalog with conditional rendering
- Built-in template functions
- Manifest-based dependency management
- GitHub Actions CI/CD pipeline
- Release automation with multi-platform builds
- npm wrapper package publishing

### Changed
- Added .kilo and .specify directories for configuration

---

## Template for Future Releases

## [VERSION] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Vulnerability fixes