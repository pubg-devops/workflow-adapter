---
description: Install multi-agent workflow system with N agents + orchestrator
argument-hint: <number_of_agents>
allowed-tools: Read, Write, Bash, Glob
---

# Install Workflow Adapter

Install a multi-agent workflow system in the current project.

## Arguments

- `$1`: Number of worker agents to create (e.g., 3 creates alpha, beta, gamma + orchestrator)

## Agent Naming Convention

Agents are named using Greek letters in order:
1. alpha
2. beta
3. gamma
4. delta
5. epsilon
6. zeta
7. eta
8. theta
9. iota
10. kappa

## Instructions

You are setting up a multi-agent workflow system. Follow these steps:

### Step 0: Check Git Repository

Check if the current directory is a git repository:

```bash
git rev-parse --is-inside-work-tree 2>/dev/null
```

**If NOT a git repository**, ask the user:

```
이 프로젝트는 Git 저장소가 아닙니다.
Git 저장소를 초기화할까요? (yes/no)

Git 초기화 시:
- git init 실행
- 적절한 .gitignore 생성
- 초기 커밋 생성
```

**If user says yes**, initialize git:

```bash
git init
```

Then create `.gitignore` based on the project type (detect from existing files or ask user):

```gitignore
# Logs
logs/
*.log

# Runtime
*.pid
*.seed

# Dependencies (Node.js)
node_modules/

# Dependencies (Python)
__pycache__/
*.py[cod]
.venv/
venv/
.env

# Dependencies (Go)
vendor/

# Build outputs
dist/
build/
*.exe
*.dll
*.so
*.dylib

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Secrets (IMPORTANT)
.env
.env.local
*.pem
*.key
credentials.json
secrets.yaml

# Agent logs
logs/*.log

# Temporary files
tmp/
temp/
*.tmp
```

After creating .gitignore, make initial commit:

```bash
git add .gitignore
git commit -m "chore: initialize repository with gitignore

🤖 Generated with workflow-adapter"
```

**If user says no**, continue without git (warn that git features will be disabled).

### Step 1: Validate Input

Check that `$1` is a valid number between 1 and 10. If not provided or invalid, ask the user how many agents they want.

### Step 2: Read Templates

Read all template files from the plugin:
- `${CLAUDE_PLUGIN_ROOT}/templates/agent-template.md`
- `${CLAUDE_PLUGIN_ROOT}/templates/orchestrator-template.md`
- `${CLAUDE_PLUGIN_ROOT}/templates/principle-template.md`
- `${CLAUDE_PLUGIN_ROOT}/templates/brainstorming-template.md`
- `${CLAUDE_PLUGIN_ROOT}/templates/spec-template.md`
- `${CLAUDE_PLUGIN_ROOT}/templates/plan-template.md`
- `${CLAUDE_PLUGIN_ROOT}/scripts/run-agent.sh`

### Step 3: Create Directory Structure

Create the following directories:
```
.claude/commands/
doc/agents/messages/
doc/feature_example/
```

### Step 4: Generate Agent Commands

For each agent (based on $1), create `.claude/commands/{agent_name}.md`:

Replace placeholders in agent-template.md:
- `{{AGENT_NAME}}` → Agent name (capitalized, e.g., "Alpha")
- `{{AGENT_NAME_LOWER}}` → Agent name (lowercase, e.g., "alpha")
- `{{FEATURE_NAME}}` → Will be replaced at runtime with `$1` argument
- `{{TIMESTAMP}}` → Will be replaced at runtime with current timestamp

The command file should have this frontmatter:
```yaml
---
description: Run as {Agent Name} agent for a feature
argument-hint: <feature_name>
allowed-tools: Read, Write, Edit, Bash, Glob, Grep, TodoWrite
---
```

### Step 5: Generate Orchestrator Command

Create `.claude/commands/orchestrator.md` using orchestrator-template.md.

Replace `{{AGENT_CHECKLIST}}` with a checklist for all installed agents:
```
- [ ] alpha: All `## Alpha` tasks completed and meet spec
- [ ] beta: All `## Beta` tasks completed and meet spec
...
```

### Step 6: Generate Principle Document

Create `doc/agents/principle.md` using principle-template.md.

Replace placeholders:
- `{{AGENT_LIST}}` → Comma-separated list of agent names (e.g., "alpha, beta, delta")
- `{{AGENT_SECTIONS}}` → List of agent sections (e.g., "- **Alpha agent**: `## Alpha` section")
- `{{AGENT_ROLE_DEFINITIONS}}` → Role definitions for each agent

### Step 7: Generate Feature Example

Create example templates in `doc/feature_example/`:
- `example_brainstorming.md` (from brainstorming-template.md)
- `example_spec.md` (from spec-template.md)
- `example_plan.md` (from plan-template.md)

For `example_plan.md`, generate agent sections:
```markdown
## Alpha

Alpha agent tasks:

- [ ] Task 1: Description
  - Details
  - Expected output

## Beta

Beta agent tasks:
...
```

### Step 8: Create .gitkeep

Create `doc/agents/messages/.gitkeep` to ensure the messages folder is tracked.

### Step 9: Report Success

After completing all steps, report:
1. Number of agents installed
2. List of created files
3. How to use the system:
   - `/workflow-adapter:feature` to plan a new feature
   - `/{agent_name} {feature}` to run as specific agent
   - `/orchestrator {feature}` to verify work

## Example Output

When user runs `/workflow-adapter:install 3`:

```
Workflow adapter installed successfully!

Git:
- Repository: initialized ✓
- .gitignore: created ✓
- Initial commit: created ✓

Created agents: alpha, beta, delta + orchestrator

Files created:
- .gitignore
- .claude/commands/alpha.md
- .claude/commands/beta.md
- .claude/commands/delta.md
- .claude/commands/orchestrator.md
- doc/agents/principle.md
- doc/agents/messages/.gitkeep
- doc/feature_example/example_brainstorming.md
- doc/feature_example/example_spec.md
- doc/feature_example/example_plan.md

Next steps:
1. Run /workflow-adapter:feature to plan a new feature
2. Run /workflow-adapter:execute {feature_name} to start agent work
3. Run /workflow-adapter:validate {feature_name} to check progress
4. Run /workflow-adapter:pr {feature_name} to create pull request

Available commands:
- /workflow-adapter:install   - Install agents (already done)
- /workflow-adapter:feature   - Plan a new feature
- /workflow-adapter:execute   - Run agents on a feature
- /workflow-adapter:validate  - Check progress
- /workflow-adapter:stop      - Stop running agents
- /workflow-adapter:pr        - Create pull request
```
