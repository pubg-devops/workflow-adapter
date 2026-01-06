package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// mockSQSClient implements a mock SQS client for testing
type mockSQSClient struct {
	listQueuesOutput       *sqs.ListQueuesOutput
	listQueuesError        error
	getQueueAttributesFunc func(queueURL string) (*sqs.GetQueueAttributesOutput, error)
	receiveMessageOutput   *sqs.ReceiveMessageOutput
	receiveMessageError    error
	sendMessageOutput      *sqs.SendMessageOutput
	sendMessageError       error
	deleteMessageOutput    *sqs.DeleteMessageOutput
	deleteMessageError     error
}

func (m *mockSQSClient) ListQueues(ctx context.Context, params *sqs.ListQueuesInput, optFns ...func(*sqs.Options)) (*sqs.ListQueuesOutput, error) {
	return m.listQueuesOutput, m.listQueuesError
}

func (m *mockSQSClient) GetQueueAttributes(ctx context.Context, params *sqs.GetQueueAttributesInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueAttributesOutput, error) {
	if m.getQueueAttributesFunc != nil {
		return m.getQueueAttributesFunc(aws.ToString(params.QueueUrl))
	}
	return &sqs.GetQueueAttributesOutput{
		Attributes: map[string]string{
			"ApproximateNumberOfMessages": "5",
			"QueueArn":                    "arn:aws:sqs:us-east-1:123456789012:test-queue",
		},
	}, nil
}

func (m *mockSQSClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	return m.receiveMessageOutput, m.receiveMessageError
}

func (m *mockSQSClient) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	return m.sendMessageOutput, m.sendMessageError
}

func (m *mockSQSClient) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	return m.deleteMessageOutput, m.deleteMessageError
}

func TestExtractQueueName(t *testing.T) {
	tests := []struct {
		name     string
		queueURL string
		expected string
	}{
		{
			name:     "standard queue URL",
			queueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue",
			expected: "my-queue",
		},
		{
			name:     "DLQ queue URL",
			queueURL: "https://sqs.us-west-2.amazonaws.com/987654321098/my-queue-dlq",
			expected: "my-queue-dlq",
		},
		{
			name:     "queue URL with dead-letter",
			queueURL: "https://sqs.eu-west-1.amazonaws.com/111222333444/orders-dead-letter-queue",
			expected: "orders-dead-letter-queue",
		},
		{
			name:     "simple name",
			queueURL: "simple-queue",
			expected: "simple-queue",
		},
		{
			name:     "empty string",
			queueURL: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractQueueName(tt.queueURL)
			if result != tt.expected {
				t.Errorf("extractQueueName(%q) = %q, want %q", tt.queueURL, result, tt.expected)
			}
		})
	}
}

func TestQueueInfo(t *testing.T) {
	// Test QueueInfo struct initialization
	info := QueueInfo{
		URL:                 "https://sqs.us-east-1.amazonaws.com/123456789012/test-dlq",
		Name:                "test-dlq",
		ApproximateMessages: 10,
		SourceQueueURL:      "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
		IsDLQ:               true,
	}

	if info.URL == "" {
		t.Error("QueueInfo.URL should not be empty")
	}
	if info.Name != "test-dlq" {
		t.Errorf("QueueInfo.Name = %q, want %q", info.Name, "test-dlq")
	}
	if info.ApproximateMessages != 10 {
		t.Errorf("QueueInfo.ApproximateMessages = %d, want %d", info.ApproximateMessages, 10)
	}
	if !info.IsDLQ {
		t.Error("QueueInfo.IsDLQ should be true")
	}
}

func TestMessageInfo(t *testing.T) {
	// Test MessageInfo struct initialization
	info := MessageInfo{
		MessageID:                        "msg-123",
		ReceiptHandle:                    "receipt-handle-xyz",
		Body:                             `{"event":"test"}`,
		ApproximateReceiveCount:          "3",
		SentTimestamp:                    "1704067200000",
		ApproximateFirstReceiveTimestamp: "1704067201000",
		Attributes:                       map[string]string{"attr1": "value1"},
		MessageAttributes: map[string]types.MessageAttributeValue{
			"customAttr": {
				DataType:    aws.String("String"),
				StringValue: aws.String("customValue"),
			},
		},
	}

	if info.MessageID != "msg-123" {
		t.Errorf("MessageInfo.MessageID = %q, want %q", info.MessageID, "msg-123")
	}
	if info.Body != `{"event":"test"}` {
		t.Errorf("MessageInfo.Body = %q, want %q", info.Body, `{"event":"test"}`)
	}
	if info.ApproximateReceiveCount != "3" {
		t.Errorf("MessageInfo.ApproximateReceiveCount = %q, want %q", info.ApproximateReceiveCount, "3")
	}
}

func TestIsDLQDetection(t *testing.T) {
	tests := []struct {
		name     string
		queueURL string
		isDLQ    bool
	}{
		{
			name:     "queue with dlq suffix",
			queueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue-dlq",
			isDLQ:    true,
		},
		{
			name:     "queue with DLQ uppercase",
			queueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue-DLQ",
			isDLQ:    true,
		},
		{
			name:     "queue with dead-letter",
			queueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-dead-letter",
			isDLQ:    true,
		},
		{
			name:     "queue with deadletter",
			queueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/orders-deadletter",
			isDLQ:    true,
		},
		{
			name:     "regular queue",
			queueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue",
			isDLQ:    false,
		},
		{
			name:     "queue with similar name",
			queueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/daily-queue",
			isDLQ:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queueName := extractQueueName(tt.queueURL)
			isDLQ := containsDLQKeyword(queueName)
			if isDLQ != tt.isDLQ {
				t.Errorf("isDLQ(%q) = %v, want %v", tt.queueURL, isDLQ, tt.isDLQ)
			}
		})
	}
}

// containsDLQKeyword checks if queue name contains DLQ-related keywords
func containsDLQKeyword(queueName string) bool {
	lowerName := lowerString(queueName)
	return contains(lowerName, "dlq") ||
		contains(lowerName, "dead-letter") ||
		contains(lowerName, "deadletter")
}

// lowerString converts string to lowercase (inline to avoid import)
func lowerString(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

// contains checks if s contains substr
func contains(s, substr string) bool {
	return len(substr) <= len(s) && findSubstring(s, substr) >= 0
}

// findSubstring returns index of substr in s, or -1 if not found
func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
