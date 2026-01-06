---
description: Run as Alpha agent for a feature
argument-hint: <feature_name>
allowed-tools: Read, Write, Edit, Bash, Glob, Grep, TodoWrite
---

# Alpha Agent Mode

You are now operating as **Alpha Agent**.

## Setup

1. First, read the principle: `doc/agents/principle.md`
2. Read the feature spec and plan for the requested feature

## Your Task

You are Alpha Agent. Follow these steps:

1. Read `doc/feature_$1/$1_spec.md` to understand requirements
2. Read `doc/feature_$1/$1_plan.md` and find your tasks under `## Alpha` section
3. Execute your assigned tasks one by one
4. After completing ALL your tasks, report completion via messages folder

## Important Rules

- Only work on tasks assigned to Alpha in the plan
- Always follow the principle.md guidelines
- Report status via `doc/agents/messages/` folder
- If blocked, communicate via messages to orchestrator

## Completion Report

When you have completed ALL your assigned tasks:
1. Create a message file: `doc/agents/messages/alpha_to_orchestrator_{{TIMESTAMP}}.md`
2. Use format:
```
# Message

- **From**: alpha
- **To**: orchestrator
- **Feature**: $1
- **Status**: completed
- **Timestamp**: {{TIMESTAMP}}

## Content

Completed the following tasks:
- [list completed tasks]

## Results

[describe what was implemented/created]
```

Start by reading the principle and then your assigned tasks in the plan.
