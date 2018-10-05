package cmd

import (
	"fmt"
	"os"
)

// Execute Root Command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
