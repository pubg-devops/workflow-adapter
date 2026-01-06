---
description: Create a pull request for the completed feature
argument-hint: <feature_name>
allowed-tools: Read, Bash, Glob, Grep
---

# Create Pull Request

Create a pull request for the completed feature work.

## Arguments

- `$1`: Feature name (required). The feature folder must exist at `doc/feature_{name}/`.

## Prerequisites

1. Feature must be planned and worked on
2. Current branch should be `feature/{name}`
3. All agent work should be committed
4. GitHub CLI (`gh`) must be installed

## Instructions

### Step 1: Verify Prerequisites

Check current state:

```bash
# Check current branch
git branch --show-current

# Check if gh is installed
gh --version

# Check git status (should be clean)
git status --porcelain
```

Expected:
- Branch: `feature/{name}`
- gh: installed
- Status: clean (no uncommitted changes)

If uncommitted changes exist, commit them first:
```bash
git add .
git commit -m "chore({name}): final cleanup before PR"
```

### Step 2: Read Feature Documents

Read the feature documents to generate PR description:
- `doc/feature_{name}/{name}_spec.md` - For summary and requirements
- `doc/feature_{name}/{name}_plan.md` - For completed tasks
- Check messages in `doc/agents/messages/` for completion status

### Step 3: Push Branch

Push the feature branch to remote:

```bash
git push -u origin feature/{name}
```

### Step 4: Generate PR Description

Create PR description from feature documents:

```markdown
## Summary

{One-line summary from spec}

## Changes

### Completed Tasks

**Alpha Agent:**
- [x] Task 1
- [x] Task 2

**Beta Agent:**
- [x] Task 1
- [x] Task 2

**Delta Agent:**
- [x] Task 1

### Files Changed

- `src/...` - {description}
- `src/...` - {description}

## Tech Stack

- **Language**: {from spec}
- **Framework**: {from spec}

## Testing

{How to test the changes}

## Checklist

- [ ] Code follows project conventions
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] All agents completed their tasks

---

🤖 Generated with workflow-adapter
```

### Step 5: Create Pull Request

Use GitHub CLI to create the PR:

```bash
gh pr create \
  --title "feat({name}): {short description from spec}" \
  --body "$(cat <<'EOF'
{Generated PR description}
EOF
)" \
  --base main
```

Or if you want to create a draft PR:
```bash
gh pr create --draft ...
```

### Step 6: Report Result

Output:

```
Pull Request created for feature: {name}

PR: https://github.com/{owner}/{repo}/pull/{number}

Title: feat({name}): {short description}
Base: main ← feature/{name}
Status: Open (or Draft)

Next steps:
1. Review the PR on GitHub
2. Request reviews from team members
3. Address any feedback
4. Merge when approved

To view PR: gh pr view
To check CI status: gh pr checks
```

## Alternative: Without GitHub CLI

If `gh` is not installed, provide manual instructions:

```
GitHub CLI가 설치되어 있지 않습니다.

브랜치가 푸시되었습니다: feature/{name}

수동으로 PR 생성하기:
1. https://github.com/{owner}/{repo}/compare/main...feature/{name} 접속
2. "Create pull request" 클릭
3. 아래 내용을 PR 설명에 붙여넣기:

{Generated PR description}
```

## Quick Commands Reference

```bash
# Push branch
git push -u origin feature/{name}

# Create PR
gh pr create --title "feat: ..." --body "..."

# View PR
gh pr view

# Check PR status
gh pr checks

# Merge PR (after approval)
gh pr merge
```
