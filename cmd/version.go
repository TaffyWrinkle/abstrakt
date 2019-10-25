package cmd

import (
	"github.com/microsoft/abstrakt/internal/tool"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "The version of Abstrakt being used",
	Long:  "The version of Abstrakt being used",
	Run: func(cmd *cobra.Command, args []string) {
		PrintVersion()
	},
}

// PrintVersion prints the current version of Abstrakt being used.
func PrintVersion() {
	logger.Info("abstrakt version 0.0.1")
}

func init() {
	rootCmd.AddCommand(versionCmd)
}