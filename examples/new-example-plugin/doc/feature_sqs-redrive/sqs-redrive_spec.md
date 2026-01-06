# SQS Redrive CLI - Specification

## Summary

AWS SQS Dead Letter Queue(DLQ)에 쌓인 메시지를 원본 큐로 쉽게 redrive할 수 있는 CLI 도구

## Goals

- DLQ 메시지를 빠르고 쉽게 원본 큐로 redrive
- 대량의 메시지 일괄 처리 지원
- 메시지 내용 확인 후 선택적 redrive 가능
- 자동화 스크립트에서 사용 가능

## Non-Goals (v1)

- 메시지 필터링 (JSON path/regex)
- 인터랙티브 모드
- 메시지 수정 후 redrive
- 설정 파일 지원

## Requirements

### Functional Requirements

1. **FR-1**: DLQ 목록 조회
   - 계정 내 모든 DLQ 나열
   - 각 DLQ의 메시지 수 표시

2. **FR-2**: 메시지 목록 조회
   - 지정한 DLQ의 메시지 목록 조회
   - 메시지 ID, 수신 시간, 대략적인 내용 미리보기

3. **FR-3**: 메시지 상세 보기 (peek)
   - 특정 메시지의 전체 내용 조회
   - JSON 포맷팅 지원

4. **FR-4**: 단일 메시지 redrive
   - 메시지 ID로 특정 메시지를 원본 큐로 이동

5. **FR-5**: 전체 메시지 일괄 redrive
   - DLQ의 모든 메시지를 원본 큐로 이동
   - 진행률 표시

6. **FR-6**: dry-run 모드
   - 실제 실행 없이 수행될 작업 미리보기

7. **FR-7**: AWS 설정 지원
   - `--profile` 플래그로 AWS profile 선택
   - `--region` 플래그로 region 선택
   - 환경변수 fallback (AWS_PROFILE, AWS_REGION)

### Non-Functional Requirements

1. **NFR-1**: AWS credentials 필요
   - 환경변수 또는 AWS profile 지원
   
2. **NFR-2**: 필요 IAM 권한
   - sqs:ListQueues
   - sqs:GetQueueAttributes
   - sqs:ReceiveMessage
   - sqs:SendMessage
   - sqs:DeleteMessage

3. **NFR-3**: 단일 바이너리 배포
   - 의존성 없이 실행 가능

## CLI Structure

```
sqs-redrive
├── list                    # DLQ 목록 조회
├── messages <queue-url>    # 메시지 목록 조회
├── peek <queue-url> <msg-id>  # 메시지 상세 보기
├── redrive <queue-url>     # 메시지 redrive
│   ├── --message-id <id>   # 특정 메시지만
│   ├── --all               # 전체 메시지
│   └── --dry-run           # 미리보기
└── Global flags
    ├── --profile <name>    # AWS profile
    ├── --region <region>   # AWS region
    └── --help              # 도움말
```

## Acceptance Criteria

- [ ] `sqs-redrive list` 실행 시 DLQ 목록과 메시지 수 출력
- [ ] `sqs-redrive messages <url>` 실행 시 메시지 목록 출력
- [ ] `sqs-redrive peek <url> <id>` 실행 시 메시지 전체 내용 출력
- [ ] `sqs-redrive redrive <url> --message-id <id>` 실행 시 해당 메시지 redrive
- [ ] `sqs-redrive redrive <url> --all` 실행 시 모든 메시지 redrive
- [ ] `--dry-run` 플래그 사용 시 실제 작업 없이 예상 결과 출력
- [ ] `--profile`, `--region` 플래그 정상 동작
- [ ] 에러 발생 시 명확한 에러 메시지 출력

## Tech Stack

- **Language**: Go
- **Framework**: Cobra (github.com/spf13/cobra)
- **AWS SDK**: aws-sdk-go-v2
- **Dependencies**:
  - github.com/fatih/color (컬러 출력)
  - github.com/schollz/progressbar (진행률 표시)

## Project Structure

```
sqs-redrive/
├── cmd/
│   ├── root.go         # Root command, global flags
│   ├── list.go         # list subcommand
│   ├── messages.go     # messages subcommand
│   ├── peek.go         # peek subcommand
│   └── redrive.go      # redrive subcommand
├── internal/
│   ├── aws/
│   │   └── sqs.go      # SQS client wrapper
│   └── output/
│       └── printer.go  # Output formatting
├── main.go
├── go.mod
└── go.sum
```
