package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/dalpark/sqs-redrive/internal/aws"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var peekCmd = &cobra.Command{
	Use:   "peek <queue-url> <message-id>",
	Short: "View the full content of a specific message",
	Long: `View the full content of a specific message in an SQS queue.

This command retrieves and displays the complete message body and attributes
for the specified message ID. The message content is pretty-printed if it's
valid JSON.

Note: This command receives messages from the queue, so the visibility timeout
will be applied. The message will become visible again after the timeout expires.`,
	Args: cobra.ExactArgs(2),
	RunE: runPeek,
}

func init() {
	rootCmd.AddCommand(peekCmd)
}

func runPeek(cmd *cobra.Command, args []string) error {
	queueURL := args[0]
	targetMessageID := args[1]

	ctx := context.Background()
	client := GetSQSClient()

	// We need to receive messages to find the one with the specified ID
	// SQS doesn't support getting a specific message by ID directly
	fmt.Printf("Searching for message %s...\n\n", targetMessageID)

	// Receive messages in batches to find the target message
	// We'll keep receiving until we find it or exhaust visible messages
	maxAttempts := 10
	messagesChecked := 0

	for attempt := 0; attempt < maxAttempts; attempt++ {
		messages, err := client.ReceiveMessages(ctx, queueURL, 10, 30)
		if err != nil {
			return fmt.Errorf("failed to receive messages: %w", err)
		}

		if len(messages) == 0 {
			break
		}

		for _, msg := range messages {
			messagesChecked++
			if msg.MessageID == targetMessageID {
				return displayMessage(msg)
			}
		}
	}

	return fmt.Errorf("message with ID %s not found (checked %d messages)", targetMessageID, messagesChecked)
}

func displayMessage(msg aws.MessageInfo) error {
	bold := color.New(color.Bold)
	cyan := color.New(color.FgCyan)
	yellow := color.New(color.FgYellow)

	bold.Println("=== Message Details ===")
	fmt.Println()

	cyan.Print("Message ID: ")
	fmt.Println(msg.MessageID)

	cyan.Print("Receipt Handle: ")
	fmt.Printf("%.50s...\n", msg.ReceiptHandle)

	if msg.SentTimestamp != "" {
		cyan.Print("Sent: ")
		ts, err := strconv.ParseInt(msg.SentTimestamp, 10, 64)
		if err == nil {
			t := time.UnixMilli(ts)
			fmt.Println(t.Format(time.RFC3339))
		} else {
			fmt.Println(msg.SentTimestamp)
		}
	}

	if msg.ApproximateFirstReceiveTimestamp != "" {
		cyan.Print("First Received: ")
		ts, err := strconv.ParseInt(msg.ApproximateFirstReceiveTimestamp, 10, 64)
		if err == nil {
			t := time.UnixMilli(ts)
			fmt.Println(t.Format(time.RFC3339))
		} else {
			fmt.Println(msg.ApproximateFirstReceiveTimestamp)
		}
	}

	if msg.ApproximateReceiveCount != "" {
		cyan.Print("Receive Count: ")
		fmt.Println(msg.ApproximateReceiveCount)
	}

	// Display system attributes
	if len(msg.Attributes) > 0 {
		fmt.Println()
		yellow.Println("--- System Attributes ---")
		for key, value := range msg.Attributes {
			// Skip ones we already displayed
			if key == "SentTimestamp" || key == "ApproximateFirstReceiveTimestamp" || key == "ApproximateReceiveCount" {
				continue
			}
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	// Display message attributes
	if len(msg.MessageAttributes) > 0 {
		fmt.Println()
		yellow.Println("--- Message Attributes ---")
		for key, attr := range msg.MessageAttributes {
			if attr.StringValue != nil {
				fmt.Printf("  %s (%s): %s\n", key, *attr.DataType, *attr.StringValue)
			} else if attr.BinaryValue != nil {
				fmt.Printf("  %s (%s): [binary data, %d bytes]\n", key, *attr.DataType, len(attr.BinaryValue))
			}
		}
	}

	// Display message body
	fmt.Println()
	bold.Println("--- Message Body ---")

	// Try to pretty-print if JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(msg.Body), &jsonData); err == nil {
		prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
		if err == nil {
			fmt.Println(string(prettyJSON))
		} else {
			fmt.Println(msg.Body)
		}
	} else {
		fmt.Println(msg.Body)
	}

	return nil
}
