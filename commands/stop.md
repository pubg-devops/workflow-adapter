---
description: Stop all running agents or a specific agent
argument-hint: [agent_name]
allowed-tools: Bash, Glob, Read
---

# Stop Agents

Stop running agent processes.

## Arguments

- `$1` (optional): Specific agent to stop (e.g., "alpha", "beta"). If not provided, stops all agents.

## Instructions

### Option 1: Stop All Agents

If no argument provided, stop all running agent scripts:

```bash
# Find and kill all agent scripts
pkill -f "run-.*\.sh" 2>/dev/null

# Also kill any claude processes spawned by agents
pkill -f "claude.*--print" 2>/dev/null
```

### Option 2: Stop Specific Agent

If agent name provided (e.g., `$1` = "alpha"):

```bash
# Kill specific agent script
pkill -f "run-alpha\.sh" 2>/dev/null

# Kill associated claude process
pkill -f "claude.*alpha" 2>/dev/null
```

### Step 2: Verify Processes Stopped

Check if any agent processes are still running:

```bash
ps aux | grep -E "run-.*\.sh|claude.*--print" | grep -v grep
```

### Step 3: Report Status

**If all stopped:**
```
All agent processes stopped.

Stopped:
- Alpha agent
- Beta agent
- Delta agent

To resume work:
  /workflow-adapter:execute {feature_name}
```

**If specific agent stopped:**
```
Alpha agent stopped.

Still running:
- Beta agent (PID: XXXX)
- Delta agent (PID: XXXX)

To stop all: /workflow-adapter:stop
```

**If no processes were running:**
```
No agent processes were running.
```

### Step 4: Clean Up (Optional)

Ask user if they want to clean up log files:

```
로그 파일을 삭제할까요? (logs/*.log)
- yes: 로그 삭제
- no: 로그 유지 (나중에 확인 가능)
```

If yes:
```bash
rm -f logs/*.log
echo "Log files cleaned up."
```

## Quick Commands Reference

```bash
# Stop all agents
pkill -f "run-.*\.sh"

# Stop specific agent
pkill -f "run-alpha\.sh"

# Check running agents
ps aux | grep "run-.*\.sh" | grep -v grep

# View agent logs
tail -f logs/alpha.log
```
