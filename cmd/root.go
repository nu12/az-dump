/*
Copyright © 2025 nu12
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var subscriptionID string
var rgName string

var rootCmd = &cobra.Command{
	Use:   "az-dump",
	Short: "Export Azure ARM templates to create local backups with backup restore capabilities",
	Long:  `Export Azure ARM templates to create local backups with backup restore capabilities.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(versionCmd)
}
