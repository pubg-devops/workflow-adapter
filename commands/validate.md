---
description: Validate feature progress and calculate completion rate
argument-hint: <feature_name>
allowed-tools: Read, Write, Bash, Glob, Grep
---

# Validate Feature Progress

Calculate completion rate and provide feedback for a feature.

## Arguments

- `$1`: Feature name (required). The feature folder must exist at `doc/feature_{name}/`.

## Prerequisites

1. Feature must be planned: `doc/feature_$1/` must exist
2. Workflow adapter must be installed

## Instructions

### Step 1: Read Feature Documents

Read and analyze:
- `doc/feature_$1/$1_spec.md` - Requirements and acceptance criteria
- `doc/feature_$1/$1_plan.md` - Task checklist for each agent

### Step 2: Identify Agents

Scan `.claude/commands/` to find installed agents.

### Step 3: Analyze Plan

Parse `$1_plan.md` to extract:
- Total tasks per agent
- Task descriptions
- Checkbox status (checked vs unchecked)

### Step 4: Check Messages

Read all messages in `doc/agents/messages/` related to this feature:
- Look for `*_$1_*.md` or messages mentioning feature name
- Extract status (completed, in_progress, blocked)
- Note any issues reported

### Step 5: Calculate Progress

For each agent:
1. Count total tasks in their section
2. Count completed tasks (checked checkboxes `[x]`)
3. Check if completion message exists
4. Calculate percentage: (completed / total) * 100

Overall progress = average of all agents

### Step 6: Verify Against Spec

Check acceptance criteria in spec:
- List each criterion
- Determine if met based on messages and code changes
- Note unmet criteria

### Step 7: Generate Feedback Messages

For agents with incomplete work:
1. Create message: `doc/agents/messages/orchestrator_to_{agent}_{timestamp}.md`
2. Include:
   - Incomplete tasks
   - Issues found
   - Specific action items

### Step 8: Report Results

Output progress report with visual progress bars:

```
╔══════════════════════════════════════════════════════════════╗
║           Feature Progress Report: {feature_name}            ║
╚══════════════════════════════════════════════════════════════╝

Overall Progress
────────────────────────────────────────────────────────────────
[████████████████████░░░░░░░░░░░░░░░░░░░░] 52% (13/25 tasks)


Agent Progress
────────────────────────────────────────────────────────────────

Alpha    [████████████████████████████████████████] 100%  ✅
         8/8 tasks completed
         Status: completed

Beta     [████████████████████░░░░░░░░░░░░░░░░░░░░]  50%  🔄
         4/8 tasks completed
         Status: in_progress
         Remaining:
           - [ ] Implement error handling
           - [ ] Add retry logic
           - [ ] Write unit tests
           - [ ] Update documentation

Delta    [████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  20%  ⏸️
         2/10 tasks completed
         Status: blocked
         Issue: Waiting for Beta to complete API implementation


Task Breakdown
────────────────────────────────────────────────────────────────

## Alpha (8/8) ✅
[x] Set up project structure
[x] Implement CLI framework
[x] Add list command
[x] Add redrive command
[x] Add config management
[x] Handle AWS credentials
[x] Add --dry-run flag
[x] Add --verbose flag

## Beta (4/8) 🔄
[x] Create API client
[x] Implement authentication
[x] Add request/response logging
[x] Handle pagination
[ ] Implement error handling
[ ] Add retry logic
[ ] Write unit tests
[ ] Update documentation

## Delta (2/10) ⏸️
[x] Set up test environment
[x] Create test fixtures
[ ] Write integration tests
[ ] Add CI/CD pipeline
[ ] Create deployment script
[ ] Write README
[ ] Add usage examples
[ ] Create demo video
[ ] Performance testing
[ ] Security audit


Acceptance Criteria
────────────────────────────────────────────────────────────────
[x] CLI can list DLQ messages
[x] CLI can redrive single message
[ ] CLI can redrive all messages
[ ] Dry-run mode works correctly
[ ] Error handling is robust


Git Status
────────────────────────────────────────────────────────────────
Branch: feature/{name}
Commits: 12 commits ahead of main
Last commit: 2 hours ago by Alpha agent


Recommendations
────────────────────────────────────────────────────────────────
1. ⚡ Beta: Focus on error handling and tests
2. ⏳ Delta: Wait for Beta API completion, then proceed
3. 📝 Consider adding more test coverage

Next Actions:
- /workflow-adapter:execute {name}  → Resume agent work
- /workflow-adapter:stop            → Stop all agents
- /workflow-adapter:pr {name}       → Create PR (when 100%)
```

## Progress Bar Rendering

Use these characters for progress bars:

```
Full block:  █ (U+2588)
Empty block: ░ (U+2591)

Width: 40 characters

Calculation:
- filled = (percentage / 100) * 40
- empty = 40 - filled
- bar = "█" * filled + "░" * empty
```

**Examples:**
```
  0% [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]
 25% [██████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]
 50% [████████████████████░░░░░░░░░░░░░░░░░░░░]
 75% [██████████████████████████████░░░░░░░░░░]
100% [████████████████████████████████████████]
```

**Status Icons:**
- ✅ completed (100%)
- 🔄 in_progress (1-99%)
- ⏸️ blocked
- ⏳ pending (0%)

## Message Format for Feedback

```markdown
# Message

- **From**: orchestrator
- **To**: {agent_name}
- **Feature**: $1
- **Status**: feedback
- **Timestamp**: {current_timestamp}

## Content

Progress check for feature "$1".

### Incomplete Tasks

The following tasks in your section are not yet complete:
- [ ] Task description 1
- [ ] Task description 2

### Action Required

Please complete these tasks and update your completion status.

When done, create a message:
`{agent}_to_orchestrator_{timestamp}.md` with Status: completed
```
