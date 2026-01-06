---
description: Execute agents to work on a feature (generates and runs scripts)
argument-hint: <feature_name> [max_iterations]
allowed-tools: Read, Write, Bash, Glob
---

# Execute Feature

Generate agent scripts and run them in background.

## Arguments

- `$1`: Feature name (required). The feature folder must exist at `doc/feature_{name}/`.
- `$2`: Max iterations per agent (optional, default: 5). Each agent will run up to this many times before stopping.

## Prerequisites

1. Workflow adapter must be installed
2. Feature must be planned: `doc/feature_$1/` must exist with:
   - `$1_spec.md`
   - `$1_plan.md`

If feature doesn't exist, tell user to run `/workflow-adapter:feature $1` first.

## Instructions

### Step 1: Validate Feature Exists

Check that `doc/feature_$1/` folder exists and contains required files.
If not found, tell user to run `/workflow-adapter:feature` first.

### Step 2: Identify Installed Agents

Scan `.claude/commands/` to find installed agents:
- List all .md files except orchestrator.md
- Extract agent names (alpha, beta, delta, etc.)

### Step 3: Read Script Template

Read the script template from:
`${CLAUDE_PLUGIN_ROOT}/scripts/run-agent.sh`

### Step 4: Create Scripts Directory

Create `scripts/` folder in project root if it doesn't exist.

### Step 5: Parse Max Iterations

Get max iterations from `$2`, default to 5 if not provided:
```bash
MAX_ITERATIONS=${2:-5}
```

### Step 6: Generate Agent Scripts

For each agent, create `scripts/run-{agent}.sh`:

Replace in template:
- `{{AGENT_NAME}}` → Capitalized agent name (e.g., "Alpha")
- `{{AGENT_NAME_LOWER}}` → Lowercase agent name (e.g., "alpha")
- `{{MAX_ITERATIONS}}` → Value from $2 or default 5

Make the script executable: `chmod +x scripts/run-{agent}.sh`

### Step 7: Generate Orchestrator Script

Create `scripts/run-orchestrator.sh` with similar replacements.

### Step 8: Run Scripts in Background

Execute all agent scripts in background:

```bash
cd /path/to/project

# Run each agent in background, redirect output to log files
nohup ./scripts/run-alpha.sh $1 > logs/alpha.log 2>&1 &
echo "Alpha agent started (PID: $!)"

nohup ./scripts/run-beta.sh $1 > logs/beta.log 2>&1 &
echo "Beta agent started (PID: $!)"

nohup ./scripts/run-delta.sh $1 > logs/delta.log 2>&1 &
echo "Delta agent started (PID: $!)"
```

Create `logs/` directory if it doesn't exist.

### Step 9: Report Status

Output:
```
Feature agents executing: $1
Max iterations per agent: {max_iterations}

Scripts created:
- scripts/run-alpha.sh
- scripts/run-beta.sh
- scripts/run-delta.sh
- scripts/run-orchestrator.sh

Background processes started:
- Alpha agent (PID: XXXX) → logs/alpha.log (max {max_iterations} iterations)
- Beta agent (PID: XXXX) → logs/beta.log (max {max_iterations} iterations)
- Delta agent (PID: XXXX) → logs/delta.log (max {max_iterations} iterations)

Monitor commands:
- View logs: tail -f logs/alpha.log
- Check processes: ps aux | grep run-
- Stop all: /workflow-adapter:stop

Check progress:
- /workflow-adapter:validate $1
```

## Examples

```bash
# Default: 5 iterations per agent
/workflow-adapter:execute sqs-redrive

# Custom: 10 iterations per agent
/workflow-adapter:execute sqs-redrive 10

# Quick test: 2 iterations
/workflow-adapter:execute sqs-redrive 2
```

## Notes

- Scripts run with `nohup` so they continue after terminal closes
- Each agent logs to `logs/{agent}.log`
- Max 3 iterations per agent (configurable in script)
- Agents stop when they report "Status: completed" in messages folder
- Use `pkill -f "run-.*\.sh"` to stop all agents
