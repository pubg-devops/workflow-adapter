package cmd

import (
	"bytes"
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Test that root command exists and has correct use string
	if rootCmd.Use != "sqs-redrive" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "sqs-redrive")
	}

	// Test that Short description is set
	if rootCmd.Short == "" {
		t.Error("rootCmd.Short should not be empty")
	}

	// Test that Long description is set
	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}
}

func TestGlobalFlags(t *testing.T) {
	// Test --profile flag exists
	profileFlag := rootCmd.PersistentFlags().Lookup("profile")
	if profileFlag == nil {
		t.Error("--profile flag should exist")
	} else {
		if profileFlag.Usage == "" {
			t.Error("--profile flag should have usage text")
		}
	}

	// Test --region flag exists
	regionFlag := rootCmd.PersistentFlags().Lookup("region")
	if regionFlag == nil {
		t.Error("--region flag should exist")
	} else {
		if regionFlag.Usage == "" {
			t.Error("--region flag should have usage text")
		}
	}
}

func TestRootCommandHelp(t *testing.T) {
	// Test that help doesn't return an error
	rootCmd.SetArgs([]string{"--help"})
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetErr(new(bytes.Buffer))

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("help command should not return error, got: %v", err)
	}
}

func TestExitWithError(t *testing.T) {
	// This function calls os.Exit, so we can't easily test it directly
	// We just verify the function exists and has the correct signature
	// In a real scenario, you might use a mock or dependency injection
	_ = exitWithError
}

func TestGetSQSClient(t *testing.T) {
	// Before initialization, client should be nil
	// This tests the getter function exists
	_ = GetSQSClient
}
