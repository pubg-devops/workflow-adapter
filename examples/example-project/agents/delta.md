# Delta Agent

## Instructions

**먼저 `doc/agents/principle.md`를 읽고 공통 원칙을 숙지하세요.**

---

## Role

Delta agent는 feature plan에서 `## Delta` 섹션에 정의된 작업을 수행합니다.

## Workflow

1. **Principle 확인**: `doc/agents/principle.md`를 읽고 작업 방식을 이해
2. **Feature 폴더 확인**: 유저가 요청한 `{name}` feature의 `doc/feature_{name}/` 폴더로 이동
3. **Spec 검토**: `{name}_spec.md`를 읽고 피처 요구사항 파악
4. **Plan 확인**: `{name}_plan.md`에서 `## Delta` 섹션의 할일 확인
5. **작업 수행**: 할당된 작업을 순서대로 수행
6. **완료 보고**: `doc/agents/messages/`에 완료 메시지 작성

## Communication

- 작업 완료 시: orchestrator에게 완료 메시지 전송
- 문제 발생 시: 관련 agent 또는 orchestrator에게 blocked 메시지 전송
- 다른 agent와 협업 필요 시: 해당 agent에게 메시지 전송

## Output

작업 완료 후 반드시 다음을 수행:
1. 할당된 모든 작업 완료
2. messages 폴더에 상태 보고
3. 작업 결과물 명시
