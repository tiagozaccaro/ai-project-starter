---
name: git-flow-finish
description: Create PR, merge with squash, and cleanup local branch
---

# Git Flow Finish - Create PR, Merge, and Cleanup

## When to use
Run this when your feature branch is pushed and ready to be merged to main.

## Workflow Steps

### 1. Create PR
```bash
gh pr create --title "<type>: <description>" --body "<details>" --base main
```

### 2. Merge with squash and delete branch
```bash
gh pr merge <pr-number> --squash --delete-branch
```

### 3. Cleanup local
```bash
git checkout main
git pull
git fetch --prune
```

## PR Title Format
Use the same conventional commit format as your commit:
```
<type>: <description>
```

## PR Body Template
```markdown
## Summary
Brief description of changes

## Changes
- Item 1
- Item 2

## Testing
How to verify the changes
```

## Example
```bash
gh pr create --title "feat: add user authentication module" --body "## Summary\nImplements JWT-based auth with refresh tokens\n\n## Changes\n- Added auth middleware\n- Added login/register endpoints\n- Added token refresh\n\n## Testing\nRun `npm test` and verify auth flow" --base main
gh pr merge 212 --squash --delete-branch
git checkout main && git pull && git fetch --prune
```

## Important Notes
- Repository is configured for squash-only merges
- Always delete branch after merge (local and remote)
- Keep main branch history linear
- Use descriptive PR titles and bodies