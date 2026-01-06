package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	redriveMessageID string
	redriveAll       bool
	redriveDryRun    bool
)

var redriveCmd = &cobra.Command{
	Use:   "redrive <queue-url>",
	Short: "Redrive messages from a DLQ back to the source queue",
	Long: `Redrive messages from a Dead Letter Queue (DLQ) back to their source queue.

You can redrive a single message by specifying --message-id, or redrive all
messages using --all. Use --dry-run to preview what would happen without
actually moving any messages.

The command automatically detects the source queue by scanning queue
redrive policies.`,
	Args: cobra.ExactArgs(1),
	RunE: runRedrive,
}

func init() {
	redriveCmd.Flags().StringVar(&redriveMessageID, "message-id", "", "Redrive a specific message by ID")
	redriveCmd.Flags().BoolVar(&redriveAll, "all", false, "Redrive all messages in the DLQ")
	redriveCmd.Flags().BoolVar(&redriveDryRun, "dry-run", false, "Preview actions without executing them")

	rootCmd.AddCommand(redriveCmd)
}

func runRedrive(cmd *cobra.Command, args []string) error {
	dlqURL := args[0]

	// Validate flags - must specify either --message-id or --all
	if redriveMessageID == "" && !redriveAll {
		return fmt.Errorf("must specify either --message-id or --all")
	}
	if redriveMessageID != "" && redriveAll {
		return fmt.Errorf("cannot use both --message-id and --all")
	}

	ctx := context.Background()
	client := GetSQSClient()

	// Find the source queue
	fmt.Println("Finding source queue...")
	sourceQueueURL, err := client.GetSourceQueueURL(ctx, dlqURL)
	if err != nil {
		return fmt.Errorf("failed to find source queue: %w", err)
	}

	cyan := color.New(color.FgCyan)
	cyan.Print("Source Queue: ")
	fmt.Println(sourceQueueURL)
	fmt.Println()

	if redriveMessageID != "" {
		return redriveSingleMessage(ctx, dlqURL, sourceQueueURL, redriveMessageID, redriveDryRun)
	}

	return redriveAllMessages(ctx, dlqURL, sourceQueueURL, redriveDryRun)
}

func redriveSingleMessage(ctx context.Context, dlqURL, sourceQueueURL, messageID string, dryRun bool) error {
	client := GetSQSClient()

	fmt.Printf("Searching for message %s...\n", messageID)

	// Receive messages to find the target
	maxAttempts := 10
	messagesChecked := 0

	for attempt := 0; attempt < maxAttempts; attempt++ {
		messages, err := client.ReceiveMessages(ctx, dlqURL, 10, 30)
		if err != nil {
			return fmt.Errorf("failed to receive messages: %w", err)
		}

		if len(messages) == 0 {
			break
		}

		for _, msg := range messages {
			messagesChecked++
			if msg.MessageID == messageID {
				if dryRun {
					color.Yellow("[DRY RUN] Would redrive message %s to %s", messageID, sourceQueueURL)
					return nil
				}

				// Redrive: send to source queue, then delete from DLQ
				err := client.RedriveMessage(ctx, dlqURL, sourceQueueURL, msg)
				if err != nil {
					return fmt.Errorf("failed to redrive message: %w", err)
				}

				color.Green("Successfully redrove message %s", messageID)
				return nil
			}
		}
	}

	return fmt.Errorf("message with ID %s not found (checked %d messages)", messageID, messagesChecked)
}

func redriveAllMessages(ctx context.Context, dlqURL, sourceQueueURL string, dryRun bool) error {
	client := GetSQSClient()

	// Get queue attributes to know approximate message count
	queueInfo, err := client.GetQueueAttributes(ctx, dlqURL)
	if err != nil {
		return fmt.Errorf("failed to get queue attributes: %w", err)
	}

	totalMessages := queueInfo.ApproximateMessages
	if totalMessages == 0 {
		fmt.Println("No messages in the queue.")
		return nil
	}

	fmt.Printf("Found approximately %d messages to redrive.\n\n", totalMessages)

	if dryRun {
		color.Yellow("[DRY RUN] Would redrive %d messages from DLQ to source queue", totalMessages)
		color.Yellow("[DRY RUN] DLQ: %s", dlqURL)
		color.Yellow("[DRY RUN] Source: %s", sourceQueueURL)
		return nil
	}

	// Create progress bar
	bar := progressbar.NewOptions(totalMessages,
		progressbar.OptionSetDescription("Redriving messages"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	successCount := 0
	errorCount := 0
	emptyReceives := 0
	maxEmptyReceives := 3

	for emptyReceives < maxEmptyReceives {
		messages, err := client.ReceiveMessages(ctx, dlqURL, 10, 30)
		if err != nil {
			return fmt.Errorf("failed to receive messages: %w", err)
		}

		if len(messages) == 0 {
			emptyReceives++
			continue
		}
		emptyReceives = 0

		for _, msg := range messages {
			err := client.RedriveMessage(ctx, dlqURL, sourceQueueURL, msg)
			if err != nil {
				errorCount++
				fmt.Fprintf(os.Stderr, "\nError redriving message %s: %v\n", msg.MessageID, err)
			} else {
				successCount++
			}
			bar.Add(1)
		}
	}

	fmt.Println()
	fmt.Println()

	color.Green("Redrive completed!")
	fmt.Printf("  Successful: %d\n", successCount)
	if errorCount > 0 {
		color.Red("  Failed: %d", errorCount)
	}

	return nil
}
