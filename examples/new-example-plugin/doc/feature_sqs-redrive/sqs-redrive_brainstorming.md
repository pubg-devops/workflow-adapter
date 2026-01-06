# SQS Redrive CLI - Brainstorming

> 아래 내용을 검토하고 수정/추가해주세요. 맞으면 그대로 두셔도 됩니다.

## Overview

AWS SQS Dead Letter Queue(DLQ)에 쌓인 메시지를 원본 큐로 쉽게 redrive할 수 있는 CLI 도구입니다.

**검토/수정:**


## Problem

AWS SQS Dead Letter Queue(DLQ)에 쌓인 메시지를 원본 큐로 다시 보내는 작업이
AWS Console에서는 번거롭고, 대량 처리가 어렵습니다.

- Console에서는 한 번에 10개씩만 처리 가능
- 메시지 내용을 확인하며 선택적으로 redrive하기 어려움
- 반복 작업 자동화 불가
- 여러 DLQ를 순회하며 처리하기 불편

**검토/수정:**


## Users

- DevOps 엔지니어
- Backend 개발자
- SRE 팀
- 운영팀

**검토/수정:**


## Key Requirements

- [O] DLQ 목록 조회 (계정 내 DLQ 찾기)
- [O] DLQ에서 메시지 목록 조회
- [O] 메시지 내용 미리보기 (peek)
- [O] 선택한 메시지를 원본 큐로 redrive
- [O] 전체 메시지 일괄 redrive
- [O] dry-run 모드 지원
- [O] AWS profile 선택 지원
- [O] Region 선택 지원

**추가할 요구사항:**
-
-

## Nice-to-Have

- [ ] 메시지 필터링 (JSON path나 regex로 특정 메시지만 redrive)
- [ ] 진행률 표시 (progress bar)
- [ ] 결과 리포트 출력 (성공/실패 카운트)
- [ ] 메시지 삭제 기능 (redrive 없이 삭제)
- [ ] 인터랙티브 모드 (메시지 하나씩 확인하며 처리)
- [ ] 메시지 수정 후 redrive
- [ ] 설정 파일 지원 (~/.sqs-redrive.yaml)

**추가:**
-

## Tech Stack

- Language: Go
- Framework: Cobra (github.com/spf13/cobra)
- AWS SDK: aws-sdk-go-v2
- Other tools:
  - github.com/fatih/color (컬러 출력)
  - github.com/schollz/progressbar (진행률 표시)

**확정됨**


## Constraints

- AWS credentials 필요 (환경변수 또는 AWS profile)
- SQS 접근 권한 필요:
  - sqs:ListQueues
  - sqs:GetQueueAttributes
  - sqs:ReceiveMessage
  - sqs:SendMessage
  - sqs:DeleteMessage
- DLQ와 연결된 원본 큐 정보 조회 필요

**추가:**


## References

- AWS SQS DLQ: https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-dead-letter-queues.html
- AWS CLI sqs: https://awscli.amazonaws.com/v2/documentation/api/latest/reference/sqs/index.html
- StartMessageMoveTask API: https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_StartMessageMoveTask.html

**추가:**


## Open Questions

- 메시지 redrive 실패 시 어떻게 처리할까요? (재시도? 로그만?)
- 동시에 여러 DLQ를 처리할 수 있어야 할까요?
- AWS의 native StartMessageMoveTask API를 사용할까요, 아니면 직접 receive/send/delete 할까요?
- FIFO 큐도 지원해야 할까요?

**추가:**


---

내용을 검토/수정한 후 "done" 또는 "완료"라고 말씀해주세요.
