package main

import (
	"os"

	"github.com/dalpark/sqs-redrive/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
