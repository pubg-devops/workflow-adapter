# Agent Principle

This document defines the common principles that all agents (alpha, beta, delta, gamma and orchestrator) must follow.

---

## 1. Feature-based Workflow

### 1.1 Check Feature Folder
When a user requests work on `{name}` feature:
1. Navigate to `doc/feature_{name}/` folder
2. Check the following files:
   - `{name}_brainstorming.md` - Ideas and discussions
   - `{name}_spec.md` - Feature specification
   - `{name}_plan.md` - Implementation plan with tasks for each agent

### 1.2 Find Your Tasks
In `{name}_plan.md`, find the section assigned to you:
- **Alpha agent**: `## Alpha` section
- **Beta agent**: `## Beta` section
- **Delta agent**: `## Delta` section
- **Gamma agent**: `## Gamma` section
- **Orchestrator**: Review all sections and verify completion

---

## 2. Agent Roles

### Alpha
- Primary implementation agent
- Handles core feature development
- Creates main functionality

### Beta
- Secondary implementation agent
- Handles supporting features
- Creates auxiliary functionality

### Delta
- Testing and quality agent
- Handles test creation
- Verifies implementation quality

### Gamma
- Documentation and integration agent
- Handles documentation updates
- Ensures integration with existing code

### Orchestrator
- Review spec and plan to understand requirements
- Verify each agent's work meets the spec
- Report overall progress and issues

---

## 3. Inter-Agent Communication

### 3.1 Message Location
Create message files in `doc/agents/messages/` folder.

### 3.2 File Naming Convention
```
{from_agent}_to_{to_agent}_{YYYYMMDD}_{HHMM}.md
```

**Examples:**
- `alpha_to_orchestrator_20240115_1430.md`
- `beta_to_delta_20240115_1500.md`
- `orchestrator_to_all_20240115_1600.md` (for announcements)

### 3.3 Message Format
```markdown
# Message

- **From**: {agent_name}
- **To**: {target_agent_name | all}
- **Feature**: {feature_name}
- **Status**: {completed | in_progress | blocked}
- **Timestamp**: {YYYY-MM-DD HH:MM}

## Content

{Write message content here}

## Action Required (optional)

{Actions requested from the recipient}
```

---

## 4. Completion Report

Each agent MUST report completion:
1. Create message in messages folder
2. Set Status to `completed`
3. List completed tasks in Content
4. Describe results/outputs

---

## 5. Handling Issues

If you encounter problems:
1. Create message in messages folder
2. Set Status to `blocked`
3. Describe the issue in Content
4. Specify who should help in To field
