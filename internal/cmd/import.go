// Package cmd implements the command line interface for Neuron CLI.
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/soyomarvaldezg/neuron-cli/internal/db"
	"github.com/soyomarvaldezg/neuron-cli/internal/note"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [path]",
	Short: "Import and sync notes from a directory",
	Long: `Imports notes from a specified directory of Markdown files.
The command will intelligently sync your notes, adding new ones,
and updating modified ones based on filename.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		notesPath := args[0]
		fmt.Printf("Starting import from directory: %s\n", notesPath)

		// Get a database connection
		database, err := db.GetDB("neuron.db")
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		importedCount := 0
		// Walk the directory
		err = filepath.Walk(notesPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// We only care about markdown files, not directories or other files
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
				// Parse the file
				parsedNote, err := note.ParseFile(path)
				if err != nil {
					log.Printf("Error parsing %s: %v. Skipping.", path, err)
					return nil // Continue walking
				}

				// Insert into database
				err = db.InsertNote(database, parsedNote)
				if err != nil {
					log.Printf("Error inserting %s into DB: %v. Skipping.", path, err)
					return nil // Continue walking
				}
				fmt.Printf("✓ Synced: %s\n", parsedNote.Title)
				importedCount++
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("error walking the path %q: %w", notesPath, err)
		}

		fmt.Printf("\nSync complete. Processed %d notes.\n", importedCount)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
