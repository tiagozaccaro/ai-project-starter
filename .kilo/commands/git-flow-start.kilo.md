---
name: git-flow-start
description: Start a new feature branch, commit staged changes, and push to origin
---

# Git Flow Start - Create Branch and Push

## When to use
Run this when you have staged changes ready to commit and want to start a new feature branch.

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

## Example
```bash
git checkout -b feat/user-authentication
git commit -m "feat: add user authentication module"
git push -u origin feat/user-authentication
```