package cmd

import (
	"context"
	"fmt"

	"github.com/dalpark/sqs-redrive/internal/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Dead Letter Queues (DLQs)",
	Long: `List all Dead Letter Queues (DLQs) in the AWS account.

Displays each DLQ with its name, URL, and approximate message count.
DLQs are identified by having "dlq", "dead-letter", or "deadletter" in their name.`,
	Example: `  # List all DLQs
  sqs-redrive list

  # List DLQs using a specific AWS profile
  sqs-redrive list --profile production

  # List DLQs in a specific region
  sqs-redrive list --region us-west-2`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client := GetSQSClient()
	printer := output.NewPrinter()

	// List all DLQs
	dlqs, err := client.ListDLQs(ctx)
	if err != nil {
		return fmt.Errorf("failed to list DLQs: %w", err)
	}

	if len(dlqs) == 0 {
		printer.Println("No Dead Letter Queues found.")
		printer.Println("")
		printer.Println("DLQs are identified by having 'dlq', 'dead-letter', or 'deadletter' in their name.")
		return nil
	}

	// Prepare table data
	headers := []string{"NAME", "MESSAGES", "URL"}
	rows := make([][]string, 0, len(dlqs))

	for _, dlq := range dlqs {
		rows = append(rows, []string{
			dlq.Name,
			fmt.Sprintf("%d", dlq.ApproximateMessages),
			dlq.URL,
		})
	}

	// Print table
	printer.Table(headers, rows)

	// Print summary
	printer.Printf("\nFound %d Dead Letter Queue(s)\n", len(dlqs))

	return nil
}
