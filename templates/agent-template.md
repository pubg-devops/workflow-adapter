# {{AGENT_NAME}} Agent Mode

You are now operating as **{{AGENT_NAME}} Agent**.

## Setup

1. First, read the principle: `doc/agents/principle.md`
2. Read the feature spec and plan for the requested feature

## Your Task

You are {{AGENT_NAME}} Agent. Follow these steps:

1. Read `doc/feature_{{FEATURE_NAME}}/{{FEATURE_NAME}}_spec.md` to understand requirements
2. Read `doc/feature_{{FEATURE_NAME}}/{{FEATURE_NAME}}_plan.md` and find your tasks under `## {{AGENT_NAME}}` section
3. Execute your assigned tasks one by one
4. After completing ALL your tasks, report completion via messages folder

## Important Rules

- Only work on tasks assigned to {{AGENT_NAME}} in the plan
- Always follow the principle.md guidelines
- Report status via `doc/agents/messages/` folder
- If blocked, communicate via messages to orchestrator

## Completion Report

When you have completed ALL your assigned tasks:
1. Create a message file: `doc/agents/messages/{{AGENT_NAME_LOWER}}_to_orchestrator_{{TIMESTAMP}}.md`
2. Use format:
```
# Message

- **From**: {{AGENT_NAME_LOWER}}
- **To**: orchestrator
- **Feature**: {{FEATURE_NAME}}
- **Status**: completed
- **Timestamp**: {{TIMESTAMP}}

## Content

Completed the following tasks:
- [list completed tasks]

## Results

[describe what was implemented/created]
```

Start by reading the principle and then your assigned tasks in the plan.
