# SQS Redrive CLI - Implementation Plan

## Overview

Go + Cobra 기반의 SQS Redrive CLI 도구 구현. 4개의 agent가 병렬로 작업을 수행합니다.

## Alpha

프로젝트 초기 설정 및 핵심 AWS SQS 클라이언트 구현

- [ ] Go 모듈 초기화 및 의존성 설정 (go.mod)
- [ ] main.go 및 root command 구현 (global flags: --profile, --region)
- [ ] internal/aws/sqs.go - SQS 클라이언트 wrapper 구현
  - NewClient() 함수
  - ListQueues() - DLQ 필터링 포함
  - GetQueueAttributes() - 메시지 수, 원본 큐 URL 조회

## Beta

list 및 messages 커맨드 구현

- [ ] cmd/list.go - DLQ 목록 조회 커맨드
  - DLQ 목록 출력
  - 각 DLQ의 메시지 수 표시
  - 테이블 형식 출력
- [ ] cmd/messages.go - 메시지 목록 조회 커맨드
  - queue-url 인자 처리
  - 메시지 ID, 수신 시간, 미리보기 출력
- [ ] internal/output/printer.go - 출력 포맷팅 유틸리티

## Delta

peek 및 redrive 커맨드 구현

- [ ] cmd/peek.go - 메시지 상세 보기 커맨드
  - queue-url, message-id 인자 처리
  - JSON pretty print
- [ ] cmd/redrive.go - redrive 커맨드
  - --message-id 플래그 (단일 메시지)
  - --all 플래그 (전체 메시지)
  - --dry-run 플래그
  - 진행률 표시 (progressbar)
- [ ] internal/aws/sqs.go에 redrive 관련 메서드 추가
  - RedriveMessage()
  - RedriveAllMessages()

## Gamma

테스트, 문서화, 빌드 설정

- [ ] 단위 테스트 작성
  - internal/aws/sqs_test.go
  - cmd/ 각 커맨드 테스트
- [ ] README.md 작성
  - 설치 방법
  - 사용 예시
  - 필요 IAM 권한
- [ ] Makefile 작성
  - build, test, lint 타겟
- [ ] .goreleaser.yml (선택) - 크로스 컴파일 빌드

## Dependencies

```
Alpha 완료 → Beta, Delta 시작 가능
Beta, Delta 완료 → Gamma 테스트 작성 가능
```

## Notes

- AWS StartMessageMoveTask API 대신 직접 receive/send/delete 방식 사용
  - 더 세밀한 제어 가능 (메시지별 선택)
  - 진행률 표시 가능
- FIFO 큐는 v1에서 지원하지 않음 (표준 큐만)
