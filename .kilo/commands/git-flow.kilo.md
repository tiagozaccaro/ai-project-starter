---
name: git-flow
description: Standard git flow for feature branches with squash merge
---

# Git Flow Command for Kilo

When the user asks you to create a branch, commit changes, create a PR, and merge it, follow this exact workflow:

## Workflow Steps

### 1. Check git status
```bash
git status
```

### 2. Create feature branch from main
```bash
git checkout -b <branch-name>
```
Use conventional branch naming: `feat/`, `fix/`, `docs/`, `refactor/`, `test/`, `chore/`

### 3. Commit staged changes
```bash
git commit -m "<type>: <description>"
```
Use conventional commits format.

### 4. Push branch to origin
```bash
git push -u origin <branch-name>
```

### 5. Create PR
```bash
gh pr create --title "<type>: <description>" --body "<details>" --base main
```

### 6. Merge with squash and delete branch
```bash
gh pr merge <pr-number> --squash --delete-branch
```

### 7. Cleanup local
```bash
git checkout main
git pull
git fetch --prune
```

## Branch Naming Convention
- `feat/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation only
- `style/` - Formatting, no code change
- `refactor/` - Code restructuring
- `test/` - Adding tests
- `chore/` - Maintenance tasks

## Commit Message Format
```
<type>: <subject>

<body> (optional)
```
Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

## Important Notes
- Repository is configured for squash-only merges
- Always delete branch after merge (local and remote)
- Keep main branch history linear
- Use descriptive PR titles and bodies