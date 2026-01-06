---
description: Start brainstorming and planning for a new feature
argument-hint: <feature_name>
allowed-tools: Read, Write, Edit, Bash, Glob, Grep, AskUserQuestion
---

# Feature Planning

Create documentation for a new feature through file-based brainstorming.

## Arguments

- `$1`: Feature name (required). Use short kebab-case name (e.g., 'login', 'sqs-redrive', 'user-dashboard')

## Instructions

### Step 0: Understand Project Context

**First, read project documentation to understand the codebase:**

1. Check for and read these files (if they exist):
   ```
   CLAUDE.md
   .claude/CLAUDE.md
   README.md
   agent.md
   .claude/agent.md
   doc/README.md
   ```

2. Scan project structure:
   ```bash
   ls -la
   ls -la src/ 2>/dev/null
   ls -la doc/ 2>/dev/null
   ```

3. Detect project type and tech stack:
   - Check for `package.json` → Node.js/TypeScript
   - Check for `go.mod` → Go
   - Check for `requirements.txt` or `pyproject.toml` → Python
   - Check for `Cargo.toml` → Rust
   - Check for `pom.xml` or `build.gradle` → Java

4. Extract key information:
   - Project name and description
   - Existing tech stack
   - Coding conventions
   - Project structure
   - Existing features/modules

**Store this context for use in brainstorming pre-fill.**

### Step 1: Parse Feature Name

Extract feature name from `$ARGUMENTS`.
- If user provided description along with name, extract the short name
- Convert to kebab-case if needed (e.g., "aws sqs redrive" → "sqs-redrive")
- If unclear, ask user for a short name

### Step 1.5: Create Git Feature Branch

Check if the project is a git repository:
```bash
git rev-parse --is-inside-work-tree 2>/dev/null
```

If it's a git repo:
1. Check current branch and status
2. Create and checkout feature branch:
```bash
git checkout -b feature/{name}
```

Branch naming: `feature/{feature_name}` (e.g., `feature/sqs-redrive`)

If not a git repo, skip this step and continue.

### Step 2: Create Feature Folder

Create `doc/feature_{name}/` directory.

### Step 3: Create Brainstorming File

**IMPORTANT**: Use the project context from Step 0 and user's input to pre-fill the brainstorming file:

1. **Project context** (from CLAUDE.md, README.md, etc.):
   - Existing architecture and patterns
   - Current tech stack (don't suggest different stack unless needed)
   - Coding conventions to follow
   - Related existing features

2. **User's description** (`$ARGUMENTS`):
   - What they want to build
   - Any specific requirements mentioned

3. **Your knowledge**:
   - Best practices for this type of feature
   - Common pitfalls to avoid
   - Suggested implementation approach

4. **Integration points**:
   - How this feature connects to existing code
   - Which modules/files will be affected

Create `doc/feature_{name}/{name}_brainstorming.md`:

```markdown
# {Feature Name} - Brainstorming

> 아래 내용을 검토하고 수정/추가해주세요. 맞으면 그대로 두셔도 됩니다.

## Project Context

{PRE-FILL: Information gathered from CLAUDE.md, README.md, project structure}

- **Project**: {project name from README or package.json}
- **Tech Stack**: {detected from project files}
- **Architecture**: {from CLAUDE.md or inferred}
- **Related Modules**: {existing code that relates to this feature}

## Overview

{PRE-FILL: Based on user's description, write 1-2 sentences about the feature}

**검토/수정:**


## Problem

{PRE-FILL: Infer the problem from context. Example for SQS redrive:
"AWS SQS Dead Letter Queue(DLQ)에 쌓인 메시지를 원본 큐로 다시 보내는 작업이
AWS Console에서는 번거롭고, 대량 처리가 어렵습니다."}

**검토/수정:**


## Users

{PRE-FILL: Infer target users. Example:
"- DevOps 엔지니어
- Backend 개발자
- SRE 팀"}

**검토/수정:**


## Key Requirements

{PRE-FILL: Suggest requirements based on feature type. Example for CLI:
"- [ ] DLQ에서 메시지 목록 조회
- [ ] 선택한 메시지를 원본 큐로 redrive
- [ ] 전체 메시지 일괄 redrive
- [ ] 메시지 내용 미리보기
- [ ] dry-run 모드 지원"}

**추가할 요구사항:**
-
-

## Nice-to-Have

{PRE-FILL: Suggest optional features. Example:
"- [ ] 메시지 필터링 (특정 패턴만 redrive)
- [ ] 진행률 표시
- [ ] 결과 리포트 출력"}

**추가:**
-

## Tech Stack

{PRE-FILL: Use detected project tech stack. Only suggest changes if there's a good reason.}

- Language: {from project context, e.g., "TypeScript (기존 프로젝트와 동일)"}
- Framework: {from project context}
- Database: {if applicable}
- New Dependencies: {any new libraries needed for this feature}

**기존 프로젝트 스택 유지 권장. 변경이 필요하면:**


## Constraints

{PRE-FILL: Known constraints. Example:
"- AWS credentials 필요 (환경변수 또는 AWS profile)
- SQS 접근 권한 필요"}

**추가:**


## References

{PRE-FILL: Relevant documentation}
- AWS SQS DLQ: https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-dead-letter-queues.html
- AWS CLI sqs: https://awscli.amazonaws.com/v2/documentation/api/latest/reference/sqs/index.html

**추가:**


## Open Questions

{PRE-FILL: Questions to consider. Example:
"- 메시지 redrive 실패 시 어떻게 처리할까요?
- 동시에 여러 DLQ를 처리할 수 있어야 할까요?"}

**추가:**


---

내용을 검토/수정한 후 "done" 또는 "완료"라고 말씀해주세요.
```

**Pre-fill Guidelines:**

For **CLI tools**, suggest:
- list/get/create/delete 같은 CRUD 커맨드
- --dry-run, --verbose, --output-format 같은 공통 플래그
- 설정 파일 지원 여부

For **AWS-related features**, include:
- 필요한 AWS 권한
- Region 설정 방법
- Credentials 처리 방식

For **Web APIs**, suggest:
- 주요 엔드포인트
- 인증 방식
- 에러 핸들링

The goal is to give the user a **70% complete draft** that they can review and refine, not an empty template.

### Step 4: Notify User

Tell the user:

```
브레인스토밍 파일을 생성했습니다:
doc/feature_{name}/{name}_brainstorming.md

파일을 열어서 각 섹션을 채워주세요.
완료되면 "done" 또는 "brainstorming 완료"라고 말씀해주세요.
```

### Step 5: Wait for User

Wait for user to respond that they're done filling the brainstorming file.

### Step 6: Read and Process Brainstorming

Read the filled brainstorming file and extract:
- Overview → Summary for spec
- Problem + Users → Context
- Key Requirements → Functional requirements
- Nice-to-Have → Non-goals or optional requirements
- Tech Stack → Technology choices
- Constraints → Non-functional requirements

### Step 6.5: Tech Stack Confirmation

**If project already has a defined tech stack** (detected in Step 0):
- Confirm using the existing stack
- Only suggest additions (new dependencies) needed for this feature
- Example: "기존 프로젝트가 Go + Cobra를 사용하므로 동일한 스택을 사용합니다. AWS SDK v2를 추가해야 합니다."

**If Tech Stack section is empty or this is a new project**, provide recommendations with reasoning.

Analyze the feature type and present options in this format:

```
Tech Stack이 비어있습니다. Feature 분석 결과를 바탕으로 추천드립니다:

## Language 선택지

| 옵션 | 추천 | 이유 |
|------|------|------|
| Go | ⭐ 추천 | 단일 바이너리 배포, 빠른 실행 속도, cobra로 CLI 구조화 용이 |
| Python | 대안 | 빠른 개발, boto3로 AWS 연동 쉬움, 하지만 배포 시 의존성 관리 필요 |
| Rust | 고려 | 최고 성능, 하지만 학습 곡선 높음 |

## Framework 선택지

| 옵션 | 추천 | 이유 |
|------|------|------|
| Cobra | ⭐ 추천 | Go CLI 표준, subcommand 지원, 자동 help 생성 |
| Click | 대안 | Python용, 데코레이터 기반으로 직관적 |

## 최종 추천

**Go + Cobra + aws-sdk-go-v2**

이유:
1. CLI 도구는 단일 바이너리가 배포/사용에 유리
2. AWS SDK v2는 모듈화되어 필요한 것만 import 가능
3. Cobra는 kubectl, docker 등 유명 CLI에서 사용하는 검증된 프레임워크

이 추천을 사용하시겠습니까? (yes/no/다른 선택)
```

**Recommendation Templates by Feature Type:**

**CLI Tool:**
- ⭐ Go + Cobra: 단일 바이너리, 크로스 컴파일 쉬움, 빠른 시작 시간
- Python + Click: 빠른 프로토타이핑, AWS boto3 연동 쉬움, 런타임 필요
- Rust + Clap: 최고 성능, 메모리 안전, 학습 곡선 높음

**Web API:**
- ⭐ TypeScript + Fastify: 타입 안전, 빠른 성능, 생태계 풍부
- Python + FastAPI: 자동 문서화, 타입 힌트, ML 연동 좋음
- Go + Gin: 고성능, 간단한 구조, 컴파일 타임 체크

**Web Frontend:**
- ⭐ React + TypeScript: 큰 생태계, 채용 쉬움, 타입 안전
- Vue + TypeScript: 학습 쉬움, 문서 좋음, 중소규모에 적합
- Svelte: 번들 작음, 빠름, 생태계 작음

**AWS Integration:**
- ⭐ Go + aws-sdk-go-v2: 프로덕션 CLI용, AWS 공식 지원
- Python + boto3: 스크립트/자동화용, 빠른 개발
- TypeScript + AWS SDK v3: CDK 연동, 풀스택 TypeScript

Wait for user response and update the brainstorming file with chosen tech stack.

### Step 7: Create Spec File

Based on brainstorming, create `doc/feature_{name}/{name}_spec.md`:

```markdown
# {Feature Name} - Specification

## Summary

{One-line summary from brainstorming overview}

## Goals

- {Goal 1 from requirements}
- {Goal 2}

## Non-Goals

- {Nice-to-have items that won't be in initial version}

## Requirements

### Functional Requirements
1. FR-1: {requirement}
2. FR-2: {requirement}

### Non-Functional Requirements
1. NFR-1: {constraint}

## Acceptance Criteria

- [ ] {Criterion derived from requirements}
- [ ] {Criterion 2}

## Tech Stack

- **Language**: {chosen language}
- **Framework**: {if applicable}
- **Dependencies**: {key libraries/tools}

## Dependencies

- {From constraints section}
```

### Step 8: Create Plan File

Check installed agents in `.claude/commands/` (exclude orchestrator.md).

Create `doc/feature_{name}/{name}_plan.md` with sections for each agent:

```markdown
# {Feature Name} - Implementation Plan

## Overview

{Brief implementation overview}

## Alpha

Alpha agent tasks:

- [ ] Task 1: {description}
- [ ] Task 2: {description}

## Beta

Beta agent tasks:

- [ ] Task 1: {description}
- [ ] Task 2: {description}

## Delta

Delta agent tasks:

- [ ] Task 1: {description}

## Dependencies

- {Task dependencies if any}
```

Distribute tasks based on:
- Break down requirements into individual tasks
- Assign related tasks to same agent
- Balance workload across agents

### Step 9: Git Commit

If in a git repository, commit the feature planning files:

```bash
git add doc/feature_{name}/
git commit -m "feat({name}): add feature planning documents

- Add brainstorming document
- Add specification document
- Add implementation plan

🤖 Generated with workflow-adapter"
```

### Step 10: Summary

Report to user:

```
Feature planning complete: {name}

Git:
- Branch: feature/{name} ✓
- Commit: feat({name}): add feature planning documents ✓

Created files:
- doc/feature_{name}/{name}_brainstorming.md ✓
- doc/feature_{name}/{name}_spec.md ✓
- doc/feature_{name}/{name}_plan.md ✓

Next steps:
1. Review the spec and plan files
2. Run /workflow-adapter:execute {name} to start agent work
3. Run /workflow-adapter:validate {name} to check progress
4. Run /workflow-adapter:pr {name} to create pull request when done
```

## Your Task

User input: `$ARGUMENTS`

Start by extracting the feature name and creating the brainstorming file.
