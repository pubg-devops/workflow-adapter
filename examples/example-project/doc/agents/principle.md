# Agent Principle

이 문서는 모든 agent (alpha, beta, delta, orchestrator)가 따라야 할 공통 원칙을 정의합니다.

---

## 1. Feature 기반 작업 흐름

### 1.1 Feature 폴더 확인
유저가 `{name}` feature 작업을 요청하면:
1. `doc/feature_{name}/` 폴더로 이동
2. 다음 파일들을 확인:
   - `{name}_brainstorming.md` - 아이디어 및 논의 내용
   - `{name}_spec.md` - 피처 명세서
   - `{name}_plan.md` - 구현 계획 및 각 agent별 할일

### 1.2 자신의 작업 확인
`{name}_plan.md`에서 자신에게 할당된 섹션을 찾아 작업을 수행합니다:
- **Alpha agent**: `## Alpha` 섹션의 작업
- **Beta agent**: `## Beta` 섹션의 작업
- **Delta agent**: `## Delta` 섹션의 작업
- **Orchestrator**: 전체 plan과 spec을 검토하여 작업 완료 여부 검증

---

## 2. Agent 역할 정의

### Alpha
- plan의 `## Alpha` 섹션에 정의된 작업 수행
- 작업 완료 시 messages 폴더에 상태 보고

### Beta
- plan의 `## Beta` 섹션에 정의된 작업 수행
- 작업 완료 시 messages 폴더에 상태 보고

### Delta
- plan의 `## Delta` 섹션에 정의된 작업 수행
- 작업 완료 시 messages 폴더에 상태 보고

### Orchestrator
- `{name}_spec.md`와 `{name}_plan.md`를 기준으로 각 agent의 작업 완료 여부 검증
- messages 폴더의 메시지를 확인하여 진행 상황 파악
- 모든 작업이 spec을 충족하는지 최종 검토

---

## 3. Agent 간 메시지 통신

### 3.1 메시지 저장 위치
`doc/agents/messages/` 폴더에 메시지 파일을 생성합니다.

### 3.2 파일 네이밍 규칙
```
{from_agent}_to_{to_agent}_{YYYYMMDD}_{HHMM}.md
```

**예시:**
- `alpha_to_orchestrator_20240115_1430.md`
- `beta_to_delta_20240115_1500.md`
- `orchestrator_to_all_20240115_1600.md` (전체 공지 시)

### 3.3 메시지 포맷
```markdown
# Message

- **From**: {agent_name}
- **To**: {target_agent_name | all}
- **Feature**: {feature_name}
- **Status**: {completed | in_progress | blocked}
- **Timestamp**: {YYYY-MM-DD HH:MM}

## Content

{메시지 내용을 여기에 작성}

## Action Required (optional)

{상대 agent에게 요청하는 액션이 있다면 작성}
```

---

## 4. 작업 완료 보고

각 agent는 자신의 작업을 완료하면 반드시:
1. messages 폴더에 완료 메시지 작성
2. Status를 `completed`로 설정
3. 수행한 작업 내용을 Content에 명시

---

## 5. 문제 발생 시

작업 중 문제가 발생하면:
1. messages 폴더에 메시지 작성
2. Status를 `blocked`로 설정
3. 문제 상황과 필요한 도움을 Content에 명시
4. 관련 agent나 orchestrator를 To에 지정
