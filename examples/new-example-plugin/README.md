# sqs-redrive

AWS SQS Dead Letter Queue(DLQ) redrive CLI tool. Easily redrive messages from DLQ back to the original queue.

## Features

- List all DLQs in your AWS account with message counts
- View messages in a DLQ with preview
- Peek at full message content with JSON formatting
- Redrive single message or all messages at once
- Dry-run mode to preview actions before execution
- Progress bar for bulk redrive operations

## Installation

### From Binary

Download the latest release from the [Releases](https://github.com/dalpark/sqs-redrive/releases) page.

```bash
# Linux (amd64)
curl -Lo sqs-redrive https://github.com/dalpark/sqs-redrive/releases/latest/download/sqs-redrive_linux_amd64
chmod +x sqs-redrive
sudo mv sqs-redrive /usr/local/bin/

# macOS (arm64)
curl -Lo sqs-redrive https://github.com/dalpark/sqs-redrive/releases/latest/download/sqs-redrive_darwin_arm64
chmod +x sqs-redrive
sudo mv sqs-redrive /usr/local/bin/
```

### From Source

```bash
go install github.com/dalpark/sqs-redrive@latest
```

Or build locally:

```bash
git clone https://github.com/dalpark/sqs-redrive.git
cd sqs-redrive
make build
```

## Usage

### List DLQs

List all Dead Letter Queues in your account:

```bash
sqs-redrive list
```

Output:
```
QUEUE NAME                          MESSAGES  SOURCE QUEUE
my-queue-dlq                        42        my-queue
another-queue-dlq                   7         another-queue
```

### View Messages

List messages in a specific DLQ:

```bash
sqs-redrive messages https://sqs.us-east-1.amazonaws.com/123456789012/my-queue-dlq
```

Output:
```
MESSAGE ID                              RECEIVED AT           PREVIEW
abc123-def456-...                       2024-01-15 10:30:00   {"event":"user.created"...
xyz789-uvw012-...                       2024-01-15 10:31:00   {"event":"order.placed"...
```

### Peek at Message

View full content of a specific message:

```bash
sqs-redrive peek https://sqs.us-east-1.amazonaws.com/123456789012/my-queue-dlq abc123-def456-ghi789
```

Output:
```json
{
  "event": "user.created",
  "userId": "12345",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Redrive Messages

Redrive a single message by ID:

```bash
sqs-redrive redrive https://sqs.us-east-1.amazonaws.com/123456789012/my-queue-dlq --message-id abc123-def456-ghi789
```

Redrive all messages in the DLQ:

```bash
sqs-redrive redrive https://sqs.us-east-1.amazonaws.com/123456789012/my-queue-dlq --all
```

Preview actions with dry-run:

```bash
sqs-redrive redrive https://sqs.us-east-1.amazonaws.com/123456789012/my-queue-dlq --all --dry-run
```

### Global Flags

```
--profile string   AWS profile to use (default: AWS_PROFILE env or "default")
--region string    AWS region (default: AWS_REGION env or from profile)
--help             Show help for command
```

Examples with AWS profile and region:

```bash
sqs-redrive list --profile production --region us-west-2
sqs-redrive redrive <queue-url> --all --profile staging
```

## Required IAM Permissions

The following IAM permissions are required:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "sqs:ListQueues",
        "sqs:GetQueueAttributes",
        "sqs:ReceiveMessage",
        "sqs:SendMessage",
        "sqs:DeleteMessage"
      ],
      "Resource": "*"
    }
  ]
}
```

For more restrictive permissions, replace `"Resource": "*"` with specific queue ARNs:

```json
{
  "Resource": [
    "arn:aws:sqs:us-east-1:123456789012:my-queue",
    "arn:aws:sqs:us-east-1:123456789012:my-queue-dlq"
  ]
}
```

## AWS Credentials

The tool uses standard AWS credential chain:

1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. Shared credentials file (`~/.aws/credentials`)
3. IAM role (when running on EC2/ECS/Lambda)

Use `--profile` flag to specify a named profile from `~/.aws/credentials`.

## Limitations

- FIFO queues are not supported in v1
- Standard queues only
- No message filtering (planned for v2)
- No interactive mode (planned for v2)

## License

MIT License
