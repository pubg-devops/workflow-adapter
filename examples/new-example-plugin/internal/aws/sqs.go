package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// SQSClient wraps the AWS SQS client with helper methods
type SQSClient struct {
	client *sqs.Client
}

// QueueInfo contains information about an SQS queue
type QueueInfo struct {
	URL                   string
	Name                  string
	ApproximateMessages   int
	SourceQueueURL        string
	IsDLQ                 bool
}

// MessageInfo contains information about an SQS message
type MessageInfo struct {
	MessageID              string
	ReceiptHandle          string
	Body                   string
	ApproximateReceiveCount string
	SentTimestamp          string
	ApproximateFirstReceiveTimestamp string
	Attributes             map[string]string
	MessageAttributes      map[string]types.MessageAttributeValue
}

// NewClient creates a new SQS client with the given configuration options
func NewClient(ctx context.Context, optFns ...func(*config.LoadOptions) error) (*SQSClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	return &SQSClient{
		client: sqs.NewFromConfig(cfg),
	}, nil
}

// ListQueues returns all queues, optionally filtered by prefix
func (c *SQSClient) ListQueues(ctx context.Context, prefix string) ([]string, error) {
	var queueURLs []string
	var nextToken *string

	for {
		input := &sqs.ListQueuesInput{}
		if prefix != "" {
			input.QueueNamePrefix = aws.String(prefix)
		}
		if nextToken != nil {
			input.NextToken = nextToken
		}

		result, err := c.client.ListQueues(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list queues: %w", err)
		}

		queueURLs = append(queueURLs, result.QueueUrls...)

		if result.NextToken == nil {
			break
		}
		nextToken = result.NextToken
	}

	return queueURLs, nil
}

// ListDLQs returns only queues that are Dead Letter Queues
func (c *SQSClient) ListDLQs(ctx context.Context) ([]QueueInfo, error) {
	queueURLs, err := c.ListQueues(ctx, "")
	if err != nil {
		return nil, err
	}

	var dlqs []QueueInfo

	for _, queueURL := range queueURLs {
		info, err := c.GetQueueAttributes(ctx, queueURL)
		if err != nil {
			// Skip queues we can't get attributes for
			continue
		}

		// Check if this queue is a DLQ by looking for queues that reference it
		// A DLQ typically has "dlq", "dead-letter", or "deadletter" in its name
		queueName := extractQueueName(queueURL)
		isDLQ := strings.Contains(strings.ToLower(queueName), "dlq") ||
			strings.Contains(strings.ToLower(queueName), "dead-letter") ||
			strings.Contains(strings.ToLower(queueName), "deadletter")

		if isDLQ {
			info.IsDLQ = true
			dlqs = append(dlqs, *info)
		}
	}

	return dlqs, nil
}

// GetQueueAttributes retrieves attributes for a specific queue
func (c *SQSClient) GetQueueAttributes(ctx context.Context, queueURL string) (*QueueInfo, error) {
	result, err := c.client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(queueURL),
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameAll,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get queue attributes: %w", err)
	}

	info := &QueueInfo{
		URL:  queueURL,
		Name: extractQueueName(queueURL),
	}

	// Parse approximate number of messages
	if msgCount, ok := result.Attributes["ApproximateNumberOfMessages"]; ok {
		var count int
		fmt.Sscanf(msgCount, "%d", &count)
		info.ApproximateMessages = count
	}

	// Parse redrive policy to find source queue
	if redrivePolicy, ok := result.Attributes["RedriveAllowPolicy"]; ok {
		// The redrive allow policy contains info about which queues can use this as DLQ
		_ = redrivePolicy // For now, we note it exists
	}

	// Check if this queue has a redrive policy (meaning it sends to a DLQ)
	if redrivePolicy, ok := result.Attributes["RedrivePolicy"]; ok {
		// Parse the redrive policy to extract DLQ ARN
		// Format: {"deadLetterTargetArn":"arn:aws:sqs:region:account:queue","maxReceiveCount":"5"}
		_ = redrivePolicy // For now, we note it exists
	}

	return info, nil
}

// GetSourceQueueURL attempts to find the source queue for a DLQ
// It does this by scanning all queues and finding ones that have this DLQ in their redrive policy
func (c *SQSClient) GetSourceQueueURL(ctx context.Context, dlqURL string) (string, error) {
	queueURLs, err := c.ListQueues(ctx, "")
	if err != nil {
		return "", err
	}

	dlqARN, err := c.getQueueARN(ctx, dlqURL)
	if err != nil {
		return "", err
	}

	for _, queueURL := range queueURLs {
		if queueURL == dlqURL {
			continue
		}

		result, err := c.client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl: aws.String(queueURL),
			AttributeNames: []types.QueueAttributeName{
				types.QueueAttributeNameRedrivePolicy,
			},
		})
		if err != nil {
			continue
		}

		if redrivePolicy, ok := result.Attributes["RedrivePolicy"]; ok {
			if strings.Contains(redrivePolicy, dlqARN) {
				return queueURL, nil
			}
		}
	}

	return "", fmt.Errorf("no source queue found for DLQ: %s", dlqURL)
}

// getQueueARN retrieves the ARN for a queue
func (c *SQSClient) getQueueARN(ctx context.Context, queueURL string) (string, error) {
	result, err := c.client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(queueURL),
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameQueueArn,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to get queue ARN: %w", err)
	}

	if arn, ok := result.Attributes["QueueArn"]; ok {
		return arn, nil
	}

	return "", fmt.Errorf("QueueArn not found in attributes")
}

// ReceiveMessages receives messages from a queue
func (c *SQSClient) ReceiveMessages(ctx context.Context, queueURL string, maxMessages int32, visibilityTimeout int32) ([]MessageInfo, error) {
	result, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(queueURL),
		MaxNumberOfMessages:   maxMessages,
		VisibilityTimeout:     visibilityTimeout,
		AttributeNames:        []types.QueueAttributeName{types.QueueAttributeNameAll},
		MessageAttributeNames: []string{"All"},
		WaitTimeSeconds:       0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to receive messages: %w", err)
	}

	var messages []MessageInfo
	for _, msg := range result.Messages {
		info := MessageInfo{
			MessageID:         aws.ToString(msg.MessageId),
			ReceiptHandle:     aws.ToString(msg.ReceiptHandle),
			Body:              aws.ToString(msg.Body),
			Attributes:        msg.Attributes,
			MessageAttributes: msg.MessageAttributes,
		}

		if v, ok := msg.Attributes["ApproximateReceiveCount"]; ok {
			info.ApproximateReceiveCount = v
		}
		if v, ok := msg.Attributes["SentTimestamp"]; ok {
			info.SentTimestamp = v
		}
		if v, ok := msg.Attributes["ApproximateFirstReceiveTimestamp"]; ok {
			info.ApproximateFirstReceiveTimestamp = v
		}

		messages = append(messages, info)
	}

	return messages, nil
}

// SendMessage sends a message to a queue
func (c *SQSClient) SendMessage(ctx context.Context, queueURL string, body string, attributes map[string]types.MessageAttributeValue) error {
	_, err := c.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:          aws.String(queueURL),
		MessageBody:       aws.String(body),
		MessageAttributes: attributes,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// DeleteMessage deletes a message from a queue
func (c *SQSClient) DeleteMessage(ctx context.Context, queueURL string, receiptHandle string) error {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

// RedriveMessage moves a single message from DLQ to source queue
// It sends the message to the source queue, then deletes it from the DLQ
func (c *SQSClient) RedriveMessage(ctx context.Context, dlqURL, sourceQueueURL string, msg MessageInfo) error {
	// Send message to source queue (preserve message attributes)
	err := c.SendMessage(ctx, sourceQueueURL, msg.Body, msg.MessageAttributes)
	if err != nil {
		return fmt.Errorf("failed to send message to source queue: %w", err)
	}

	// Delete message from DLQ
	err = c.DeleteMessage(ctx, dlqURL, msg.ReceiptHandle)
	if err != nil {
		return fmt.Errorf("failed to delete message from DLQ: %w", err)
	}

	return nil
}

// RedriveAllMessages moves all messages from DLQ to source queue
// Returns the number of successfully redriven messages and any error
func (c *SQSClient) RedriveAllMessages(ctx context.Context, dlqURL, sourceQueueURL string) (int, error) {
	successCount := 0
	emptyReceives := 0
	maxEmptyReceives := 3

	for emptyReceives < maxEmptyReceives {
		messages, err := c.ReceiveMessages(ctx, dlqURL, 10, 30)
		if err != nil {
			return successCount, fmt.Errorf("failed to receive messages: %w", err)
		}

		if len(messages) == 0 {
			emptyReceives++
			continue
		}
		emptyReceives = 0

		for _, msg := range messages {
			err := c.RedriveMessage(ctx, dlqURL, sourceQueueURL, msg)
			if err != nil {
				return successCount, fmt.Errorf("failed to redrive message %s: %w", msg.MessageID, err)
			}
			successCount++
		}
	}

	return successCount, nil
}

// extractQueueName extracts the queue name from a queue URL
func extractQueueName(queueURL string) string {
	parts := strings.Split(queueURL, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return queueURL
}
