package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/dalpark/sqs-redrive/internal/aws"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	awsProfile string
	awsRegion  string

	// SQS client (initialized in PersistentPreRunE)
	sqsClient *aws.SQSClient
)

var rootCmd = &cobra.Command{
	Use:   "sqs-redrive",
	Short: "A CLI tool to redrive messages from SQS Dead Letter Queues",
	Long: `sqs-redrive is a CLI tool that helps you manage and redrive messages
from AWS SQS Dead Letter Queues (DLQ) back to their source queues.

It supports listing DLQs, viewing messages, and redriving messages
either individually or in bulk.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip client initialization for help commands
		if cmd.Name() == "help" || cmd.Name() == "completion" {
			return nil
		}

		// Initialize AWS config options
		var configOpts []func(*config.LoadOptions) error

		if awsProfile != "" {
			configOpts = append(configOpts, config.WithSharedConfigProfile(awsProfile))
		}

		if awsRegion != "" {
			configOpts = append(configOpts, config.WithRegion(awsRegion))
		}

		// Create SQS client
		client, err := aws.NewClient(context.Background(), configOpts...)
		if err != nil {
			return fmt.Errorf("failed to initialize AWS client: %w", err)
		}
		sqsClient = client

		return nil
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&awsProfile, "profile", "", "AWS profile to use (overrides AWS_PROFILE)")
	rootCmd.PersistentFlags().StringVar(&awsRegion, "region", "", "AWS region to use (overrides AWS_REGION)")
}

// GetSQSClient returns the initialized SQS client
func GetSQSClient() *aws.SQSClient {
	return sqsClient
}

// exitWithError prints an error message and exits with code 1
func exitWithError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}
