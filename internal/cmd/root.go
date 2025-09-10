// Package cmd implements the command line interface for Neuron CLI.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "neuron",
	Short: "Neuron CLI is your smart study partner for Markdown notes.",
	Long: `A powerful, evidence-based learning tool for the command line.
Neuron CLI helps you learn and retain knowledge from your notes
by using spaced repetition, active recall, and AI-powered questioning.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// The init function in root.go should be empty.
// Each command file (e.g., review.go, import.go) is responsible
// for adding itself to the rootCmd in its own init() function.
func init() {
	// Intentionally empty.
}
