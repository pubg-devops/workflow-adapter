package cmd

import (
	"context"
	"fmt"

	"github.com/dalpark/sqs-redrive/internal/output"
	"github.com/spf13/cobra"
)

var (
	messagesMaxCount int32
)

var messagesCmd = &cobra.Command{
	Use:   "messages <queue-url>",
	Short: "List messages in a queue",
	Long: `List messages in the specified SQS queue.

Displays message ID, sent timestamp, receive count, and a preview of the message body.
Messages are received but not deleted, so they will return to the queue after the visibility timeout.`,
	Example: `  # List messages in a DLQ
  sqs-redrive messages https://sqs.us-east-1.amazonaws.com/123456789012/my-dlq

  # List up to 5 messages
  sqs-redrive messages https://sqs.us-east-1.amazonaws.com/123456789012/my-dlq --max 5`,
	Args: cobra.ExactArgs(1),
	RunE: runMessages,
}

func init() {
	rootCmd.AddCommand(messagesCmd)

	messagesCmd.Flags().Int32VarP(&messagesMaxCount, "max", "m", 10, "Maximum number of messages to retrieve (1-10)")
}

func runMessages(cmd *cobra.Command, args []string) error {
	queueURL := args[0]
	ctx := context.Background()
	client := GetSQSClient()
	printer := output.NewPrinter()

	// Validate max count
	if messagesMaxCount < 1 {
		messagesMaxCount = 1
	}
	if messagesMaxCount > 10 {
		messagesMaxCount = 10
	}

	// Get queue info for display
	queueInfo, err := client.GetQueueAttributes(ctx, queueURL)
	if err != nil {
		return fmt.Errorf("failed to get queue attributes: %w", err)
	}

	printer.Printf("Queue: %s\n", queueInfo.Name)
	printer.Printf("Approximate Messages: %d\n\n", queueInfo.ApproximateMessages)

	if queueInfo.ApproximateMessages == 0 {
		printer.Println("No messages in queue.")
		return nil
	}

	// Receive messages with a short visibility timeout (30 seconds)
	// This allows messages to be quickly returned to the queue
	messages, err := client.ReceiveMessages(ctx, queueURL, messagesMaxCount, 30)
	if err != nil {
		return fmt.Errorf("failed to receive messages: %w", err)
	}

	if len(messages) == 0 {
		printer.Println("No messages available to receive.")
		printer.Println("(Messages may be in-flight or the queue may be empty)")
		return nil
	}

	// Prepare table data
	headers := []string{"MESSAGE ID", "SENT", "RECEIVES", "PREVIEW"}
	rows := make([][]string, 0, len(messages))

	for _, msg := range messages {
		sentTime := output.FormatTimestamp(msg.SentTimestamp)
		preview := output.FormatMessagePreview(msg.Body, 50)

		rows = append(rows, []string{
			output.TruncateString(msg.MessageID, 36),
			sentTime,
			msg.ApproximateReceiveCount,
			preview,
		})
	}

	// Print table
	printer.Table(headers, rows)

	// Print summary and tips
	printer.Printf("\nShowing %d message(s)\n", len(messages))
	printer.Println("")
	printer.Println("Tips:")
	printer.Println("  - Use 'sqs-redrive peek <queue-url> <message-id>' to see full message content")
	printer.Println("  - Use 'sqs-redrive redrive <queue-url> --message-id <id>' to redrive a specific message")
	printer.Println("  - Use 'sqs-redrive redrive <queue-url> --all' to redrive all messages")

	return nil
}
