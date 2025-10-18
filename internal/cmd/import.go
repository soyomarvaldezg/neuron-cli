// Package cmd implements the command line interface for Neuron CLI.
package cmd

import (
	"database/sql"
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
updating modified ones, and removing deleted ones based on filename.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		notesPath := args[0]
		fmt.Printf("Starting import from directory: %s\n", notesPath)

		// Get a database connection
		database, err := db.GetDB()
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		// Track which files we found during this import
		foundFiles := make(map[string]bool)
		importedCount := 0

		// Walk the directory
		err = filepath.Walk(notesPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// We only care about markdown files, not directories or other files
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
				// Mark this file as found
				foundFiles[path] = true

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

		// Now clean up deleted notes
		deletedCount, err := cleanupDeletedNotes(database, foundFiles)
		if err != nil {
			return fmt.Errorf("error cleaning up deleted notes: %w", err)
		}

		fmt.Printf("\nSync complete. Processed %d notes.", importedCount)
		if deletedCount > 0 {
			fmt.Printf(" Removed %d deleted notes.", deletedCount)
		}
		fmt.Println()

		return nil
	},
}

// cleanupDeletedNotes removes database entries for files that no longer exist
func cleanupDeletedNotes(database *sql.DB, foundFiles map[string]bool) (int, error) {
	// Get all filenames currently in the database
	query := `SELECT filename FROM notes;`
	rows, err := database.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var toDelete []string
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return 0, err
		}

		// If this file wasn't found during our walk, it's been deleted
		if !foundFiles[filename] {
			toDelete = append(toDelete, filename)
		}
	}

	// Delete the orphaned entries
	deletedCount := 0
	for _, filename := range toDelete {
		deleteQuery := `DELETE FROM notes WHERE filename = ?;`
		_, err := database.Exec(deleteQuery, filename)
		if err != nil {
			log.Printf("Error deleting %s from database: %v", filename, err)
			continue
		}
		fmt.Printf("✗ Removed: %s\n", filepath.Base(filename))
		deletedCount++
	}

	return deletedCount, nil
}

func init() {
	rootCmd.AddCommand(importCmd)
}
