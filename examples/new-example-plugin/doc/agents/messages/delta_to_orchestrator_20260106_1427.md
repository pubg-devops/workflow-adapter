# Message

- **From**: delta
- **To**: orchestrator
- **Feature**: sqs-redrive
- **Status**: completed
- **Timestamp**: 2026-01-06 14:27

## Content

Completed the following tasks:
- [x] cmd/peek.go - Message detail view command
- [x] cmd/redrive.go - Redrive command with flags (--message-id, --all, --dry-run)
- [x] internal/aws/sqs.go - Added RedriveMessage() and RedriveAllMessages() methods

## Results

### 1. cmd/peek.go
Implemented message detail view command with:
- Accepts `<queue-url>` and `<message-id>` arguments
- Searches through messages to find the target by ID
- Displays message ID, receipt handle, timestamps
- Shows system attributes and message attributes
- Pretty-prints JSON message body

### 2. cmd/redrive.go
Implemented redrive command with:
- `--message-id <id>` flag for single message redrive
- `--all` flag for bulk redrive of all messages
- `--dry-run` flag to preview actions without execution
- Automatic source queue detection via redrive policy scanning
- Progress bar display using progressbar library for bulk operations
- Success/error count reporting

### 3. internal/aws/sqs.go additions
- `RedriveMessage(ctx, dlqURL, sourceQueueURL, msg)` - Moves single message from DLQ to source queue (send then delete)
- `RedriveAllMessages(ctx, dlqURL, sourceQueueURL)` - Bulk redrive with automatic batching

### Build Verification
- All code compiles successfully with `go build ./...`

## Dependencies Met
- Used Alpha's completed SQS client wrapper (NewClient, ReceiveMessages, SendMessage, DeleteMessage, GetSourceQueueURL)
- Used progressbar v3 and fatih/color libraries as specified in tech stack
