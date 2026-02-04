---
name: git-status
description: Explains current repository status including branch, uncommitted changes, and associated PR state. Triggered when the user asks "git status", "branch status", "PR status", "what's the status", "where am I", "current state", or wants to understand the current working context before resuming work.
---

# Git Status

Provide a comprehensive overview of the current repository state, including branch information, local changes, and associated pull request status.

## Workflow

### 1. Gather Git Information

```bash
# Current branch (empty if detached HEAD)
git branch --show-current

# If detached HEAD, get the commit hash
git rev-parse --short HEAD

# Short status (staged, unstaged, untracked)
git status --short

# Stashed changes
git stash list

# Recent commits (last 5)
git log --oneline -5

# Determine base branch
BASE_BRANCH=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@' || echo "main")
echo "Base branch: $BASE_BRANCH"

# Relationship with remote (ahead/behind)
git rev-list --left-right --count origin/$(git branch --show-current)...HEAD 2>/dev/null || echo "No remote tracking"
```

### 2. Compare with Base Branch

```bash
# Use the detected base branch
BASE_BRANCH=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@' || echo "main")

# Commits ahead of base
git rev-list --count origin/$BASE_BRANCH..HEAD 2>/dev/null || echo "0"

# Summary of changes from base
git diff --stat origin/$BASE_BRANCH...HEAD 2>/dev/null

# Check if rebase is needed
git fetch origin --quiet
git merge-base --is-ancestor origin/$BASE_BRANCH HEAD && echo "Up to date with $BASE_BRANCH" || echo "May need rebase"
```

### 3. Check PR Status

```bash
# Get PR information (all relevant fields in one call)
gh pr view --json number,title,state,url,isDraft,reviewDecision,statusCheckRollup,comments,reviews 2>/dev/null || echo "No PR found"

# Get CI/check status
gh pr checks 2>/dev/null || echo "No checks"
```

## Output Format

Present the gathered information in a structured summary:

```
## Repository Status

### Branch
- **Current Branch**: feature/xyz
- **Base Branch**: main
- **Commits Ahead**: 3
- **Commits Behind**: 0

### Local Changes
- **Staged**: 2 files
- **Unstaged**: 1 file
- **Untracked**: 0 files
- **Stashed**: 1 stash

### Recent Commits
1. abc1234 Add feature X
2. def5678 Fix bug in Y
3. ghi9012 Update tests

### Pull Request
- **PR #123**: "Add feature X"
- **Status**: Open (Draft)
- **URL**: https://github.com/owner/repo/pull/123
- **CI Status**: Passing (3/3 checks)
- **Reviews**: 1 approved, 1 changes requested
- **Comments**: 5

### Recommendations
- Consider addressing the requested changes before merging
- Branch is up to date with main
```

## Error Handling

- If not in a git repository, inform the user
- If in detached HEAD state, show the commit hash and explain the state
- If no remote is configured, skip remote-related checks
- If gh CLI is not authenticated, skip PR-related checks and suggest `gh auth login`
- If on the default branch (main/master), note that PR checks are skipped
